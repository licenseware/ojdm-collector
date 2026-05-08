package main

import (
	"flag"
	"fmt"
	ojdmc "ojdmcollector/ojdmcollector"
	"os"
	"strings"
)

func main() {

	fmt.Print("\n\nLicenseware OJDM Collector - Gather all java info in one place\n\n")

	csvReportPath := flag.String("output-path", "report.csv", "Optional: Path to csv report.")
	searchPaths := flag.String("search-paths", "", "Optional: List of paths separated by comma where to search for java info.")
	searchPathsFile := flag.String("search-paths-file", "", "Optional: Path to a file containing additional search paths, one path per line.")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("OJDMCollector - Utility to find JVMs/JDKs report their versions and related running processes")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("     $ ojdm-collector")
		fmt.Println("     $ ojdm-collector -output-path=/path/to/csvreport.csv")
		fmt.Println("     $ ojdm-collector -search-paths=/home,/oracle,/opt")
		fmt.Println("     $ ojdm-collector -search-paths-file=/path/to/search-paths.txt")
		fmt.Println("     $ ojdm-collector -search-paths=/home,/usr,/opt -output-path=/path/to/csvreport.csv")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	if !strings.HasSuffix(*csvReportPath, ".csv") {
		fmt.Println("Error: Invalid output report path. The report path must be a csv file.")
		return
	}

	trimSpaths := parseSearchPaths(*searchPaths)

	fileSearchPaths, err := parseSearchPathsFile(*searchPathsFile)
	if err != nil {
		fmt.Println("Error reading search paths file:", err)
		return
	}
	trimSpaths = append(trimSpaths, fileSearchPaths...)

	javaInfoRunningProcs := ojdmc.CollectJavaInfo(trimSpaths)

	fmt.Println("\nJava Info with Running Processes:")
	// ojdmc.Println(javaInfoRunningProcs)

	ojdmc.CreateCSVReport(*csvReportPath, javaInfoRunningProcs)

}

func parseSearchPaths(searchPaths string) []string {
	spaths := strings.Split(searchPaths, ",")
	trimSpaths := []string{}
	for _, path := range spaths {
		path = strings.TrimSpace(path)
		if path != "" {
			trimSpaths = append(trimSpaths, path)
		}
	}

	return trimSpaths
}

func parseSearchPathsFile(searchPathsFile string) ([]string, error) {
	if strings.TrimSpace(searchPathsFile) == "" {
		return nil, nil
	}

	fileContent, err := os.ReadFile(searchPathsFile)
	if err != nil {
		return nil, err
	}

	searchPaths := []string{}
	lines := strings.Split(string(fileContent), "\n")
	for _, line := range lines {
		path := strings.TrimSpace(line)
		if path == "" || strings.HasPrefix(path, "#") {
			continue
		}
		searchPaths = append(searchPaths, path)
	}

	return searchPaths, nil
}
