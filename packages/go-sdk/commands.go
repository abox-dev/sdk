package agentbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"slices"
	"sync"

	"connectrpc.com/connect"
	process "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process/processconnect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// CommandOptions configures a command process.
type CommandOptions struct {
	// User selects the process owner. Empty uses the template default.
	User     string
	Args     []string
	Env      map[string]string
	Cwd      string
	Tag      string
	Stdin    bool
	OnStdout func([]byte)
	OnStderr func([]byte)
	// Streaming opts out of output capture; nil preserves collecting behavior.
	Streaming *CommandStreamingOptions
}

// CommandStreamingOptions enables output delivery without retaining command output.
// Each output uses its callback, its explicitly enabled channel, or is discarded.
// A callback and channel for the same output are mutually exclusive.
// Callbacks run synchronously and must return promptly. Their slices are read-only
// and valid only until the callback returns; copy bytes that must outlive it.
// Channel slices belong to the receiver and contain at most 32 KiB each. Channels
// are unbuffered: all enabled channels must be consumed concurrently with Wait.
// A stalled reader pauses transport reads and may eventually stall the process.
// The SDK retains no output queue and at most one pending 32 KiB channel slice,
// plus one transport event (limited to 4 MiB encoded and decompressed). Transport
// decoding and HTTP buffers add overhead; caller-retained bytes are not bounded.
// Oversized events fail the local attachment; bytes are never silently dropped.
// Reattachment adds no replay, retry, restart, or exactly-once guarantee.
type CommandStreamingOptions struct {
	StdoutChannel bool
	StderrChannel bool
	PTYChannel    bool
}

// CommandConnectOptions configures output delivery when attaching to a process.
type CommandConnectOptions struct {
	OnStdout func([]byte)
	OnStderr func([]byte)
	OnPTY    func([]byte)
	// Streaming opts out of output capture; nil preserves collecting behavior.
	Streaming *CommandStreamingOptions
}

const commandStreamChunkBytes = 32 << 10
const commandStreamMessageBytes = 4 << 20

// CommandResult contains process metadata and, in collecting mode, output.
type CommandResult struct {
	PID      uint32
	ExitCode int
	Stdout   []byte
	Stderr   []byte
	Status   string
}

// CommandExitError reports a process that completed with a non-zero exit code.
type CommandExitError struct {
	Result  CommandResult
	Message string
}

// Error describes the non-zero command exit code.
func (e *CommandExitError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("agentbox: command exited with code %d: %s", e.Result.ExitCode, e.Message)
	}
	return fmt.Sprintf("agentbox: command exited with code %d", e.Result.ExitCode)
}

// ProcessInfo describes a running envd process.
type ProcessInfo struct {
	PID     uint32
	Tag     string
	Command string
	Args    []string
	Env     map[string]string
	Cwd     string
}

// CommandService executes and manages sandbox processes.
type CommandService struct {
	sandbox *Sandbox
	client  processconnect.ProcessClient
}

// CommandHandle represents a process attachment. In collecting mode, Wait returns
// complete output without requiring channel reads. In streaming mode, Wait returns
// only metadata and requires enabled channels to be consumed. Close detaches
// locally without killing the process or closing its stdin.
type CommandHandle struct {
	service    *CommandService
	ready      chan struct{}
	mu         sync.RWMutex
	closeReady func()
	pid        uint32
	tag        string

	Stdout <-chan []byte
	Stderr <-chan []byte
	PTY    <-chan []byte
	Done   <-chan struct{}

	stdout    *outputStream
	stderr    *outputStream
	pty       *outputStream
	done      chan struct{}
	result    CommandResult
	err       error
	streaming bool
	cancel    context.CancelFunc
}

func newCommandService(sandbox *Sandbox) *CommandService {
	client := processconnect.NewProcessClient(sandbox.client.envdClient, sandbox.envdURL(envdPort, false), connect.WithCodec(tolerantJSONCodec{}), connect.WithAcceptCompression("gzip", nil, nil))
	return &CommandService{sandbox: sandbox, client: client}
}

