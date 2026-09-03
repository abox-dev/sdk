package agentbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
	filesystem "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/filesystem"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/filesystem/filesystemconnect"
	process "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process/processconnect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type testFilesystemServer struct {
	filesystemconnect.UnimplementedFilesystemHandler
}

type delayedReader struct {
	delay time.Duration
	data  []byte
}

type closeTrackingTransport struct {
	base   http.RoundTripper
	closed atomic.Int64
}

func (transport *closeTrackingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := transport.base.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	switch request.URL.Path {
	case processconnect.ProcessStartProcedure, processconnect.ProcessConnectProcedure, filesystemconnect.FilesystemWatchDirProcedure:
		response.Body = &closeTrackingBody{ReadCloser: response.Body, closedCount: &transport.closed}
	}
	return response, nil
}

type closeTrackingBody struct {
	io.ReadCloser
	closed      atomic.Bool
	closedCount *atomic.Int64
}

func (body *closeTrackingBody) Close() error {
	if body.closed.CompareAndSwap(false, true) {
		body.closedCount.Add(1)
	}
	return body.ReadCloser.Close()
}

func (reader *delayedReader) Read(buffer []byte) (int, error) {
	if reader.delay > 0 {
		time.Sleep(reader.delay)
		reader.delay = 0
	}
	if len(reader.data) == 0 {
		return 0, io.EOF
	}
	count := copy(buffer, reader.data)
	reader.data = reader.data[count:]
	return count, nil
}

func testEntry(path string) *filesystem.EntryInfo {
	return &filesystem.EntryInfo{Name: "file.txt", Path: path, Type: filesystem.FileType_FILE_TYPE_FILE, Size: 5, Mode: 0o644, Permissions: "rw-r--r--", Owner: "user", Group: "user", ModifiedTime: timestamppb.Now(), Metadata: map[string]string{"kind": "test"}}
}
func (testFilesystemServer) Stat(ctx context.Context, request *connect.Request[filesystem.StatRequest]) (*connect.Response[filesystem.StatResponse], error) {
	if request.Msg.Path == "missing" {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("missing"))
	}
	if request.Msg.Path == "slow" {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return connect.NewResponse(&filesystem.StatResponse{Entry: testEntry(request.Msg.Path)}), nil
}
func (testFilesystemServer) MakeDir(_ context.Context, request *connect.Request[filesystem.MakeDirRequest]) (*connect.Response[filesystem.MakeDirResponse], error) {
	entry := testEntry(request.Msg.Path)
	entry.Type = filesystem.FileType_FILE_TYPE_DIRECTORY
	return connect.NewResponse(&filesystem.MakeDirResponse{Entry: entry}), nil
}
func (testFilesystemServer) Move(_ context.Context, request *connect.Request[filesystem.MoveRequest]) (*connect.Response[filesystem.MoveResponse], error) {
	return connect.NewResponse(&filesystem.MoveResponse{Entry: testEntry(request.Msg.Destination)}), nil
}
func (testFilesystemServer) ListDir(context.Context, *connect.Request[filesystem.ListDirRequest]) (*connect.Response[filesystem.ListDirResponse], error) {
	return connect.NewResponse(&filesystem.ListDirResponse{Entries: []*filesystem.EntryInfo{testEntry("/file.txt")}}), nil
}
func (testFilesystemServer) Remove(context.Context, *connect.Request[filesystem.RemoveRequest]) (*connect.Response[filesystem.RemoveResponse], error) {
	return connect.NewResponse(&filesystem.RemoveResponse{}), nil
}
func (testFilesystemServer) WatchDir(ctx context.Context, request *connect.Request[filesystem.WatchDirRequest], stream *connect.ServerStream[filesystem.WatchDirResponse]) error {
	count := 1
	if request.Msg.Path == "/flood" {
		count = 40
	}
	for range count {
		if err := stream.Send(&filesystem.WatchDirResponse{Event: &filesystem.WatchDirResponse_Filesystem{Filesystem: &filesystem.FilesystemEvent{Name: "file.txt", Type: filesystem.EventType_EVENT_TYPE_WRITE, Entry: testEntry("/file.txt")}}}); err != nil {
			return err
		}
	}
	if request.Msg.Path == "/flood" {
		<-ctx.Done()
		return ctx.Err()
	}
	return nil
}

type testProcessServer struct {
	processconnect.UnimplementedProcessHandler
}

