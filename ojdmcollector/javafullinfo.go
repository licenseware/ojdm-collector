package ojdmcollector

import (
	"fmt"
	"sort"
)

func CollectJavaInfo(searchPaths []string) []JavaInfoRunningProcs {

	var javaInfos []JavaInfoRunningProcs
	var versionInfos []JavaVersionInfo
	javaLibPaths := getJavaSharedLibPaths(nil) // Assuming this function is defined in javasearch.go

	// hostInfo := GetHostInfo()

	versionInfos = GetJavaVersionInfos(javaLibPaths)
	sort.Slice(versionInfos, func(i, j int) bool {
		return versionInfos[i].JavaVersion > versionInfos[j].JavaVersion
	})

	var toolFound *JavaVersionInfo
	for _, info := range versionInfos {
		if info.JpsJinfoPresent {
			toolFound = &info
			break
		}
	}

	if toolFound != nil {
		javaProcesses := GetJavaProcessInfo(toolFound.JavaHome)
		javaInfos = append(javaInfos, javaProcesses...)
	} else {
		fmt.Println("Did not find the jinfo and jps binaries, running processes will not be identified.")
	}

	return javaInfos
}
