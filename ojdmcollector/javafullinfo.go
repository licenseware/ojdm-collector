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

func getJavaBinFileName() string {
	if runtime.GOOS == "windows" {
		return "java.exe"
	}
	return "java"
}

func getJavaCBinFileName() string {
	if runtime.GOOS == "windows" {
		return "javac.exe"
	}
	return "javac"
}

func getJavaBasePath(jpath string) string {

	jSplitedPath := strings.Split(jpath, string(os.PathSeparator))
	libLastIdx := 0
	for dix, dir := range jSplitedPath {
		if dir == "lib" {
			libLastIdx = dix
		}
	}
	jbasePath := filepath.Join(jSplitedPath[:libLastIdx]...)

	if runtime.GOOS != "windows" {
		return filepath.Join(string(os.PathSeparator), jbasePath)
	}

	return jbasePath

}

func getJavaPath(jbasePath, javaFileName string) string {
	javaPath := filepath.Join(jbasePath, "bin", javaFileName)
	if fileExists(javaPath) {
		return javaPath
	}
	// NOTE: try up one dir join with bin
	jbasePath = upDir(jbasePath, 1)
	javaPath = filepath.Join(jbasePath, "bin", javaFileName)
	if fileExists(javaPath) {
		return javaPath
	}
	return ""
}

func getJavaFullVersionSettings(javaBinPath string) string {

	cmdSettingsAllVersion := exec.Command(javaBinPath, "-XshowSettings:all", "-version")
	fullOutput, err := cmdSettingsAllVersion.CombinedOutput()
	if err == nil {
		return string(fullOutput)
	}

	if strings.Contains(string(fullOutput), "Unrecognized option:") {
		cmdVersion := exec.Command(javaBinPath, "-version")
		partialOutput, err := cmdVersion.CombinedOutput()
		if err == nil {
			return javaBinPath + "\n" + string(partialOutput)
		}
	}

	fmt.Printf("Failed to retrieve java settings info: %v\n", err)
	return ""

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

func GetFullJavaInfo() []JavaInfoRunningProcs {

	hostName, _ := os.Hostname()
	javaBinFileName := getJavaBinFileName()
	javaCBinFileName := getJavaCBinFileName()

	javaSharedLibPaths := getJavaSharedLibPaths()
	runningProcs := getRunningProcCommands()

	Pprint(runningProcs)

	jInfoProcs := []JavaInfoRunningProcs{}
	for _, jpath := range javaSharedLibPaths {

		jbasePath := getJavaBasePath(jpath)
		javaBinPath := getJavaPath(jbasePath, javaBinFileName)
		javaCBinPath := getJavaPath(jbasePath, javaCBinFileName)

		isJDK := javaCBinPath != ""

		vinfo := JavaInfoRunningProcs{}
		if javaBinPath != "" {
			javaBinSettingsOutput := getJavaFullVersionSettings(javaBinPath)
			if javaBinSettingsOutput != "" {
				vinfo = extractInfoFromFullVersionSettings(javaBinSettingsOutput)
			}
		}

		procRunning := false
		procPath := ""
		cmdLine := ""
		for _, rProc := range runningProcs {
			if javaBinPath == rProc.ProcDir {
				procRunning = true
				procPath = rProc.ProcDir
				cmdLine = rProc.CommandLine
				break
			}
		}

		jinfo := JavaInfoRunningProcs{
			HostName:      hostName,
			DynLibBinPath: jpath,
			JavaBinPath:   javaBinPath,
			JavaCBinPath:  javaCBinPath,
			IsJDK:         isJDK,

			JavaHome:           vinfo.JavaHome,
			JavaRuntimeName:    vinfo.JavaRuntimeName,
			JavaRuntimeVersion: vinfo.JavaRuntimeVersion,
			JavaVendor:         vinfo.JavaVendor,
			JavaVersion:        vinfo.JavaVersion,
			JavaVersionDate:    vinfo.JavaVersionDate,
			JavaVMName:         vinfo.JavaVMName,
			JavaVMVendor:       vinfo.JavaVMVendor,
			JavaVMVersion:      vinfo.JavaVMVersion,

			ProcessRunning: procRunning,
			ProcessPath:    procPath,
			CommandLine:    cmdLine,
		}

		jInfoProcs = append(jInfoProcs, jinfo)

	}

	return jInfoProcs

}
