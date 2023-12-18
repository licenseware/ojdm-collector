package ojdmcollector

import (
	"runtime"
)

type HostInfo struct {
	Processors int
	Cores      int
}

func GetHostInfo() HostInfo {

	return HostInfo{
		Processors: runtime.NumCPU(),
		Cores:      runtime.GOMAXPROCS(0),
	}
}
