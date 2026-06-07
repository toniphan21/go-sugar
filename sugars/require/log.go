//go:build !dev

package require

import "log/slog"

func SetLogger(l *slog.Logger) {
}

func debug(msg string, args ...any) {
}