// Run executes a foreground command, draining channels and waiting for completion.
// It collects output unless CommandOptions.Streaming is set.
func (service *CommandService) Run(ctx context.Context, command string, options *CommandOptions) (CommandResult, error) {
	handle, err := service.Start(ctx, command, options)
	if err != nil {
		return CommandResult{}, err
	}
	stdout, stderr, pty := handle.Stdout, handle.Stderr, handle.PTY
	for stdout != nil || stderr != nil || pty != nil {
		select {
		case <-ctx.Done():
			return CommandResult{}, ctx.Err()
		case _, ok := <-stdout:
			if !ok {
				stdout = nil
			}
		case _, ok := <-stderr:
			if !ok {
				stderr = nil
			}
		case _, ok := <-pty:
			if !ok {
				pty = nil
			}
		}
	}
	return handle.Wait(ctx)
}

// Start starts a process and streams output through the returned handle.
func (service *CommandService) Start(ctx context.Context, command string, options *CommandOptions) (*CommandHandle, error) {
	if command == "" {
		return nil, &InvalidArgumentError{Message: "command cannot be empty"}
	}
	if options == nil {
		options = &CommandOptions{}
	}
	callbacks := outputCallbacks{stdout: options.OnStdout, stderr: options.OnStderr}
	if err := validateCommandStreaming(options.Streaming, callbacks); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	config := &process.ProcessConfig{Cmd: command, Args: options.Args, Envs: options.Env}
	if options.Cwd != "" {
		config.Cwd = &options.Cwd
	}
	stdin := options.Stdin
	request := connect.NewRequest(&process.StartRequest{Process: config, Stdin: &stdin})
	if options.Tag != "" {
		request.Msg.Tag = &options.Tag
	}
	service.addHeaders(request.Header())
	if err := service.sandbox.addUserHeader(request.Header(), options.User); err != nil {
		cancel()
		return nil, err
	}
	stream, err := service.outputClient(options.Streaming).Start(ctx, request)
	if err != nil {
		cancel()
		return nil, connectError(err)
	}
	handle := newCommandHandle(service, options.Tag, options.Streaming)
	handle.cancel = cancel
	go handle.receive(ctx, func() (*process.ProcessEvent, bool) {
		if !stream.Receive() {
			return nil, false
		}
		return stream.Msg().GetEvent(), true
	}, stream.Err, stream.Close, callbacks)
	return handle, nil
}

// Connect attaches to an existing process by PID or tag.
func (service *CommandService) Connect(ctx context.Context, pid uint32, tag string) (*CommandHandle, error) {
	return service.ConnectWithOptions(ctx, pid, tag, nil)
}

// ConnectWithOptions attaches by PID or tag with an explicit output policy.
// Canceling ctx or calling Close detaches locally without sending a signal or EOF.
func (service *CommandService) ConnectWithOptions(ctx context.Context, pid uint32, tag string, options *CommandConnectOptions) (*CommandHandle, error) {
	if options == nil {
		options = &CommandConnectOptions{}
	}
	callbacks := outputCallbacks{stdout: options.OnStdout, stderr: options.OnStderr, pty: options.OnPTY}
	if err := validateCommandStreaming(options.Streaming, callbacks); err != nil {
		return nil, err
	}
	selector, err := processSelector(pid, tag)
	if err != nil {
		return nil, err
	}
	request := connect.NewRequest(&process.ConnectRequest{Process: selector})
	service.addHeaders(request.Header())
	ctx, cancel := context.WithCancel(ctx)
	stream, err := service.outputClient(options.Streaming).Connect(ctx, request)
	if err != nil {
		cancel()
		return nil, connectError(err)
	}
	handle := newCommandHandle(service, tag, options.Streaming)
	handle.cancel = cancel
	if pid != 0 {
		handle.pid = pid
		handle.result.PID = pid
		handle.closeReady()
	}
	go handle.receive(ctx, func() (*process.ProcessEvent, bool) {
		if !stream.Receive() {
			return nil, false
		}
		return stream.Msg().GetEvent(), true
	}, stream.Err, stream.Close, callbacks)
	return handle, nil
}

