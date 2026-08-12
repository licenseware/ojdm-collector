package ojdmcollector

import "testing"

// A JVM running from a location no search path reaches used to be merged with
// nothing, so its row carried a java home and a command line but no version, no
// IsJDK and no shared library. Its own java.home is an installation the scan
// missed, and has to be collected like any other.
func TestUnscannedProcessHomesFindsWhatTheScanMissed(t *testing.T) {
	processInfo := []JavaInfoRunningProcs{
		{JavaHome: "/opt/jdk-21"},
		{JavaHome: "/Users/runner/hostedtoolcache/jdk-25/Contents/Home"},
	}
	versionInfo := []JavaInfoRunningProcs{
		{JavaHome: "/opt/jdk-21"},
	}

	got := unscannedProcessHomes(processInfo, versionInfo)
	want := []string{"/Users/runner/hostedtoolcache/jdk-25/Contents/Home"}

	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("unscannedProcessHomes() = %v, want %v", got, want)
	}
}

// Two processes out of the same installation must not make the collector run
// the java binary twice.
func TestUnscannedProcessHomesDeduplicates(t *testing.T) {
	processInfo := []JavaInfoRunningProcs{
		{JavaHome: "/opt/jdk-21"},
		{JavaHome: "/opt/jdk-21"},
	}

	if got := unscannedProcessHomes(processInfo, nil); len(got) != 1 {
		t.Fatalf("unscannedProcessHomes() = %v, want one entry", got)
	}
}

// A process jinfo could not describe has no home to collect, and an empty base
// path would be walked as the working directory.
func TestUnscannedProcessHomesSkipsEmptyHomes(t *testing.T) {
	processInfo := []JavaInfoRunningProcs{{JavaHome: ""}}

	if got := unscannedProcessHomes(processInfo, nil); len(got) != 0 {
		t.Fatalf("unscannedProcessHomes() = %v, want none", got)
	}
}

// The scan and jinfo can spell the same installation differently, and a
// spelling difference must not present it as a second installation.
func TestUnscannedProcessHomesMatchesOnNormalisedPaths(t *testing.T) {
	processInfo := []JavaInfoRunningProcs{{JavaHome: `C:\Program Files\Java\jdk-21`}}
	versionInfo := []JavaInfoRunningProcs{{JavaHome: "C:/Program Files/Java/jdk-21"}}

	if got := unscannedProcessHomes(processInfo, versionInfo); len(got) != 0 {
		t.Fatalf("unscannedProcessHomes() = %v, want none", got)
	}
}
