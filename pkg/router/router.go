package router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/gcontext"
	"github.com/wendisx/puzzle/pkg/log"
)

const (
	_ctx_echo    = "_ctx_echo"
	_defaultRoot = ""
)

type (
	RouteRule[G any] interface {
		Block() bool
		Pattern() string
		IsEnd() bool
		RegisterRule(prefix string, g G)
		RegisterEndpoint(prefix string, g G)
	}
	router struct {
		routeCnt    int
		activeCnt   int
		rootPattern string // _default = "" replace root
	}
	Endpoint[HF any] struct {
		Method string
		Path   string
		Hf     HF
	}
	echoRouter struct {
		router
		g *echo.Group
	}
)

func NewEchoRouter() *echoRouter {
	r := &echoRouter{
		router: router{
			rootPattern: _defaultRoot,
			routeCnt:    0,
			activeCnt:   0,
		},
	}
	return r
}

func (r *echoRouter) RootPattern(rp string) *echoRouter {
	e, ok := gcontext.GetGlobalContext().Get(_ctx_echo).(*echo.Echo)
	if !ok {
		panic(fmt.Sprintf("invalid gcontext key with [%s]", _ctx_echo))
	}
	r.g = e.Group(rp)
	return r
}

func (r *echoRouter) Apply(rr RouteRule[*echo.Group]) *echoRouter {
	if r.g == nil {
		panic("router can't find root pattern")
	}
	EchoApply(r.rootPattern, r.g, rr)
	r.routeCnt += 1
	if !rr.Block() {
		r.activeCnt += 1
	}
	return r
}

func EchoApply(prefix string, g *echo.Group, rr RouteRule[*echo.Group]) {
	apply(prefix, g, rr)
}

func EchoEndpoint(prefix string, g *echo.Group, ep *Endpoint[echo.HandlerFunc], mf ...echo.MiddlewareFunc) {
	log.PlainLog.Info(fmt.Sprintf("endpoint at %s in %s", prefix+ep.Path, ep.Method))
	g.Add(ep.Method, ep.Path, ep.Hf, mf...)
}

func apply[G any](prefix string, g G, rr RouteRule[G]) {
	if rr.Block() {
		return
	}
	rr.RegisterEndpoint(prefix, g)
	if !rr.IsEnd() {
		rr.RegisterRule(prefix, g)
	}
}