// List returns currently running processes.
func (service *CommandService) List(ctx context.Context) ([]ProcessInfo, error) {
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	request := connect.NewRequest(&process.ListRequest{})
	service.addHeaders(request.Header())
	response, err := service.client.List(requestCtx, request)
	if err != nil {
		return nil, connectError(err)
	}
	items := make([]ProcessInfo, 0, len(response.Msg.GetProcesses()))
	for _, item := range response.Msg.GetProcesses() {
		cfg := item.GetConfig()
		info := ProcessInfo{PID: item.GetPid(), Tag: item.GetTag()}
		if cfg != nil {
			info.Command = cfg.GetCmd()
			info.Args = cfg.GetArgs()
			info.Env = cfg.GetEnvs()
			info.Cwd = cfg.GetCwd()
		}
		items = append(items, info)
	}
	return items, nil
}

// Kill sends SIGKILL to a process.
func (service *CommandService) Kill(ctx context.Context, pid uint32, tag string) error {
	return service.signal(ctx, pid, tag, process.Signal_SIGNAL_SIGKILL)
}

// Terminate sends SIGTERM to a process.
func (service *CommandService) Terminate(ctx context.Context, pid uint32, tag string) error {
	return service.signal(ctx, pid, tag, process.Signal_SIGNAL_SIGTERM)
}
func (service *CommandService) signal(ctx context.Context, pid uint32, tag string, signal process.Signal) error {
	selector, err := processSelector(pid, tag)
	if err != nil {
		return err
	}
	request := connect.NewRequest(&process.SendSignalRequest{Process: selector, Signal: signal})
	service.addHeaders(request.Header())
	requestCtx, cancel := service.sandbox.unaryContext(ctx)
	defer cancel()
	_, err = service.client.SendSignal(requestCtx, request)
	return connectError(err)
}

func (service *CommandService) addHeaders(header http.Header) {
	for key, values := range service.sandbox.envdHeaders(envdPort) {
		header[key] = slices.Clone(values)
	}
	header.Set("Keepalive-Ping-Interval", "50")
}

func newCommandHandle(service *CommandService, tag string, policy *CommandStreamingOptions) *CommandHandle {
	ready := make(chan struct{})
	var stdout, stderr, pty *outputStream
	if policy == nil {
		stdout, stderr, pty = newOutputStream(), newOutputStream(), newOutputStream()
	} else {
		stdout = newDirectOutputStream(policy.StdoutChannel)
		stderr = newDirectOutputStream(policy.StderrChannel)
		pty = newDirectOutputStream(policy.PTYChannel)
	}
	done := make(chan struct{})
	handle := &CommandHandle{
		service: service, ready: ready, closeReady: sync.OnceFunc(func() { close(ready) }),
		tag: tag, streaming: policy != nil,
		Stdout: stdout.output, Stderr: stderr.output, PTY: pty.output, Done: done,
		stdout: stdout, stderr: stderr, pty: pty, done: done,
	}
	runtime.AddCleanup(handle, cleanupCommandOutputs, commandOutputs{stdout: stdout, stderr: stderr, pty: pty})
	return handle
}

type outputCallbacks struct {
	stdout func([]byte)
	stderr func([]byte)
	pty    func([]byte)
}

type commandOutputs struct{ stdout, stderr, pty *outputStream }

func cleanupCommandOutputs(outputs commandOutputs) {
	outputs.stdout.abort()
	outputs.stderr.abort()
	outputs.pty.abort()
}

