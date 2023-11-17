package ojdmcollector

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/process"
)

func getExeNamesOfJavaProgramsOnWindows(javaBinInfo []JavaInfoRunningProcs) []string {

	javaExeNames := []string{"java.exe", "javac.exe", "javaw.exe"}

	for _, jinfo := range javaBinInfo {
		err := filepath.Walk(jinfo.BaseDir, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				if os.IsPermission(err) {
					return filepath.SkipDir
				}
				return err
			}

			if !info.IsDir() && strings.HasSuffix(info.Name(), ".exe") {
				javaExeNames = append(javaExeNames, info.Name())
			}

			return nil
		})

		if err != nil {
			fmt.Println("Encountered error", err)
		}
	}

	return javaExeNames

}

func isJavaProcess(procName string, javaProcNames []string) bool {
	for _, jname := range javaProcNames {
		if strings.EqualFold(procName, jname) {
			return true
		}
	}
	return false
}

func getJavaRunningProcsCommands(javaBinInfo []JavaInfoRunningProcs) []string {

	javaProcNames := []string{"java"}

	if runtime.GOOS == "windows" {
		javaProcNames = getExeNamesOfJavaProgramsOnWindows(javaBinInfo)
	}

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

		if isJavaProcess(name, javaProcNames) {
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
	re := regexp.MustCompile(regexp.QuoteMeta(regex))
	match := re.FindStringSubmatch(text)
	return len(match) > 0
}
