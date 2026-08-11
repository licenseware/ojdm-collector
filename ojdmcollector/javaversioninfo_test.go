package ojdmcollector

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Windows keeps jvm.dll under bin\server; only the POSIX runtimes use
// lib/server. Looking in the POSIX location alone left DynLibBinPath empty in
// every report produced on Windows.
func TestGetJavaDLLPathFindsPlatformLayouts(t *testing.T) {
	dllName := map[string]string{
		"windows": "jvm.dll",
		"darwin":  "libjvm.dylib",
	}[runtime.GOOS]
	if dllName == "" {
		dllName = "libjvm.so"
	}

	for _, relativeDir := range getJavaDLLRelativeDirs() {
		layout := filepath.Join(relativeDir...)
		t.Run(layout, func(t *testing.T) {
			base := t.TempDir()
			dir := filepath.Join(base, layout)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatal(err)
			}
			want := filepath.Join(dir, dllName)
			if err := os.WriteFile(want, []byte("stub"), 0o644); err != nil {
				t.Fatal(err)
			}

			got, err := getJavaDLLPath(base)
			if err != nil {
				t.Fatalf("getJavaDLLPath() error = %v", err)
			}
			if got != want {
				t.Fatalf("getJavaDLLPath() = %q, want %q", got, want)
			}
		})
	}
}

func TestGetJavaDLLPathErrorsWhenAbsent(t *testing.T) {
	if _, err := getJavaDLLPath(t.TempDir()); err == nil {
		t.Fatal("expected an error when no VM library is present")
	}
}
