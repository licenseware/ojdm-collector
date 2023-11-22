package ojdmcollector

import (
	"os"
	"path/filepath"
	"runtime"
)

func getSearchPaths() []string {

	oracleHomePath := os.Getenv("ORACLE_HOME")
	paths := []string{}
	if oracleHomePath != "" {
		paths = append(paths, oracleHomePath)
	}

	switch runtime.GOOS {

	case "darwin":
		macPaths := []string{"/Applications"}
		paths = append(paths, macPaths...)
		return paths

	case "linux":
		homeDir, _ := os.UserHomeDir()
		localSharePath := filepath.Join(homeDir, ".local/share")
		linuxPaths := []string{
			"/home",
			"/usr/bin",
			"/usr/local",
			"/usr/lib",
			"/usr/share",
			"/opt",
			"/snap",
			localSharePath,
		}
		paths = append(paths, linuxPaths...)
		return paths

	case "windows":
		appDataDir := os.Getenv("LocalAppData")
		winPaths := []string{"C:\\Program Files", "C:\\Program Files (x86)", appDataDir}
		paths = append(paths, winPaths...)
		return paths

	default:
		return []string{"/"}
	}
}
