package watch

import (
	"log/slog"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/sugar"
)

type Logger interface {
	Error(msg string, target sugar.Target, err error)
	Skip(target sugar.Target)
	Watching(target sugar.Target)
	Stop(target sugar.Target)
}

type defaultLogger struct {
	log *slog.Logger
}

func (l *defaultLogger) Error(msg string, target sugar.Target, err error) {
	l.log.Error(cli.ColorRed(msg), slog.String("target", target.Path), slog.Any("error", err))
}

func (l *defaultLogger) Skip(target sugar.Target) {
	l.log.Info(cli.ColorYellow("skip"), slog.String("target", target.Path))
}

func (l *defaultLogger) Watching(target sugar.Target) {
	l.startStop(target, "watching")
}

func (l *defaultLogger) Stop(target sugar.Target) {
	l.startStop(target, "stop watcher")
}

func (l *defaultLogger) startStop(target sugar.Target, msg string) {
	switch {
	case target.IsDir && target.Recursive:
		l.log.Info(cli.ColorGreen(msg) + cli.ColorYellow(" dir recursively ") + cli.ColorCyan(target.DisplayPath()))
	case target.IsDir:
		l.log.Info(cli.ColorGreen(msg) + cli.ColorYellow(" dir (non-recursive) ") + cli.ColorCyan(target.DisplayPath()))
	default:
		l.log.Info(cli.ColorGreen(msg) + cli.ColorYellow(" file ") + cli.ColorCyan(target.DisplayPath()))
	}
}

var _ Logger = (*defaultLogger)(nil)
