package ojdmcollector

import (
	"os"
	"regexp"
	"runtime"
	"strings"
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

func findRegexInText(regex, text string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
