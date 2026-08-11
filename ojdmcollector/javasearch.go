package ojdmcollector

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
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

func getJavaSharedLibPaths(searchPaths []string, log *zerolog.Logger) []string {
	javaSharedLibFilenames := getJavaSharedLibFileName()

	searchPaths = append(searchPaths, getSearchPaths(log)...)

	log.Info().Strs("search_paths", searchPaths).Msg("scanning for java installations")

	javaFilesMap := make(map[string]bool)
	var javaFiles []string
	for _, searchPath := range searchPaths {
		javaFiles = append(javaFiles, walkForJavaFiles(searchPath, javaSharedLibFilenames, javaFilesMap, log)...)
	}

	log.Info().Int("count", len(javaFiles)).Msg("finished gathering java related paths")
	return javaFiles
}

func walkForJavaFiles(searchPath string, javaSharedLibFilenames []string, javaFilesMap map[string]bool, log *zerolog.Logger) []string {
	var javaFiles []string

	filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Returning the error aborts the walk for this whole search root,
			// silently dropping every installation that sorts after the failing
			// entry. Broken junctions, cloud placeholders and files removed
			// mid-scan by security agents all reach here, so skip and carry on.
			log.Warn().Err(err).Str("path", path).Bool("permission_denied", os.IsPermission(err)).
				Msg("skipping unreadable path")
			return nil
		}

		if !info.IsDir() && isInTargetSubfolder(path) {
			for _, javaSharedLibFilename := range javaSharedLibFilenames {
				if info.Name() == javaSharedLibFilename {
					cleanPath := processPath(path)
					if _, exists := javaFilesMap[cleanPath]; !exists {
						log.Debug().Str("file", info.Name()).Str("path", path).Msg("found java file")
						javaFilesMap[cleanPath] = true
						javaFiles = append(javaFiles, cleanPath)
					}
				}
			}
		}

		return nil
	})

	return javaFiles
}

func processPath(path string) string {
	normalizedPath := normalizeSearchPath(path)
	// The last separator wins: an installation under a root that itself
	// contains "bin" (C:\bin-tools\jdk-21\bin\java.exe) would otherwise be
	// truncated to the root, and every lookup below it then fails.
	if idx := strings.LastIndex(normalizedPath, "/bin/"); idx != -1 {
		return normalizedPath[:idx]
	}
	if idx := strings.LastIndex(normalizedPath, "/lib/server/"); idx != -1 {
		return normalizedPath[:idx]
	}
	return normalizedPath
}

func isInTargetSubfolder(path string) bool {
	normalizedPath := normalizeSearchPath(path)
	return strings.Contains(normalizedPath, "/bin/") || strings.Contains(normalizedPath, "/lib/server/")
}

func normalizeSearchPath(path string) string {
	return strings.ReplaceAll(filepath.ToSlash(path), "\\", "/")
}
