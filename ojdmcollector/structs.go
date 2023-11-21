package ojdmcollector

type JavaInfoRunningProcs struct {
	HostName           string
	DynLibBinPath      string
	JavaBinPath        string
	JavaCBinPath       string
	IsJDK              bool
	JavaHome           string
	JavaRuntimeName    string
	JavaRuntimeVersion string
	JavaVendor         string
	JavaVersion        string
	JavaVersionDate    string
	JavaVMName         string
	JavaVMVendor       string
	JavaVMVersion      string
	ProcessRunning     bool
	ProcessPath        string
	CommandLine        string
}
