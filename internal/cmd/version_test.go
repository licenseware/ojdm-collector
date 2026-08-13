package cmd

import (
	"runtime"
	"strings"
	"testing"
)

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
	build := Build{Version: "v1.2.3", Commit: "abc1234", Date: "2026-08-11T00:00:00Z", BuiltBy: "goreleaser"}

	got := versionString(build)
	for _, want := range []string{"ojdm-collector v1.2.3", "commit:    abc1234", "built:     2026-08-11T00:00:00Z"} {
		if !strings.Contains(got, want) {
			t.Fatalf("versionString() = %q, want it to contain %q", got, want)
		}
	}
}

func TestVersionStringFallsBackToBuildInfo(t *testing.T) {
	build := Build{Version: "dev", Commit: "none", Date: "unknown", BuiltBy: "unknown"}

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

	got := versionString(build)
	if !strings.Contains(got, "commit:    "+wantCommit) {
		t.Fatalf("versionString() = %q, want commit %q", got, wantCommit)
	}
	if !strings.Contains(got, "built:     "+wantDate) {
		t.Fatalf("versionString() = %q, want date %q", got, wantDate)
	}
}
