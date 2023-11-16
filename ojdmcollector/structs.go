package ojdmcollector

type JavaInfoRunningProcs struct {
	HostName           string
	AppDirName         string
	DynLibBinPath      string
	JavaBinPath        string
	JavaCBinPath       string
	IsJDK              bool
	BaseDir            string
	JavaHome           string
	JavaRuntimeName    string
	JavaRuntimeVersion string
	JavaVendor         string
	JavaVersion        string
	JavaVersionDate    string
	JavaVMName         string
	JavaVMVendor       string
	JavaVMVersion      string
	CommandLine        string
	ProcessRunning     bool
}
