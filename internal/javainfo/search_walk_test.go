package javainfo

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

// makeUnstatablePath nests directories until the absolute path exceeds
// PATH_MAX, so filepath.Walk gets ENAMETOOLONG from Lstat. Each level is
// created relative to the previous one, otherwise the creation itself would
// fail with the same error. This stands in for the non-permission errors a
// Windows scan hits in practice (broken junctions, cloud placeholders, files
// removed mid-walk by an EDR agent).
func makeUnstatablePath(t *testing.T, parent string) {
	t.Helper()

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(cwd); err != nil {
			t.Fatal(err)
		}
	}()

	if err := os.Chdir(parent); err != nil {
		t.Fatal(err)
	}

	segment := strings.Repeat("d", 200)
	for i := 0; i < 30; i++ {
		if err := os.Mkdir(segment, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(segment); err != nil {
			t.Fatal(err)
		}
	}
}

func writeJavaBinary(t *testing.T, root, installDir string) string {
	t.Helper()

	binDir := filepath.Join(root, installDir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}

	name := "java"
	if runtime.GOOS == "windows" {
		name = "java.exe"
	}
	if err := os.WriteFile(filepath.Join(binDir, name), []byte("stub"), 0o755); err != nil {
		t.Fatal(err)
	}

	return normalizeSearchPath(filepath.Join(root, installDir))
}

// A single unreadable entry must not stop the scan: the walk callback returns
// the error, which aborts filepath.Walk for the whole search root, so every
// installation sorting after it is silently dropped from the report.
func TestWalkForJavaFilesContinuesAfterStatError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("long-path construction is POSIX specific")
	}

	root := t.TempDir()

	// "aaa_broken" sorts before "zzz_java", so the walk hits the error first.
	broken := filepath.Join(root, "aaa_broken")
	if err := os.Mkdir(broken, 0o755); err != nil {
		t.Fatal(err)
	}
	makeUnstatablePath(t, broken)

	want := writeJavaBinary(t, root, "zzz_java")

	log := zerolog.Nop()
	got := walkForJavaFiles(root, getJavaSharedLibFileName(), map[string]bool{}, &log)

	if len(got) != 1 || got[0] != want {
		t.Fatalf("walkForJavaFiles() = %v, want [%s]", got, want)
	}
}

func TestWalkForJavaFilesContinuesAfterPermissionError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits do not apply")
	}
	if os.Geteuid() == 0 {
		t.Skip("root bypasses permission bits")
	}

	root := t.TempDir()

	denied := filepath.Join(root, "aaa_denied", "bin")
	if err := os.MkdirAll(denied, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Dir(denied), 0o000); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(filepath.Dir(denied), 0o755) })

	want := writeJavaBinary(t, root, "zzz_java")

	log := zerolog.Nop()
	got := walkForJavaFiles(root, getJavaSharedLibFileName(), map[string]bool{}, &log)

	if len(got) != 1 || got[0] != want {
		t.Fatalf("walkForJavaFiles() = %v, want [%s]", got, want)
	}
}

// The walk has no measurement today, so widening a platform's default roots is
// guesswork. Each root reports what it cost.
func TestWalkForJavaFilesReportsItsCost(t *testing.T) {
	root := t.TempDir()
	writeJavaBinary(t, root, "jdk-21")

	var sink bytes.Buffer
	log := zerolog.New(&sink)

	walkForJavaFiles(root, getJavaSharedLibFileName(), map[string]bool{}, &log)

	var found bool
	for _, line := range strings.Split(strings.TrimSpace(sink.String()), "\n") {
		var entry map[string]any
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Fatalf("log line is not json: %v", err)
		}
		if entry["message"] != "finished walking search path" {
			continue
		}
		found = true

		if entry["search_path"] != root {
			t.Errorf("search_path = %v, want %q", entry["search_path"], root)
		}
		if entries, _ := entry["entries"].(float64); entries < 1 {
			t.Errorf("entries = %v, want at least 1", entry["entries"])
		}
		if _, ok := entry["duration_ms"]; !ok {
			t.Error("duration_ms missing")
		}
	}

	if !found {
		t.Fatalf("no completion entry in the log:\n%s", sink.String())
	}
}
