package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/garymjr/pi-go/wire"
)

type Client struct {
	opts Options

	mu          sync.Mutex
	writeMu     sync.Mutex
	started     bool
	closed      bool
	cmd         *exec.Cmd
	stdin       io.WriteCloser
	stdout      io.ReadCloser
	stderrPipe  io.ReadCloser
	encoder     *wire.Encoder
	pending     map[string]chan wire.Response
	handlers    map[uint64]NotificationHandler
	nextHandler uint64
	nextReq     uint64
	procErr     error
	notifyCh    chan Notification

	stderrMu     sync.Mutex
	stderrBuffer []byte

	readDone   chan struct{}
	stderrDone chan struct{}
	waitDone   chan struct{}
	closeOnce  sync.Once
}

func NewClient(opts Options) *Client {
	return &Client{
		opts:       opts.withDefaults(),
		pending:    map[string]chan wire.Response{},
		handlers:   map[uint64]NotificationHandler{},
		notifyCh:   make(chan Notification, 32),
		readDone:   make(chan struct{}),
		stderrDone: make(chan struct{}),
		waitDone:   make(chan struct{}),
	}
}

func (c *Client) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.started {
		c.mu.Unlock()
		return ErrAlreadyStarted
	}
	c.started = true
	c.mu.Unlock()
	launched := false
	defer func() {
		if launched {
			return
		}
		c.mu.Lock()
		c.started = false
		c.mu.Unlock()
	}()

	args := []string{"--mode", "rpc"}
	if c.opts.Provider != "" {
		args = append(args, "--provider", c.opts.Provider)
	}
	if c.opts.Model != "" {
		args = append(args, "--model", c.opts.Model)
	}
	if c.opts.NoSession {
		args = append(args, "--no-session")
	}
	if c.opts.SessionDir != "" {
		args = append(args, "--session-dir", c.opts.SessionDir)
	}
	args = append(args, c.opts.ExtraArgs...)

	cmd := exec.CommandContext(ctx, c.opts.Executable, args...)
	cmd.Dir = c.opts.Dir
	if len(c.opts.Env) > 0 {
		cmd.Env = append(os.Environ(), c.opts.Env...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	launched = true

	c.mu.Lock()
	c.cmd = cmd
	c.stdin = stdin
	c.stdout = stdout
	c.stderrPipe = stderrPipe
	c.encoder = wire.NewEncoder(stdin)
	c.mu.Unlock()

	go c.readStdout()
	go c.runNotifications()
	go c.captureStderr()
	go c.waitProcess()

	t := time.NewTimer(c.opts.StartupDelay)
	defer t.Stop()
	select {
	case <-ctx.Done():
		_ = c.Close()
		return ctx.Err()
	case <-c.waitDone:
		return fmt.Errorf("startup failed: %w: %s", ErrProcessExited, c.Stderr())
	case <-t.C:
		c.mu.Lock()
		err := c.procErr
		cmd := c.cmd
		c.mu.Unlock()
		if err != nil {
			return fmt.Errorf("startup failed: %w: %s", err, c.Stderr())
		}
		if cmd != nil && cmd.Process != nil {
			if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
				return fmt.Errorf("startup failed: %w: %s", ErrProcessExited, c.Stderr())
			}
		}
		return nil
	}
}

func (c *Client) Close() error {
	var ret error
	c.closeOnce.Do(func() {
		c.mu.Lock()
		alreadyClosed := c.closed
		c.closed = true
		cmd := c.cmd
		pending := c.pending
		c.pending = map[string]chan wire.Response{}
		c.mu.Unlock()
		if alreadyClosed {
			ret = nil
			return
		}
		for _, ch := range pending {
			close(ch)
		}
		if cmd == nil || cmd.Process == nil {
			ret = nil
			return
		}
		_ = cmd.Process.Signal(syscall.SIGTERM)
		timer := time.NewTimer(c.opts.ShutdownTimeout)
		defer timer.Stop()
		select {
		case <-c.waitDone:
		case <-timer.C:
			_ = cmd.Process.Kill()
			<-c.waitDone
		}
	})
	return ret
}

