package router

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/palette"
)

/*
	echo route impl
*/

const (
	_default_echo_gateway = ""
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
	EchoGateway struct {
		EchoRoute
	}
)

func NewEchoPack(p Pack, g *echo.Group) EchoPack {
	return EchoPack{
		P: p,
		G: g,
	}
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

// return the single instance of EchoPack
func DefaultEchoPack(e *echo.Echo) EchoPack {
	// 确保所有 group from default root
	if _default_echopack == nil {
		_default_echopack = &EchoPack{
			P: Pack{
				Prefix: _default_echo_gateway,
			},
			G: e.Group(_default_echo_gateway),
		}
	}
	return *_default_echopack
}

// return a new EchoRoot to start all
func NewEchoGateway(e *echo.Echo) Route[EchoPack] {
	er := &EchoGateway{
		*NewEchoRoute(_default_echo_gateway),
	}
	er.Inbound(DefaultEchoPack(e))
	return er
}

// can override to block route and peer
func (r *EchoRoute) Active() bool {
	// clog.Warn("<pkg.router.echo> Call default Active function, this method should be explicitly overridden.")
	return true
}

func (r *EchoRoute) Path() string {
	return r.path
}

// should override and call it in Inbound()
func (r *EchoRoute) Handle(p EchoPack) {
	clog.Warn("<pkg.router.echo> Call default Handle function, this method should be explicitly overridden.")
}

// transform EchoPack and update path prefix
func (r *EchoRoute) Inbound(p EchoPack) {
	if !r.Active() {
		return
	}
	// new outbound EchoPack to next route or just peer
	prefix := p.P.Prefix + r.path
	clog.Info("<pkg.router.echo> router cover " + palette.SkyBlue(prefix))
	r.ep = EchoPack{
		P: Pack{
			Prefix: prefix,
		},
		G: p.G.Group(r.path),
	}
}

func (r *EchoRoute) ToRoute(rr Route[EchoPack]) {
	r.qroute = append(r.qroute, rr)
}

func (r *EchoRoute) ToPeer(p Peer[EchoPack]) {
	r.qpeer = append(r.qpeer, p)
}

// go out route and routing EchoPack
func (r *EchoRoute) Outbound() {
	if !r.Active() {
		return
	}
	if r.ep.G == nil {
		clog.Panic("<pkg.router.echo> Outbound routing link terminate.")
	}
	for i := range r.qroute {
		// if override Active() to false, then just skip all.
		if !strings.Contains(r.ep.P.Prefix, r.qroute[i].Path()) && r.qroute[i].Active() {
			r.qroute[i].Inbound(r.ep)
			r.qroute[i].Outbound()
		}
	}
	for i := range r.qpeer {
		r.qpeer[i].Parse(r.ep)
	}
}

func (p *EchoPeer) ToEndpoint(ep Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]) {
	if p.qendpoint == nil {
		p.qendpoint = make([]Endpoint[echo.HandlerFunc, echo.MiddlewareFunc], 0)
	}
	p.qendpoint = append(p.qendpoint, ep)
}

// mount all endpoints
func (p EchoPeer) Parse(pp EchoPack) {
	if len(p.qendpoint) == 0 {
		clog.Warn("<pkg.router.echo> Call default Mount function, this method should be explicitly overridden and is eventually called.")
	}
	for _, ep := range p.qendpoint {
		pp.G.Add(ep.Method, ep.Path, ep.Handler, ep.PreHandlers...)
		clog.Info(fmt.Sprintf("<pkg.router.echo> %s#%s", ep.Method, pp.P.Prefix+ep.Path))
	}
}
