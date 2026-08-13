package ojdmcollector

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// ReportMeta describes the run that produced a report rather than any single
// java installation, so every row carries the same values.
type ReportMeta struct {
	CollectedAt time.Time
	Version     string
}

func CreateCSVReport(csvPath string, javaFullInfo []JavaInfoRunningProcs, meta ReportMeta, log *zerolog.Logger) error {

	collectedAt := meta.CollectedAt.UTC().Format(time.RFC3339)

	log.Info().Str("path", csvPath).Str("collected_at", collectedAt).Msg("creating csv report")

	file, err := os.Create(csvPath)
	if err != nil {
		return fmt.Errorf("creating csv report: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"CollectedAtUTC",
		"CollectorVersion",
		"HostName",
		"DynLibBinPath",
		"JavaBinPath",
		"JavaCBinPath",
		"IsJDK",
		"JavaHome",
		"JavaRuntimeName",
		"JavaRuntimeVersion",
		"JavaVendor",
		"JavaVersion",
		"JavaVersionDate",
		"JavaVMName",
		"JavaVMVendor",
		"JavaVMVersion",
		"ProcessPath",
		"ProcessRunning",
		"CommandLine",
		"HostLogicalProcessors",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("writing csv header: %w", err)
	}

	for _, value := range javaFullInfo {

		stringData := []string{
			collectedAt,
			meta.Version,
			value.HostName,
			value.DynLibBinPath,
			value.JavaBinPath,
			value.JavaCBinPath,
			strconv.FormatBool(value.IsJDK),
			value.JavaHome,
			value.JavaRuntimeName,
			value.JavaRuntimeVersion,
			value.JavaVendor,
			value.JavaVersion,
			value.JavaVersionDate,
			value.JavaVMName,
			value.JavaVMVendor,
			value.JavaVMVersion,
			value.ProcessPath,
			strconv.FormatBool(value.ProcessRunning),
			value.CommandLine,
			strconv.Itoa(value.HostLogicalProcessors),
		}

		if err := writer.Write(stringData); err != nil {
			return fmt.Errorf("writing csv row: %w", err)
		}
	}

	log.Info().Int("rows", len(javaFullInfo)).Msg("csv report written")

	return nil
}
