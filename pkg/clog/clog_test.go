package clog

import (
	"os"
	"testing"
)

// basic logger [pass]
func Test_clog(t *testing.T) {
	template := `{_temp_shortpath}:{_temp_linenum} - [{_temp_level}]`
	h := NewPlainTextHandler(os.Stderr, DEBUG, template)
	l := NewLogger(h)
	l.Debug("debug message")
	l.Info("info message")
	l.Warn("warn message")
	l.Error("error message")
	// l.Panic("panic message")
}

// template extension [pass]
func Test_clog_extension(t *testing.T) {
	// template := `{_temp_timestamp} [{_temp_level}] {_temp_shortpath}:{_temp_linenum} {_temp_prefix}`
	// _format_timestamp := func() string {
	// 	return time.Now().UTC().Format("2006/01/02 15:04:05")
	// }
	// h := NewPlainTextHandler(os.Stderr, DEBUG, template)
	// h.With(TEMP_TIMESTAMP, _format_timestamp)
	// l := NewLogger(h)
	// SetDefault(l)
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
}

// test change log level [passed]
func Test_change_level(t *testing.T) {
	DefaultLevel(DEBUG)
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
	DefaultLevel(WARN)
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
}
