package router

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/clog"
)

/*
router test with echo framework.
eg: usi=/agent/[api|openapi]/[v..]/<ns>/...
*/
type (
	VersionRoute struct {
		EchoRoute
	}
	ApiRoute struct {
		EchoRoute
	}
	UserPeer struct {
		EchoPeer
	}
)

func (up *UserPeer) Mount(p EchoPack) {
	up.ToEndpoint(Endpoint[echo.HandlerFunc, echo.MiddlewareFunc]{
		Method: http.MethodGet,
		Path:   "/info",
		Handler: func(c echo.Context) error {
			return c.String(http.StatusOK, "get info successfully.")
		},
		PreHandlers: []echo.MiddlewareFunc{},
	})
	up.EchoPeer.Mount(p)
}

// test not avtive route [passed]
func (vr *VersionRoute) Active() bool {
	return true
}

var (
	e = echo.New()
)

// test single instance of echopack [passed]
func Test_single_instance(t *testing.T) {
	a := DefaultEchoPack(e)
	b := DefaultEchoPack(e)
	clog.Debug(fmt.Sprintf("the same instance of _default_echopack: %v", a == b))
}

// test route [passed]
func Test_route(t *testing.T) {
	rootRoute := NewEchoRoot(e)
	// _ = NewEchoRoot(e)
	apiRoute := &ApiRoute{
		EchoRoute: *NewEchoRoute("/api"),
	}
	vRoute := &VersionRoute{
		EchoRoute: *NewEchoRoute("/v1"),
	}
	up := &UserPeer{}
	rootRoute.ToRoute(apiRoute)
	apiRoute.ToRoute(vRoute)
	vRoute.ToPeer(up)
	rootRoute.Outbound()
	e.Server.Addr = "127.0.0.1:3333"
	if err := e.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		clog.Panic(err.Error())
	}
}
