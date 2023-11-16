package ojdmcollector

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func CreateCSVReport(csvPath string, javaFullInfo []JavaInfoRunningProcs) {

	fmt.Println("Creating csv report...")

	file, err := os.Create(csvPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"HostName",
		"AppDirName",
		"DynLibBinPath",
		"JavaBinPath",
		"JavaCBinPath",
		"IsJDK",
		"BaseDir",
		"JavaHome",
		"JavaRuntimeName",
		"JavaRuntimeVersion",
		"JavaVendor",
		"JavaVersion",
		"JavaVersionDate",
		"JavaVMName",
		"JavaVMVendor",
		"JavaVMVersion",
		"CommandLine",
		"ProcessRunning",
	}

	writer.Write(header)

	for _, value := range javaFullInfo {

		stringData := []string{
			value.HostName,
			value.AppDirName,
			value.DynLibBinPath,
			value.JavaBinPath,
			value.JavaCBinPath,
			strconv.FormatBool(value.IsJDK),
			value.BaseDir,
			value.JavaHome,
			value.JavaRuntimeName,
			value.JavaRuntimeVersion,
			value.JavaVendor,
			value.JavaVersion,
			value.JavaVersionDate,
			value.JavaVMName,
			value.JavaVMVendor,
			value.JavaVMVersion,
			value.CommandLine,
			strconv.FormatBool(value.ProcessRunning),
		}

		err := writer.Write(stringData)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Done!")

}
