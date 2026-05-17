//go:build dev

package sugar

import "log/slog"

var logger = slog.Default()

func SetLogger(l *slog.Logger) {
	logger = l
}

func debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}
