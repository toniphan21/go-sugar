//go:build !dev

package sugar

import "log/slog"

func SetLogger(l *slog.Logger) {
}

func debug(msg string, args ...any) {
}
