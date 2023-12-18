package ojdmcollector

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
	javaInfo.ProcessPath = normalizePath(findRegexInText(`user.dir=(.*)`, jinfoOutput))
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

	return javaInfo
}

func GetJavaProcessInfo(javaHome string) []JavaInfoRunningProcs {
	var runningProcsJavaInfos []JavaInfoRunningProcs

	jpsOutput, err := runJps(javaHome)
	if err != nil {
		fmt.Println("Could not run jps utility to identify running JVM instances")
		return runningProcsJavaInfos
	}

	for _, process := range jpsOutput {
		jinfoOutput, err := runJinfo(javaHome, process)
		if err != nil {
			continue // Optionally handle or log the error
		}
		parsedJInfo := updateJavaInfoWithJinfoData(jinfoOutput, process)
		runningProcsJavaInfos = append(runningProcsJavaInfos, parsedJInfo)
	}

	return runningProcsJavaInfos
}
