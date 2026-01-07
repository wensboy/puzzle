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

	"github.com/valyala/fasttemplate"
)

/*
	discard log inplement but useful
*/

const (
	_ts_rfc3339     = "ts_rfc3339"
	_ts_rfc3339nano = "ts_rfc3339nano"
	_prefix         = "prefix"
	_level          = "level"
	_from           = "from"
	_message        = "message"

	_defaultPrefix   = "-"
	_defaultTemplate = `${ts_rfc3339} [${level}] ${from} ${prefix} ${message}`
	_defaultBufSize  = 1 << 10
	_defaultLevel    = slog.LevelDebug
)

var (
	_defaultOutput = os.Stderr

	PlainLog = slog.New(NewPlainHandler(&LogOption{
		Prefix:   _defaultPrefix,
		Template: _defaultTemplate,
		Level:    slog.LevelDebug,
		Out:      os.Stderr,
	}))
)

type (
	groupOrAttrs struct {
		group string
		attrs []slog.Attr
	}
	LogOption struct {
		Prefix   string
		Template string
		Level    slog.Leveler // log level interface
		Out      io.Writer    // output
	}
	PlainHandler struct {
		opt     *LogOption
		bufPool sync.Pool
		temp    *fasttemplate.Template
		mx      *sync.Mutex
	}
)

// plain hanlder
// non structed
func NewPlainHandler(opt *LogOption) slog.Handler {
	h := &PlainHandler{
		opt: opt,
		bufPool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 0, _defaultBufSize))
			},
		},
		temp: fasttemplate.New(_defaultTemplate, "${", "}"),
		mx:   &sync.Mutex{},
	}
	h.setup()
	return h
}

func (h *PlainHandler) setup() {
	if t, err := fasttemplate.NewTemplate(h.opt.Template, "${", "}"); h.opt.Template != "" && err == nil {
		h.temp = t
	}
	if h.opt.Prefix == "" {
		h.opt.Prefix = _defaultPrefix
	}
	if h.opt.Level == nil {
		h.opt.Level = _defaultLevel
	}
	if h.opt.Out == nil {
		h.opt.Out = _defaultOutput
	}
}

func (h *PlainHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opt.Level.Level()
}

func (h *PlainHandler) Handle(_ context.Context, r slog.Record) error {
	buf := h.bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer h.bufPool.Put(buf)
	message := h.temp.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		switch tag {
		case _ts_rfc3339:
			return w.Write([]byte(r.Time.Format(time.RFC3339)))
		case _ts_rfc3339nano:
			return w.Write([]byte(r.Time.Format(time.RFC3339Nano)))
		case _level:
			return w.Write([]byte(r.Level.Level().String()))
		case _from:
			f := runtime.CallersFrames([]uintptr{r.PC})
			fm, _ := f.Next()
			return fmt.Fprintf(w, "%s:%d", filepath.Base(fm.File), fm.Line)
		case _prefix:
			return w.Write([]byte(h.opt.Prefix))
		case _message:
			return w.Write([]byte(r.Message))
		default:
			return fmt.Fprintf(w, "?%s?", tag)
		}
	})
	buf.WriteString(message)
	buf.WriteByte('\n')
	h.mx.Lock()
	defer h.mx.Unlock()
	h.opt.Out.Write(buf.Bytes())
	return nil
}

func (h *PlainHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PlainHandler) WithGroup(name string) slog.Handler {
	return h
}