func (testProcessServer) List(context.Context, *connect.Request[process.ListRequest]) (*connect.Response[process.ListResponse], error) {
	cwd := "/app"
	tag := "tag"
	return connect.NewResponse(&process.ListResponse{Processes: []*process.ProcessInfo{{Pid: 7, Tag: &tag, Config: &process.ProcessConfig{Cmd: "echo", Args: []string{"ok"}, Envs: map[string]string{"A": "B"}, Cwd: &cwd}}}}), nil
}
func (testProcessServer) Start(_ context.Context, request *connect.Request[process.StartRequest], stream *connect.ServerStream[process.StartResponse]) error {
	if err := stream.Send(&process.StartResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_Start{Start: &process.ProcessEvent_StartEvent{Pid: 7}}}}); err != nil {
		return err
	}
	output := &process.ProcessEvent_DataEvent{Output: &process.ProcessEvent_DataEvent_Stdout{Stdout: []byte("ok\n")}}
	if request.Msg.Pty != nil {
		output.Output = &process.ProcessEvent_DataEvent_Pty{Pty: []byte("pty\n")}
	} else if request.Msg.Process.GetCmd() == "stderr" {
		output.Output = &process.ProcessEvent_DataEvent_Stderr{Stderr: []byte("bad\n")}
	}
	count := 1
	if request.Msg.Process.GetCmd() == "flood" {
		count = 40
	}
	for range count {
		if err := stream.Send(&process.StartResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_Data{Data: output}}}); err != nil {
			return err
		}
	}
	exitCode := int32(0)
	message := ""
	if request.Msg.Process.GetCmd() == "false" {
		exitCode, message = 2, "failed"
	}
	return stream.Send(&process.StartResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_End{End: &process.ProcessEvent_EndEvent{Exited: true, ExitCode: exitCode, Status: "exited", Error: &message}}}})
}
func (testProcessServer) Connect(_ context.Context, _ *connect.Request[process.ConnectRequest], stream *connect.ServerStream[process.ConnectResponse]) error {
	if err := stream.Send(&process.ConnectResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_Data{Data: &process.ProcessEvent_DataEvent{Output: &process.ProcessEvent_DataEvent_Stderr{Stderr: []byte("connected")}}}}}); err != nil {
		return err
	}
	return stream.Send(&process.ConnectResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_End{End: &process.ProcessEvent_EndEvent{Exited: true}}}})
}
func (testProcessServer) SendInput(context.Context, *connect.Request[process.SendInputRequest]) (*connect.Response[process.SendInputResponse], error) {
	return connect.NewResponse(&process.SendInputResponse{}), nil
}
func (testProcessServer) SendSignal(context.Context, *connect.Request[process.SendSignalRequest]) (*connect.Response[process.SendSignalResponse], error) {
	return connect.NewResponse(&process.SendSignalResponse{}), nil
}
func (testProcessServer) CloseStdin(context.Context, *connect.Request[process.CloseStdinRequest]) (*connect.Response[process.CloseStdinResponse], error) {
	return connect.NewResponse(&process.CloseStdinResponse{}), nil
}
func (testProcessServer) Update(context.Context, *connect.Request[process.UpdateRequest]) (*connect.Response[process.UpdateResponse], error) {
	return connect.NewResponse(&process.UpdateResponse{}), nil
}

