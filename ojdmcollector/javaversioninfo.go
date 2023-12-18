package ojdmcollector

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func getJavaBinaryPath(javaDllPath string) string {
	// Assuming javaDllPath is the path to the jvm.dll or equivalent file
	dir := filepath.Dir(javaDllPath)
	if strings.HasSuffix(dir, "server") {
		dir = filepath.Dir(dir) // Go up one level to the bin directory
	}
	javaBin := "java"
	if runtime.GOOS == "windows" {
		javaBin = "java.exe"
	}
	return filepath.Join(dir, javaBin)
}

func executeJavaBinary(javaBinPath string) (string, error) {
	cmd := exec.Command(javaBinPath, "-XshowSettings:all", "-version")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func findRegexInText(regex, text string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

// checkToolExists checks if a given tool (jps or jinfo) exists in the Java installation's bin directory.
func checkToolExists(javaBinPath, toolName string) bool {
	binDir := filepath.Dir(javaBinPath)
	var toolPath string
	if runtime.GOOS == "windows" {
		toolPath = filepath.Join(binDir, toolName+".exe")
	} else {
		toolPath = filepath.Join(binDir, toolName)
	}

	if _, err := os.Stat(toolPath); err == nil {
		return true
	}
	return false
}

func parseJavaVersionOutput(output string) JavaInfoRunningProcs {
	return JavaInfoRunningProcs{
		JavaHome:           normalizePath(findRegexInText(`java.home\s=\s(.*)`, output)),
		JavaRuntimeName:    findRegexInText(`java.runtime.name\s=\s(.*)`, output),
		JavaRuntimeVersion: findRegexInText(`java.runtime.version\s=\s(.*)`, output),
		JavaVersion:        findRegexInText(`java.version\s=\s(.*)`, output),
		JavaVersionDate:    findRegexInText(`java.version.date\s=\s(.*)`, output),
		JavaVMName:         findRegexInText(`java.vm.name\s=\s(.*)`, output),
		JavaVendor:         findRegexInText(`java.vendor\s=\s(.*)`, output),
		JavaVMVendor:       findRegexInText(`java.vm.vendor\s=\s(.*)`, output),
		JavaVMVersion:      findRegexInText(`java.vm.version\s=\s(.*)`, output),
	}
}

func GetJavaVersionInfos(javaDllPaths []string) []JavaInfoRunningProcs {
	var versionInfos []JavaInfoRunningProcs
	for _, dllPath := range javaDllPaths {
		javaBinPath := getJavaBinaryPath(dllPath)
		output, err := executeJavaBinary(javaBinPath)
		if err != nil {
			continue // Handle error or log as needed
		}
		info := parseJavaVersionOutput(output)
		if checkToolExists(javaBinPath, "jps") && checkToolExists(javaBinPath, "jinfo") {
			info.JpsJinfoPresent = true
		}
		info.JavaBinPath = javaBinPath
		if checkToolExists(javaBinPath, "javac") {
			info.JavaCBinPath = normalizePath(filepath.Join(filepath.Dir(javaBinPath), "javac"))
			info.IsJDK = true
		} else {
			info.IsJDK = false
		}
		info.HostName = getHostName()
		versionInfos = append(versionInfos, info)
	}
	return versionInfos
}