func (c *Client) Stderr() string {
	c.stderrMu.Lock()
	defer c.stderrMu.Unlock()
	return string(c.stderrBuffer)
}

func (c *Client) OnNotification(fn NotificationHandler) func() {
	id := atomic.AddUint64(&c.nextHandler, 1)
	c.mu.Lock()
	c.handlers[id] = fn
	c.mu.Unlock()
	return func() {
		c.mu.Lock()
		delete(c.handlers, id)
		c.mu.Unlock()
	}
}

func (c *Client) SendUIResponse(ctx context.Context, resp any) error {
	return c.write(ctx, resp)
}

func (c *Client) dispatch(n Notification) {
	c.notifyCh <- n
}

func (c *Client) runNotifications() {
	for n := range c.notifyCh {
		c.mu.Lock()
		handlers := make([]NotificationHandler, 0, len(c.handlers))
		for _, h := range c.handlers {
			handlers = append(handlers, h)
		}
		c.mu.Unlock()
		for _, h := range handlers {
			h(n)
		}
	}
}

func (c *Client) readStdout() {
	defer close(c.notifyCh)
	defer close(c.readDone)
	dec := wire.NewDecoder(c.stdout)
	for {
		frame, err := dec.Decode()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			c.fail(fmt.Errorf("%w: %v", ErrProtocol, err))
			return
		}
		switch f := frame.(type) {
		case wire.Response:
			c.resolveResponse(f)
		case wire.Event:
			c.dispatch(EventNotification{Event: f})
		case wire.UIRequest:
			c.dispatch(UIRequestNotification{Request: f})
		default:
			if _, ok := frame.(wire.UnknownFrame); ok {
				continue
			}
		}
	}
}

func (c *Client) captureStderr() {
	defer close(c.stderrDone)
	buf := make([]byte, 4096)
	for {
		n, err := c.stderrPipe.Read(buf)
		if n > 0 {
			c.appendStderr(buf[:n])
		}
		if err != nil {
			return
		}
	}
}

func (c *Client) waitProcess() {
	defer close(c.waitDone)
	err := c.cmd.Wait()
	if err != nil {
		c.fail(fmt.Errorf("%w: %v", ErrProcessExited, err))
		return
	}
	c.fail(ErrProcessExited)
}

func (c *Client) fail(err error) {
	c.mu.Lock()
	if c.procErr != nil {
		c.mu.Unlock()
		return
	}
	c.procErr = err
	pending := c.pending
	c.pending = map[string]chan wire.Response{}
	c.mu.Unlock()
	for _, ch := range pending {
		close(ch)
	}
}

func (c *Client) appendStderr(p []byte) {
	if len(p) == 0 {
		return
	}
	if w := c.opts.Stderr; w != nil {
		_, _ = w.Write(p)
	}
	c.stderrMu.Lock()
	defer c.stderrMu.Unlock()
	c.stderrBuffer = append(c.stderrBuffer, p...)
	if extra := len(c.stderrBuffer) - c.opts.StderrBufferBytes; extra > 0 {
		c.stderrBuffer = append([]byte(nil), c.stderrBuffer[extra:]...)
	}
}

func (c *Client) resolveResponse(resp wire.Response) {
	id := responseID(resp)
	if id == "" {
		return
	}
	c.mu.Lock()
	ch := c.pending[id]
	delete(c.pending, id)
	c.mu.Unlock()
	if ch == nil {
		return
	}
	ch <- resp
	close(ch)
}

