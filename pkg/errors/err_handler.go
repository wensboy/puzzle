package errors

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/wendisx/puzzle/pkg/log"
)

// from Echo
func EchoErrHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}
	var e error
	httpCode := http.StatusInternalServerError
	var message any
	message = map[string]any{
		"code":    1,
		"message": "unknown error",
	}
	// BErr应当调用fail()正常执行后返回
	if he, ok := err.(*echo.HTTPError); ok {
		httpCode = he.Code
		message = he.Message
	} else if he, ok := err.(IErr); ok {
		httpCode = http.StatusOK
		message = map[string]any{
			"code":    he.Code,
			"message": he.Msg,
		}
	}
	e = c.JSON(httpCode, message)
	if e != nil {
		log.PlainLog.Error(fmt.Sprintf("error handler throw %s", e.Error()))
	}
}