func newEnvdTestSandbox(t *testing.T, options ...ClientOption) (*Sandbox, func()) {
	t.Helper()
	mux := http.NewServeMux()
	fsPath, fsHandler := filesystemconnect.NewFilesystemHandler(testFilesystemServer{}, connect.WithCodec(tolerantJSONCodec{}))
	mux.Handle(fsPath, fsHandler)
	processPath, processHandler := processconnect.NewProcessHandler(testProcessServer{}, connect.WithCodec(tolerantJSONCodec{}))
	mux.Handle(processPath, processHandler)
	mux.HandleFunc("/health", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/headers", func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("X-Test") != "value" || request.Header.Get("Content-Type") != "text/plain" {
			http.Error(writer, "missing custom headers", http.StatusBadRequest)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("/files", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodGet {
			if request.URL.Query().Get("path") == "missing-http" {
				http.NotFound(writer, request)
				return
			}
			if request.URL.Query().Get("path") == "old-user" && request.URL.Query().Get("username") != "user" {
				http.Error(writer, "missing legacy user", http.StatusBadRequest)
				return
			}
			writer.Write([]byte("hello"))
			return
		}
		if err := request.ParseMultipartForm(1 << 20); err != nil {
			http.Error(writer, err.Error(), 400)
			return
		}
		file, _, err := request.FormFile("file")
		if err != nil {
			http.Error(writer, err.Error(), 400)
			return
		}
		data, _ := io.ReadAll(file)
		file.Close()
		if string(data) != "hello" {
			http.Error(writer, "bad file", 400)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprint(writer, `[{"name":"file.txt","type":"file","path":"/file.txt","metadata":{"kind":"test"}}]`)
	})
	server := httptest.NewServer(mux)
	clientOptions := []ClientOption{WithAPIURL(server.URL), WithSandboxURL(server.URL)}
	clientOptions = append(clientOptions, options...)
	client, err := NewClient(clientOptions...)
	if err != nil {
		server.Close()
		t.Fatal(err)
	}
	sandbox := client.sandboxFromAPI(api.Sandbox{SandboxID: "sbx", TemplateID: "base", EnvdVersion: "1", EnvdAccessToken: ptr("token")})
	return sandbox, server.Close
}

func TestFilesystem(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	ctx := context.Background()
	if text, err := sandbox.Files.ReadText(ctx, "/file.txt", ""); err != nil || text != "hello" {
		t.Fatalf("read: %q %v", text, err)
	}
	if _, err := sandbox.Files.Read(ctx, "", ""); err == nil {
		t.Fatal("expected empty path validation")
	}
	var output strings.Builder
	if count, err := sandbox.Files.ReadTo(ctx, "/file.txt", "", &output); err != nil || count != 5 {
		t.Fatalf("read to: %d %v", count, err)
	}
	if entry, err := sandbox.Files.WriteText(ctx, "/file.txt", "hello", &WriteFileOptions{Metadata: map[string]string{"kind": "test"}}); err != nil || entry.Path != "/file.txt" {
		t.Fatalf("write: %#v %v", entry, err)
	}
	if _, err := sandbox.Files.WriteBytes(ctx, "/file.txt", []byte("hello"), nil); err != nil {
		t.Fatal(err)
	}
	if entries, err := sandbox.Files.WriteBatch(ctx, []WriteFile{{Path: "/file.txt", Data: strings.NewReader("hello")}}, ""); err != nil || len(entries) != 1 {
		t.Fatalf("batch: %v", err)
	}
	if entry, err := sandbox.Files.Stat(ctx, "/file.txt"); err != nil || entry.Metadata["kind"] != "test" {
		t.Fatalf("stat: %#v %v", entry, err)
	}
	if exists, err := sandbox.Files.Exists(ctx, "missing"); err != nil || exists {
		t.Fatalf("exists: %v %v", exists, err)
	}
	if exists, err := sandbox.Files.Exists(ctx, "/file.txt"); err != nil || !exists {
		t.Fatalf("existing file: %v %v", exists, err)
	}
	if entries, err := sandbox.Files.WriteBatch(ctx, []WriteFile{{Path: "", Data: nil}}, ""); err == nil || len(entries) != 0 {
		t.Fatal("expected batch failure")
	}
	if entries, err := sandbox.Files.List(ctx, "/", 1); err != nil || len(entries) != 1 {
		t.Fatalf("list: %v", err)
	}
	if _, err := sandbox.Files.MakeDir(ctx, "/dir"); err != nil {
		t.Fatal(err)
	}
	if _, err := sandbox.Files.Rename(ctx, "/file.txt", "/new.txt"); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.Files.Remove(ctx, "/new.txt"); err != nil {
		t.Fatal(err)
	}
	watcher, err := sandbox.Files.Watch(ctx, "/", &WatchOptions{IncludeEntry: true})
	if err != nil {
		t.Fatal(err)
	}
	event := <-watcher.Events
	if event.Type != "write" || event.Entry == nil {
		t.Fatalf("event: %#v", event)
	}
	if err := watcher.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := sandbox.Files.WriteText(ctx, "/x", "x", &WriteFileOptions{Metadata: map[string]string{"bad key": "x"}}); err == nil {
		t.Fatal("expected metadata validation")
	}
	if _, err := sandbox.Files.ReadText(ctx, "missing-http", ""); err == nil {
		t.Fatal("expected HTTP file error")
	} else {
		var missing *FileNotFoundError
		if !errors.As(err, &missing) {
			t.Fatalf("expected FileNotFoundError, got %T", err)
		}
	}
	response, err := sandbox.Request(ctx, envdPort, http.MethodGet, "/files?path=/file.txt", nil, false)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	response, err = sandbox.RequestWithOptions(ctx, envdPort, http.MethodPost, "/headers", strings.NewReader("body"), &SandboxRequestOptions{Headers: http.Header{"X-Test": {"value"}}, ContentType: "text/plain"})
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if running, err := sandbox.IsRunning(ctx); err != nil || !running {
		t.Fatalf("is running: %v %v", running, err)
	}
}

func TestFileWriteStreamingTimeout(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	sandbox.client.config.requestTimeout = 10 * time.Millisecond

	entry, err := sandbox.Files.Write(t.Context(), "/file.txt", &delayedReader{delay: 30 * time.Millisecond, data: []byte("hello")}, nil)
	if err != nil || entry.Path != "/file.txt" {
		t.Fatalf("streaming upload inherited unary timeout: %#v %v", entry, err)
	}

	_, err = sandbox.Files.Write(t.Context(), "/file.txt", &delayedReader{delay: 30 * time.Millisecond, data: []byte("hello")}, &WriteFileOptions{RequestTimeout: 10 * time.Millisecond})
	var timeout *TimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("upload timeout error = %T %v", err, err)
	}
	if _, err := sandbox.Files.WriteText(t.Context(), "/file.txt", "hello", &WriteFileOptions{RequestTimeout: -time.Second}); err == nil {
		t.Fatal("expected negative upload timeout validation")
	}
}

