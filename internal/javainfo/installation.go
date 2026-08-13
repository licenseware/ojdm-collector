package javainfo

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

func recordFromVersionOutput(output string) Record {
	return Record{
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

func inspect(javaBasePaths []string, log *zerolog.Logger) []Record {
	var versionInfos []Record
	indexByHome := make(map[string]int)
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
		info := recordFromVersionOutput(output)
		javaDllPath, err := getJavaDLLPath(basePath)
		if err != nil && info.JavaHome != "" && normalizePath(basePath) != info.JavaHome {
			// A java binary reached through a symlink (/usr/bin/java) yields a
			// base path that holds no runtime. The home the JVM reports itself
			// is the authoritative location.
			javaDllPath, err = getJavaDLLPath(info.JavaHome)
		}
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

		// Several search paths can reach one installation (a symlinked
		// /usr/bin/java and the real /usr/lib/jvm/... entry). Reporting it once
		// per route would double count it.
		key := info.JavaHome
		if key == "" {
			key = normalizePath(basePath)
		}
		if index, seen := indexByHome[key]; seen {
			if isRicherRecord(info, versionInfos[index]) {
				log.Debug().Str("java_home", key).Str("base_path", basePath).
					Msg("replacing duplicate installation with a more complete record")
				versionInfos[index] = info
			} else {
				log.Debug().Str("java_home", key).Str("base_path", basePath).
					Msg("skipping duplicate installation")
			}
			continue
		}

		indexByHome[key] = len(versionInfos)
		versionInfos = append(versionInfos, info)
	}
	return versionInfos
}

// isRicherRecord reports whether candidate describes an installation more
// completely than existing, deciding which route to keep for a java home that
// was reached more than once.
func isRicherRecord(candidate, existing Record) bool {
	if (candidate.DynLibBinPath != "") != (existing.DynLibBinPath != "") {
		return candidate.DynLibBinPath != ""
	}
	if candidate.IsJDK != existing.IsJDK {
		return candidate.IsJDK
	}
	return false
}
