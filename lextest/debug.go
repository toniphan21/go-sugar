package lextest

import (
	"log/slog"
	"os"
)

type SetLogger func(logger *slog.Logger)

func Debug(setters ...SetLogger) {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)
	for _, set := range setters {
		set(logger)
	}
}
