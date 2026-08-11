//go:build acceptance && !windows

package acceptance

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// scanRoot is a writable location for planted installations.
func scanRoot() string { return "/opt" }

// searchableSystemDir is covered by the default search paths, so anything
// planted here is met without passing extra flags.
func searchableSystemDir() string { return "/opt" }

func removeStubbornly(path string) { os.RemoveAll(path) }

// unprivilegedUser returns an account that cannot read root-owned directories,
// which is the only way to exercise the permission branch: root reads
// everything, so running the suite as root would silently skip it.
func unprivilegedUser() string {
	if name := os.Getenv("SUDO_USER"); name != "" && name != "root" {
		return name
	}

	homes, err := os.ReadDir("/home")
	if err != nil || len(homes) == 0 {
		return ""
	}
	return homes[0].Name()
}

// A directory the scan may not read must be reported with its cause and must
// not stop the scan.
func TestPermissionFailureIsFlaggedAndSurvived(t *testing.T) {
	jdk := requireEnv(t, envJDK)
	binary := requireEnv(t, envBinary)

	denied := filepath.Join(scanRoot(), "aaa_denied")
	if err := os.MkdirAll(filepath.Join(denied, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(denied, 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		os.Chmod(denied, 0o755)
		os.RemoveAll(denied)
	})

	// The report and log must be writable by whoever runs the collector.
	workDir, err := os.MkdirTemp("", "ojdm-unpriv")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(workDir) })
	if err := os.Chmod(workDir, 0o777); err != nil {
		t.Fatal(err)
	}

	reportPath := filepath.Join(workDir, "report.csv")
	logPath := filepath.Join(workDir, "logs", "debug.log")
	args := []string{
		"-search-paths=" + denied + "," + jdk,
		"-output-path=" + reportPath,
	}

	var cmd *exec.Cmd
	if os.Geteuid() == 0 {
		user := unprivilegedUser()
		if user == "" {
			t.Skip("no unprivileged account available to exercise the permission branch")
		}
		cmd = exec.Command("su", user, "-c", binary+" "+strings.Join(args, " "))
	} else {
		cmd = exec.Command(binary, args...)
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("collector failed: %v\n%s", err, output)
	}

	var flagged bool
	for _, entry := range logEntriesWithMessage(readLog(t, logPath), "skipping unreadable path") {
		path, _ := entry["path"].(string)
		refused, _ := entry["permission_denied"].(bool)
		if strings.Contains(path, "aaa_denied") && refused {
			flagged = true
		}
	}
	if !flagged {
		t.Error("the unreadable directory was not reported as a permission failure")
	}

	if len(readReport(t, reportPath)) == 0 {
		t.Error("the scan stopped at the permission failure instead of continuing")
	}
}
