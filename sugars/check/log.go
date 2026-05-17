//go:build !dev

package check

import "log/slog"

func SetLogger(l *slog.Logger) {
}

func debug(msg string, args ...any) {
}
