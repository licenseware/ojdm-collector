package ojdmcollector

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func getJavaSharedLibFileName() string {
	switch runtime.GOOS {
	case "darwin":
		return "libjvm.dylib"
	case "linux":
		return "libjvm.so"
	case "windows":
		return "jvm.dll"
	default:
		return "libjvm.so"
	}
}

func getJavaSharedLibPaths(searchPaths []string) []string {
	javaSharedLibFilename := getJavaSharedLibFileName()

	searchPaths = append(searchPaths, getSearchPaths()...)

	fmt.Println("Java Search Paths: ", searchPaths)

	javaFilesMap := make(map[string]bool)
	var javaFiles []string
	for _, searchPath := range searchPaths {
		filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if !info.IsDir() {
				serverFolder := filepath.Base(filepath.Dir(path))
				if info.Name() == javaSharedLibFilename && serverFolder == "server" {
					if _, exists := javaFilesMap[path]; !exists {
						fmt.Printf("Found %s in path %s\n", info.Name(), path)
						javaFilesMap[path] = true
						javaFiles = append(javaFiles, path)
					}
				}
			}

			return nil
		})
	}

	fmt.Printf("Finished gathering all java related paths!\n")
	return javaFiles
}
