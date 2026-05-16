// Proxy is an LSP message pump with structural awareness.
//
// It understands two layers of structure:
//   - Wire framing: Content-Length headers
//   - Message envelope: JSON-RPC fields (jsonrpc, id, method)
//
// It does NOT understand meaning. It never inspects Params, Result, or Error.
// Those stay as raw bytes until the bridge layer opens them.
//
// Think of it as a post office: it reads the address on the envelope to route
// messages through hooks, but never opens the letter inside.
//
// All semantic concerns (method dispatch, position translation, source map)
// belong in the bridge layer, which produces the hooks.
//
// Lifecycle (starting/stopping gopls, handling OS signals, closing streams)
// belongs in the cmd layer. The proxy does not own its streams and never
// closes them. When the caller closes a stream, the blocked read returns
// an error, the pump exits, and Done() fires.

package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
)

// Proxy pumps LSP messages between two streams.
// It parses each message into an Envelope, applies a EnvelopeHook, and forwards it.
type Proxy struct {
	logger *slog.Logger

	r  io.Reader // source: read from here
	w  io.Writer // source: write back here
	tr io.Reader // target: read from here
	tw io.Writer // target: write to here

	toTarget   EnvelopeHook
	fromTarget EnvelopeHook

	done     chan struct{}
	doneOnce sync.Once
}

func NewProxy(
	logger *slog.Logger,
	r io.Reader, w io.Writer,
	tr io.Reader, tw io.Writer,
	toTarget, fromTarget EnvelopeHook,
) *Proxy {
	return &Proxy{
		logger:     logger,
		r:          r,
		w:          w,
		tr:         tr,
		tw:         tw,
		toTarget:   toTarget,
		fromTarget: fromTarget,
		done:       make(chan struct{}),
	}
}

// Start launches both pump goroutines and returns immediately.
func (p *Proxy) Start() {
	logger := p.logger.WithGroup("proxy")
	go p.run(logger, "editor->gopls", p.r, p.tw, p.toTarget)
	go p.run(logger, "gopls->editor", p.tr, p.w, p.fromTarget)
}

// Done returns a channel that is closed when either pump exits.
// The caller is responsible for closing the underlying streams
// to ensure both goroutines exit cleanly.
func (p *Proxy) Done() <-chan struct{} {
	return p.done
}

func (p *Proxy) run(logger *slog.Logger, direction string, r io.Reader, w io.Writer, hook EnvelopeHook) {
	defer p.doneOnce.Do(func() {
		close(p.done)
		logger.Info("run done once do")
	})
	p.pump(logger, direction, r, w, hook)
}

// pump is a pure function with no state. It reads frames, parses them
// into envelopes, passes them through a hook, and writes the result
// until an error occurs.
func (p *Proxy) pump(log *slog.Logger, direction string, r io.Reader, w io.Writer, hook EnvelopeHook) {
	br := bufio.NewReader(r)
	for {
		env, err := p.readEnvelope(br)
		if err != nil {
			log.Error("read envelope", "direction", direction, "error", err)
			return
		}

		log.Info("pump", env.logArgs(direction)...)

		out, err := hook(*env)
		if err != nil {
			log.Error("hook", "direction", direction, "error", err)
			return
		}

		if err = p.writeEnvelope(w, &out); err != nil {
			log.Error("write envelope", "direction", direction, "error", err)
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Wire framing — Content-Length: N  CRLF  CRLF  <N bytes of JSON>
// ---------------------------------------------------------------------------

func (p *Proxy) readFrame(r *bufio.Reader) ([]byte, error) {
	contentLength := -1

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")

		if line == "" {
			break
		}

		const prefix = "Content-Length: "
		if strings.HasPrefix(line, prefix) {
			n, err := strconv.Atoi(strings.TrimPrefix(line, prefix))
			if err != nil {
				return nil, fmt.Errorf("bad Content-Length: %w", err)
			}
			contentLength = n
		}
	}

	if contentLength < 0 {
		return nil, fmt.Errorf("no Content-Length header")
	}

	body := make([]byte, contentLength)
	if _, err := io.ReadFull(r, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (p *Proxy) writeFrame(w io.Writer, body []byte) error {
	_, err := fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(body))
	if err != nil {
		return err
	}
	_, err = w.Write(body)
	return err
}

// ---------------------------------------------------------------------------
// Envelope — JSON-RPC structure parsed, inner fields left as raw bytes
// ---------------------------------------------------------------------------

// Envelope is a JSON-RPC message with the inner fields left as raw bytes.
// The proxy parses the structure; the bridge interprets the content.
type Envelope struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   json.RawMessage `json:"error,omitempty"`
}

func (e *Envelope) logArgs(direction string) []any {
	return []any{
		slog.String("direction", direction),
		slog.String("id", string(e.ID)),
		slog.String("method", e.Method),
		slog.Any("params", e.Params),
		slog.Any("result", e.Result),
		slog.String("error", string(e.Error)),
	}
}

// EnvelopeHook can inspect and rewrite a message envelope.
// Return the same envelope to pass through, a modified envelope to rewrite,
// or nil to drop the message.
type EnvelopeHook func(Envelope) (Envelope, error)

func (p *Proxy) readEnvelope(r *bufio.Reader) (*Envelope, error) {
	body, err := p.readFrame(r)
	if err != nil {
		return nil, err
	}
	var env Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("bad JSON-RPC message: %w", err)
	}
	return &env, nil
}

func (p *Proxy) writeEnvelope(w io.Writer, env *Envelope) error {
	body, err := json.Marshal(env)
	if err != nil {
		return err
	}
	return p.writeFrame(w, body)
}
