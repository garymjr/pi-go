package rpc

import (
	"context"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/garymjr/pi-go/wire"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func buildTestProc(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "testproc")
	cmd := exec.Command("go", "build", "-o", bin, "./internal/testproc")
	cmd.Dir = repoRoot(t)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build testproc: %v\n%s", err, out)
	}
	return bin
}

func newClientFor(t *testing.T, scenario string) *Client {
	t.Helper()
	return NewClient(Options{
		Executable:      buildTestProc(t),
		ExtraArgs:       []string{"--scenario", scenario},
		StartupDelay:    500 * time.Millisecond,
		ResponseTimeout: time.Second,
	})
}

func TestStartupSuccess(t *testing.T) {
	c := newClientFor(t, "startup_success")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
}

func TestStartupImmediateExit(t *testing.T) {
	c := newClientFor(t, "startup_exit")
	err := c.Start(context.Background())
	if err == nil || c.Stderr() == "" {
		t.Fatalf("err=%v stderr=%q", err, c.Stderr())
	}
}

func TestStartFailureAllowsRetry(t *testing.T) {
	c := NewClient(Options{Executable: filepath.Join(t.TempDir(), "missing-bin")})
	err := c.Start(context.Background())
	if err == nil {
		t.Fatal("expected startup error")
	}
	if !strings.Contains(err.Error(), "no such file") {
		t.Fatalf("unexpected error: %v", err)
	}
	c.opts.Executable = buildTestProc(t)
	c.opts.ExtraArgs = []string{"--scenario", "startup_success"}
	c.opts.StartupDelay = 500 * time.Millisecond
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
}

func TestGetState(t *testing.T) {
	c := newClientFor(t, "get_state")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	resp, err := c.GetState(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Data.SessionID != "s1" {
		t.Fatalf("session id = %q", resp.Data.SessionID)
	}
}

func TestConcurrentOutOfOrder(t *testing.T) {
	c := newClientFor(t, "concurrent")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	errCh := make(chan error, 2)
	go func() { defer wg.Done(); errCh <- c.Abort(context.Background()) }()
	go func() { defer wg.Done(); errCh <- c.AbortRetry(context.Background()) }()
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestConcurrentLargeRequests(t *testing.T) {
	c := newClientFor(t, "concurrent_large")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	msg := strings.Repeat("x", 128*1024)
	var wg sync.WaitGroup
	wg.Add(2)
	errCh := make(chan error, 2)
	go func() {
		defer wg.Done()
		errCh <- c.Prompt(context.Background(), wire.PromptCommand{Message: msg})
	}()
	go func() {
		defer wg.Done()
		errCh <- c.FollowUp(context.Background(), wire.FollowUpCommand{Message: msg})
	}()
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestPromptAndNotificationsContinue(t *testing.T) {
	c := newClientFor(t, "prompt_stream")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	seen := 0
	unsub := c.OnNotification(func(n Notification) {
		if _, ok := n.(EventNotification); ok {
			seen++
		}
	})
	defer unsub()
	if err := c.Prompt(context.Background(), wire.PromptCommand{Message: "hi"}); err != nil {
		t.Fatal(err)
	}
	if err := c.WaitForIdle(context.Background()); err != nil {
		t.Fatal(err)
	}
	if seen == 0 {
		t.Fatal("expected notifications")
	}
}

func TestCollectEventsUntilIdle(t *testing.T) {
	c := newClientFor(t, "prompt_stream")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	go func() { _ = c.Prompt(context.Background(), wire.PromptCommand{Message: "hi"}) }()
	events, err := c.CollectEventsUntilIdle(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(events) < 2 {
		t.Fatalf("events=%d", len(events))
	}
}

func TestExtensionUIRoundTrip(t *testing.T) {
	c := newClientFor(t, "ui_roundtrip")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	gotUI := make(chan struct{}, 1)
	c.OnNotification(func(n Notification) {
		ui, ok := n.(UIRequestNotification)
		if !ok {
			return
		}
		if _, ok := ui.Request.(*wire.InputUIRequest); ok {
			gotUI <- struct{}{}
			_ = c.SendUIResponse(context.Background(), wire.ExtensionUIResponseValue{Type: "extension_ui_response", ID: "ui1", Value: "ok"})
		}
	})
	if err := c.Prompt(context.Background(), wire.PromptCommand{Message: "hi"}); err != nil {
		t.Fatal(err)
	}
	select {
	case <-gotUI:
	case <-time.After(time.Second):
		t.Fatal("no ui request")
	}
	if err := c.WaitForIdle(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestNotificationHandlerCanRequestState(t *testing.T) {
	c := newClientFor(t, "notification_get_state")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	done := make(chan error, 1)
	var once sync.Once
	unsub := c.OnNotification(func(n Notification) {
		evt, ok := n.(EventNotification)
		if !ok {
			return
		}
		if _, ok := evt.Event.(*wire.AgentStartEvent); !ok {
			return
		}
		once.Do(func() {
			_, err := c.GetState(ctx)
			done <- err
		})
	})
	defer unsub()
	if err := c.Prompt(ctx, wire.PromptCommand{Message: "hi"}); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
}

func TestCommandFailure(t *testing.T) {
	c := newClientFor(t, "command_failure")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	err := c.Abort(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*CommandError); !ok {
		t.Fatalf("err=%T", err)
	}
}

func TestInvalidJSONCausesProtocolFailure(t *testing.T) {
	c := newClientFor(t, "invalid_json")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	time.Sleep(100 * time.Millisecond)
	err := c.Abort(context.Background())
	if err == nil {
		t.Fatal("expected protocol error")
	}
}

func TestProcessExitFailsPending(t *testing.T) {
	c := newClientFor(t, "process_exit_pending")
	if err := c.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	_, err := c.GetState(context.Background())
	if err == nil {
		t.Fatal("expected exit error")
	}
}
