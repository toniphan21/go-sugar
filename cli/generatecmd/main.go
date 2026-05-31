package generatecmd

import (
	"io"
	"log/slog"

	"nhatp.com/go/gen-lib/cli"
	"nhatp.com/go/sugar/internal/util"
)

const cmdName = "generate"

type Arguments struct {
	Args     []string
	Watch    bool
	DryRun   bool
	JSON     bool
	Log      string
	LogLevel slog.Level
}

func (a *Arguments) inputs() []string {
	if len(a.Args) == 0 {
		return []string{"./"}
	}
	return a.Args
}

func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer, args Arguments) error {
	log, logFile, err := util.NewCLILogger(args.Log, stderr, args.LogLevel)
	if err != nil {
		return err
	}
	log = log.With("cmd", cmdName).WithGroup(cmdName)
	defer func() {
		msg := "done"
		if args.Watch {
			msg = "end"
		}
		log.Info(cli.ColorGreen(msg))
		if logFile != nil {
			if ce := logFile.Close(); ce != nil {
				err = ce
			}
		}
	}()

	if args.Watch {
		return runWatch(stdin, stdout, stderr, args, log)
	}

	if err = runGenerate(stdin, stdout, stderr, args, log); err != nil {
		log.Error(cli.ColorRed(err.Error()), slog.Any("error", err))
		return err
	}
	return nil
}
