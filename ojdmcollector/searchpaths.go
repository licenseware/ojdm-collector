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
		winPaths := []string{"C:\\Program Files", "C:\\Program Files (x86)"}
		
		// Get a list of user profiles
		userProfiles, err := filepath.Glob(filepath.Join(filepath.Dir(userProfileDir), "*"))
		if err != nil {
        		fmt.Println("Error getting user profiles:", err)
        		return paths
		}

		// Iterate over user profiles and add local app data paths
    		for _, userProfile := range userProfiles {
       			if userProfile != userProfileDir {
            			appDataPath := filepath.Join(userProfile, "AppData", "Local")
            			if _, err := os.Stat(appDataPath); err == nil {
	                		winPaths = append(winPaths, appDataPath)
        	   		}
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
		fmt.Println("Windows Java Search Paths: ", paths)
		return paths

	default:
		fmt.Println("Default Java Search Paths /")
		return []string{"/"}
	}
}
