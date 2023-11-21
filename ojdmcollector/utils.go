package ojdmcollector

import (
	"fmt"
	"os"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func sliceContains(str string, strSlice []string) bool {
	for _, name := range strSlice {
		if strings.EqualFold(str, name) {
			return true
		}
	}
	return false
}

func fileExists(fp string) bool {
	if _, err := os.Stat(fp); err == nil {
		return true
	}
	return false
}

func getRunningProcCommands() []string {

	runningProcs := []string{}

	processes, procerr := process.Processes()
	if procerr != nil {
		return runningProcs
	}

	for _, p := range processes {
		procRunning, runerr := p.IsRunning()
		if runerr != nil || !procRunning {
			continue
		}

		name, err := p.Name()
		if err != nil {
			continue
		}

		cmdline, cmderr := p.Cmdline()
		if cmderr != nil {
			continue
		}

		procCmd := "CommandLine: " + cmdline + " ProcName: " + name

		fmt.Println(procCmd)

		runningProcs = append(runningProcs, procCmd)
	}

	return runningProcs

}
