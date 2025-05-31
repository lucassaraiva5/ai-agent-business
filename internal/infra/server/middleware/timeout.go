package middleware

import (
	"time"

	"ai-agent-business/internal/infra/logger"
	"ai-agent-business/internal/infra/logger/attributes"
	"ai-agent-business/internal/infra/variables"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	DefaultTimeoutConfig = middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "{\"error\":\"Request Timeout\"}",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Warn(c.Request().Context(), "Request Timeout", attributes.Attributes{
				"uri": c.Request().RequestURI,
			}.WithError(err))
		},
		Timeout: time.Second * time.Duration(variables.ServerTimeout()),
	}
)

// ConfigTimeout middleware adds a `timeout`
func ConfigTimeout() echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(DefaultTimeoutConfig)
}
