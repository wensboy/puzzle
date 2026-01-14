// Package clod provide a simple colorable log.
package clog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/valyala/fasttemplate"
	"github.com/wendisx/puzzle/pkg/palette"
)

const (
	_format_debug   = "DEBUG"
	_format_info    = "INFO"
	_format_warn    = "WARN"
	_format_error   = "ERROR"
	_format_panic   = "PANIC"
	_format_fatal   = "FATAL"
	_format_default = "???"

	_max_level    = 1 << 5
	_offset_level = 1 << 2

	// caller(exactly need to show) -> caller(clog) -> clog.Log() -> logger.log() -> runtime.Caller -> extern low level
	_default_skip_step = 5
	_default_template  = `{_temp_timestamp} {_temp_shortpath}:{_temp_linenum} [{_temp_level}] {_temp_prefix}`
)

const (
	/* internal record info */
	TEMP_LONGPATH  = "_temp_longpath"  // long file path
	TEMP_SHORTPATH = "_temp_shortpath" // short file path
	TEMP_LINENUM   = "_temp_linenum"   // line number
	TEMP_PREFIX    = "_temp_prefix"    //  log prefix
	TEMP_LEVEL     = "_temp_level"     // log level
	/* external record info */
	TEMP_TIMESTAMP = "_temp_timestamp" // log time stamp
)

const (
	DEBUG          = slog.LevelDebug // same as slog
	INFO           = slog.LevelInfo  // same as slog
	WARN           = slog.LevelWarn  // same as slog
	ERROR          = slog.LevelError // same as slog
	PANIC LogLevel = 9               // slog.Level(9)
	FATAL LogLevel = 10              // slog.Level(10)
)

var (
	_fg_debug = palette.RGB_BLUE   // base blue
	_fg_info  = palette.RGB_GREEN  // base green
	_fg_warn  = palette.RGB_YELLOW // base yellow
	_fg_error = palette.RGB_RED    // base red
	_fg_panic = palette.RGB_PURPLE // base purple
	_fg_fatal = palette.RGB_GREY   // base grey

	_default_level  = INFO
	_default_logger *Logger
)

type (
	LogLevel = slog.Level
	// TempFunc is a function type to define a template function
	TempFunc func() string
	// PlainTextHandler is colorable plain text handler.
	// no attrs store here, it's useless.
	PlainTextHandler struct {
		colorable bool
		color     *color.Color           // color instance
		colordict []palette.RGB          // color dict
		temparser *fasttemplate.Template // parser to parse the log template
		tempdict  map[string]TempFunc    // get template string from here
		level     *slog.LevelVar         // switch log level dynamically
		out       io.Writer
		bufs      *sync.Pool // buffer pool
	}
	// Logger provide public api to record message.
	Logger struct {
		h        slog.Handler
		skipstep int
	}
)

// DefaultLevel update global log level dynamically.
func DefaultLevel(lv LogLevel) {
	_default_level = lv
	_default_logger = NewLogger(NewPlainTextHandler(os.Stderr, _default_level, _default_template))
}

func init() {
	_default_logger = NewLogger(NewPlainTextHandler(os.Stderr, _default_level, _default_template))
}

func _format_timestamp() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

func _format_prefix() string {
	return "-"
}

// SetDefault update global logger with l.
func SetDefault(l *Logger) {
	_default_logger = l
}

// NewPlainTextHandler return a new hanlder of PlainTextHandler. Is's an internal and default handler.
// minLevel controls the minmum output log level.
func NewPlainTextHandler(out io.Writer, minLevel LogLevel, template string) *PlainTextHandler {
	td := make(map[string]TempFunc) // custom string from outer
	td[TEMP_PREFIX] = _format_prefix
	td[TEMP_TIMESTAMP] = _format_timestamp
	if minLevel < DEBUG || minLevel > ERROR {
		minLevel = INFO
	}
	var lva slog.LevelVar
	// color list
	colorList := make([]palette.RGB, _max_level)
	colorList[int(DEBUG)+_offset_level] = _fg_debug
	colorList[int(INFO)+_offset_level] = _fg_info
	colorList[int(WARN)+_offset_level] = _fg_warn
	colorList[int(ERROR)+_offset_level] = _fg_error
	colorList[int(PANIC)+_offset_level] = _fg_panic
	colorList[int(FATAL)+_offset_level] = _fg_fatal
	lva.Set(minLevel)
	if template == "" {
		template = _default_template
	}
	tparser, err := fasttemplate.NewTemplate(template, "{", "}")
	if err != nil {
		panic(err.Error())
	}
	return &PlainTextHandler{
		colorable: true,
		colordict: colorList,
		tempdict:  td,
		temparser: tparser, level: &lva,
		out: out,
		bufs: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

// With record the template building process and wait for log parsing to be called.
// It's not concurrency safe.
func (h *PlainTextHandler) With(key string, value TempFunc) {
	h.tempdict[key] = value
}

// SetLogLevel set log level for specific handler.
func (h *PlainTextHandler) SetLogLevel(newLevel LogLevel) {
	h.level.Set(newLevel)
}

// from slog, same as slog.
func (h *PlainTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *PlainTextHandler) replaceTemp(tag string) []byte {
	if v, found := h.tempdict[tag]; found {
		return []byte(v())
	}
	// not found to ???
	return []byte(_format_default)
}

func (h *PlainTextHandler) replaceLevel(level *slog.LevelVar) []byte {
	levelStr := ""
	lv := level.Level()
	switch lv {
	case PANIC:
		levelStr = _format_panic
	case FATAL:
		levelStr = _format_fatal
	default:
		levelStr = level.Level().String()
	}
	if h.colorable {
		rgb := h.colordict[int(lv)+_offset_level]
		h.color = color.RGB(int(rgb.R), int(rgb.G), int(rgb.B))
		sf := h.color.SprintfFunc()
		return []byte(sf(levelStr))
	}
	return []byte(levelStr)
}

// from slog, same as slog.
func (h *PlainTextHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.Level < h.level.Level() {
		return nil
	}
	buf := h.bufs.Get().(*bytes.Buffer)
	buf.Reset()
	defer func() {
		if buf != nil {
			buf.Reset()
			h.bufs.Put(buf)
		}
	}()
	// 尝试解析 template, 如果没找到指定模板替换字符串, 采用???替换
	fm := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fm.Next()
	longPath, lineNum := f.File, f.Line
	shortPath := filepath.Base(longPath)
	s := h.temparser.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		switch tag {
		case TEMP_LONGPATH:
			return w.Write([]byte(longPath))
		case TEMP_SHORTPATH:
			return w.Write([]byte(shortPath))
		case TEMP_LINENUM:
			return fmt.Fprintf(w, "%d", lineNum)
		case TEMP_LEVEL:
			lv := &slog.LevelVar{}
			lv.Set(r.Level)
			return w.Write(h.replaceLevel(lv))
		default:
			return w.Write(h.replaceTemp(tag))
		}
	})
	buf.WriteString(s)
	if s[len(s)-1] != ' ' {
		buf.WriteByte(' ')
	}
	buf.WriteString(r.Message)
	if r.Message[len(r.Message)-1] != '\n' {
		buf.WriteByte('\n')
	}
	h.out.Write(buf.Bytes())
	return nil
}

