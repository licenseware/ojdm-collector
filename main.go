package main

import (
	"encoding/json"
	"flag"
	"fmt"
	ojdmc "ojdmcollector/ojdmcollector"
)

func pprint(any interface{}) {
	empJSON, _ := json.MarshalIndent(any, "", "  ")
	fmt.Printf("\n%s\n", string(empJSON))
}

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

	fmt.Println("Java Info with Running Processes:")
	pprint(javaInfoRunningProcs)

	ojdmc.CreateCSVReport(*csvReportPath, javaInfoRunningProcs)

}