func responseID(resp wire.Response) string {
	switch r := resp.(type) {
	case *wire.PromptResponse:
		return r.ID
	case *wire.SteerResponse:
		return r.ID
	case *wire.FollowUpResponse:
		return r.ID
	case *wire.AbortResponse:
		return r.ID
	case *wire.NewSessionResponse:
		return r.ID
	case *wire.GetStateResponse:
		return r.ID
	case *wire.SetModelResponse:
		return r.ID
	case *wire.CycleModelResponse:
		return r.ID
	case *wire.GetAvailableModelsResponse:
		return r.ID
	case *wire.SetThinkingLevelResponse:
		return r.ID
	case *wire.CycleThinkingLevelResponse:
		return r.ID
	case *wire.SetSteeringModeResponse:
		return r.ID
	case *wire.SetFollowUpModeResponse:
		return r.ID
	case *wire.CompactResponse:
		return r.ID
	case *wire.SetAutoCompactionResponse:
		return r.ID
	case *wire.SetAutoRetryResponse:
		return r.ID
	case *wire.AbortRetryResponse:
		return r.ID
	case *wire.BashResponse:
		return r.ID
	case *wire.AbortBashResponse:
		return r.ID
	case *wire.GetSessionStatsResponse:
		return r.ID
	case *wire.ExportHTMLResponse:
		return r.ID
	case *wire.SwitchSessionResponse:
		return r.ID
	case *wire.ForkResponse:
		return r.ID
	case *wire.GetForkMessagesResponse:
		return r.ID
	case *wire.GetLastAssistantTextResponse:
		return r.ID
	case *wire.SetSessionNameResponse:
		return r.ID
	case *wire.GetMessagesResponse:
		return r.ID
	case *wire.GetCommandsResponse:
		return r.ID
	case *wire.ErrorResponse:
		return r.ID
	default:
		return ""
	}
}

func (c *Client) write(ctx context.Context, v any) error {
	c.mu.Lock()
	if !c.started {
		c.mu.Unlock()
		return ErrNotStarted
	}
	if c.closed {
		c.mu.Unlock()
		return ErrClientClosed
	}
	if c.procErr != nil {
		err := c.procErr
		c.mu.Unlock()
		return err
	}
	enc := c.encoder
	c.mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		c.writeMu.Lock()
		defer c.writeMu.Unlock()
		errCh <- enc.Encode(v)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			c.fail(err)
			return err
		}
		return nil
	}
}

func (c *Client) send(ctx context.Context, cmdType string, cmd any) (wire.Response, error) {
	id := fmt.Sprintf("req_%d", atomic.AddUint64(&c.nextReq, 1))
	line, err := injectID(cmd, id)
	if err != nil {
		return nil, err
	}
	respCh := make(chan wire.Response, 1)
	c.mu.Lock()
	if !c.started {
		c.mu.Unlock()
		return nil, ErrNotStarted
	}
	if c.closed {
		c.mu.Unlock()
		return nil, ErrClientClosed
	}
	if c.procErr != nil {
		err := c.procErr
		c.mu.Unlock()
		return nil, err
	}
	c.pending[id] = respCh
	c.mu.Unlock()

	if err := c.writeLine(ctx, line); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, err
	}

	timeout := c.opts.ResponseTimeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)
		if remaining < timeout {
			timeout = remaining
		}
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, ctx.Err()
	case <-timer.C:
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return nil, fmt.Errorf("%w: %s", ErrResponseTimeout, cmdType)
	case resp, ok := <-respCh:
		if !ok {
			c.mu.Lock()
			err := c.procErr
			closed := c.closed
			c.mu.Unlock()
			if err != nil {
				return nil, err
			}
			if closed {
				return nil, ErrClientClosed
			}
			return nil, ErrProcessExited
		}
		if errResp, ok := resp.(*wire.ErrorResponse); ok && !errResp.Success {
			return nil, &CommandError{Command: errResp.Command, Message: errResp.Error}
		}
		return resp, nil
	}
}

func injectID(v any, id string) ([]byte, error) {
	data, err := wire.EncodeLine(v)
	if err != nil {
		return nil, err
	}
	data = bytes.TrimSuffix(data, []byte("\n"))
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	m["id"] = id
	return wire.EncodeLine(m)
}

func (c *Client) writeLine(ctx context.Context, line []byte) error {
	c.mu.Lock()
	if !c.started {
		c.mu.Unlock()
		return ErrNotStarted
	}
	if c.closed {
		c.mu.Unlock()
		return ErrClientClosed
	}
	if c.procErr != nil {
		err := c.procErr
		c.mu.Unlock()
		return err
	}
	stdin := c.stdin
	c.mu.Unlock()

	errCh := make(chan error, 1)
	go func() {
		c.writeMu.Lock()
		defer c.writeMu.Unlock()
		_, err := stdin.Write(line)
		errCh <- err
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errCh:
		if err != nil {
			c.fail(err)
		}
		return err
	}
}
