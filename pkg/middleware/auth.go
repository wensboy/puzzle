package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/clog"
	"github.com/wendisx/puzzle/pkg/server"
	"github.com/wendisx/puzzle/pkg/util"
)

/*  auth middleware for echo */
func (m EchoMiddleware) SimpleJwtAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := c.Request().Header.Get("Authorization")
			clog.Debug(tokenStr)
			if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
				return server.WithEchoRes().Err(http.StatusUnauthorized, "unauthorized")
			}
			tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
			jwtClaim, err := util.ParseToken(tokenStr)
			if err != nil || !jwtClaim.Check() {
				return server.WithEchoRes().Err(http.StatusUnauthorized, "invalid token")
			}
			// TODO: to be optimized...
			clog.Info(fmt.Sprintf("context.user{id=%s,name=%s}", string(jwtClaim.ExternId), jwtClaim.Name))
			c.Set("userId", jwtClaim.ExternId)
			c.Set("name", jwtClaim.Name)
			return next(c)
		}
	}
}

/* check middleware for echo */
func (m EchoMiddleware) ParseAndCheckBody(s interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := util.Bind(c, s); err != nil {
				return server.WithEchoRes().Err(http.StatusBadRequest, "Invalid parameter passed")
			}
			err := util.GetGlobalValidator().Check(s)
			if err != nil {
				return server.WithEchoRes().Err(http.StatusBadRequest, "Parameter verification failed")
			}
			clog.Debug(fmt.Sprintf("%#v", s))
			c.Set("body", s)
			return next(c)
		}
	}
}
