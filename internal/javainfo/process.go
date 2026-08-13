package javainfo

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

type process struct {
	ProcessID   string
	CommandLine string
}

func parseJpsOutput(output string) []process {
	var processes []process
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if line == "" || strings.Contains(line, "jps.Jps -mvl") {
			continue // Skip empty lines and the jps command itself
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			processes = append(processes, process{
				ProcessID:   fields[0],
				CommandLine: strings.Join(fields[1:], " "),
			})
		}
	}

	return processes
}

func runJps(javaHome string) ([]process, error) {
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

func runJinfo(javaHome string, proc process) (string, error) {
	jinfoPath := filepath.Join(javaHome, "bin", "jinfo")
	if runtime.GOOS == "windows" {
		jinfoPath += ".exe"
	}

	cmd := exec.Command(jinfoPath, "-sysprops", proc.ProcessID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func recordFromJinfo(jinfoOutput string, proc process) Record {
	// Update JavaInfo with data extracted from jinfoOutput
	var javaInfo Record

	javaInfo.HostName = getHostName()
	javaInfo.ProcessPath = normalizePath(findRegexInText(`(?m)^user.dir=(.*)`, jinfoOutput))
	javaInfo.CommandLine = proc.CommandLine
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

func runningProcesses(javaHome string, log *zerolog.Logger) []Record {
	var runningProcsJavaInfos []Record

	jpsOutput, err := runJps(javaHome)
	if err != nil {
		log.Warn().Err(err).Str("java_home", javaHome).
			Msg("could not run jps, running jvm instances will not be identified")
		return runningProcsJavaInfos
	}

	log.Debug().Int("count", len(jpsOutput)).Msg("jps reported running jvm instances")
	for _, proc := range jpsOutput {
		jinfoOutput, err := runJinfo(javaHome, proc)
		if err != nil {
			log.Warn().Err(err).Str("pid", proc.ProcessID).Msg("could not inspect running jvm with jinfo")
			continue
		}
		parsedJInfo := recordFromJinfo(jinfoOutput, proc)
		runningProcsJavaInfos = append(runningProcsJavaInfos, parsedJInfo)
	}

	return runningProcsJavaInfos
}
