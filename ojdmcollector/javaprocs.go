package ojdmcollector

import (
	"fmt"
	"strings"
)

func getJavaProcInfo(javaBinInfo []JavaInfoRunningProcs) []JavaInfoRunningProcs {

	procCommands := getJavaRunningProcsCommands(javaBinInfo)

	javaBinProcs := []JavaInfoRunningProcs{}
	for _, procCmd := range procCommands {
		for _, jinfo := range javaBinInfo {

			if !subStringExistsInText(jinfo.BaseDir, procCmd) {
				continue
			}

			fmt.Println("Found running process ", jinfo.AppDirName)

			jproc := JavaInfoRunningProcs{
				HostName:           jinfo.HostName,
				AppDirName:         jinfo.AppDirName,
				DynLibBinPath:      jinfo.DynLibBinPath,
				JavaBinPath:        jinfo.JavaBinPath,
				JavaCBinPath:       jinfo.JavaCBinPath,
				IsJDK:              jinfo.IsJDK,
				BaseDir:            jinfo.BaseDir,
				JavaHome:           jinfo.JavaHome,
				JavaRuntimeName:    jinfo.JavaRuntimeName,
				JavaRuntimeVersion: jinfo.JavaRuntimeVersion,
				JavaVendor:         jinfo.JavaVendor,
				JavaVersion:        jinfo.JavaVersion,
				JavaVersionDate:    jinfo.JavaVersionDate,
				JavaVMName:         jinfo.JavaVMName,
				JavaVMVendor:       jinfo.JavaVMVendor,
				JavaVMVersion:      jinfo.JavaVMVersion,
				CommandLine:        procCmd,
				ProcessPath:        strings.Split(procCmd, " ")[0],
				ProcessRunning:     true,
			}

			javaBinProcs = append(javaBinProcs, jproc)
		}
	}

	return javaBinProcs
}
