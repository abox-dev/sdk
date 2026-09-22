# Go command streaming without output retention

## Contract and owner

The handwritten Go command transport in
[`commands.go`](../packages/go-sdk/commands.go) supports explicit streaming on
start and reattachment. [`pty.go`](../packages/go-sdk/pty.go) shares the delivery
implementation. No backend, wire or generated binding changes are required.
The minimum Go version remains 1.24.0.

A non-nil `CommandOptions.Streaming`, `CommandConnectOptions.Streaming` or
`PTYOptions.Streaming` selects streaming. Use `Commands.ConnectWithOptions` or
`PTY.ConnectWithOptions` for attachment by PID or tag. Existing `Connect` calls
and nil options retain collecting behavior.

| Concern | Collecting mode | Streaming mode |
| --- | --- | --- |
| Result | Full stdout/stderr and exit metadata | PID, status, exit code and errors; nil stdout/stderr |
| Callbacks | Existing callback and channel delivery | Synchronous callback only for that output |
| Channels | Unbounded queues with a preserved natural-completion tail | Explicitly enabled, unbuffered, at most 32 KiB per delivered slice |
| Unselected output | Captured and queued | Discarded; its public channel is closed |
| Callback plus channel | Existing fan-out | Invalid for the same output |
| `Wait` without reads | Supported | Supported only when no channels are enabled; callbacks must return |
| Local detach | `Close` or attachment context cancellation | Same; no signal or stdin EOF |

`CommandStreamingOptions` has `StdoutChannel`, `StderrChannel` and `PTYChannel`
flags. Existing start callbacks remain in `CommandOptions`/`PTYOptions`;
`CommandConnectOptions` supplies `OnStdout`, `OnStderr` and `OnPTY`.
`Run` continues to drain channels itself and collects output by default.

## Memory and delivery bounds

Streaming creates no queue goroutines or output history. An enabled channel uses
synchronous delivery, splitting events into independent slices of at most 32 KiB.
The SDK retains at most one pending delivery slice plus the current decoded
transport event. A streaming Connect client imposes a 4 MiB message-size limit;
this also bounds messages selected for discard. Decode buffers and HTTP transport
buffers add overhead. This is a bound independent of total emitted bytes, not a
promise of a 4 MiB total heap or RSS. Memory retained by the consumer is outside
the SDK bound.

A reader that stops blocks subsequent event processing, including exit events and
other outputs. Backpressure can propagate to the remote process. Read every
enabled channel concurrently, or select callbacks/discard for unused outputs.
There is no silent overflow, truncation, or goroutine per chunk. An oversized
transport message fails and closes the local attachment.

Callbacks execute serially on the receiver. Their byte slices are read-only and
valid only until the callback returns; copy bytes needed later. Channel receivers
own their independent slices. Channel chunk boundaries may differ from transport
boundaries, but byte order and content are preserved. Slow callbacks also apply
backpressure. Callbacks must not call `Wait` on their own handle.

## Completion, cancellation and errors

Natural streaming completion hands every enabled-channel byte to a reader before
closing the channels and `Done`; application processing may finish later.
Collecting-mode completion preserves its queued channel tail, even after `Wait`.

`CommandHandle.Close` cancels the local attachment and releases queued output.
Canceling the context supplied to start/connect does the same while the receiver
is active. Both close the HTTP response without relying on GC and without killing
the process or closing stdin. `Close` is idempotent and does not wait for arbitrary
callback code. `Done`/`Wait` complete when that callback returns and the receiver
exits. Canceling only `Wait`'s context does not detach.

Local cancellation reports `context.Canceled` or `context.DeadlineExceeded`.
A confirmed process end retains its result, including `CommandExitError` for a
nonzero exit. In streaming mode that error contains no stdout/stderr or server
error message. Transport failures preserve Connect error classification but use
generic diagnostic text, avoiding output-bearing server error details.

Reattachment adds no stdout replay, retry, restart or exactly-once promise.
Consumers must select the streaming policy on both launch and reconnect.

## Verification

Hermetic regressions in
[`command_streaming_test.go`](../packages/go-sdk/command_streaming_test.go) cover:

- Multiple MiB of ordered stdout/stderr through callbacks and channels, including
  inspection of a still-live handle at a synchronized callback barrier.
- No capture or hidden channel queue; unused stderr discard and closed channels.
- Channel blocks larger than 32 KiB, byte ownership and ordered splitting.
- Unread output preventing completion, context cancellation, repeated detach,
  response-body closure, and no signal/EOF side effect.
- Start, attach by PID/tag, nonzero exit, transport failure and oversized events.
- Cancellation while a callback is blocked, collecting tail release on `Close`,
  PTY callback/channel delivery and PTY detach.

Existing collecting-mode no-drain, queued-tail, command/PTY and response-closure
regressions remain required. The supported Go matrix is 1.24–1.27; run
`make go-check` and generation/artifact gates described in
[`RELEASING.md`](../RELEASING.md) using the pinned containers.

Validation passed: `make go-check` (91.2% handwritten statement coverage), builds
and hermetic tests on Go 1.24.13, 1.25.14, 1.26.8 and 1.27.1, race checks on
1.24.13 and 1.27.1, format/lint/type checks, workspace tests, release-artifact
builds and clean installation checks. `make generate` and reference-contract
checks passed; a second generation produced identical tracked files. This does
not replace the full multi-language KVM release suite required before publication.

`BenchmarkCommandOutputRetention` compares callback consumers at fixed 1 KiB
chunks, retaining the handle until observation and leaving channels unread.
A Go 1.24.13 linux/arm64 sample (`-benchtime=1x -benchmem`) measured:

| Emitted stdout | Collecting retained output storage | Streaming retained output storage | Streaming allocated bytes |
| --- | --- | --- | --- |
| 1 MiB | 2,161,664 bytes | 0 bytes | 1,904 bytes |
| 16 MiB | 37,428,224 bytes | 0 bytes | 1,904 bytes |

Retained storage counts capture capacity and queued slices, excluding queue
metadata and any collecting slice already in flight. The benchmark isolates SDK
delivery from transport decoding. Allocation counts are supporting evidence;
structural live-retention assertions are the deterministic gate.

[`TestCommandStreamingAttachDetachKVM`](../packages/go-sdk/integration/command_streaming_test.go)
starts `cat`, exchanges bytes, detaches/reconnects three times using PID and tag,
exchanges further bytes, then explicitly closes stdin and checks exit metadata.
It creates and deletes only its own sandbox. This runtime smoke has passed.

## Delivery

The implementation must follow the coordinated SDK release procedure in
[`RELEASING.md`](../RELEASING.md). No ad-hoc fork or replacement module is needed.
Publication and the consumer dependency upgrade remain release steps; no released
module version is assigned by this source change. The consuming orchestrator must
use the released version and explicitly enable streaming on launch and reconnect.
