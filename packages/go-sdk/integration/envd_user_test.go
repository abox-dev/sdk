//go:build integration

package integration_test

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/abox-dev/sdk/packages/go-sdk"
)

func TestEnvdUsersKVM(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()
	client, err := agentbox.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	sandbox, err := client.Sandboxes.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 3 * time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if _, err := sandbox.Kill(cleanup); err != nil {
			t.Errorf("sandbox cleanup: %v", err)
		}
	})
	for _, user := range []string{"root", "user"} {
		t.Run(user, func(t *testing.T) {
			result, err := sandbox.Commands.Run(ctx, "id", &agentbox.CommandOptions{User: user, Args: []string{"-un"}})
			if err != nil || strings.TrimSpace(string(result.Stdout)) != user {
				t.Fatalf("identity: %q %v", result.Stdout, err)
			}
			var output strings.Builder
			h, err := sandbox.Commands.Start(ctx, "id", &agentbox.CommandOptions{User: user, Args: []string{"-un"}, Streaming: &agentbox.CommandStreamingOptions{}, OnStdout: func(p []byte) { output.Write(p) }})
			if err != nil {
				t.Fatal(err)
			}
			defer h.Close()
			if _, err := h.Wait(ctx); err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(output.String()) != user {
				t.Fatalf("streaming identity: %q", output.String())
			}
			var terminal strings.Builder
			pty, err := sandbox.PTY.Create(ctx, "id", &agentbox.PTYOptions{User: user, Args: []string{"-un"}, Streaming: &agentbox.CommandStreamingOptions{}, OnPTY: func(p []byte) { terminal.Write(p) }})
			if err != nil {
				t.Fatal(err)
			}
			defer pty.Close()
			if _, err := pty.Wait(ctx); err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(terminal.String()) != user {
				t.Fatalf("PTY identity: %q", terminal.String())
			}

			options := &agentbox.FileOptions{User: user}
			directory := fmt.Sprintf("~/go-sdk-user-%d", time.Now().UnixNano())
			entry, err := sandbox.Files.MakeDir(ctx, directory, options)
			if err != nil {
				t.Fatal(err)
			}
			if entry.Owner != user {
				t.Fatalf("directory owner: %q", entry.Owner)
			}
			watchCtx, stop := context.WithTimeout(ctx, 15*time.Second)
			defer stop()
			watcher, err := sandbox.Files.Watch(watchCtx, directory, &agentbox.WatchOptions{User: user, IncludeEntry: true})
			if err != nil {
				t.Fatal(err)
			}
			defer watcher.Close()
			file := directory + "/file"
			if _, err := sandbox.Files.WriteText(ctx, file, user, &agentbox.WriteFileOptions{User: user}); err != nil {
				t.Fatal(err)
			}
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					t.Fatalf("watch ended: %v", watcher.Close())
				}
				if event.Entry == nil || event.Entry.Path != path.Join(entry.Path, "file") {
					t.Fatalf("watch path: %#v", event)
				}
			case <-watchCtx.Done():
				t.Fatal("watch did not observe the selected user's directory")
			}
			if err := watcher.Close(); err != nil {
				t.Fatal(err)
			}
			info, err := sandbox.Files.Stat(ctx, file, options)
			if err != nil || info.Owner != user {
				t.Fatalf("file owner: %#v %v", info, err)
			}
			text, err := sandbox.Files.ReadText(ctx, file, options)
			if err != nil || text != user {
				t.Fatalf("read: %q %v", text, err)
			}
			entries, err := sandbox.Files.List(ctx, directory, &agentbox.ListFilesOptions{User: user, Depth: 1})
			if err != nil || len(entries) != 1 {
				t.Fatalf("list: %#v %v", entries, err)
			}
			if _, err := sandbox.Files.Rename(ctx, file, directory+"/renamed", options); err != nil {
				t.Fatal(err)
			}
			if exists, err := sandbox.Files.Exists(ctx, directory+"/renamed", options); err != nil || !exists {
				t.Fatalf("exists: %v %v", exists, err)
			}
			if err := sandbox.Files.Remove(ctx, directory, options); err != nil {
				t.Fatal(err)
			}
		})
	}
	_, err = sandbox.Commands.Run(ctx, "id", &agentbox.CommandOptions{User: "agentbox-user-does-not-exist"})
	var auth *agentbox.AuthenticationError
	if !errors.As(err, &auth) {
		t.Fatalf("unknown user: %T %v", err, err)
	}
}