func (handle *CommandHandle) receive(ctx context.Context, next func() (*process.ProcessEvent, bool), streamErr, closeStream func() error, callbacks outputCallbacks) {
	defer close(handle.done)
	closeResponse := sync.OnceFunc(func() { _ = closeStream() })
	stopCancel := context.AfterFunc(ctx, func() {
		cleanupCommandOutputs(commandOutputs{handle.stdout, handle.stderr, handle.pty})
		closeResponse()
	})
	defer func() {
		stopCancel()
		if ctx.Err() != nil {
			cleanupCommandOutputs(commandOutputs{handle.stdout, handle.stderr, handle.pty})
		}
		closeResponse()
		if handle.cancel != nil {
			handle.cancel()
		}
	}()
	defer handle.stdout.close()
	defer handle.stderr.close()
	defer handle.pty.close()
	for {
		event, ok := next()
		if !ok {
			break
		}
		if event == nil {
			continue
		}
		if start := event.GetStart(); start != nil {
			handle.mu.Lock()
			handle.pid = start.GetPid()
			handle.result.PID = start.GetPid()
			handle.mu.Unlock()
			handle.closeReady()
			continue
		}
		if data := event.GetData(); data != nil {
			var chunk []byte
			var output *outputStream
			var callback func([]byte)
			switch data := data.GetOutput().(type) {
			case *process.ProcessEvent_DataEvent_Stdout:
				chunk, output, callback = data.Stdout, handle.stdout, callbacks.stdout
				if !handle.streaming {
					handle.result.Stdout = append(handle.result.Stdout, chunk...)
				}
			case *process.ProcessEvent_DataEvent_Stderr:
				chunk, output, callback = data.Stderr, handle.stderr, callbacks.stderr
				if !handle.streaming {
					handle.result.Stderr = append(handle.result.Stderr, chunk...)
				}
			case *process.ProcessEvent_DataEvent_Pty:
				chunk, output, callback = data.Pty, handle.pty, callbacks.pty
			}
			if output != nil && !handle.deliver(ctx, output, callback, chunk) {
				break
			}
		}
		if end := event.GetEnd(); end != nil {
			handle.closeReady()
			handle.result.ExitCode = int(end.GetExitCode())
			handle.result.Status = end.GetStatus()
			if end.GetExited() && end.GetExitCode() != 0 {
				message := end.GetError()
				if handle.streaming {
					message = ""
				}
				handle.err = &CommandExitError{Result: handle.result, Message: message}
			}
			return
		}
	}
	handle.closeReady()
	// Closing the response can surface as a transport read error. The attachment
	// context is authoritative when local cancellation caused delivery to stop.
	if ctx.Err() != nil {
		handle.err = ctx.Err()
		return
	}
	if err := streamErr(); err != nil {
		if handle.streaming && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			// Preserve the error code without retaining server-provided output in errors.
			err = connect.NewError(connect.CodeOf(err), errors.New("command attachment failed"))
		}
		handle.err = connectError(err)
	}
}

func validateCommandStreaming(policy *CommandStreamingOptions, callbacks outputCallbacks) error {
	if policy != nil && ((policy.StdoutChannel && callbacks.stdout != nil) ||
		(policy.StderrChannel && callbacks.stderr != nil) || (policy.PTYChannel && callbacks.pty != nil)) {
		return &InvalidArgumentError{Message: "streaming output must use either a callback or a channel"}
	}
	return nil
}

func (service *CommandService) outputClient(policy *CommandStreamingOptions) processconnect.ProcessClient {
	if policy == nil {
		return service.client
	}
	return processconnect.NewProcessClient(service.sandbox.client.envdClient, service.sandbox.envdURL(envdPort, false),
		connect.WithCodec(tolerantJSONCodec{}), connect.WithAcceptCompression("gzip", nil, nil),
		connect.WithReadMaxBytes(commandStreamMessageBytes))
}

