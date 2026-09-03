//go:build runtime

package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/abox-dev/sdk/packages/go-sdk"
	"github.com/abox-dev/sdk/packages/go-sdk/codeinterpreter"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()
	core, err := agentbox.NewClient()
	must(err)
	sandbox, err := core.Sandboxes.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 5 * time.Minute})
	must(err)
	result, err := sandbox.Commands.Run(ctx, "printf", &agentbox.CommandOptions{Args: []string{"go-runtime-smoke"}})
	must(err)
	if string(result.Stdout) != "go-runtime-smoke" {
		panic(fmt.Sprintf("unexpected stdout %q", result.Stdout))
	}
	must(sandbox.Kill(ctx))

	interpreter, err := codeinterpreter.NewClient()
	must(err)
	codeSandbox, err := interpreter.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 5 * time.Minute})
	must(err)
	execution, err := codeSandbox.RunCode(ctx, "40 + 2", nil)
	must(err)
	if !strings.Contains(execution.Text(), "42") {
		panic(fmt.Sprintf("unexpected execution %#v", execution))
	}
	must(codeSandbox.Kill(ctx))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
