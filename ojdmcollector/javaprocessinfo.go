package ojdmcollector

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type JavaProcess struct {
	ProcessID   string
	CommandLine string
}

func parseJpsOutput(output string) []JavaProcess {
	var processes []JavaProcess
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if line == "" || strings.Contains(line, "jps.Jps -mvl") {
			continue // Skip empty lines and the jps command itself
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			process := JavaProcess{
				ProcessID:   fields[0],
				CommandLine: strings.Join(fields[1:], " "),
			}
			processes = append(processes, process)
		}
	}

	return processes
}

func runJps(javaHome string) ([]JavaProcess, error) {
	jpsPath := filepath.Join(javaHome, "bin", "jps")
	if runtime.GOOS == "windows" {
		jpsPath += ".exe"
	}

	cmd := exec.Command(jpsPath, "-mvl")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	processes := parseJpsOutput(string(output))
	return processes, nil
}

func runJinfo(javaHome string, javaProcess JavaProcess) (string, error) {
	jinfoPath := filepath.Join(javaHome, "bin", "jinfo")
	if runtime.GOOS == "windows" {
		jinfoPath += ".exe"
	}

	cmd := exec.Command(jinfoPath, "-sysprops", javaProcess.ProcessID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func updateJavaInfoWithJinfoData(jinfoOutput string, javaProcess JavaProcess) JavaInfoRunningProcs {
	// Update JavaInfo with data extracted from jinfoOutput
	var javaInfo JavaInfoRunningProcs

	javaInfo.HostName = getHostName()
	javaInfo.ProcessPath = normalizePath(findRegexInText(`(?m)^user.dir=(.*)`, jinfoOutput))
	javaInfo.CommandLine = javaProcess.CommandLine
	javaInfo.JavaHome = normalizePath(findRegexInText(`java.home=(.*)`, jinfoOutput))
	javaInfo.JavaRuntimeName = findRegexInText(`java.runtime.name=(.*)`, jinfoOutput)
	javaInfo.JavaRuntimeVersion = findRegexInText(`java.runtime.version\s=\s(.*)`, jinfoOutput)
	javaInfo.JavaVersion = findRegexInText(`java.version=(.*)`, jinfoOutput)
	javaInfo.JavaVMName = findRegexInText(`java.vm.name=(.*)`, jinfoOutput)
	javaInfo.JavaVendor = findRegexInText(`java.vendor=(.*)`, jinfoOutput)
	javaInfo.JavaVMVendor = findRegexInText(`java.vm.vendor=(.*)`, jinfoOutput)
	javaInfo.JavaVMVersion = findRegexInText(`java.vm.version=(.*)`, jinfoOutput)
	javaInfo.ProcessRunning = true
	javaInfo.HostLogicalProcessors = runtime.NumCPU()

	return javaInfo
}

func GetJavaProcessInfo(javaHome string, log *zerolog.Logger) []JavaInfoRunningProcs {
	var runningProcsJavaInfos []JavaInfoRunningProcs

	jpsOutput, err := runJps(javaHome)
	if err != nil {
		log.Warn().Err(err).Str("java_home", javaHome).
			Msg("could not run jps, running jvm instances will not be identified")
		return runningProcsJavaInfos
	}

	log.Debug().Int("count", len(jpsOutput)).Msg("jps reported running jvm instances")
	for _, process := range jpsOutput {
		jinfoOutput, err := runJinfo(javaHome, process)
		if err != nil {
			log.Warn().Err(err).Str("pid", process.ProcessID).Msg("could not inspect running jvm with jinfo")
			continue
		}
		parsedJInfo := updateJavaInfoWithJinfoData(jinfoOutput, process)
		runningProcsJavaInfos = append(runningProcsJavaInfos, parsedJInfo)
	}

	return runningProcsJavaInfos
}
