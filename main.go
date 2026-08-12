package main

import (
	"flag"
	"fmt"
	"ojdmcollector/internal/logging"
	ojdmc "ojdmcollector/ojdmcollector"
	"os"
	"runtime"
	"strings"
	"time"
)

// Stamped at build time by GoReleaser via -ldflags -X main.<name>=...
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func versionString() string {
	return fmt.Sprintf("ojdm-collector %s (commit %s, built %s by %s)", version, commit, date, builtBy)
}

func main() {

	// One stamp for the whole run: the csv report, the log lines and the
	// support hand-over all have to refer to the same instant.
	runStartedAt := time.Now().UTC()

	csvReportPath := flag.String("output-path", "report.csv", "Optional: Path to csv report.")
	searchPaths := flag.String("search-paths", "", "Optional: List of paths separated by comma where to search for java info.")
	searchPathsFile := flag.String("search-paths-file", "", "Optional: Path to a file containing additional search paths, one path per line.")
	logPath := flag.String("log-path", "", "Optional: Path to the debug log. Defaults to a logs/ directory beside the csv report.")
	logLevel := flag.String("log-level", "info", "Optional: Console verbosity (debug, info, warn, error). The debug log always records everything.")
	showVersion := flag.Bool("version", false, "Print the collector version and exit.")

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
		fmt.Println("     $ ojdm-collector -log-path=/path/to/debug.log -log-level=debug")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVersion {
		fmt.Println(versionString())
		return
	}

	fmt.Print("\n\nLicenseware OJDM Collector - Gather all java info in one place\n\n")

	if !strings.HasSuffix(*csvReportPath, ".csv") {
		fmt.Println("Error: Invalid output report path. The report path must be a csv file.")
		return
	}

	resolvedLogPath := *logPath
	if resolvedLogPath == "" {
		resolvedLogPath = logging.DefaultLogPath(*csvReportPath)
	}

	logFile, archived, err := logging.OpenFileSink(resolvedLogPath)
	if err != nil {
		fmt.Println("Error preparing the debug log:", err)
		return
	}
	defer logFile.Close()

	runLog := logging.New(*logLevel, true, logFile).With().
		Time("run_started_at", runStartedAt).Logger()
	log := &runLog
	log.Info().
		Str("log_path", resolvedLogPath).
		Str("csv_report_path", *csvReportPath).
		Str("version", version).
		Str("commit", commit).
		Str("os", runtime.GOOS).
		Str("arch", runtime.GOARCH).
		Strs("args", os.Args[1:]).
		Msg("ojdm collector starting")
	if archived != "" {
		log.Info().Str("archived_log", archived).Msg("rotated previous debug log")
	}

	trimSpaths := parseSearchPaths(*searchPaths)

	fileSearchPaths, err := parseSearchPathsFile(*searchPathsFile)
	if err != nil {
		log.Error().Err(err).Str("path", *searchPathsFile).Msg("could not read search paths file")
		return
	}
	trimSpaths = append(trimSpaths, fileSearchPaths...)

	javaInfoRunningProcs := ojdmc.CollectJavaInfo(trimSpaths, log)

	reportMeta := ojdmc.ReportMeta{CollectedAt: runStartedAt, Version: version}

	if err := ojdmc.CreateCSVReport(*csvReportPath, javaInfoRunningProcs, reportMeta, log); err != nil {
		log.Error().Err(err).Str("path", *csvReportPath).Msg("could not write csv report")
		return
	}

	log.Info().Str("csv_report_path", *csvReportPath).Str("log_path", resolvedLogPath).
		Msg("done, attach the log file if support asks for it")
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
