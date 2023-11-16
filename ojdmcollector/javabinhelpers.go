package ojdmcollector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
)

var IGNORE_LINUX_DIRS []string = []string{
	"boot", "cdrom", "dev", "etc", "lost+found", "media",
	"mnt", "proc", "root", "run", "sbin", "srv", "sys",
	"tmp", "var",
}

func getDYNLIBFileName() string {
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

func getOsAppsInstalledRootPaths() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"/Applications"}
	case "linux":
		return []string{"/"}
	case "windows":
		return []string{"C:\\Program Files", "C:\\Program Files (x86)"}
	default:
		return []string{"/"}
	}
}

func getDYNLIBFullPaths() ([]string, error) {

	libjvmFilename := getDYNLIBFileName()
	rootPaths := getOsAppsInstalledRootPaths()

	fmt.Printf("Gathering all %s binary filepaths...\n", libjvmFilename)

	var libjvmPaths []string
	for _, rootPath := range rootPaths {

		err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if runtime.GOOS == "linux" {
				if info.IsDir() {
					for _, dir := range IGNORE_LINUX_DIRS {
						if info.Name() == dir {
							return filepath.SkipDir
						}
					}
				}
			}

			if !info.IsDir() && info.Name() == libjvmFilename {
				fmt.Printf("Found %s in path %s\n", libjvmFilename, path)
				libjvmPaths = append(libjvmPaths, path)
			}

			return nil
		})

		if err != nil {
			fmt.Println("Encountered error", err)
			return libjvmPaths, err
		}

	}

	fmt.Printf("Finished gathering all %s binary paths!\n", libjvmFilename)
	return libjvmPaths, nil

}

func getBaseDYNLIBPaths(libjvmPaths []string) map[string]string {

	baseMaplibjvmPaths := map[string]string{}
	for _, path := range libjvmPaths {
		dir := filepath.Dir(path)

		osSep := string(os.PathSeparator)
		dirSlice := strings.Split(dir, osSep)

		if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
			dirBaseSlice := dirSlice[:3]
			baseDirPath := filepath.Join(dirBaseSlice...)
			baseMaplibjvmPaths[baseDirPath] = path
		} else {
			dirBaseSlice := dirSlice[:len(dirSlice)-3]
			baseDirPath := filepath.Join(dirBaseSlice...)
			baseMaplibjvmPaths[osSep+baseDirPath] = path
		}
	}

	return baseMaplibjvmPaths

}

func getJavaBinPaths() []JavaInfoRunningProcs {

	libjvmPaths, err := getDYNLIBFullPaths()
	if err != nil {
		panic("Cannot find all usage of java dynamic binaries")
	}

	baseJVMSOPaths := getBaseDYNLIBPaths(libjvmPaths)

	fmt.Println("Gather related java binaries...")

	javaBins := []string{"java", "javac", "java.exe", "javac.exe"}

	javaMappedPaths := map[string]JavaInfoRunningProcs{}
	for jvmlibRootPath, jvmlibBinPath := range baseJVMSOPaths {

		err := filepath.Walk(jvmlibRootPath, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if !info.IsDir() && slices.Contains(javaBins, info.Name()) {

				javaInfoRunningProcs, ok := javaMappedPaths[jvmlibRootPath]

				isJDK := false
				javaBinPath := ""
				javaCBinPath := ""

				if ok {

					if javaInfoRunningProcs.IsJDK {
						isJDK = javaInfoRunningProcs.IsJDK
					}

					if javaInfoRunningProcs.JavaBinPath != "" {
						javaBinPath = javaInfoRunningProcs.JavaBinPath
					}

					if javaInfoRunningProcs.JavaCBinPath != "" {
						javaCBinPath = javaInfoRunningProcs.JavaCBinPath
					}

				}

				if info.Name() == "java" || info.Name() == "java.exe" {
					fmt.Println("Found java in path ", path)
					javaBinPath = path
				}

				if info.Name() == "javac" || info.Name() == "javac.exe" {
					fmt.Println("Found javac in path ", path)
					javaCBinPath = path
					isJDK = true
				}

				javaMappedPaths[jvmlibRootPath] = JavaInfoRunningProcs{
					AppDirName:    filepath.Base(jvmlibRootPath),
					DynLibBinPath: jvmlibBinPath,
					JavaBinPath:   javaBinPath,
					JavaCBinPath:  javaCBinPath,
					IsJDK:         isJDK,
					BaseDir:       jvmlibRootPath,
				}
			}

			return nil
		})

		if err != nil {
			fmt.Println("Encountered error", err)
		}

	}

	javaLibsInfo := []JavaInfoRunningProcs{}
	for _, values := range javaMappedPaths {
		javaLibsInfo = append(javaLibsInfo, values)
	}

	fmt.Println("Finished gathering related java binaries to java dynamic library!")

	return javaLibsInfo

}

func getJavaFullVersionSettings(javaBinPath string) (string, error) {

	cmd := exec.Command(javaBinPath, "-XshowSettings:all", "-version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to execute command: %v\n", err)
		return "", err
	}

	return string(output), nil

}

func findRegexInText(regex, text string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractInfoFromFullVersionSettings(versionSettings string) JavaInfoRunningProcs {

	javaHome := findRegexInText(`java.home\s=\s(.*)`, versionSettings)
	javaRuntimeName := findRegexInText(`java.runtime.name\s=\s(.*)`, versionSettings)
	javaRuntimeVersion := findRegexInText(`java.runtime.version\s=\s(.*)`, versionSettings)
	javaVendor := findRegexInText(`java.vendor\s=\s(.*)`, versionSettings)
	javaVersion := findRegexInText(`java.version\s=\s(.*)`, versionSettings)
	javaVersionDate := findRegexInText(`java.version.date\s=\s(.*)`, versionSettings)
	javaVMName := findRegexInText(`java.vm.name\s=\s(.*)`, versionSettings)
	javaVMVendor := findRegexInText(`java.vm.vendor\s=\s(.*)`, versionSettings)
	javaVMVersion := findRegexInText(`java.vm.version\s=\s(.*)`, versionSettings)

	versionInfo := JavaInfoRunningProcs{
		JavaHome:           javaHome,
		JavaRuntimeName:    javaRuntimeName,
		JavaRuntimeVersion: javaRuntimeVersion,
		JavaVendor:         javaVendor,
		JavaVersion:        javaVersion,
		JavaVersionDate:    javaVersionDate,
		JavaVMName:         javaVMName,
		JavaVMVendor:       javaVMVendor,
		JavaVMVersion:      javaVMVersion,
	}

	return versionInfo

}
