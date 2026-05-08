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
