package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wendisx/puzzle/pkg/errors"
	"github.com/wendisx/puzzle/pkg/gcontext"
	"github.com/wendisx/puzzle/pkg/log"
)

const (
	_ctx_echo = "_ctx_echo"

	_defaultDelay = 1 * time.Millisecond
	_defaultAddr  = "0.0.0.0:2333"
)

type (
	WebOption struct {
		Addr string
	}
	webServer[H http.Handler] struct {
		opt  *WebOption
		h    H
		s    *http.Server
		quit chan os.Signal
		exit chan struct{}
	}
	EchoServer struct {
		webServer[*echo.Echo]
	}
)

func NewEchoServer(opt *WebOption) *EchoServer {
	h, ok := gcontext.GetGlobalContext().Get(_ctx_echo).(*echo.Echo)
	if !ok {
		panic(fmt.Sprintf("invalid gcontext key [%s]", _ctx_echo))
	}
	// setup handler here
	h.HTTPErrorHandler = errors.EchoErrHandler
	skipper := func(c echo.Context) bool {
		return false
	}
	h.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: skipper,
		Format:  `${time_rfc3339} [${method}] network - from ${remote_ip} hit ${uri} with (${status},${latency_human}) in ${bytes_in} out ${bytes_out}` + "\n",
	}))
	es := &EchoServer{
		webServer[*echo.Echo]{
			opt:  opt,
			h:    h,
			quit: make(chan os.Signal, 1),
			exit: make(chan struct{}),
			s:    &http.Server{},
		},
	}
	return es
}

func (es *EchoServer) Setup() *EchoServer {
	es.setupServer()
	return es
}

func (es *EchoServer) Start() {
	es.startServer()
}

func (es *EchoServer) Stop() {
	es.quit <- syscall.SIGQUIT
}

func (ws *webServer[H]) setupServer() {
	// setup server here
	ws.s.Handler = ws.h
	if ws.opt.Addr == "" {
		ws.opt.Addr = _defaultAddr
	}
	ws.s.Addr = ws.opt.Addr
}

func (ws *webServer[H]) startServer() {
	go ws.stopServer(_defaultDelay)
	log.PlainLog.Info(fmt.Sprintf("web server listen %s", ws.s.Addr))
	if err := ws.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.PlainLog.Error(fmt.Sprintf("web server start fail for %s", err.Error()))
		close(ws.exit)
	}
	<-ws.exit
	log.PlainLog.Info("web server exit")
}

func (ws *webServer[H]) stopServer(delay time.Duration) {
	signal.Notify(ws.quit, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT)
	<-ws.quit
	ctx, cancle := context.WithTimeout(context.Background(), delay)
	defer cancle()
	if err := ws.s.Shutdown(ctx); err != nil {
		log.PlainLog.Error(fmt.Sprintf("web server shutdown fail after %s", delay))
		os.Exit(1)
	}
	log.PlainLog.Info(fmt.Sprintf("web server shutdown success after %s", delay))
	close(ws.exit)
}
