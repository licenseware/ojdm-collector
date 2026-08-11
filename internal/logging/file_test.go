package logging

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDefaultLogPathSitsBesideTheReport(t *testing.T) {
	got := DefaultLogPath(filepath.Join("reports", "report.csv"))
	want := filepath.Join("reports", "logs", "debug.log")

	if got != want {
		t.Fatalf("DefaultLogPath() = %q, want %q", got, want)
	}
}

// A rerun must not destroy the log of the run that went wrong, which is the
// one support actually needs.
func TestOpenFileSinkArchivesPreviousLogByModTime(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "logs", "debug.log")

	first, _, err := OpenFileSink(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := first.WriteString("previous run\n"); err != nil {
		t.Fatal(err)
	}
	first.Close()

	modTime := time.Date(2026, 8, 11, 14, 3, 22, 0, time.UTC)
	if err := os.Chtimes(logPath, modTime, modTime); err != nil {
		t.Fatal(err)
	}

	second, archived, err := OpenFileSink(logPath)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()

	wantArchive := filepath.Join(filepath.Dir(logPath), "debug-20260811T140322Z.log")
	if archived != wantArchive {
		t.Fatalf("archived = %q, want %q", archived, wantArchive)
	}

	preserved, err := os.ReadFile(archived)
	if err != nil {
		t.Fatalf("archived log unreadable: %v", err)
	}
	if !strings.Contains(string(preserved), "previous run") {
		t.Fatalf("archived log lost its content: %q", preserved)
	}

	current, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(current) != 0 {
		t.Fatalf("new log should start empty, got %q", current)
	}
}

func TestOpenFileSinkReportsNoArchiveOnFirstRun(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "logs", "debug.log")

	file, archived, err := OpenFileSink(logPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if archived != "" {
		t.Fatalf("archived = %q, want empty on the first run", archived)
	}
}
