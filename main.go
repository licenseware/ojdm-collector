package main

import (
	"flag"
	"fmt"
	ojdmc "ojdmcollector/ojdmcollector"
	"strings"
)

func main() {

	fmt.Print("\n\nLicenseware OJDM Collector - Gather all java info in one place\n\n")

	csvReportPath := flag.String("csv", "report.csv", "Optional: Path to csv report.")
	searchPaths := flag.String("paths", "", "Optional: List of paths separated by comma where to search for java info.")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("OJDMCollector - Utility to find JVMs/JDKs report their versions and related running processes")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("     $ ojdm-collector")
		fmt.Println("     $ ojdm-collector -csv=/path/to/csvreport.csv")
		fmt.Println("     $ ojdm-collector -paths=/home,/oracle,/opt")
		fmt.Println("     $ ojdm-collector -paths=/home,/usr,/opt -csv=/home/alin/Downloads/ubuntu_report.csv")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	spaths := strings.Split(*searchPaths, ",")
	trimSpaths := []string{}
	for _, path := range spaths {
		path = strings.TrimSpace(path)
		if path != "" {
			trimSpaths = append(trimSpaths, path)
		}
	}

	javaInfoRunningProcs := ojdmc.CollectJavaInfo(trimSpaths)

	fmt.Println("\nJava Info with Running Processes:")
	// ojdmc.Println(javaInfoRunningProcs)

	ojdmc.CreateCSVReport(*csvReportPath, javaInfoRunningProcs)

}
