package middleware

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

type contextKey string

// AddXRequestIdToCtx will (if client supplies X-Request-Id) insert into context
func AddXRequestIdToCtx() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rid := c.Request().Header.Get(echo.HeaderXRequestID)
			if rid != "" {
				ctx := context.WithValue(c.Request().Context(), contextKey(echo.HeaderXRequestID), rid)
				request := c.Request().Clone(ctx)
				c.SetRequest(request)
			}
			return next(c)
		}
	}
}

// RequestLogging logs appropriate info about the request
func RequestLogging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			slog.InfoContext(
				req.Context(),
				"REQUEST",
				slog.String("uri", req.RequestURI),
				slog.String("method", req.Method),
				slog.String(echo.HeaderXRequestID, req.Header.Get(echo.HeaderXRequestID)),
			)
			return next(c)
		}
	}
}

// Authenticated confirms user is authenticated
func Authenticated() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			// Normally we would appropriately authenticate

			slog.InfoContext(
				req.Context(),
				"Authenticated",
				slog.Bool("isAuthenticated", true),
			)
			return next(c)
		}
	}
}
