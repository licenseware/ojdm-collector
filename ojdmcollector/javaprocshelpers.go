package ojdmcollector

import (
	"regexp"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func getJavaRunningProcsCommands() []string {

	javaRunningProcs := []string{}
	processes, procerr := process.Processes()
	if procerr != nil {
		return javaRunningProcs
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

		if strings.ToLower(name) == "java" {
			cmdline, cmderr := p.Cmdline()
			if cmderr != nil {
				continue
			}

			javaRunningProcs = append(javaRunningProcs, cmdline)
		}
	}

	return javaRunningProcs

}

func subStringExistsInText(regex, text string) bool {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	return len(match) > 0
}
