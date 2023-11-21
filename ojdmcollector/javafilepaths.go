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

func getJavaSharedLibPaths() []string {

	javaSharedLibFilename := getJavaSharedLibFileName()
	searchPaths := getSearchPaths()

	var javaFiles []string
	for _, searchPath := range searchPaths {

		err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if !info.IsDir() {
				serverFolder := filepath.Base(filepath.Base(filepath.Dir(path)))
				if info.Name() == javaSharedLibFilename && serverFolder == "server" {
					fmt.Printf("Found %s in path %s\n", info.Name(), path)
					javaFiles = append(javaFiles, path)
				}
			}

			return nil
		})

		if err != nil {
			fmt.Println("Encountered error", err)
			return javaFiles
		}

	}

	fmt.Printf("Finished gathering all java related paths!\n")
	return javaFiles

}
