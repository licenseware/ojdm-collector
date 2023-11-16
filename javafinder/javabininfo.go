package javafinder

import (
	"fmt"
	"os"
)

func getJavaBinInfo() []JavaInfoRunningProcs {

	javaBinPaths := getJavaBinPaths()

	hostname, _ := os.Hostname()

	JavaInfoRunningProcsSlice := []JavaInfoRunningProcs{}
	for _, jp := range javaBinPaths {

		versionSettings, err := getJavaFullVersionSettings(jp.JavaBinPath)
		if err != nil {
			fmt.Println("Cannot get full version info from java binary at path ", jp.JavaBinPath)
			continue
		}

		vInfo := extractInfoFromFullVersionSettings(versionSettings)

		jinfo := JavaInfoRunningProcs{
			HostName:           hostname,
			AppDirName:         jp.AppDirName,
			DynLibBinPath:      jp.DynLibBinPath,
			JavaBinPath:        jp.JavaBinPath,
			JavaCBinPath:       jp.JavaCBinPath,
			IsJDK:              jp.IsJDK,
			BaseDir:            jp.BaseDir,
			JavaHome:           vInfo.JavaHome,
			JavaRuntimeName:    vInfo.JavaRuntimeName,
			JavaRuntimeVersion: vInfo.JavaRuntimeVersion,
			JavaVendor:         vInfo.JavaVendor,
			JavaVersion:        vInfo.JavaVersion,
			JavaVersionDate:    vInfo.JavaVersionDate,
			JavaVMName:         vInfo.JavaVMName,
			JavaVMVendor:       vInfo.JavaVMVendor,
			JavaVMVersion:      vInfo.JavaVMVersion,
		}

		JavaInfoRunningProcsSlice = append(JavaInfoRunningProcsSlice, jinfo)

	}

	return JavaInfoRunningProcsSlice

}
