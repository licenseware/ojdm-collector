package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// DefaultLogPath returns the debug log location for a given csv report path,
// which is a logs/ directory beside the report. Operators are told where the
// report lands, so the log is findable with the same instruction.
func DefaultLogPath(csvReportPath string) string {
	return filepath.Join(filepath.Dir(csvReportPath), "logs", "debug.log")
}

// rotationStamp names the archive of a previous log by when it was last
// written, falling back to the current UTC time when that is unavailable.
func rotationStamp(info os.FileInfo, err error) string {
	if err == nil && !info.ModTime().IsZero() {
		return info.ModTime().UTC().Format("20060102T150405Z")
	}
	return time.Now().UTC().Format("20060102T150405Z")
}

// rotate moves an existing log aside so each run starts with a clean file
// without destroying the evidence of the run before it. A failed rerun is the
// most likely thing an operator does, and it must not overwrite the failure.
func rotate(logPath string) (string, error) {
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		return "", nil
	}

	extension := filepath.Ext(logPath)
	archived := fmt.Sprintf("%s-%s%s",
		logPath[:len(logPath)-len(extension)], rotationStamp(info, err), extension)

	if err := os.Rename(logPath, archived); err != nil {
		return "", fmt.Errorf("rotating previous log: %w", err)
	}

	return archived, nil
}

// OpenFileSink prepares the debug log, rotating any previous one, and returns
// the open file along with the path of the archive it created, if any.
func OpenFileSink(logPath string) (*os.File, string, error) {
	if err := os.MkdirAll(filepath.Dir(logPath), 0o755); err != nil {
		return nil, "", fmt.Errorf("creating log directory: %w", err)
	}

	archived, err := rotate(logPath)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Create(logPath)
	if err != nil {
		return nil, "", fmt.Errorf("creating log file: %w", err)
	}

	return file, archived, nil
}
