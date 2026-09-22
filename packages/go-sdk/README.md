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

API reference for this release: [core SDK on pkg.go.dev](https://pkg.go.dev/github.com/abox-dev/sdk/packages/go-sdk@v0.1.8) and
[Code Interpreter on pkg.go.dev](https://pkg.go.dev/github.com/abox-dev/sdk/packages/go-sdk/codeinterpreter@v0.1.8).

Documentation: [core SDK](https://docs.agentbox.ru/en/sdk/),
[sandboxes](https://docs.agentbox.ru/en/sdk/sandboxes/),
[templates](https://docs.agentbox.ru/en/sdk/templates/), and
[Code Interpreter](https://docs.agentbox.ru/en/sdk/code-interpreter/).

`Sandbox.Kill` returns `false, nil` when the sandbox no longer exists.

## Command output

By default, command handles collect complete stdout/stderr for `Wait`, which can
be called without draining their output channels. Channel tails may finish
draining after `Wait` returns. Callbacks also receive output in this mode.

For long-lived processes, opt into streaming on both start and reconnect:

```go
policy := &agentbox.CommandStreamingOptions{}
handle, err := sandbox.Commands.Start(ctx, "cat", &agentbox.CommandOptions{
    Stdin: true,
    Streaming: policy,
    OnStdout: func(chunk []byte) { consume(chunk) },
    // No stderr callback or channel: discard stderr without retaining it.
})
if err != nil {
    return err
}
defer handle.Close()
pid, err := handle.PID(ctx)
if err != nil {
    return err
}
// Exchange messages, then detach locally. The remote process and stdin stay open.
_ = handle.Close()
_, _ = handle.Wait(ctx) // reports context.Canceled for an active attachment
attached, err := sandbox.Commands.ConnectWithOptions(ctx, pid, "", &agentbox.CommandConnectOptions{
    Streaming: policy,
    OnStdout: func(chunk []byte) { consume(chunk) },
})
if err != nil {
    return err
}
defer attached.Close()
```

A non-nil `Streaming` selects exactly one delivery path per output: its callback,
its explicitly enabled channel (`StdoutChannel`, `StderrChannel`, `PTYChannel`),
or discard. Selecting both a callback and a channel for the same output returns
`InvalidArgumentError`. Unselected channels are already closed. Callbacks execute
serially on the receiver; their read-only slices are valid until the callback
returns. Copy data that must be retained and return promptly. Cancellation closes
the transport but cannot interrupt callback code; `Done`/`Wait` finish after it
returns. Do not call `Wait` from a callback.

Streaming channels are unbuffered. Consume every enabled channel concurrently
with `Wait`; a stalled reader applies backpressure and can eventually stall the
remote process. Byte order is preserved, but transport chunks may be split into
slices of at most **32 KiB**. The receiver owns those slices. There is no output
queue or captured result, only one pending channel slice plus one in-flight
transport event. Each encoded/decompressed Connect message is limited to **4 MiB**;
an oversized message fails the attachment with a resource-exhausted error. HTTP
buffers, decoding allocations and caller-retained slices are additional memory.
The limit is independent of total process output, not an exact heap/RSS cap.

`Wait` returns PID, status, exit code and errors, with nil stdout/stderr in streaming
mode. Successful completion means all enabled channel bytes were handed to their
readers; their processing may still be running. Nonzero exits return
`CommandExitError` with an empty output result and empty server error message.
Transport failures retain their error code but omit server-provided diagnostic
text/details that might contain output. Neither path silently truncates output
and reports success.

Cancel the context passed to `Start`/`ConnectWithOptions`, or call `Close`, to
release the local attachment without killing the process or sending EOF.
Cancellation returns `context.Canceled`/`context.DeadlineExceeded`, unless a process
end event was already confirmed. `Close` also discards an unread collecting-mode
channel tail; natural completion preserves that tail. Canceling only the context
passed to `Wait` stops waiting without detaching. Reconnecting by PID or tag adds
no stdout replay, delivery retry, process restart or exactly-once guarantee.

PTY callers can use `PTYOptions.Streaming` with `OnPTY` or `PTYChannel`, and
`PTY.ConnectWithOptions` to reattach with the same policy. Default PTY behavior
remains unchanged. See the [streaming contract](../../docs/go-command-streaming.md)
for bounds and validation, and `examples_test.go` for a complete example.

Streaming uploads have no SDK deadline by default. Set
`WriteFileOptions.RequestTimeout` to limit a complete upload. A client supplied
with `WithHTTPClient` remains authoritative, so its `http.Client.Timeout` also
applies to command, watch, download, upload, and Code Interpreter streams. Keep
that timeout at zero and use contexts or operation-specific options when streams
may be long-lived.
