package ojdmcollector

import (
	"sort"

	"github.com/rs/zerolog"
)

func CollectJavaInfo(searchPaths []string, log *zerolog.Logger) []JavaInfoRunningProcs {

	var javaInfos []JavaInfoRunningProcs
	var versionInfos []JavaInfoRunningProcs
	javaBasePAths := getJavaSharedLibPaths(searchPaths, log)

	versionInfos = GetJavaVersionInfos(javaBasePAths, log)
	sort.Slice(versionInfos, func(i, j int) bool {
		return versionInfos[i].JavaVersion > versionInfos[j].JavaVersion
	})

	var toolFound *JavaInfoRunningProcs
	for _, info := range versionInfos {
		if info.JpsJinfoPresent {
			toolFound = &info
			break
		}
	}

	if toolFound != nil {
		javaProcesses := GetJavaProcessInfo(toolFound.JavaHome, log)
		javaInfos = append(javaInfos, javaProcesses...)
	} else {
		log.Warn().Msg("did not find the jinfo and jps binaries, running processes will not be identified")
	}

	mergedJavaInfo := mergeSlices(javaInfos, versionInfos)

	log.Info().Int("installations", len(versionInfos)).Int("running_processes", len(javaInfos)).
		Int("rows", len(mergedJavaInfo)).Msg("collection complete")

	return mergedJavaInfo
}

func mergeSlices(processInfo, versionInfo []JavaInfoRunningProcs) []JavaInfoRunningProcs {
	mergedSlice := make([]JavaInfoRunningProcs, 0)
	versionMap := make(map[string]*JavaInfoRunningProcs)
	processMap := make(map[string]bool)

	// Create a map from versionInfo
	for i, vInfoItem := range versionInfo {
		normalizedJavaHome := normalizePath(vInfoItem.JavaHome)
		versionMap[normalizedJavaHome] = &versionInfo[i]
	}

	// Iterate over processInfo and merge or add unique items
	for _, pInfoItem := range processInfo {
		normalizedJavaHome := normalizePath(pInfoItem.JavaHome)
		if vInfoItem, exists := versionMap[normalizedJavaHome]; exists {
			// Merge with versionInfo item
			mergedItem := *vInfoItem
			mergedItem.ProcessRunning = pInfoItem.ProcessRunning
			mergedItem.ProcessPath = pInfoItem.ProcessPath
			mergedItem.CommandLine = pInfoItem.CommandLine
			mergedSlice = append(mergedSlice, mergedItem)

			// Mark as processed
			processMap[normalizedJavaHome] = true
		} else {
			// Add unique processInfo item
			mergedSlice = append(mergedSlice, pInfoItem)
		}
	}

	// Add any versionInfo items that weren't merged
	for _, vInfoItem := range versionInfo {
		normalizedJavaHome := normalizePath(vInfoItem.JavaHome)
		if !processMap[normalizedJavaHome] {
			mergedSlice = append(mergedSlice, vInfoItem)
		}
	}

	return mergedSlice
}
