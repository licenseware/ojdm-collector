package main

import (
	"flag"
	"fmt"
	ojdmc "ojdmcollector/ojdmcollector"
)

func main() {

	fmt.Println("Licenseware OJDM Collector - Gather all java info in one place")

	csvReportPath := flag.String("csv", "report.csv", "Path to csv report.")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("OJDMCollector - Utility to find JVMs/JDKs report their versions and related running processes")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("     $ ojdm-collector")
		fmt.Println("     $ ojdm-collector -csv=/path/to/csvreport.csv")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	javaInfoRunningProcs := ojdmc.GetFullJavaInfo()

	fmt.Println("\nJava Info with Running Processes:")
	ojdmc.Pprint(javaInfoRunningProcs)

	ojdmc.CreateCSVReport(*csvReportPath, javaInfoRunningProcs)

}
