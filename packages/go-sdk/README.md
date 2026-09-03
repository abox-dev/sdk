<p align="center"><a href="https://agentbox.ru"><img src="https://raw.githubusercontent.com/abox-dev/sdk/main/readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240"></a></p>

# AgentBox Go SDK

The official Go client for AgentBox sandboxes, templates, and Code Interpreter.
Go 1.24 or newer is required.

```bash
go get github.com/abox-dev/sdk/packages/go-sdk
```

```go
package main

import (
    "context"
    "log"

    agentbox "github.com/abox-dev/sdk/packages/go-sdk"
)

func main() {
    ctx := context.Background()
    client, err := agentbox.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    sandbox, err := client.Sandboxes.Create(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer sandbox.Kill(context.Background())

    result, err := sandbox.Commands.Run(ctx, "echo", &agentbox.CommandOptions{
        Args: []string{"Hello from AgentBox"},
    })
    if err != nil {
        log.Fatal(err)
    }
    log.Print(string(result.Stdout))
}
```

Code Interpreter is available from
`github.com/abox-dev/sdk/packages/go-sdk/codeinterpreter`.

API reference for this release: [core SDK on pkg.go.dev](https://pkg.go.dev/github.com/abox-dev/sdk/packages/go-sdk@v0.1.2) and
[Code Interpreter on pkg.go.dev](https://pkg.go.dev/github.com/abox-dev/sdk/packages/go-sdk/codeinterpreter@v0.1.2).

Documentation: [core SDK](https://docs.agentbox.ru/en/sdk/),
[sandboxes](https://docs.agentbox.ru/en/sdk/sandboxes/),
[templates](https://docs.agentbox.ru/en/sdk/templates/), and
[Code Interpreter](https://docs.agentbox.ru/en/sdk/code-interpreter/).

`Sandbox.Kill` returns `false, nil` when the sandbox no longer exists. Command
handles can be waited on without draining their live output channels; `Wait`
always returns complete stdout and stderr collected by the SDK. Output channels
may finish draining and close after `Wait` returns. PTY callers can consume
`CommandHandle.PTY` or use `PTYOptions.OnPTY`.

Streaming uploads have no SDK deadline by default. Set
`WriteFileOptions.RequestTimeout` to limit a complete upload. A client supplied
with `WithHTTPClient` remains authoritative, so its `http.Client.Timeout` also
applies to command, watch, download, upload, and Code Interpreter streams. Keep
that timeout at zero and use contexts or operation-specific options when streams
may be long-lived.
