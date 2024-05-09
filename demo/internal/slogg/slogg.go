package slogg

import (
	"context"
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}

// ErrorWithOSExit replicates the behaviour of `log.Fatal`,
// and logs the underlying error via slogg.ErrorAttr
func ErrorWithOSExit(msg string, err error) {
	slog.Error(msg, ErrorAttr(err))
	os.Exit(1)
}

// ErrorContextWithOSExit replicates the behaviour of `log.Fatal`,
// and logs the underlying error via slogg.ErrorAttr, whilst also accepting
// a `context.Context`
func ErrorContextWithOSExit(ctx context.Context, msg string, err error) {
	slog.ErrorContext(ctx, msg, ErrorAttr(err))
	os.Exit(1)
}

// ErrorAttr returns an Attr for an error, with a key of `err`
func ErrorAttr(err error) slog.Attr {
	return slog.String("err", err.Error())
}
