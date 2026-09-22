package agentbox

import (
	"context"

	"connectrpc.com/connect"
	process "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process"
)

// PTYOptions configures an interactive terminal.
type PTYOptions struct {
	// User selects the process owner. Empty uses the template default (user on envd < 0.4.0).
	User  string
	Args  []string
	Env   map[string]string
	Cwd   string
	Tag   string
	Cols  uint32
	Rows  uint32
	OnPTY func([]byte)
	// Streaming selects bounded delivery; nil preserves collecting-mode PTY queues.
	Streaming *CommandStreamingOptions
}

// PTYService manages pseudo-terminal processes.
type PTYService struct{ commands *CommandService }

// Create starts a process attached to a pseudo-terminal.
func (service *PTYService) Create(ctx context.Context, command string, options *PTYOptions) (*CommandHandle, error) {
	if command == "" {
		return nil, &InvalidArgumentError{Message: "command cannot be empty"}
	}
	if options == nil {
		options = &PTYOptions{}
	}
	callbacks := outputCallbacks{pty: options.OnPTY}
	if err := validateCommandStreaming(options.Streaming, callbacks); err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	cols, rows := options.Cols, options.Rows
	if cols == 0 {
		cols = 80
	}
	if rows == 0 {
		rows = 24
	}
	config := &process.ProcessConfig{Cmd: command, Args: options.Args, Envs: options.Env}
	if options.Cwd != "" {
		config.Cwd = &options.Cwd
	}
	request := connect.NewRequest(&process.StartRequest{Process: config, Pty: &process.PTY{Size: &process.PTY_Size{Cols: cols, Rows: rows}}})
	if options.Tag != "" {
		request.Msg.Tag = &options.Tag
	}
	service.commands.addHeaders(request.Header())
	if err := service.commands.sandbox.addUserHeader(request.Header(), options.User); err != nil {
		cancel()
		return nil, err
	}
	stream, err := service.commands.outputClient(options.Streaming).Start(ctx, request)
	if err != nil {
		cancel()
		return nil, connectError(err)
	}
	handle := newCommandHandle(service.commands, options.Tag, options.Streaming)
	handle.cancel = cancel
	go handle.receive(ctx, func() (*process.ProcessEvent, bool) {
		if !stream.Receive() {
			return nil, false
		}
		return stream.Msg().GetEvent(), true
	}, stream.Err, stream.Close, callbacks)
	return handle, nil
}

// Connect attaches to an existing PTY process.
func (service *PTYService) Connect(ctx context.Context, pid uint32, tag string) (*CommandHandle, error) {
	return service.commands.Connect(ctx, pid, tag)
}

// ConnectWithOptions attaches to an existing PTY using the given output policy.
// Use OnPTY or Streaming.PTYChannel for terminal output.
func (service *PTYService) ConnectWithOptions(ctx context.Context, pid uint32, tag string, options *CommandConnectOptions) (*CommandHandle, error) {
	return service.commands.ConnectWithOptions(ctx, pid, tag, options)
}

// Input sends terminal input.
func (service *PTYService) Input(ctx context.Context, handle *CommandHandle, data []byte) error {
	pid, err := handle.PID(ctx)
	if err != nil {
		return err
	}
	request := connect.NewRequest(&process.SendInputRequest{Process: &process.ProcessSelector{Selector: &process.ProcessSelector_Pid{Pid: pid}}, Input: &process.ProcessInput{Input: &process.ProcessInput_Pty{Pty: data}}})
	service.commands.addHeaders(request.Header())
	requestCtx, cancel := service.commands.sandbox.unaryContext(ctx)
	defer cancel()
	_, err = service.commands.client.SendInput(requestCtx, request)
	return connectError(err)
}

// Resize changes terminal dimensions.
func (service *PTYService) Resize(ctx context.Context, handle *CommandHandle, cols, rows uint32) error {
	if cols == 0 || rows == 0 {
		return &InvalidArgumentError{Message: "PTY columns and rows must be positive"}
	}
	pid, err := handle.PID(ctx)
	if err != nil {
		return err
	}
	request := connect.NewRequest(&process.UpdateRequest{Process: &process.ProcessSelector{Selector: &process.ProcessSelector_Pid{Pid: pid}}, Pty: &process.PTY{Size: &process.PTY_Size{Cols: cols, Rows: rows}}})
	service.commands.addHeaders(request.Header())
	requestCtx, cancel := service.commands.sandbox.unaryContext(ctx)
	defer cancel()
	_, err = service.commands.client.Update(requestCtx, request)
	return connectError(err)
}

// Kill stops the terminal process.
func (service *PTYService) Kill(ctx context.Context, handle *CommandHandle) error {
	return handle.Kill(ctx)
}
