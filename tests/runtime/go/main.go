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
	defer func() { _, _ = sandbox.Kill(context.Background()) }()
	result, err := sandbox.Commands.Run(ctx, "printf", &agentbox.CommandOptions{Args: []string{"go-runtime-smoke"}})
	must(err)
	if string(result.Stdout) != "go-runtime-smoke" {
		panic(fmt.Sprintf("unexpected stdout %q", result.Stdout))
	}
	killed, err := sandbox.Kill(ctx)
	must(err)
	if !killed {
		panic("core sandbox was not killed")
	}

	interpreter, err := codeinterpreter.NewClient()
	must(err)
	codeSandbox, err := interpreter.Create(ctx, &agentbox.CreateSandboxOptions{Timeout: 5 * time.Minute})
	must(err)
	defer func() { _, _ = codeSandbox.Kill(context.Background()) }()
	execution, err := codeSandbox.RunCode(ctx, "40 + 2", nil)
	must(err)
	if !strings.Contains(execution.Text(), "42") {
		panic(fmt.Sprintf("unexpected execution %#v", execution))
	}
	killed, err = codeSandbox.Kill(ctx)
	must(err)
	if !killed {
		panic("code interpreter sandbox was not killed")
	}
	templateUserSmoke(ctx, core)
}

func templateUserSmoke(ctx context.Context, client *agentbox.Client) {
	name := fmt.Sprintf("sdk-go-user-smoke-%d", time.Now().UnixNano())
	builder := agentbox.NewTemplate("").FromTemplate("base").User("user").Workdir("/home/user").
		RunAs("", "test \"$(id -un)\" = user", "mkdir -p sdk-user-smoke", "printf '{\"name\":\"sdk-user-smoke\",\"version\":\"1.0.0\"}' >sdk-user-smoke/package.json", "git init --bare sdk-user-smoke/source.git").
		Workdir("/home/user/sdk-user-smoke").
		NPMInstall(agentbox.PackageInstallOptions{}).
		GitClone("/home/user/sdk-user-smoke/source.git", &agentbox.GitCloneOptions{Path: "checkout"}).
		RunAs("root", "test \"$(id -un)\" = root").
		Run("test \"$(id -un)\" = user", "test -f package-lock.json", "test -d checkout/.git")
	reference, err := client.Templates.BuildInBackground(ctx, builder, name, nil)
	must(err)
	defer func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		must(client.Templates.Delete(cleanupCtx, reference.TemplateID))
	}()
	for {
		status, err := client.Templates.BuildStatus(ctx, reference.TemplateID, reference.BuildID, 0)
		must(err)
		if status.Status == agentbox.BuildReady {
			fmt.Println("Go template user inheritance smoke passed")
			return
		}
		if status.Status == agentbox.BuildFailed {
			panic(fmt.Sprintf("template user smoke failed: %+v", status.Reason))
		}
		timer := time.NewTimer(time.Second)
		select {
		case <-ctx.Done():
			timer.Stop()
			panic(ctx.Err())
		case <-timer.C:
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
