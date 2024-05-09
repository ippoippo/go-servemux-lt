package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	mware "github.com/ippoippo/go-servemux-lt/demo/internal/middleware/std"
	routes "github.com/ippoippo/go-servemux-lt/demo/internal/routehandlers/std"
	"github.com/ippoippo/go-servemux-lt/demo/internal/slogg"
)

func main() {
	logger := slogg.NewLogger()
	slog.SetDefault(logger)

	//
	srv := &http.Server{
		Addr:    ":1323",
		Handler: setupMux(),
		// Error handling, needs to use older *log.logger (bridge old APIs to new)
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := start(srv); err != nil && err != http.ErrServerClosed {
			slogg.ErrorWithOSExit("failed to e.Start", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := shutdown(ctx, srv); err != nil {
		slogg.ErrorContextWithOSExit(ctx, "failed to e.Shutdown", err)
	}
	slog.Info("graceful shutdown complete")
}

func start(srv *http.Server) error {
	return srv.ListenAndServe()
}

func shutdown(ctx context.Context, srv *http.Server) error {
	slog.InfoContext(ctx, "Closing server")
	return srv.Close()
}

func setupMux() *http.ServeMux {
	slog.Info("configuring mux")
	mux := http.NewServeMux()
	return withRoutes(mux)
}

func withMiddlware(f http.Handler) http.Handler {
	return mware.JsonContentType(
		mware.AddXRequestIdToCtx(
			mware.RequestLogging(
				mware.RecoverPanic(f))))
}

func withRoutes(mux *http.ServeMux) *http.ServeMux {
	slog.Info("configuring routes")

	// Setup a group for an entity
	notesGroup := http.NewServeMux()
	// Lets assume these endpoints are not authenticated
	notesGroup.HandleFunc("GET /notes", routes.GetAllNotes)
	notesGroup.HandleFunc("GET /notes/{id}", routes.GetNote)

	// Lets assume these endpoints need to be authenticated
	notesGroup.Handle("POST /notes", mware.Authenticated(http.HandlerFunc(routes.CreateNote)))
	notesGroup.Handle("DELETE /notes/{id}", mware.Authenticated(http.HandlerFunc(routes.DeleteNote)))

	// Demo malformed
	// notesGroup.HandleFunc("GER /notes", routes.GetAllNotes)
	// The above does not cause a compile time problem

	// Setup API version
	v1Api := http.NewServeMux()
	v1Api.Handle("/", notesGroup)

	mux.Handle("/v1/", withMiddlware(http.Handler(http.StripPrefix("/v1", v1Api))))
	return mux
}
