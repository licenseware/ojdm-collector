package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"javafinder/javafinder"
)

func pprint(any interface{}) {
	empJSON, _ := json.MarshalIndent(any, "", "  ")
	fmt.Printf("\n%s\n", string(empJSON))
}

func main() {

	csvReportPath := flag.String("csv", "report.csv", "Path to csv report.")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("JavaFinder - Utility to find JVMs/JDKs report their versions and related running processes")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("     $ javafinder")
		fmt.Println("     $ javafinder -csv=/path/to/csvreport.csv")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	javaInfoRunningProcs := javafinder.GetFullJavaInfo()

	fmt.Println("Java Info with Running Processes:")
	pprint(javaInfoRunningProcs)

	javafinder.CreateCSVReport(*csvReportPath, javaInfoRunningProcs)

}
