package ojdmcollector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/rs/zerolog"
)

func getJavaBinaryPath(basePath string) (string, error) {
	javaExec := "java"
	if runtime.GOOS == "windows" {
		javaExec += ".exe"
	}

	javaPath := filepath.Join(basePath, "bin", javaExec)

	if _, err := os.Stat(javaPath); !os.IsNotExist(err) {
		return javaPath, nil // Executable found
	}

	return "", fmt.Errorf("java executable not found at %s", javaPath)
}

// getJavaDLLRelativeDirs lists the directories holding the VM shared library,
// relative to the installation root, in priority order. Windows keeps it under
// bin, never lib. The jre/ prefixed entries cover JDK 8, where the runtime is
// nested inside the JDK, and the arch subdirectory covers Linux JRE 8.
func getJavaDLLRelativeDirs() [][]string {
	switch runtime.GOOS {
	case "windows":
		return [][]string{
			{"bin", "server"},
			{"bin", "client"},
			{"jre", "bin", "server"},
			{"jre", "bin", "client"},
		}
	case "darwin":
		return [][]string{
			{"lib", "server"},
			{"jre", "lib", "server"},
		}
	default:
		return [][]string{
			{"lib", "server"},
			{"lib", "amd64", "server"},
			{"jre", "lib", "server"},
			{"jre", "lib", "amd64", "server"},
		}
	}
}

func getJavaDLLPath(basePath string) (string, error) {
	var dllFileName string
	switch runtime.GOOS {
	case "windows":
		dllFileName = "jvm.dll"
	case "darwin":
		dllFileName = "libjvm.dylib"
	default:
		dllFileName = "libjvm.so"
	}

	var tried []string
	for _, relativeDir := range getJavaDLLRelativeDirs() {
		dllPath := filepath.Join(append([]string{basePath}, append(relativeDir, dllFileName)...)...)
		if _, err := os.Stat(dllPath); err == nil {
			return dllPath, nil // File found
		}
		tried = append(tried, dllPath)
	}

	return "", fmt.Errorf("%s not found in any of %v", dllFileName, tried)
}

func executeJavaBinary(javaBinPath string) (string, error) {
	cmd := exec.Command(javaBinPath, "-XshowSettings:all", "-version")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// getToolPath builds the path to a tool (jps, jinfo, javac) sitting next to the
// java binary, including the platform executable suffix.
func getToolPath(javaBinPath, toolName string) string {
	if runtime.GOOS == "windows" {
		toolName += ".exe"
	}
	return filepath.Join(filepath.Dir(javaBinPath), toolName)
}

// checkToolExists checks if a given tool (jps or jinfo) exists in the Java installation's bin directory.
func checkToolExists(javaBinPath, toolName string) bool {
	if _, err := os.Stat(getToolPath(javaBinPath, toolName)); err == nil {
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

func GetJavaVersionInfos(javaBasePaths []string, log *zerolog.Logger) []JavaInfoRunningProcs {
	var versionInfos []JavaInfoRunningProcs
	for _, basePath := range javaBasePaths {
		javaBinPath, err := getJavaBinaryPath(basePath)
		if err != nil {
			log.Warn().Err(err).Str("base_path", basePath).Msg("skipping installation, java binary not found")
			continue
		}
		output, err := executeJavaBinary(javaBinPath)
		if err != nil {
			log.Warn().Err(err).Str("base_path", basePath).Str("java_bin_path", javaBinPath).
				Msg("skipping installation, java binary could not be executed")
			continue
		}
		info := parseJavaVersionOutput(output)
		javaDllPath, err := getJavaDLLPath(basePath)
		if err != nil {
			log.Warn().Err(err).Str("base_path", basePath).Msg("vm shared library not found")
		} else {

			info.DynLibBinPath = javaDllPath
		}
		if checkToolExists(javaBinPath, "jps") && checkToolExists(javaBinPath, "jinfo") {
			info.JpsJinfoPresent = true
		}
		info.JavaBinPath = javaBinPath
		if checkToolExists(javaBinPath, "javac") {
			info.JavaCBinPath = normalizePath(getToolPath(javaBinPath, "javac"))
			info.IsJDK = true
		} else {
			info.IsJDK = false
		}
		info.HostName = getHostName()
		info.HostLogicalProcessors = runtime.NumCPU()
		log.Debug().Str("java_home", info.JavaHome).Str("java_version", info.JavaVersion).
			Bool("is_jdk", info.IsJDK).Msg("collected java installation")
		versionInfos = append(versionInfos, info)
	}
	return versionInfos
}