func (handle *CommandHandle) deliver(ctx context.Context, output *outputStream, callback func([]byte), chunk []byte) bool {
	if !handle.streaming {
		chunk = bytes.Clone(chunk)
		output.send(chunk)
		if callback != nil {
			callback(chunk)
		}
		return true
	}
	if ctx.Err() != nil {
		return false
	}
	if callback != nil {
		callback(chunk)
	} else if output.enabled {
		for len(chunk) > 0 {
			n := min(len(chunk), commandStreamChunkBytes)
			part := bytes.Clone(chunk[:n])
			select {
			case output.output <- part:
			case <-ctx.Done():
				return false
			}
			chunk = chunk[n:]
		}
	}
	return ctx.Err() == nil
}

// Close cancels this local attachment and releases queued output, including an
// unread collecting-mode tail. It does not kill the process or send stdin EOF.
// Close does not wait for user callbacks; Done closes after the receiver exits.
// A callback must return before Wait can complete. Close is safe to call repeatedly.
func (handle *CommandHandle) Close() error {
	if handle.cancel != nil {
		handle.cancel()
	}
	cleanupCommandOutputs(commandOutputs{handle.stdout, handle.stderr, handle.pty})
	return nil
}

type outputStream struct {
	output    chan []byte
	wake      chan struct{}
	aborted   chan struct{}
	abortOnce func()
	mu        sync.Mutex
	queue     [][]byte
	closed    bool
	direct    bool
	enabled   bool
}

func newOutputStream() *outputStream {
	stream := &outputStream{output: make(chan []byte), wake: make(chan struct{}, 1), aborted: make(chan struct{})}
	stream.abortOnce = sync.OnceFunc(func() {
		stream.mu.Lock()
		clear(stream.queue)
		stream.queue = nil
		stream.closed = true
		stream.mu.Unlock()
		close(stream.aborted)
	})
	go stream.run()
	return stream
}

func newDirectOutputStream(enabled bool) *outputStream {
	stream := &outputStream{output: make(chan []byte), direct: true, enabled: enabled}
	stream.abortOnce = func() {} // Direct delivery owns no queue or background goroutine.
	if !enabled {
		close(stream.output)
	}
	return stream
}

func (stream *outputStream) send(chunk []byte) {
	stream.mu.Lock()
	if !stream.closed {
		stream.queue = append(stream.queue, chunk)
	}
	stream.mu.Unlock()
	stream.notify()
}

func (stream *outputStream) close() {
	if stream.direct {
		if stream.enabled {
			close(stream.output)
		}
		return
	}
	stream.mu.Lock()
	stream.closed = true
	stream.mu.Unlock()
	stream.notify()
}

func (stream *outputStream) abort() { stream.abortOnce() }

func (stream *outputStream) notify() {
	select {
	case stream.wake <- struct{}{}:
	default:
	}
}

func (stream *outputStream) run() {
	defer close(stream.output)
	for {
		stream.mu.Lock()
		if len(stream.queue) > 0 {
			chunk := stream.queue[0]
			stream.queue[0] = nil
			stream.queue = stream.queue[1:]
			stream.mu.Unlock()
			select {
			case stream.output <- chunk:
			case <-stream.aborted:
				return
			}
			continue
		}
		closed := stream.closed
		stream.mu.Unlock()
		if closed {
			return
		}
		select {
		case <-stream.wake:
		case <-stream.aborted:
			return
		}
	}
}

// PID waits for and returns the process identifier.
func (handle *CommandHandle) PID(ctx context.Context) (uint32, error) {
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-handle.ready:
		handle.mu.RLock()
		defer handle.mu.RUnlock()
		if handle.pid == 0 {
			return 0, errors.New("agentbox: process did not start")
		}
		return handle.pid, nil
	}
}

// Wait waits for receiver completion. Collecting mode preserves the queued channel
// tail and returns full output. Streaming mode returns metadata with nil output;
// all enabled channels must be read concurrently, and are closed before Done.
// Attachment cancellation returns a cancellation error unless a process end event
// was already confirmed. CommandExitError preserves confirmed nonzero exits.
// Canceling only Wait's context stops waiting; use Close to cancel the attachment.
func (handle *CommandHandle) Wait(ctx context.Context) (CommandResult, error) {
	select {
	case <-ctx.Done():
		return CommandResult{}, ctx.Err()
	case <-handle.done:
		return handle.result, handle.err
	}
}

