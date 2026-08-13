// Package javainfo finds the java installations on a host, inspects each one
// and reports the running JVM instances it can attribute to them.
package javainfo

// A Record describes one java installation, optionally carrying the details of
// a JVM instance observed running out of it.
type Record struct {
	HostName              string
	DynLibBinPath         string
	JavaBinPath           string
	JavaCBinPath          string
	IsJDK                 bool
	JavaHome              string
	JavaRuntimeName       string
	JavaRuntimeVersion    string
	JavaVendor            string
	JavaVersion           string
	JavaVersionDate       string
	JavaVMName            string
	JavaVMVendor          string
	JavaVMVersion         string
	ProcessRunning        bool
	ProcessPath           string
	CommandLine           string
	HostLogicalProcessors int
	JpsJinfoPresent       bool
}
