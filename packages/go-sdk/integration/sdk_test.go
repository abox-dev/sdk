//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/abox-dev/sdk/packages/go-sdk"
	"github.com/abox-dev/sdk/packages/go-sdk/codeinterpreter"
)

func TestCoreKVM(t *testing.T) {
	client, err := agentbox.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Minute)
	defer cancel()
	sandbox, err := client.Sandboxes.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 5 * time.Minute, Metadata: map[string]string{"sdk": "go-integration"}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = sandbox.Kill(context.Background()) })
	result, err := sandbox.Commands.Run(ctx, "sh", &agentbox.CommandOptions{Args: []string{"-lc", "printf go-sdk"}})
	if err != nil {
		t.Fatal(err)
	}
	if string(result.Stdout) != "go-sdk" {
		t.Fatalf("stdout: %q", result.Stdout)
	}
	if _, err := sandbox.Files.WriteText(ctx, "/tmp/go-sdk.txt", "content", nil); err != nil {
		t.Fatal(err)
	}
	if text, err := sandbox.Files.ReadText(ctx, "/tmp/go-sdk.txt", ""); err != nil || text != "content" {
		t.Fatalf("file: %q %v", text, err)
	}
	if info, err := sandbox.Info(ctx); err != nil || info.SandboxID != sandbox.ID {
		t.Fatalf("info: %#v %v", info, err)
	}
}

func TestCodeInterpreterKVM(t *testing.T) {
	client, err := codeinterpreter.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Minute)
	defer cancel()
	sandbox, err := client.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 5 * time.Minute})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _, _ = sandbox.Kill(context.Background()) })
	execution, err := sandbox.RunCode(ctx, "6 * 7", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(execution.Text(), "42") {
		t.Fatalf("result: %#v", execution)
	}
	codeContext, err := sandbox.CreateContext(ctx, &codeinterpreter.CreateContextOptions{Language: codeinterpreter.Python})
	if err != nil {
		t.Fatal(err)
	}
	if err := sandbox.RemoveContext(ctx, codeContext.ID); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateKVM(t *testing.T) {
	client, err := agentbox.NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Minute)
	defer cancel()
	name := fmt.Sprintf("go-sdk-integration-%d", time.Now().Unix())
	reference, err := client.Templates.Build(ctx, agentbox.NewTemplate("").FromAlpine("").Run("printf go-sdk >/tmp/go-sdk"), name, &agentbox.TemplateBuildOptions{CPUCount: 2, MemoryMB: 1024})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = client.Templates.Delete(context.Background(), reference.TemplateID) })
}
