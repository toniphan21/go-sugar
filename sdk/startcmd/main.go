package startcmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/exp/jsonrpc2"
	"nhatp.com/go/sugar"
	"nhatp.com/go/sugar/internal/sdk/transport"
	"nhatp.com/go/sugar/internal/util"
)

type Arguments struct {
	Sugar    sugar.Sugar
	Log      string
	LogLevel slog.Level
}

func Run(stdin, stdout, stderr *os.File, args Arguments) error {
	log, logFile, err := util.NewCLILogger(args.Log, stderr, args.LogLevel)
	if err != nil {
		return err
	}
	log = log.With("cmd", "start").WithGroup("server")
	defer func() {
		msg := "done"
		log.Info(msg)
		if logFile != nil {
			if ce := logFile.Close(); ce != nil {
				err = ce
			}
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	handler := &Handler{
		sugar: args.Sugar,
		log:   log,
	}

	binder := jsonrpc2.ConnectionOptions{
		Framer:  jsonrpc2.HeaderFramer(),
		Handler: handler,
	}

	listener := newStdioListener(stdin, stdout, log)
	defer func() { _ = listener.Close() }()

	srv, err := jsonrpc2.Serve(ctx, listener, binder)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		log.Info("shutting down")
	}()

	log.Info(fmt.Sprintf("jsonrpc server for %v started", args.Sugar.ID()))
	log.Info("supported node", "types", transport.Registered())
	if err := srv.Wait(); err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
	}
	return nil
}

type stdioListener struct {
	stdin    *os.File
	stdout   *os.File
	log      *slog.Logger
	accepted chan struct{}
}

func newStdioListener(stdin, stdout *os.File, log *slog.Logger) *stdioListener {
	return &stdioListener{
		stdin:    stdin,
		stdout:   stdout,
		log:      log,
		accepted: make(chan struct{}, 1),
	}
}

func (l *stdioListener) Accept(ctx context.Context) (io.ReadWriteCloser, error) {
	select {
	case l.accepted <- struct{}{}:
		return &stdio{
			stdin:  l.stdin,
			stdout: l.stdout,
			log:    l.log,
		}, nil
	default:
		<-ctx.Done()
		return nil, ctx.Err()
	}
}

func (l *stdioListener) Close() error {
	return nil
}

func (l *stdioListener) Dialer() jsonrpc2.Dialer {
	return nil
}

type stdio struct {
	stdin  *os.File
	stdout *os.File
	log    *slog.Logger
}

func (s *stdio) Read(p []byte) (int, error) {
	return s.stdin.Read(p)
}

func (s *stdio) Write(p []byte) (int, error) {
	return s.stdout.Write(p)
}

func (s *stdio) Close() error {
	err := s.stdin.Close()
	if err != nil {
		return err
	}
	return s.stdout.Close()
}
