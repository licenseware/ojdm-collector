package javainfo

import (
	"strings"
	"testing"
)

func TestIsInTargetSubfolderHandlesWindowsPaths(t *testing.T) {
	if !isInTargetSubfolder(`C:\Program Files\Java\jdk-21\bin\java.exe`) {
		t.Fatal("expected Windows bin path to match target subfolder")
	}

	if !isInTargetSubfolder(`C:\Program Files\Java\jdk-21\lib\server\jvm.dll`) {
		t.Fatal("expected Windows lib server path to match target subfolder")
	}
}

func TestProcessPathNormalizesJavaHome(t *testing.T) {
	got := processPath(`C:\Program Files\Java\jdk-21\bin\java.exe`)
	want := "C:/Program Files/Java/jdk-21"

	if got != want {
		t.Fatalf("processPath() = %q, want %q", got, want)
	}
}

// An installation root that itself contains "bin" must not be truncated to that
// root, otherwise every path derived from JavaHome misses and the whole
// installation is dropped from the report.
func TestProcessPathUsesLastBinSegment(t *testing.T) {
	cases := map[string]string{
		`C:\bin-tools\jdk-21\bin\java.exe`:                "C:/bin-tools/jdk-21",
		`C:\Program Files\Java\jre1.8.0_481\bin\java.exe`: "C:/Program Files/Java/jre1.8.0_481",
		`C:\bin-tools\jre8\bin\server\jvm.dll`:            "C:/bin-tools/jre8",
		"/opt/bin-tools/jdk-21/lib/server/libjvm.so":      "/opt/bin-tools/jdk-21",
	}

	for input, want := range cases {
		if got := processPath(input); got != want {
			t.Errorf("processPath(%q) = %q, want %q", input, got, want)
		}
	}
}

// The canonical macOS JDK location is /Library/Java/JavaVirtualMachines, where
// /usr/libexec/java_home, the Homebrew casks and the Temurin installer all put
// their installations. Searching only /Applications means a stock Mac reports
// nothing at all.
func TestDarwinSearchPathsCoverTheCanonicalLocations(t *testing.T) {
	got := darwinSearchPaths("/Users/tester")

	want := []string{
		"/Library/Java/JavaVirtualMachines",
		"/Applications",
		"/usr/local",
		"/opt",
		"/Users/tester/Library/Java",
	}

	if len(got) != len(want) {
		t.Fatalf("darwinSearchPaths() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("darwinSearchPaths()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

// An empty home directory must not contribute a bare "Library/Java", which
// would be resolved relative to the working directory and scan something
// arbitrary.
func TestDarwinSearchPathsOmitsHomeWhenUnknown(t *testing.T) {
	for _, path := range darwinSearchPaths("") {
		if strings.Contains(path, "Library/Java") && !strings.HasPrefix(path, "/Library") {
			t.Errorf("unexpected home-relative path %q", path)
		}
	}
}

// A JDK on macOS lives inside an .app-style bundle, so JavaHome is the
// Contents/Home directory rather than the .jdk directory above it.
func TestProcessPathResolvesDarwinBundleLayout(t *testing.T) {
	got := processPath("/Library/Java/JavaVirtualMachines/temurin-21.jdk/Contents/Home/bin/java")
	want := "/Library/Java/JavaVirtualMachines/temurin-21.jdk/Contents/Home"

	if got != want {
		t.Fatalf("processPath() = %q, want %q", got, want)
	}
}
