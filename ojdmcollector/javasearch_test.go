package ojdmcollector

import "testing"

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
