package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/server"
)

func (m EchoMiddleware) BindResponder(key string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			responder := server.GetEchoResponder(c)
			c.Set(key, responder)
			defer server.PutEchoResponder(responder)
			return next(c)
		}
	}
}
