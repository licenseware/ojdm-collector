//go:build acceptance

package acceptance

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// Environment contract, populated by the platform planting script.
const (
	envBinary    = "OJDM_BINARY"     // collector under test
	envJDK       = "OJDM_JDK"        // a JDK holding javac, jps and jinfo
	envJRE       = "OJDM_JRE"        // optional JRE-only installation
	envExtraRoot = "OJDM_EXTRA_ROOT" // optional second drive or mount point
)

type report []map[string]string

func requireEnv(t *testing.T, name string) string {
	t.Helper()

	value := os.Getenv(name)
	if value == "" {
		t.Fatalf("%s must be set; run the planting script for this platform first", name)
	}
	return value
}

func executable(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

// runCollector executes the collector and returns everything it wrote to the
// console, so tests can assert on what the operator actually sees.
func runCollector(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command(requireEnv(t, envBinary), args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("collector failed: %v\n%s", err, output)
	}
	return string(output)
}

func readReport(t *testing.T, path string) report {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("report unreadable: %v", err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		t.Fatalf("report is not valid csv: %v", err)
	}
	if len(records) == 0 {
		t.Fatal("report has no header")
	}

	var rows report
	for _, record := range records[1:] {
		row := map[string]string{}
		for i, header := range records[0] {
			if i < len(record) {
				row[header] = record[i]
			}
		}
		rows = append(rows, row)
	}
	return rows
}

// readLog parses the debug log, which doubles as the assertion that every line
// is machine readable for the support hand-over.
func readLog(t *testing.T, path string) []map[string]any {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("debug log unreadable: %v", err)
	}

	var entries []map[string]any
	for i, line := range strings.Split(strings.TrimSpace(string(content)), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("debug log line %d is not json: %v\n%s", i+1, err, line)
		}
		entries = append(entries, entry)
	}
	return entries
}

func logEntriesWithMessage(entries []map[string]any, message string) []map[string]any {
	var matched []map[string]any
	for _, entry := range entries {
		if entry["message"] == message {
			matched = append(matched, entry)
		}
	}
	return matched
}

func copyTree(t *testing.T, source, destination string) {
	t.Helper()

	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, relative)

		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		// A JDK is full of symlinks (lib/security/cacerts and friends) and the
		// runtime will not start without them, so recreate rather than skip.
		if info.Mode()&os.ModeSymlink != 0 {
			destination, err := os.Readlink(path)
			if err != nil {
				return err
			}
			os.Remove(target)
			return os.Symlink(destination, target)
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, content, info.Mode().Perm())
	})
	if err != nil {
		t.Fatalf("copying %s: %v", source, err)
	}
}

