package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	mware "github.com/ippoippo/go-servemux-lt/demo/internal/middleware/echo"
	routes "github.com/ippoippo/go-servemux-lt/demo/internal/routehandlers/echo"
	"github.com/ippoippo/go-servemux-lt/demo/internal/slogg"
)

func main() {
	logger := slogg.NewLogger()
	slog.SetDefault(logger)

	// Create Server
	e := newEcho() // echo.New(), but with banner and port stdout supressed

	// ********************
	// DEMO FEATURE POINT 4
	// ********************
	// Setup Middleware
	setupMiddleware(e)

	// Setup Routes
	setupRoutes(e)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			slogg.ErrorWithOSExit("failed to e.Start", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slogg.ErrorContextWithOSExit(ctx, "failed to e.Shutdown", err)
	}
	slog.Info("graceful shutdown complete")
}

func setupMiddleware(e *echo.Echo) {

	slog.Info("configuring middlware")

	// Middleware Configuration
	e.Use(middleware.Recover())       // Recover from panics (supplied by Echo)
	e.Use(mware.AddXRequestIdToCtx()) // Add the XRequestId header to the ctx
	e.Use(mware.RequestLogging())     // Log incoming requests
}

func setupRoutes(e *echo.Echo) {

	slog.Info("configuring routes")

	// Routes Configuration
	// **********************
	// DEMO FEATURE POINT 1&3
	// **********************
	// Setup API version
	v1Api := e.Group("/v1")
	// Setup a group for an entity such as notes
	notesGroup := v1Api.Group("/notes")

	// Lets assume these endpoints are not authenticated
	notesGroup.GET("", routes.GetAllNotes)
	// ********************
	// DEMO FEATURE POINT 2
	// ********************
	notesGroup.GET("/:id", routes.GetNote)

	// Lets assume these endpoints need to be authenticated, so custom middleware
	notesGroup.POST("", routes.CreateNote, mware.Authenticated())
	notesGroup.DELETE("/:id", routes.DeleteNote, mware.Authenticated())
}

func newEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return e
}
