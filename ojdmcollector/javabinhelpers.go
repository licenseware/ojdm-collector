package ojdmcollector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var IGNORE_LINUX_DIRS []string = []string{
	"boot", "cdrom", "dev", "etc", "lost+found", "media",
	"mnt", "proc", "root", "run", "sbin", "srv", "sys",
	"tmp", "var", "Trash",
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
						if strings.Contains(info.Name(), dir) {
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
			if runtime.GOOS == "windows" {
				baseDirPath = strings.Replace(baseDirPath, "C:Program", "C:\\Program", 1)
			}
			baseMaplibjvmPaths[baseDirPath] = path
		} else {
			dirBaseSlice := dirSlice[:len(dirSlice)-3]
			baseDirPath := filepath.Join(dirBaseSlice...)
			baseMaplibjvmPaths[osSep+baseDirPath] = path
		}
	}

	return baseMaplibjvmPaths

}

func isJavaBinary(procName string) bool {

	javaBins := []string{"java", "javac", "java.exe", "javac.exe"}

	for _, jname := range javaBins {
		if jname == procName {
			return true
		}
	}
	return false
}

func getJavaBinPaths() []JavaInfoRunningProcs {

	libjvmPaths, err := getDYNLIBFullPaths()
	if err != nil {
		panic("Cannot find all usage of java dynamic binaries")
	}

	baseJVMSOPaths := getBaseDYNLIBPaths(libjvmPaths)

	fmt.Println("Gather related java binaries...")

	javaMappedPaths := map[string]JavaInfoRunningProcs{}
	for jvmlibRootPath, jvmlibBinPath := range baseJVMSOPaths {

		err := filepath.Walk(jvmlibRootPath, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if !info.IsDir() && isJavaBinary(info.Name()) {

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

	cmdSettingsAllVersion := exec.Command(javaBinPath, "-XshowSettings:all", "-version")
	fullOutput, err := cmdSettingsAllVersion.CombinedOutput()
	if err == nil {
		return string(fullOutput), nil
	}

	if strings.Contains(string(fullOutput), "Unrecognized option:") {
		cmdVersion := exec.Command(javaBinPath, "-version")
		partialOutput, err := cmdVersion.CombinedOutput()
		if err == nil {
			return javaBinPath + "\n" + string(partialOutput), nil
		}
	}

	fmt.Printf("Failed to retrieve java settings info: %v\n", err)
	return "", err

}

func findRegexInText(regex, text string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractInfoFromFullVersionSettings(versionSettings string) JavaInfoRunningProcs {

	javaHome := findRegexInText(`java.home\s=\s(.*)`, versionSettings)
	javaRuntimeName := findRegexInText(`java.runtime.name\s=\s(.*)`, versionSettings)
	javaRuntimeVersion := findRegexInText(`java.runtime.version\s=\s(.*)`, versionSettings)
	javaVersion := findRegexInText(`java.version\s=\s(.*)`, versionSettings)

	javaVersionDate := findRegexInText(`java.version.date\s=\s(.*)`, versionSettings)

	javaVMName := findRegexInText(`java.vm.name\s=\s(.*)`, versionSettings)

	javaVendor := findRegexInText(`java.vendor\s=\s(.*)`, versionSettings)
	javaVMVendor := findRegexInText(`java.vm.vendor\s=\s(.*)`, versionSettings)
	javaVMVersion := findRegexInText(`java.vm.version\s=\s(.*)`, versionSettings)

	if javaHome == "" {
		javaHome = findRegexInText(`(.*)/bin/java`, versionSettings)
	}

	if javaRuntimeName == "" {
		javaRuntimeName = findRegexInText(`(.*\sRuntime\sEnvironment).*?`, versionSettings)
	}

	if javaRuntimeVersion == "" {
		javaRuntimeVersion = findRegexInText(`.*Runtime\sEnvironment\s\(build (.*)\)`, versionSettings)
	}

	if javaVersion == "" {
		javaVersion = findRegexInText(`.*\sversion\s"(.*)".*`, versionSettings)
	}

	if javaVersionDate == "" {
		javaVersionDate = findRegexInText(`.*version.*".*"\s(.*)`, versionSettings)
	}

	if javaVMName == "" {
		javaVMName = findRegexInText(`(.*Server VM).*build`, versionSettings)
	}

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
