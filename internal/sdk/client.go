package sdk

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"golang.org/x/exp/jsonrpc2"
)

const DefaultTimeout = time.Second * 5

func NewProcess(path string, log string) Process {
	return &processImpl{
		path: path,
		log:  log,
	}
}

type Process interface {
	Call(ctx context.Context, method string, params json.RawMessage, result any) error

	Close() error
}

type processImpl struct {
	mu     sync.Mutex
	path   string
	log    string
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	conn   *jsonrpc2.Connection
	done   chan struct{}
}

func (p *processImpl) Call(ctx context.Context, method string, params json.RawMessage, result any) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.ensureStarted(); err != nil {
		return err
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, DefaultTimeout)
		defer cancel()
	}

	call := p.conn.Call(ctx, method, params)
	return call.Await(ctx, result)
}

func (p *processImpl) ensureStarted() error {
	if p.conn != nil {
		return nil
	}

	conn, err := jsonrpc2.Dial(context.Background(), p, jsonrpc2.ConnectionOptions{
		Framer: jsonrpc2.HeaderFramer(),
	})
	if err != nil {
		return err
	}
	p.conn = conn
	return nil
}

func (p *processImpl) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	if p.cmd == nil {
		args := []string{"start"}
		if p.log != "" {
			args = append(args, "-log", p.log)
		}

		cmd := exec.CommandContext(ctx, p.path, args...)
		cmd.SysProcAttr = sysProcAttr()

		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			if ce := stdin.Close(); ce != nil {
				return nil, ce
			}
			return nil, err
		}
		cmd.Stderr = log.Writer()

		if err := cmd.Start(); err != nil {
			if ce := stdin.Close(); ce != nil {
				return nil, ce
			}
			if ce := stdout.Close(); ce != nil {
				return nil, ce
			}
			return nil, err
		}

		p.cmd = cmd
		p.stdin = stdin
		p.stdout = stdout
		p.done = make(chan struct{})
		done := p.done

		go func() {
			_ = cmd.Wait()
			close(done)

			p.mu.Lock()
			defer p.mu.Unlock()
			if p.cmd == cmd {
				p.cmd = nil
				p.conn = nil
			}
		}()
	}
	return p, nil
}

func (p *processImpl) Read(out []byte) (n int, err error) {
	return p.stdout.Read(out)
}

func (p *processImpl) Write(out []byte) (n int, err error) {
	return p.stdin.Write(out)
}

func (p *processImpl) Close() error {
	p.mu.Lock()
	if p.cmd == nil {
		p.mu.Unlock()
		return nil
	}

	done := p.done
	pcmd := p.cmd

	_ = p.stdin.Close()
	_ = p.stdout.Close()
	_ = p.cmd.Process.Signal(syscall.SIGTERM)
	p.mu.Unlock()

	select {
	case <-done:
		// exited cleanly
	case <-time.After(2 * time.Second):
		_ = pcmd.Process.Kill()
		<-done
	}
	return nil
}

var _ Process = (*processImpl)(nil)
