package main

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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

func TestFormatVersion(t *testing.T) {
	got := formatVersion("v1.2.3", "abc1234", "2026-08-11T00:00:00Z", "goreleaser")

	want := strings.Join([]string{
		"ojdm-collector v1.2.3",
		"",
		"commit:    abc1234",
		"built:     2026-08-11T00:00:00Z",
		"built by:  goreleaser",
		"go:        " + runtime.Version(),
		"platform:  " + runtime.GOOS + "/" + runtime.GOARCH,
	}, "\n")

	if got != want {
		t.Fatalf("formatVersion() =\n%s\nwant:\n%s", got, want)
	}
}

func TestVersionStringUsesLdflagStamps(t *testing.T) {
	origVersion, origCommit, origDate, origBuiltBy := version, commit, date, builtBy
	t.Cleanup(func() {
		version, commit, date, builtBy = origVersion, origCommit, origDate, origBuiltBy
	})

	version, commit, date, builtBy = "v1.2.3", "abc1234", "2026-08-11T00:00:00Z", "goreleaser"

	got := versionString()
	for _, want := range []string{"ojdm-collector v1.2.3", "commit:    abc1234", "built:     2026-08-11T00:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Fatalf("versionString() = %q, want it to contain %q", got, want)
		}
	}
}

func TestVersionStringFallsBackToBuildInfo(t *testing.T) {
	origCommit, origDate := commit, date
	t.Cleanup(func() { commit, date = origCommit, origDate })

	commit, date = "none", "unknown"

	revision, when, dirty := buildInfoStamps()
	wantCommit := "none"
	if revision != "" {
		wantCommit = revision
		if dirty {
			wantCommit += "-dirty"
		}
	}
	wantDate := "unknown"
	if when != "" {
		wantDate = when
	}

	got := versionString()
	if !strings.Contains(got, "commit:    "+wantCommit) {
		t.Fatalf("versionString() = %q, want commit %q", got, wantCommit)
	}
	if !strings.Contains(got, "built:     "+wantDate) {
		t.Fatalf("versionString() = %q, want date %q", got, wantDate)
	}
}

func TestVersionStringDefaults(t *testing.T) {
	if version != "dev" {
		t.Fatalf("version default = %q, want %q", version, "dev")
	}
}
