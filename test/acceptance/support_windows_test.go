//go:build acceptance && windows

package acceptance

import (
	"os"
	"path/filepath"
)

// scanRoot is a writable location for planted installations.
func scanRoot() string { return `C:\` }

// searchableSystemDir is covered by the default search paths, so anything
// planted here is met without passing extra flags.
func searchableSystemDir() string { return `C:\Program Files` }

// removeStubbornly deletes entries whose names Win32 path resolution cannot
// address, such as the trailing space planted by plantUnreadableEntry. Only
// the \\?\ prefix reaches them.
func removeStubbornly(path string) {
	if entries, err := os.ReadDir(path); err == nil {
		for _, entry := range entries {
			os.Remove(`\\?\` + filepath.Join(path, entry.Name()))
		}
	}
	os.RemoveAll(`\\?\` + path)
}
