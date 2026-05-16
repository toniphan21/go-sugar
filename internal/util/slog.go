package util

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if !h.Enabled(ctx, r.Level) {
			continue
		}

		if err := h.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: hs}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: hs}
}

func NewLogger(filePath string, writer io.Writer, level slog.Level) (*slog.Logger, io.Closer, error) {
	handlers := []slog.Handler{
		slog.NewTextHandler(writer, &slog.HandlerOptions{Level: level}),
	}

	var f *os.File
	var err error
	if filePath != "" {
		f, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, nil, err
		}
		handlers = append(handlers, slog.NewJSONHandler(f, &slog.HandlerOptions{Level: level}))
	}

	return slog.New(&multiHandler{handlers: handlers}), f, err
}
