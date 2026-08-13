package main

import "testing"

// The ldflag stamps live here because GoReleaser targets -X main.<name>. If a
// default drifts, released binaries silently report the wrong build.
func TestBuildStampDefaults(t *testing.T) {
	for _, tt := range []struct {
		name string
		got  string
		want string
	}{
		{"version", version, "dev"},
		{"commit", commit, "none"},
		{"date", date, "unknown"},
		{"builtBy", builtBy, "unknown"},
	} {
		if tt.got != tt.want {
			t.Errorf("%s default = %q, want %q", tt.name, tt.got, tt.want)
		}
	}
}
