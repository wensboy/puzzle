package router

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/gcontext"
	"github.com/wendisx/puzzle/pkg/log"
)

type (
	tr1 struct{}
	tr2 struct{}
)

func newTr1() *tr1             { return &tr1{} }
func (r *tr1) Pattern() string { return "/tr1" }
func (r *tr1) Block() bool     { return false }
func (r *tr1) IsEnd() bool     { return false }
func (r *tr1) RegisterRule(prefix string, g *echo.Group) {
	g.Group(r.Pattern())
	prefix += r.Pattern()
	EchoApply(prefix, g, newTr2())
}
func (r *tr1) RegisterEndpoint(prefix string, g *echo.Group) {
	g.Group(r.Pattern())
	prefix += r.Pattern()
	log.PlainLog.Debug(fmt.Sprintf("path: %s/t1", prefix))
	g.Add(http.MethodGet, "/t1", func(c echo.Context) error {
		return c.String(http.StatusOK, "t1")
	})
}

func newTr2() *tr2             { return &tr2{} }
func (r *tr2) Pattern() string { return "/tr2" }
func (r *tr2) Block() bool     { return false }
func (r *tr2) IsEnd() bool     { return true }
func (r *tr2) RegisterRule(prefix string, g *echo.Group) {
	g.Group(r.Pattern())
	prefix += r.Pattern()
	log.PlainLog.Debug(fmt.Sprintf("path: %s", prefix))
}
func (r *tr2) RegisterEndpoint(prefix string, g *echo.Group) {
	g.Group(r.Pattern())
	prefix += r.Pattern()
	log.PlainLog.Debug(fmt.Sprintf("path: %s/t2", prefix))
	g.Add(http.MethodGet, "/t2", func(c echo.Context) error {
		return c.String(http.StatusOK, "t2")
	})
}

func TestEchoRouter(t *testing.T) {
	e := echo.New()
	gcontext.NewGlobalContext().Set(_ctx_echo, e)
	r := NewEchoRouter().RootPattern("")
	t.Run("rule_endpoint", func(t *testing.T) {
		r.Apply(newTr1())
	})
}