// Write writes bytes to process stdin.
func (handle *CommandHandle) Write(ctx context.Context, data []byte) (int, error) {
	pid, err := handle.PID(ctx)
	if err != nil {
		return 0, err
	}
	request := connect.NewRequest(&process.SendInputRequest{Process: &process.ProcessSelector{Selector: &process.ProcessSelector_Pid{Pid: pid}}, Input: &process.ProcessInput{Input: &process.ProcessInput_Stdin{Stdin: data}}})
	handle.service.addHeaders(request.Header())
	requestCtx, cancel := handle.service.sandbox.unaryContext(ctx)
	defer cancel()
	_, err = handle.service.client.SendInput(requestCtx, request)
	if err != nil {
		return 0, connectError(err)
	}
	return len(data), nil
}

// CloseStdin signals EOF to a non-PTY process.
func (handle *CommandHandle) CloseStdin(ctx context.Context) error {
	pid, err := handle.PID(ctx)
	if err != nil {
		return err
	}
	request := connect.NewRequest(&process.CloseStdinRequest{Process: &process.ProcessSelector{Selector: &process.ProcessSelector_Pid{Pid: pid}}})
	handle.service.addHeaders(request.Header())
	requestCtx, cancel := handle.service.sandbox.unaryContext(ctx)
	defer cancel()
	_, err = handle.service.client.CloseStdin(requestCtx, request)
	return connectError(err)
}

// Kill sends SIGKILL to this process.
func (handle *CommandHandle) Kill(ctx context.Context) error {
	pid, err := handle.PID(ctx)
	if err != nil {
		return err
	}
	return handle.service.Kill(ctx, pid, "")
}

func processSelector(pid uint32, tag string) (*process.ProcessSelector, error) {
	if pid != 0 && tag != "" {
		return nil, &InvalidArgumentError{Message: "provide process PID or tag, not both"}
	}
	if pid != 0 {
		return &process.ProcessSelector{Selector: &process.ProcessSelector_Pid{Pid: pid}}, nil
	}
	if tag != "" {
		return &process.ProcessSelector{Selector: &process.ProcessSelector_Tag{Tag: tag}}, nil
	}
	return nil, &InvalidArgumentError{Message: "process PID or tag is required"}
}

type tolerantJSONCodec struct{}

func (tolerantJSONCodec) Name() string { return "json" }
func (tolerantJSONCodec) Marshal(value any) ([]byte, error) {
	message, ok := value.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("agentbox: expected protobuf message, got %T", value)
	}
	return protojson.MarshalOptions{UseProtoNames: false}.Marshal(message)
}
func (tolerantJSONCodec) Unmarshal(data []byte, value any) error {
	message, ok := value.(proto.Message)
	if !ok {
		return fmt.Errorf("agentbox: expected protobuf message, got %T", value)
	}
	return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, message)
}

func connectError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return &TimeoutError{APIError: APIError{Message: "envd request timed out", Cause: err}}
	}
	var value *connect.Error
	if errors.As(err, &value) {
		apiError := APIError{Code: value.Code().String(), Message: value.Message(), Cause: err}
		switch value.Code() {
		case connect.CodeUnauthenticated, connect.CodePermissionDenied:
			return &AuthenticationError{APIError: apiError}
		case connect.CodeResourceExhausted:
			return &RateLimitError{APIError: apiError}
		case connect.CodeCanceled, connect.CodeDeadlineExceeded, connect.CodeUnavailable:
			return &TimeoutError{APIError: apiError}
		case connect.CodeInvalidArgument:
			return &InvalidArgumentError{Message: value.Message(), Cause: err}
		default:
			return &SandboxError{APIError: apiError}
		}
	}
	return err
}
