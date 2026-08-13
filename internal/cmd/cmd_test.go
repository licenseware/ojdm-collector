package cmd

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