// plantUnreadableEntry creates something the scan cannot stat, which is not a
// permission error. This is the failure that used to abort the whole walk and
// silently truncate the report.
//
// Windows: NTFS accepts a trailing space through the \\?\ prefix, while Win32
// path resolution strips it, so the directory listing yields a name that
// cannot be opened. POSIX: a path longer than PATH_MAX, built relatively so
// that creating it succeeds.
func plantUnreadableEntry(t *testing.T, root string) {
	t.Helper()

	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatal(err)
	}

	if runtime.GOOS == "windows" {
		for _, name := range []string{"ghost ", "ghost."} {
			if err := os.WriteFile(`\\?\`+filepath.Join(root, name), []byte("x"), 0o644); err != nil {
				t.Fatalf("planting %q: %v", name, err)
			}
		}
		return
	}

	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(previous) })

	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	segment := strings.Repeat("d", 200)
	for i := 0; i < 30; i++ {
		if err := os.Mkdir(segment, 0o755); err != nil {
			break
		}
		if err := os.Chdir(segment); err != nil {
			break
		}
	}
	if err := os.Chdir(previous); err != nil {
		t.Fatal(err)
	}
}

// startLiveJVM compiles a trivial class and keeps it running, so the jps and
// jinfo path has something to discover.
func startLiveJVM(t *testing.T, jdk, workDir string) {
	t.Helper()

	source := filepath.Join(workDir, "Sleeper.java")
	program := "public class Sleeper { public static void main(String[] a) throws Exception { Thread.sleep(600000); } }"
	if err := os.WriteFile(source, []byte(program), 0o644); err != nil {
		t.Fatal(err)
	}

	javac := filepath.Join(jdk, "bin", executable("javac"))
	if output, err := exec.Command(javac, "-d", workDir, source).CombinedOutput(); err != nil {
		t.Fatalf("compiling the live jvm helper: %v\n%s", err, output)
	}

	cmd := exec.Command(filepath.Join(jdk, "bin", executable("java")), "-Xmx64m", "-cp", workDir, "Sleeper")
	if err := cmd.Start(); err != nil {
		t.Fatalf("starting the live jvm: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
	})

	time.Sleep(3 * time.Second) // let the vm register with the jvmstat directory
}

// scanFixture is the planted host state plus the result of scanning it once.
type scanFixture struct {
	workDir    string
	reportPath string
	logPath    string
	console    string
	rows       report
	entries    []map[string]any
	expected   []string // java homes that must appear in the report
	brokenRoot string
}

// newScanFixture plants every layout the fixes target, then scans once. The
// planted state lives outside t.TempDir because the collector's default search
// paths are absolute system locations.
func newScanFixture(t *testing.T) *scanFixture {
	t.Helper()

	jdk := requireEnv(t, envJDK)
	workDir := t.TempDir()

	// The scan must reach installations under a root whose own name contains
	// "bin", which used to be truncated away.
	binNamedRoot := filepath.Join(scanRoot(), "bin-tools")
	planted := filepath.Join(binNamedRoot, "jdk-planted")
	if err := os.RemoveAll(binNamedRoot); err != nil {
		t.Fatal(err)
	}
	copyTree(t, jdk, planted)
	t.Cleanup(func() { os.RemoveAll(binNamedRoot) })

	brokenRoot := filepath.Join(searchableSystemDir(), "AAA_broken")
	plantUnreadableEntry(t, brokenRoot)
	t.Cleanup(func() { removeStubbornly(brokenRoot) })

	expected := []string{planted, jdk}
	if jre := os.Getenv(envJRE); jre != "" {
		expected = append(expected, jre)
	}

	extraRoot := os.Getenv(envExtraRoot)
	searchPaths := []string{binNamedRoot}
	if extraRoot != "" {
		extraInstall := filepath.Join(extraRoot, "jdk-planted")
		if err := os.RemoveAll(extraInstall); err != nil {
			t.Fatal(err)
		}
		copyTree(t, jdk, extraInstall)
		t.Cleanup(func() { os.RemoveAll(extraInstall) })
		searchPaths = append(searchPaths, extraRoot)
		expected = append(expected, extraInstall)
	} else {
		t.Logf("%s not set: skipping the second drive or mount point", envExtraRoot)
	}

	startLiveJVM(t, jdk, workDir)

	fixture := &scanFixture{
		workDir:    workDir,
		reportPath: filepath.Join(workDir, "report.csv"),
		logPath:    filepath.Join(workDir, "logs", "debug.log"),
		expected:   expected,
		brokenRoot: brokenRoot,
	}
	fixture.console = runCollector(t,
		"-search-paths="+strings.Join(searchPaths, ","),
		"-output-path="+fixture.reportPath)
	fixture.rows = readReport(t, fixture.reportPath)
	fixture.entries = readLog(t, fixture.logPath)

	return fixture
}

func TestReportCoversEveryInstallation(t *testing.T) {
	fixture := newScanFixture(t)

	for _, javaHome := range fixture.expected {
		if !fixture.reports(javaHome) {
			t.Errorf("installation %q missing from the report; rows: %v", javaHome, fixture.javaHomes())
			// Without the collector's own account of the scan there is no way
			// to tell a discovery failure from an execution failure.
			t.Logf("collector console:\n%s", fixture.console)
		}
	}
}

// One installation reachable through several search paths, such as a symlinked
// /usr/bin/java beside the real location, must not be counted twice.
func TestReportHasNoDuplicateInstallations(t *testing.T) {
	fixture := newScanFixture(t)

	seen := map[string]int{}
	for _, home := range fixture.javaHomes() {
		seen[home]++
	}
	for home, count := range seen {
		if count > 1 {
			t.Errorf("installation %q reported %d times", home, count)
		}
	}
}

func TestEveryRowResolvesItsPaths(t *testing.T) {
	fixture := newScanFixture(t)

	for _, row := range fixture.rows {
		if row["DynLibBinPath"] == "" {
			t.Errorf("row %q has no vm shared library", row["JavaHome"])
		} else if _, err := os.Stat(row["DynLibBinPath"]); err != nil {
			t.Errorf("vm shared library %q does not exist: %v", row["DynLibBinPath"], err)
		}

		if row["IsJDK"] == "true" {
			if _, err := os.Stat(row["JavaCBinPath"]); err != nil {
				t.Errorf("javac %q does not exist: %v", row["JavaCBinPath"], err)
			}
		}

		if processors, _ := strconv.Atoi(row["HostLogicalProcessors"]); processors < 1 {
			t.Errorf("row %q reports %q logical processors", row["JavaHome"], row["HostLogicalProcessors"])
		}
		if row["HostName"] == "" {
			t.Errorf("row %q has no hostname", row["JavaHome"])
		}
	}
}

func TestRunningJVMIsIdentified(t *testing.T) {
	fixture := newScanFixture(t)

	for _, row := range fixture.rows {
		if row["ProcessRunning"] == "true" {
			return
		}
	}
	t.Error("a jvm was running during the scan but no row reports it")
}

// An entry the scan cannot read must be reported and must not stop the scan:
// this is the defect that silently produced blank reports on customer hosts.
func TestUnreadableEntryIsReportedAndSurvived(t *testing.T) {
	fixture := newScanFixture(t)

	var found bool
	for _, entry := range logEntriesWithMessage(fixture.entries, "skipping unreadable path") {
		if path, _ := entry["path"].(string); strings.Contains(path, "AAA_broken") {
			found = true
		}
	}
	if !found {
		t.Fatal("the unreadable entry was not reported in the debug log")
	}

	if len(fixture.rows) == 0 {
		t.Fatal("the scan produced no rows after meeting an unreadable entry")
	}
}

func TestDebugLogRecordsTheRun(t *testing.T) {
	fixture := newScanFixture(t)

	if got := len(logEntriesWithMessage(fixture.entries, "scanning for java installations")); got != 1 {
		t.Errorf("expected exactly one search path entry, got %d", got)
	}

	var hasDebug bool
	for _, entry := range fixture.entries {
		if entry["level"] == "debug" {
			hasDebug = true
		}
	}
	if !hasDebug {
		t.Error("the debug log holds no debug level detail")
	}

	last := fixture.entries[len(fixture.entries)-1]
	if message, _ := last["message"].(string); !strings.HasPrefix(message, "done") {
		t.Errorf("the run does not finish by pointing at the log, got %q", message)
	}

	if strings.Contains(fixture.console, "found java file") {
		t.Error("debug detail leaked onto the console at the default level")
	}
}

func TestDebugLevelSurfacesDetailOnTheConsole(t *testing.T) {
	workDir := t.TempDir()

	console := runCollector(t,
		"-search-paths="+requireEnv(t, envJDK),
		"-output-path="+filepath.Join(workDir, "report.csv"),
		"-log-level=debug")

	if !strings.Contains(console, "found java file") {
		t.Error("-log-level=debug did not surface debug detail on the console")
	}
}

func TestLogPathOverrideIsHonoured(t *testing.T) {
	workDir := t.TempDir()
	custom := filepath.Join(workDir, "custom", "my.log")

	runCollector(t,
		"-search-paths="+requireEnv(t, envJDK),
		"-output-path="+filepath.Join(workDir, "report.csv"),
		"-log-path="+custom)

	if _, err := os.Stat(custom); err != nil {
		t.Fatalf("-log-path was not honoured: %v", err)
	}
}

// A rerun must not destroy the log of the run that went wrong, which is the
// one support asks for.
func TestRerunArchivesThePreviousLog(t *testing.T) {
	workDir := t.TempDir()
	reportPath := filepath.Join(workDir, "report.csv")
	logPath := filepath.Join(workDir, "logs", "debug.log")

	runCollector(t, "-search-paths="+requireEnv(t, envJDK), "-output-path="+reportPath)
	first, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1100 * time.Millisecond) // the archive is named by whole seconds
	runCollector(t, "-search-paths="+requireEnv(t, envJDK), "-output-path="+reportPath)

	archives, err := filepath.Glob(filepath.Join(workDir, "logs", "debug-*.log"))
	if err != nil {
		t.Fatal(err)
	}
	if len(archives) != 1 {
		t.Fatalf("expected exactly one archived log, got %d", len(archives))
	}

	archived, err := os.ReadFile(archives[0])
	if err != nil {
		t.Fatal(err)
	}
	if string(archived) != string(first) {
		t.Error("the archived log does not hold the previous run's content")
	}

	current, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(current) == string(archived) {
		t.Error("the new run did not start a fresh log")
	}
}

func (f *scanFixture) javaHomes() []string {
	var homes []string
	for _, row := range f.rows {
		homes = append(homes, row["JavaHome"])
	}
	return homes
}

// reports compares java homes tolerantly: the collector normalises separators,
// and a JDK may report its home one level below the installation root.
func (f *scanFixture) reports(javaHome string) bool {
	wanted := normalise(javaHome)
	for _, home := range f.javaHomes() {
		got := normalise(home)
		if got == wanted || strings.HasPrefix(got, wanted+"/") {
			return true
		}
	}
	return false
}

func normalise(path string) string {
	return strings.TrimSuffix(strings.ToLower(filepath.ToSlash(path)), "/")
}
