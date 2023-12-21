package ojdmcollector

import (
	"fmt"
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
		fmt.Println("MacOS Java Search Paths: ", paths)
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
			"/oracle",
			"/bin",
			localSharePath,
		}
		paths = append(paths, linuxPaths...)
		fmt.Println("Linux Java Search Paths: ", paths)
		return paths

	case "windows":
		appDataDir := os.Getenv("LocalAppData")
		winPaths := []string{"C:\\Program Files", "C:\\Program Files (x86)", appDataDir}
		paths = append(paths, winPaths...)
		fmt.Println("Windows Java Search Paths: ", paths)
		return paths

	default:
		fmt.Println("Default Java Search Paths /")
		return []string{"/"}
	}
}
