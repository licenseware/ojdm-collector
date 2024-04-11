package ojdmcollector

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func getJavaSharedLibFileName() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"libjvm.dylib", "java", "javac"}
	case "windows":
		return []string{"jvm.dll", "java.exe", "javac.exe"}
	default:
		return []string{"libjvm.so", "java", "javac"}
	}
}

func getJavaSharedLibPaths(searchPaths []string) []string {
	javaSharedLibFilenames := getJavaSharedLibFileName()

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

			if !info.IsDir() && isInTargetSubfolder(path) {
				for _, javaSharedLibFilename := range javaSharedLibFilenames {
					if info.Name() == javaSharedLibFilename {
						cleanPath := processPath(path)
						if _, exists := javaFilesMap[cleanPath]; !exists {
							fmt.Printf("Found %s in path %s\n", info.Name(), path)
							javaFilesMap[cleanPath] = true
							javaFiles = append(javaFiles, cleanPath)
						}
					}
				}
			}

			return nil
		})
	}

	fmt.Printf("Finished gathering all java related paths!\n")
	return javaFiles
}

func processPath(path string) string {
	if idx := strings.Index(path, "/bin"); idx != -1 {
		return path[:idx]
	}
	if idx := strings.Index(path, "/lib/server"); idx != -1 {
		return path[:idx]
	}
	return path
}

func isInTargetSubfolder(path string) bool {
	return strings.Contains(path, "/bin/") || strings.Contains(path, "/lib/server/")
}
