// Package logging provides structured logging setup.
package logging

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// parseLevel maps a textual log level to a zerolog.Level, defaulting to InfoLevel.
func parseLevel(logLevel string) zerolog.Level {
	switch strings.ToLower(logLevel) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "critical":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

// minLevelWriter drops events below its level, so the console can stay readable
// at info while the debug file keeps everything.
type minLevelWriter struct {
	w     io.Writer
	level zerolog.Level
}

func (m minLevelWriter) Write(p []byte) (int, error) {
	return m.w.Write(p)
}

func (m minLevelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level < m.level {
		return len(p), nil
	}
	return m.w.Write(p)
}

// New creates a new zerolog logger. The console sink is human readable and
// filtered to consoleLevel; any extra writers receive every event as raw JSON,
// which is what the debug file and the support hand-over rely on.
func New(consoleLevel string, isProd bool, extra ...io.Writer) *zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	console := minLevelWriter{
		w: zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    isProd,
		},
		level: parseLevel(consoleLevel),
	}

	var w io.Writer = console
	if len(extra) > 0 {
		w = zerolog.MultiLevelWriter(append([]io.Writer{console}, extra...)...)
	}

	l := zerolog.New(w).Level(zerolog.DebugLevel).With().Caller().Timestamp().Logger()

	return &l
}

// OrNop returns l, or a no-op logger when l is nil, so callers can log
// unconditionally instead of guarding every log site.
func OrNop(l *zerolog.Logger) *zerolog.Logger {
	if l == nil {
		nop := zerolog.Nop()
		return &nop
	}
	return l
}
