package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseSearchPaths(t *testing.T) {
	got := parseSearchPaths(" /opt/java, ,/oracle,/mnt/share ")
	want := []string{"/opt/java", "/oracle", "/mnt/share"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseSearchPaths() = %v, want %v", got, want)
	}
}

func TestParseSearchPathsFile(t *testing.T) {
	searchPathsFile := filepath.Join(t.TempDir(), "search-paths.txt")
	err := os.WriteFile(searchPathsFile, []byte("\n# shared locations\n/opt/java\n  /mnt/share/java  \n"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	got, err := parseSearchPathsFile(searchPathsFile)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"/opt/java", "/mnt/share/java"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseSearchPathsFile() = %v, want %v", got, want)
	}
}

func TestVersionString(t *testing.T) {
	origVersion, origCommit, origDate, origBuiltBy := version, commit, date, builtBy
	t.Cleanup(func() {
		version, commit, date, builtBy = origVersion, origCommit, origDate, origBuiltBy
	})

	version, commit, date, builtBy = "v1.2.3", "abc1234", "2026-08-11T00:00:00Z", "goreleaser"

	got := versionString()
	want := "ojdm-collector v1.2.3 (commit abc1234, built 2026-08-11T00:00:00Z by goreleaser)"
	if got != want {
		t.Fatalf("versionString() = %q, want %q", got, want)
	}
}

func TestVersionStringDefaults(t *testing.T) {
	if version != "dev" {
		t.Fatalf("version default = %q, want %q", version, "dev")
	}
}
