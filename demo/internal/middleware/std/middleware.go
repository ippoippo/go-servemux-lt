package std

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

const (
	HeaderXRequestID = "X-Request-Id"
)

type contextKey string

func JsonContentType(next http.Handler) http.Handler {
	slog.Debug("enter JsonContentType")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("execute JsonContentType")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func AddXRequestIdToCtx(next http.Handler) http.Handler {
	slog.Debug("enter AddXRequestIdToCtx")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("execute AddXRequestIdToCtx")
		rid := r.Header.Get(HeaderXRequestID)
		request := r
		if rid != "" {
			ctx := context.WithValue(r.Context(), contextKey(HeaderXRequestID), rid)
			request = r.WithContext(ctx)
		}

		next.ServeHTTP(w, request)
	})
}

func RequestLogging(next http.Handler) http.Handler {
	slog.Debug("enter RequestLogging")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("execute RequestLogging")
		slog.InfoContext(
			r.Context(),
			"REQUEST",
			slog.String("uri", r.URL.RequestURI()),
			slog.String("method", r.Method),
			slog.String(HeaderXRequestID, r.Header.Get(HeaderXRequestID)),
		)

		next.ServeHTTP(w, r)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	slog.Debug("enter RecoverPanic")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("execute RecoverPanic")
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				http.Error(w, errors.New("something unexpected").Error(), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Authenticated(next http.Handler) http.Handler {
	slog.Debug("enter Authenticated")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("execute Authenticated")
		// Normally we would appropriately authenticate
		slog.InfoContext(
			r.Context(),
			"Authenticated",
			slog.Bool("isAuthenticated", true),
		)
		next.ServeHTTP(w, r)
	})
}
