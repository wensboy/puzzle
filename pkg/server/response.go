package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	// 2xx Success
	MsgOK        = "request successful"
	MsgCreated   = "resource created successfully"
	MsgNoContent = "no content to return"

	// 3xx Redirection
	MsgMovedPermanently = "resource has been moved permanently"
	MsgFound            = "resource found at a different URI"

	// 4xx Client Errors
	MsgBadRequest       = "bad request, please check your input"
	MsgUnauthorized     = "unauthorized, please provide valid credentials"
	MsgForbidden        = "forbidden, you do not have permission"
	MsgNotFound         = "resource not found"
	MsgMethodNotAllowed = "http method not allowed"
	MsgConflict         = "resource conflict occurred"
	MsgUnprocessable    = "unprocessable entity, validation failed"

	// 5xx Server Errors
	MsgInternalError      = "internal server error, please try again later"
	MsgBadGateway         = "bad gateway, upstream server error"
	MsgServiceUnavailable = "service unavailable, server is overloaded or down"
)

type (
	EmptyData struct{}
	EmptyList []struct{}
	// Response is abstract response.
	Response[C any] interface {
		Success(c C, code int, msg string, data any) error
		Fail(c C, code int, msg string) error
		Err(httpCode int, msg string) error
	}
	echoResponse[C echo.Context] struct{}
)

// Success response with specific logic code and message and data.
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

// Fail response with specific logic code and operation fail message.
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

// Err response with specific http code and error message.
func (r echoResponse[C]) Err(httpCode int, msg string) error {
	return echo.NewHTTPError(httpCode, map[string]any{
		"code":    -1,
		"message": msg,
	})
}

// WithEchoRes limit use of Echo context.
func WithEchoRes() Response[echo.Context] {
	return echoResponse[echo.Context]{}
}
