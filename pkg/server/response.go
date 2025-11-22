package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	EmptyData       struct{}
	Response[C any] interface {
		Success(c C, code int, msg string, data any) error
		Fail(c C, code int, msg string) error
		Err(httpCode int, msg string) error
	}
	echoResponse[C echo.Context] struct{}
)

func (r echoResponse[C]) Success(c C, code int, msg string, data any) error {
	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    code,
			"message": msg,
			"data":    data,
		},
	)
}

func (r echoResponse[C]) Fail(c C, code int, msg string) error {
	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    code,
			"message": msg,
			"data":    nil,
		},
	)
}

func (r echoResponse[C]) Err(httpCode int, msg string) error {
	return echo.NewHTTPError(httpCode, map[string]any{
		"code":    1,
		"message": msg,
	})
}

func WithEchoRes() Response[echo.Context] {
	return echoResponse[echo.Context]{}
}
