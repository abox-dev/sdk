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
	Args     []string
	Env      map[string]string
	Cwd      string
	Tag      string
	Stdin    bool
	OnStdout func([]byte)
	OnStderr func([]byte)
}

// CommandResult contains collected process output.
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

// CommandHandle represents a streaming process. Wait can be called without
// draining the output channels and always returns the complete collected output.
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

	stdout *outputStream
	stderr *outputStream
	pty    *outputStream
	done   chan struct{}
	result CommandResult
	err    error
}

func newCommandService(sandbox *Sandbox) *CommandService {
	client := processconnect.NewProcessClient(sandbox.client.envdClient, sandbox.envdURL(envdPort, false), connect.WithCodec(tolerantJSONCodec{}), connect.WithAcceptCompression("gzip", nil, nil))
	return &CommandService{sandbox: sandbox, client: client}
}

// Run executes a foreground command and collects its output.
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
	stream, err := service.client.Start(ctx, request)
	if err != nil {
		return nil, connectError(err)
	}
	handle := newCommandHandle(service, options.Tag)
	go handle.receive(ctx, func() (*process.ProcessEvent, bool) {
		if !stream.Receive() {
			return nil, false
		}
		return stream.Msg().GetEvent(), true
	}, stream.Err, stream.Close, outputCallbacks{stdout: options.OnStdout, stderr: options.OnStderr})
	return handle, nil
}

// Connect attaches to an existing process by PID or tag.
func (service *CommandService) Connect(ctx context.Context, pid uint32, tag string) (*CommandHandle, error) {
	selector, err := processSelector(pid, tag)
	if err != nil {
		return nil, err
	}
	request := connect.NewRequest(&process.ConnectRequest{Process: selector})
	service.addHeaders(request.Header())
	stream, err := service.client.Connect(ctx, request)
	if err != nil {
		return nil, connectError(err)
	}
	handle := newCommandHandle(service, tag)
	if pid != 0 {
		handle.pid = pid
		handle.closeReady()
	}
	go handle.receive(ctx, func() (*process.ProcessEvent, bool) {
		if !stream.Receive() {
			return nil, false
		}
		return stream.Msg().GetEvent(), true
	}, stream.Err, stream.Close, outputCallbacks{})
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

func newCommandHandle(service *CommandService, tag string) *CommandHandle {
	ready := make(chan struct{})
	stdout := newOutputStream()
	stderr := newOutputStream()
	pty := newOutputStream()
	done := make(chan struct{})
	handle := &CommandHandle{
		service: service, ready: ready, closeReady: sync.OnceFunc(func() { close(ready) }),
		tag:    tag,
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
	defer func() { _ = closeStream() }()
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
			switch output := data.GetOutput().(type) {
			case *process.ProcessEvent_DataEvent_Stdout:
				chunk := bytes.Clone(output.Stdout)
				handle.result.Stdout = append(handle.result.Stdout, chunk...)
				handle.stdout.send(chunk)
				if callbacks.stdout != nil {
					callbacks.stdout(chunk)
				}
			case *process.ProcessEvent_DataEvent_Stderr:
				chunk := bytes.Clone(output.Stderr)
				handle.result.Stderr = append(handle.result.Stderr, chunk...)
				handle.stderr.send(chunk)
				if callbacks.stderr != nil {
					callbacks.stderr(chunk)
				}
			case *process.ProcessEvent_DataEvent_Pty:
				chunk := bytes.Clone(output.Pty)
				handle.pty.send(chunk)
				if callbacks.pty != nil {
					callbacks.pty(chunk)
				}
			}
		}
		if end := event.GetEnd(); end != nil {
			handle.closeReady()
			handle.result.ExitCode = int(end.GetExitCode())
			handle.result.Status = end.GetStatus()
			if end.GetExited() && end.GetExitCode() != 0 {
				handle.err = &CommandExitError{Result: handle.result, Message: end.GetError()}
			}
			return
		}
	}
	handle.closeReady()
	if err := streamErr(); err != nil && !errors.Is(err, context.Canceled) {
		handle.err = connectError(err)
	}
}

type outputStream struct {
	output    chan []byte
	wake      chan struct{}
	aborted   chan struct{}
	abortOnce func()
	mu        sync.Mutex
	queue     [][]byte
	closed    bool
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

func (stream *outputStream) send(chunk []byte) {
	stream.mu.Lock()
	if !stream.closed {
		stream.queue = append(stream.queue, chunk)
	}
	stream.mu.Unlock()
	stream.notify()
}

func (stream *outputStream) close() {
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

// Wait waits for completion and returns collected output.
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
