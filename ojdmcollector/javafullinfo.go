package ojdmcollector

func GetFullJavaInfo() []JavaInfoRunningProcs {
	javaBinInfo := getJavaBinInfo()
	javaProcInfo := getJavaProcInfo(javaBinInfo)

	javaInfoRunningProcs := []JavaInfoRunningProcs{}
	for _, jinfo := range javaBinInfo {
		found := false
		for _, jproc := range javaProcInfo {
			if jinfo.JavaBinPath == jproc.JavaBinPath {
				javaInfoRunningProcs = append(javaInfoRunningProcs, jproc)
				found = true
				break
			}
		}
		if !found {
			javaInfoRunningProcs = append(javaInfoRunningProcs, jinfo)
		}
	}

	return javaInfoRunningProcs
}
