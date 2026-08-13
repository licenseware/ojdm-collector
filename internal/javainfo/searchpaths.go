package javainfo

import (
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

func getSearchPaths(log *zerolog.Logger) []string {

	oracleHomePath := os.Getenv("ORACLE_HOME")
	paths := []string{}
	if oracleHomePath != "" {
		paths = append(paths, oracleHomePath)
	}

	switch runtime.GOOS {

	case "darwin":
		homeDir, _ := os.UserHomeDir()
		paths = append(paths, darwinSearchPaths(homeDir)...)
		log.Debug().Str("platform", runtime.GOOS).Strs("paths", paths).Msg("default search paths")
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
		log.Debug().Str("platform", runtime.GOOS).Strs("paths", paths).Msg("default search paths")
		return paths

	case "windows":
		winPaths := []string{"C:\\Program Files", "C:\\Program Files (x86)"}

		userProfileDir, err := os.UserHomeDir()
		if err != nil {
			log.Warn().Err(err).Msg("could not determine the user home directory")
			return paths
		}

		// Get a list of user profiles
		userProfiles, err := filepath.Glob(filepath.Join(filepath.Dir(userProfileDir), "*"))
		if err != nil {
			log.Warn().Err(err).Msg("could not enumerate user profiles")
			return paths
		}

		// Iterate over user profiles and add local app data paths
		for _, userProfile := range userProfiles {
			appDataPath := filepath.Join(userProfile, "AppData", "Local")
			if _, err := os.Stat(appDataPath); err == nil {
				winPaths = append(winPaths, appDataPath)
			}
		}

		// Add the all users profile directory
		allUsersProfileDir := os.Getenv("ALLUSERSPROFILE")
		if allUsersProfileDir != "" {
			allUsersAppDataPath := filepath.Join(allUsersProfileDir, "AppData", "Local")
			if _, err := os.Stat(allUsersAppDataPath); err == nil {
				winPaths = append(winPaths, allUsersAppDataPath)
			}
		}

		paths = append(paths, winPaths...)
		log.Debug().Str("platform", runtime.GOOS).Strs("paths", paths).Msg("default search paths")
		return paths

	default:
		log.Debug().Str("platform", runtime.GOOS).Msg("unknown platform, searching from the filesystem root")
		return []string{"/"}
	}
}

// darwinSearchPaths lists the roots a macOS installation can appear under.
// /Library/Java/JavaVirtualMachines is the canonical one — /usr/libexec/java_home,
// the Homebrew casks and the Temurin installer all target it — and searching
// only /Applications is why a stock Mac used to report nothing.
func darwinSearchPaths(homeDir string) []string {
	paths := []string{
		"/Library/Java/JavaVirtualMachines",
		"/Applications",
		"/usr/local",
		"/opt", // includes /opt/homebrew, the arm64 prefix
	}

	// path.Join, not filepath.Join: these are macOS paths whatever host the
	// tests compile on, and filepath.Join would separate them with a backslash
	// when the unit tests run on the windows runner.
	if homeDir != "" {
		paths = append(paths, path.Join(homeDir, "Library", "Java"))
	}

	return paths
}
