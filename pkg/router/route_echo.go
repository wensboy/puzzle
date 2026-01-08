package router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/clog"
)

/*
	echo route impl
	1. root route from group("")
	2.
*/

const (
	_default_echoroot = ""
)

var (
	_default_echopack *EchoPack
)

type (
	EchoPack struct {
		P Pack
		G *echo.Group
	}
	EchoRoute struct {
		ep     EchoPack
		path   string            // for self
		qroute []Route[EchoPack] // next route queue
		qpeer  []Peer[EchoPack]  // next peer queue
	}
	EchoPeer struct {
		qendpoint []Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]
	}
	EchoRoot struct {
		EchoRoute
	}
)

// return the single instance of EchoPack
func DefaultEchoPack(e *echo.Echo) EchoPack {
	// 确保所有 group from default root
	if _default_echopack == nil {
		_default_echopack = &EchoPack{
			P: Pack{
				Prefix: _default_echoroot,
			},
			G: e.Group(_default_echoroot),
		}
	}
	return *_default_echopack
}

// return a new EchoRoot to start all
func NewEchoRoot(e *echo.Echo) Route[EchoPack] {
	er := &EchoRoot{
		*NewEchoRoute(_default_echoroot),
	}
	er.Inbound(DefaultEchoPack(e))
	return er
}

// as an embedded to implement Route[Pack]
func NewEchoRoute(path string) *EchoRoute {
	if path == "" {
		clog.Warn("<pkg.router.echo> Empty path route exists.")
	}
	return &EchoRoute{
		path: path,
	}
}

// can override to block route and peer
func (r *EchoRoute) Active() bool {
	// clog.Warn("<pkg.router.echo> Call default Active function, this method should be explicitly overridden.")
	return true
}

// should override and call it in Inbound()
func (r *EchoRoute) Handle() {
	clog.Warn("<pkg.router.echo> Call default Handle function, this method should be explicitly overridden.")
}

// transform EchoPack and update path prefix
func (r *EchoRoute) Inbound(p EchoPack) {
	if !r.Active() {
		return
	}
	// new outbound EchoPack to next route or just peer
	prefix := p.P.Prefix + r.path
	clog.Debug("<pkg.router.echo> New outbound prefix: " + prefix)
	r.ep = EchoPack{
		P: Pack{
			Prefix: prefix,
		},
		G: p.G.Group(r.path),
	}
	// r.Handle() <- just like this!
}

func (r *EchoRoute) ToRoute(rr Route[EchoPack]) {
	r.qroute = append(r.qroute, rr)
}

func (r *EchoRoute) ToPeer(p Peer[EchoPack]) {
	r.qpeer = append(r.qpeer, p)
}

// go out route and routing EchoPack
func (r *EchoRoute) Outbound() {
	if r.ep.G == nil {
		clog.Panic("<pkg.router.echo> Outbound routing link terminate.")
	}
	for i := range r.qroute {
		// if override Active() to false, then just skip all.
		if r.qroute[i].Active() {
			r.qroute[i].Inbound(r.ep)
			r.qroute[i].Outbound()
		}
	}
	for i := range r.qpeer {
		r.qpeer[i].Mount(r.ep)
	}
}

func (p *EchoPeer) ToEndpoint(ep Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]) {
	if p.qendpoint == nil {
		p.qendpoint = make([]Endpoint[echo.HandlerFunc, echo.MiddlewareFunc], 0)
	}
	p.qendpoint = append(p.qendpoint, ep)
}

// mount all endpoints
func (p EchoPeer) Mount(pp EchoPack) {
	if len(p.qendpoint) == 0 {
		clog.Warn("<pkg.router.echo> Call default Mount function, this method should be explicitly overridden and is eventually called.")
	}
	for _, ep := range p.qendpoint {
		pp.G.Add(ep.Method, ep.Path, ep.Handler, ep.PreHandlers...)
		clog.Info(fmt.Sprintf("<pkg.router.echo> %s#%s", ep.Method, pp.P.Prefix+ep.Path))
	}
}
