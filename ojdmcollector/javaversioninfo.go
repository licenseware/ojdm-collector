package ojdmcollector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func getJavaBinaryPath(javaDllPath string) (string, error) {

	javaExec := "java"
	if runtime.GOOS == "windows" {
		javaExec = "java.exe"
	}

	dir := filepath.Dir(javaDllPath)
	for i := 0; i < 4; i++ {
		dir = filepath.Dir(dir)
		binPath := filepath.Join(dir, "bin")
		javaPath := filepath.Join(binPath, javaExec)
		if _, err := os.Stat(javaPath); err == nil {
			return javaPath, nil
		}
	}
	return "", fmt.Errorf("java bin folder with executable not found")

}

func executeJavaBinary(javaBinPath string) (string, error) {
	cmd := exec.Command(javaBinPath, "-XshowSettings:all", "-version")
	output, err := cmd.CombinedOutput()
	return string(output), err
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
		javaBinPath, err := getJavaBinaryPath(dllPath)
		if err != nil {
			continue
		}
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
		info.HostLogicalProcessors = runtime.NumCPU()
		versionInfos = append(versionInfos, info)
	}
	return versionInfos
}
