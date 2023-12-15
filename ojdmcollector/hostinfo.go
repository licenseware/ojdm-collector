package ojdmcollector

import (
	"runtime"
)

type HostInfo struct {
	Processors int
	Cores      int
}

func getProcessorAndCoreCounts() (int, int) {
	return runtime.NumCPU(), runtime.GOMAXPROCS(0)
}

func GetHostInfo() HostInfo {
	processors, cores := getProcessorAndCoreCounts()

	return HostInfo{
		Processors: processors,
		Cores:      cores,
	}
}
