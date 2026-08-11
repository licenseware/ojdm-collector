package ojdmcollector

import (
	"os"
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
		macPaths := []string{"/Applications"}
		paths = append(paths, macPaths...)
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
