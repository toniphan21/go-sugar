package lspcmd

import (
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"nhatp.com/go/sugar/internal/util"
	"nhatp.com/go/sugar/lsp"
)

const cmdName = "lsp"

type Arguments struct {
	Log      string
	LogLevel slog.Level
}

func Run(stdin io.Reader, stdout io.Writer, stderr io.Writer, args Arguments) error {
	log, logFile, err := util.NewLogger(args.Log, stderr, args.LogLevel)
	if err != nil {
		return err
	}
	log = log.With("cmd", cmdName).WithGroup(cmdName)

	// start
	log.Info("start")
	defer func() {
		log.Info("stop")
		if logFile != nil {
			if ce := logFile.Close(); ce != nil {
				err = ce
			}
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// gopls
	cmd, goplsIn, goplsOut, err := startGopls(ctx)
	if err != nil {
		return err
	}
	log.Info("started gopls")

	// proxy
	proxy := lsp.NewProxy(log, stdin, stdout, goplsOut, goplsIn, lsp.DoNothing, lsp.DoNothing)
	proxy.Start()
	log.Info("start proxying gopls 2")

	// run
	select {
	case <-ctx.Done():
		log.Info("ctx done")
	case <-proxy.Done():
		log.Info("proxy done")
	}

	log.Info("shutting down")
	return cleanUp(cmd, goplsIn, goplsOut)
}

func startGopls(ctx context.Context) (*exec.Cmd, io.WriteCloser, io.ReadCloser, error) {
	cmd := exec.CommandContext(ctx, "gopls", "serve")
	cmd.Stderr = os.Stderr

	goplsIn, err := cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	goplsOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, nil, nil, err
	}
	return cmd, goplsIn, goplsOut, nil
}

func cleanUp(cmd *exec.Cmd, closers ...io.Closer) error {
	var err error
	for _, closer := range closers {
		if err = closer.Close(); err != nil {
			return err
		}
	}

	if err = cmd.Process.Kill(); err != nil {
		return err
	}
	if err = cmd.Wait(); err != nil {
		return err
	}
	return nil
}
