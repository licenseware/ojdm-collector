package ojdmcollector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func fileExists(fp string) bool {
	if _, err := os.Stat(fp); err == nil {
		return true
	}
	return false
}

func getRunningProcCommands() []ProcessInfo {

	runningProcs := []ProcessInfo{}

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
		if cmderr != nil || len(cmdline) == 0 {
			continue
		}

		procDir, pdirErr := p.Exe()
		if pdirErr != nil {
			continue
		}

		procCmd := ProcessInfo{
			Name:        name,
			ProcDir:     procDir,
			CommandLine: cmdline,
		}

		runningProcs = append(runningProcs, procCmd)

	}

	return runningProcs

}

func upDir(path string, n int) string {

	splitedPath := strings.Split(path, string(os.PathSeparator))
	splitedPath = splitedPath[:len(splitedPath)-n]
	upPath := filepath.Join(splitedPath...)

	if runtime.GOOS != "windows" {
		return filepath.Join(string(os.PathSeparator), upPath)
	}

	return upPath
}

func Pprint(any interface{}) {
	empJSON, _ := json.MarshalIndent(any, "", "  ")
	fmt.Printf("\n%s\n", string(empJSON))
}