func TestEnvdVersionCompatibility(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	sandbox.EnvdVersion = "0.3.9"
	if _, err := sandbox.Files.ReadText(t.Context(), "old-user", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := sandbox.Files.WriteText(t.Context(), "/file.txt", "hello", &WriteFileOptions{Metadata: map[string]string{"kind": "test"}}); err == nil {
		t.Fatal("expected metadata version gate")
	}
	if _, err := sandbox.Files.Watch(t.Context(), "/", &WatchOptions{IncludeEntry: true}); err == nil {
		t.Fatal("expected watch entry version gate")
	} else {
		var sandboxError *SandboxError
		if !errors.As(err, &sandboxError) {
			t.Fatalf("watch version error = %T, want SandboxError", err)
		}
	}
	sandbox.EnvdVersion = "0.1.3"
	if _, err := sandbox.Files.Watch(t.Context(), "/", &WatchOptions{Recursive: true}); err == nil {
		t.Fatal("expected recursive watch version gate")
	}
	sandbox.EnvdVersion = "0.6.3"
	if _, err := sandbox.Files.Watch(t.Context(), "/", &WatchOptions{AllowNetworkMounts: true}); err == nil {
		t.Fatal("expected network mount watch version gate")
	}
	if envdAtLeast("0.6.1", 0, 6, 2) || !envdAtLeast("v0.6.2-beta", 0, 6, 2) || !envdAtLeast("invalid", 9, 0, 0) || !envdAtLeast("1.2.3.4", 9, 0, 0) {
		t.Fatal("unexpected envd version comparison")
	}
	response, err := sandbox.RequestWithOptions(t.Context(), envdPort, http.MethodGet, "/health", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
}

func TestUnaryEnvdRequestTimeout(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	sandbox.client.config.requestTimeout = 10 * time.Millisecond
	_, err := sandbox.Files.Stat(t.Context(), "slow")
	var timeout *TimeoutError
	if !errors.As(err, &timeout) {
		t.Fatalf("expected TimeoutError, got %T: %v", err, err)
	}
}

func TestCommandsAndPTY(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	ctx := context.Background()
	stdoutCalled := false
	result, err := sandbox.Commands.Run(ctx, "echo", &CommandOptions{Args: []string{"ok"}, OnStdout: func([]byte) { stdoutCalled = true }})
	if err != nil || string(result.Stdout) != "ok\n" || !stdoutCalled {
		t.Fatalf("run: %#v %v", result, err)
	}
	if _, err := sandbox.Commands.Run(ctx, "echo", nil); err != nil {
		t.Fatal(err)
	}
	stderrCalled := false
	if result, err := sandbox.Commands.Run(ctx, "stderr", &CommandOptions{OnStderr: func([]byte) { stderrCalled = true }}); err != nil || string(result.Stderr) != "bad\n" || !stderrCalled {
		t.Fatalf("stderr: %#v %v", result, err)
	}
	if _, err := sandbox.Commands.Run(ctx, "false", nil); err == nil {
		t.Fatal("expected command exit error")
	}
	processes, err := sandbox.Commands.List(ctx)
	if err != nil || len(processes) != 1 || processes[0].Cwd != "/app" {
		t.Fatalf("list: %#v %v", processes, err)
	}
	handle, err := sandbox.Commands.Connect(ctx, 7, "")
	if err != nil {
		t.Fatal(err)
	}
	defaultPTY, err := sandbox.PTY.Create(ctx, "sh", nil)
	if err != nil {
		t.Fatal(err)
	}
	for range defaultPTY.PTY {
	}
	if result, err := defaultPTY.Wait(ctx); err != nil || len(result.Stdout) != 0 {
		t.Fatal(err)
	}
	if _, err := handle.Write(ctx, []byte("input")); err != nil {
		t.Fatal(err)
	}
	if err := handle.CloseStdin(ctx); err != nil {
		t.Fatal(err)
	}
	if err := handle.Kill(ctx); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.Commands.Terminate(ctx, 7, ""); err != nil {
		t.Fatal(err)
	}
	for range handle.Stderr {
	}
	if _, err := handle.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	var ptyOutput bytes.Buffer
	pty, err := sandbox.PTY.Create(ctx, "sh", &PTYOptions{Cols: 100, Rows: 40, OnPTY: func(chunk []byte) { ptyOutput.Write(chunk) }})
	if err != nil {
		t.Fatal(err)
	}
	if err := sandbox.PTY.Input(ctx, pty, []byte("x")); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.PTY.Resize(ctx, pty, 80, 24); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.PTY.Kill(ctx, pty); err != nil {
		t.Fatal(err)
	}
	for range pty.PTY {
	}
	if _, err := pty.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	if ptyOutput.String() != "pty\n" {
		t.Fatalf("PTY callback output = %q", ptyOutput.String())
	}
	if _, err := processSelector(1, "tag"); err == nil {
		t.Fatal("expected selector error")
	}
	attached, err := sandbox.PTY.Connect(ctx, 7, "")
	if err != nil {
		t.Fatal(err)
	}
	for range attached.Stderr {
	}
	if _, err := attached.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	withoutStart, err := sandbox.Commands.Connect(ctx, 0, "tag")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := withoutStart.PID(ctx); err == nil {
		t.Fatal("expected PID error for stream ending without a start event")
	}
}

func TestCommandWaitDoesNotRequireDrainingOutput(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	handle, err := sandbox.Commands.Start(t.Context(), "flood", nil)
	if err != nil {
		t.Fatal(err)
	}
	waitCtx, cancel := context.WithTimeout(t.Context(), time.Second)
	defer cancel()
	result, err := handle.Wait(waitCtx)
	if err != nil {
		t.Fatal(err)
	}
	if expected := bytes.Repeat([]byte("ok\n"), 40); !bytes.Equal(result.Stdout, expected) {
		t.Fatalf("stdout length = %d, want %d", len(result.Stdout), len(expected))
	}
}

func TestCommandWaitPreservesOutputChannelTail(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	handle, err := sandbox.Commands.Start(t.Context(), "flood", nil)
	if err != nil {
		t.Fatal(err)
	}
	read := make(chan []byte, 1)
	go func() {
		var output bytes.Buffer
		for chunk := range handle.Stdout {
			output.Write(chunk)
			time.Sleep(2 * time.Millisecond)
		}
		read <- output.Bytes()
	}()

	result, err := handle.Wait(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	expected := bytes.Repeat([]byte("ok\n"), 40)
	if !bytes.Equal(result.Stdout, expected) {
		t.Fatalf("result stdout length = %d, want %d", len(result.Stdout), len(expected))
	}
	select {
	case output := <-read:
		if !bytes.Equal(output, expected) {
			t.Fatalf("channel stdout length = %d, want %d", len(output), len(expected))
		}
	case <-time.After(time.Second):
		t.Fatal("stdout channel did not finish draining")
	}
}

func TestStreamingResponsesAreClosed(t *testing.T) {
	transport := &closeTrackingTransport{base: http.DefaultTransport.(*http.Transport).Clone()}
	sandbox, closeServer := newEnvdTestSandbox(t, WithHTTPClient(&http.Client{Transport: transport}))
	defer closeServer()

	const commandRuns = 30
	for range commandRuns {
		if _, err := sandbox.Commands.Run(t.Context(), "echo", nil); err != nil {
			t.Fatal(err)
		}
	}

	connected, err := sandbox.Commands.Connect(t.Context(), 7, "")
	if err != nil {
		t.Fatal(err)
	}
	for range connected.Stderr {
	}
	if _, err := connected.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}

	terminal, err := sandbox.PTY.Create(t.Context(), "sh", nil)
	if err != nil {
		t.Fatal(err)
	}
	for range terminal.PTY {
	}
	if _, err := terminal.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}

	watcher, err := sandbox.Files.Watch(t.Context(), "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if event := <-watcher.Events; event.Type != "write" {
		t.Fatalf("watch event = %#v", event)
	}
	if err := watcher.Close(); err != nil {
		t.Fatal(err)
	}

	if got, want := transport.closed.Load(), int64(commandRuns+3); got != want {
		t.Fatalf("closed streaming response bodies = %d, want %d", got, want)
	}
}

func TestWatchCloseDoesNotRequireDrainingEvents(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	watcher, err := sandbox.Files.Watch(t.Context(), "/flood", nil)
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for len(watcher.Events) < cap(watcher.Events) && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	closed := make(chan error, 1)
	go func() { closed <- watcher.Close() }()
	select {
	case err := <-closed:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("watcher Close blocked with an unread full event channel")
	}
}

func ptr[T any](value T) *T { return &value }

func TestCommandHelpersAndFileMappings(t *testing.T) {
	if (&CommandExitError{Result: CommandResult{ExitCode: 2}, Message: "bad"}).Error() == "" || (&CommandExitError{Result: CommandResult{ExitCode: 1}}).Error() == "" {
		t.Fatal("empty exit error")
	}
	if _, err := processSelector(0, ""); err == nil {
		t.Fatal("expected selector validation")
	}
	if selector, err := processSelector(0, "tag"); err != nil || selector.GetTag() != "tag" {
		t.Fatalf("tag selector: %v", err)
	}
	for value, expected := range map[filesystem.EventType]string{filesystem.EventType_EVENT_TYPE_CREATE: "create", filesystem.EventType_EVENT_TYPE_WRITE: "write", filesystem.EventType_EVENT_TYPE_REMOVE: "remove", filesystem.EventType_EVENT_TYPE_RENAME: "rename", filesystem.EventType_EVENT_TYPE_CHMOD: "chmod", filesystem.EventType_EVENT_TYPE_UNSPECIFIED: "unknown"} {
		if mapEventType(value) != expected {
			t.Fatalf("event %v", value)
		}
	}
	if mapEntry(nil) != nil {
		t.Fatal("nil entry changed")
	}
	for value, expected := range map[filesystem.FileType]FileType{filesystem.FileType_FILE_TYPE_FILE: FileTypeFile, filesystem.FileType_FILE_TYPE_DIRECTORY: FileTypeDirectory, filesystem.FileType_FILE_TYPE_SYMLINK: FileTypeSymlink} {
		entry := testEntry("/x")
		entry.Type = value
		if mapEntry(entry).Type != expected {
			t.Fatalf("type %v", value)
		}
	}
	entryWithoutTime := testEntry("/x")
	entryWithoutTime.ModifiedTime = nil
	if !mapEntry(entryWithoutTime).ModifiedAt.IsZero() {
		t.Fatal("unexpected modification time")
	}
	header := make(http.Header)
	if err := addMetadataHeaders(header, map[string]string{"ok": "bad\n"}); err == nil {
		t.Fatal("expected bad metadata value")
	}
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	ctx := context.Background()
	if _, err := sandbox.Commands.Start(ctx, "", nil); err == nil {
		t.Fatal("expected empty command")
	}
	if _, err := sandbox.PTY.Create(ctx, "", nil); err == nil {
		t.Fatal("expected empty PTY command")
	}
	handle := newCommandHandle(sandbox.Commands, "")
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := handle.PID(canceled); err == nil {
		t.Fatal("expected PID cancellation")
	}
	if _, err := handle.Wait(canceled); err == nil {
		t.Fatal("expected wait cancellation")
	}
	if err := sandbox.PTY.Resize(ctx, handle, 0, 1); err == nil {
		t.Fatal("expected invalid PTY size")
	}
	if _, err := sandbox.Files.Write(ctx, "", nil, nil); err == nil {
		t.Fatal("expected write validation")
	}
	if _, err := sandbox.Request(ctx, envdPort, "bad\nmethod", "/", nil, false); err == nil {
		t.Fatal("expected request validation")
	}
}
