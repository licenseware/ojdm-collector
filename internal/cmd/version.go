package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"text/tabwriter"
)

// A Build carries the stamps the linker writes into a released binary. The
// values live in package main so GoReleaser keeps stamping -X main.<name>.
type Build struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

// buildInfoStamps reads the VCS stamps the toolchain embeds in binaries built
// from a git checkout, so local `go build` output is not stuck on the "none"
// and "unknown" defaults. Released builds are stamped via -ldflags and win.
func buildInfoStamps() (revision, when string, dirty bool) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", "", false
	}

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			when = setting.Value
		case "vcs.modified":
			dirty = setting.Value == "true"
		}
	}

	return revision, when, dirty
}

func versionString(build Build) string {
	resolvedCommit, resolvedDate := build.Commit, build.Date

	revision, when, dirty := buildInfoStamps()
	if resolvedCommit == "none" && revision != "" {
		resolvedCommit = revision
		if dirty {
			resolvedCommit += "-dirty"
		}
	}
	if resolvedDate == "unknown" && when != "" {
		resolvedDate = when
	}

	return formatVersion(build.Version, resolvedCommit, resolvedDate, build.BuiltBy)
}

func formatVersion(version, commit, date, builtBy string) string {
	var out strings.Builder
	fmt.Fprintf(&out, "ojdm-collector %s\n\n", version)

	// tabwriter buffers the whole block to size the columns, so every row has to
	// be written before Flush and the "\t" is a column break, not a literal tab.
	w := tabwriter.NewWriter(&out, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "commit:\t%s\n", commit)
	fmt.Fprintf(w, "built:\t%s\n", date)
	fmt.Fprintf(w, "built by:\t%s\n", builtBy)
	fmt.Fprintf(w, "go:\t%s\n", runtime.Version())
	fmt.Fprintf(w, "platform:\t%s/%s\n", runtime.GOOS, runtime.GOARCH)
	w.Flush()

	return strings.TrimRight(out.String(), "\n")
}
