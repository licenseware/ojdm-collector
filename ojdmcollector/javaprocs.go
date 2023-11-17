package ojdmcollector

import (
	"fmt"
	"runtime"
)

func getJavaProcInfo(javaBinInfo []JavaInfoRunningProcs) []JavaInfoRunningProcs {

	procCommands := getJavaRunningProcsCommands(javaBinInfo)

	javaBinProcs := []JavaInfoRunningProcs{}
	for _, procCmd := range procCommands {
		for _, jinfo := range javaBinInfo {

			if runtime.GOOS == "windows" {
				if !subStringExistsInText(jinfo.BaseDir, procCmd) {
					continue
				}
			} else {
				if !subStringExistsInText(jinfo.JavaBinPath, procCmd) {
					continue
				}
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
				ProcessRunning:     true,
			}

			javaBinProcs = append(javaBinProcs, jproc)
		}
	}

	return javaBinProcs
}
