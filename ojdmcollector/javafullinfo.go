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

		// A JVM can run from a location no search path reaches: an app bundle, a
		// user directory, a container mount. The home it reports is an
		// installation the scan missed, so collect it like any other, otherwise
		// the row merges with nothing and carries no version, no IsJDK and no
		// shared library.
		if missed := unscannedProcessHomes(javaInfos, versionInfos); len(missed) > 0 {
			log.Info().Strs("java_homes", missed).
				Msg("collecting installations that only a running process revealed")
			versionInfos = append(versionInfos, GetJavaVersionInfos(missed, log)...)
		}
	} else {
		log.Warn().Msg("did not find the jinfo and jps binaries, running processes will not be identified")
	}

	mergedJavaInfo := mergeSlices(javaInfos, versionInfos)

	log.Info().Int("installations", len(versionInfos)).Int("running_processes", len(javaInfos)).
		Int("rows", len(mergedJavaInfo)).Msg("collection complete")

	return mergedJavaInfo
}

// unscannedProcessHomes lists the java homes that running processes report and
// the scan did not reach, deduplicated, so each is collected once.
func unscannedProcessHomes(processInfo, versionInfo []JavaInfoRunningProcs) []string {
	known := make(map[string]struct{}, len(versionInfo))
	for _, info := range versionInfo {
		known[normalizePath(info.JavaHome)] = struct{}{}
	}

	var missed []string
	for _, info := range processInfo {
		home := normalizePath(info.JavaHome)
		_, seen := known[home]
		if home == "" || seen {
			continue
		}
		known[home] = struct{}{}
		missed = append(missed, home)
	}

	return missed
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
