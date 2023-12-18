package ojdmcollector

import (
	"os"
	"regexp"
	"runtime"
)

func normalizePath(path string) string {
	// Specific handling for Windows drive letter
	if runtime.GOOS == "windows" {
		driveLetterRe := regexp.MustCompile(`([a-zA-Z])\\:\\`)
		path = driveLetterRe.ReplaceAllString(path, `$1:\`)
	}

	// Normalize slashes
	re := regexp.MustCompile(`[/\\]+`)
	normalizedPath := re.ReplaceAllString(path, `/`)

	return normalizedPath
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return hostname
}