// from slog, same as slog.
func (h *PlainTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// from slog, same as slog.
func (h *PlainTextHandler) WithGroup(name string) slog.Handler {
	return h
}

// NewLogger return pointer to new logger.
// It should take a valid slog.Handler.
func NewLogger(h slog.Handler) *Logger {
	if h == nil {
		panic("nil handler!")
	}
	return &Logger{
		h:        h,
		skipstep: _default_skip_step,
	}
}

// from log.slog impl
func (l *Logger) Handler() slog.Handler {
	return l.h
}

// from log.slog impl
func (l *Logger) Enabled(ctx context.Context, level LogLevel) bool {
	if ctx == nil {
		ctx = context.Background()
	}
	return l.Handler().Enabled(ctx, level)
}

// from log.slog impl
func (l *Logger) log(ctx context.Context, level LogLevel, msg string, args ...any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	// dyn skip step
	runtime.Callers(l.skipstep, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...) // for structed log
	_ = l.Handler().Handle(ctx, r)
}

func (l *Logger) Log(ctx context.Context, level LogLevel, msg string, args ...any) {
	l.log(ctx, level, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), DEBUG, msg, args...)
}

func (l *Logger) DebugX(ctx context.Context, msg string, args ...any) {
	l.log(context.Background(), DEBUG, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(context.Background(), INFO, msg, args...)
}

func (l *Logger) InfoX(ctx context.Context, msg string, args ...any) {
	l.log(context.Background(), INFO, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(context.Background(), WARN, msg, args...)
}

func (l *Logger) WarnX(ctx context.Context, msg string, args ...any) {
	l.log(context.Background(), WARN, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(context.Background(), ERROR, msg, args...)
}

func (l *Logger) ErrorX(ctx context.Context, msg string, args ...any) {
	l.log(context.Background(), ERROR, msg, args...)
}

func (l *Logger) Panic(msg string) {
	l.log(context.Background(), PANIC, msg)
	panic(msg)
}

func (l *Logger) PanicX(ctx context.Context, msg string) {
	l.log(ctx, PANIC, msg)
	panic(msg)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.log(context.Background(), FATAL, msg, args...)
	os.Exit(1)
}

func (l *Logger) FatalX(ctx context.Context, msg string, args ...any) {
	l.log(ctx, FATAL, msg, args...)
	os.Exit(1)
}

func Log(ctx context.Context, level LogLevel, msg string, args ...any) {
	_default_logger.Log(ctx, level, msg, args...)
}

func Debug(msg string, args ...any) {
	_default_logger.Debug(msg, args...)
}

func DebugX(ctx context.Context, msg string, args ...any) {
	_default_logger.DebugX(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	_default_logger.Info(msg, args...)
}

func InfoX(ctx context.Context, msg string, args ...any) {
	_default_logger.InfoX(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	_default_logger.Warn(msg, args...)
}

func WarnX(ctx context.Context, msg string, args ...any) {
	_default_logger.WarnX(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	_default_logger.Error(msg, args...)
}

func ErrorX(ctx context.Context, msg string, args ...any) {
	_default_logger.ErrorX(ctx, msg, args...)
}

func Panic(msg string) {
	_default_logger.Panic(msg)
}

func PanicX(ctx context.Context, msg string) {
	_default_logger.Panic(msg)
}

func Fatal(msg string, args ...any) {
	_default_logger.Fatal(msg, args...)
}

func FatalX(ctx context.Context, msg string, args ...any) {
	_default_logger.FatalX(ctx, msg, args...)
}
