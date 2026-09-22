//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abox-dev/sdk/packages/go-sdk"
)

func TestCommandStreamingAttachDetachKVM(t *testing.T) {
	client, err := agentbox.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	sandbox, err := client.Sandboxes.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 3 * time.Minute, Metadata: map[string]string{"sdk": "go-streaming-smoke"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if _, err := sandbox.Kill(cleanupCtx); err != nil {
			t.Errorf("sandbox cleanup: %v", err)
		}
	})
	policy := &agentbox.CommandStreamingOptions{StdoutChannel: true}
	handle, err := sandbox.Commands.Start(ctx, "cat", &agentbox.CommandOptions{Stdin: true, Tag: "streaming-smoke", Streaming: policy})
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	pid, err := handle.PID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	exchange := func(handle *agentbox.CommandHandle, message string) {
		t.Helper()
		if _, err := handle.Write(ctx, []byte(message)); err != nil {
			t.Fatal(err)
		}
		var got []byte
		for !bytes.Contains(got, []byte(message)) {
			select {
			case chunk, ok := <-handle.Stdout:
				if !ok {
					t.Fatal("process output closed")
				}
				got = append(got, chunk...)
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
		}
	}
	exchange(handle, "before detach\n")
	for i := range 3 {
		handle.Close()
		result, err := handle.Wait(ctx)
		if !errors.Is(err, context.Canceled) || result.PID != pid || len(result.Stdout) != 0 || len(result.Stderr) != 0 {
			t.Fatalf("detach: %+v %v", result, err)
		}
		attachPID, tag := pid, ""
		if i%2 == 1 {
			attachPID, tag = 0, "streaming-smoke"
		}
		handle, err = sandbox.Commands.ConnectWithOptions(ctx, attachPID, tag, &agentbox.CommandConnectOptions{Streaming: policy})
		if err != nil {
			t.Fatal(err)
		}
		defer handle.Close()
		exchange(handle, "after detach\n")
	}
	if err := handle.CloseStdin(ctx); err != nil {
		t.Fatal(err)
	}
	for range handle.Stdout {
	}
	result, err := handle.Wait(ctx)
	if err != nil || result.ExitCode != 0 || result.PID != pid || result.Stdout != nil || result.Stderr != nil {
		t.Fatalf("exit: %+v %v", result, err)
	}
}
