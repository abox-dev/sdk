package agentbox

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
	process "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process/processconnect"
)

type streamingProcessServer struct {
	testProcessServer
	chunks   int
	size     int
	failure  bool
	pty      bool
	exitCode int32
	signals  atomic.Int64
	eof      atomic.Int64
}

func (server *streamingProcessServer) events(send func(*process.ProcessEvent) error) error {
	if err := send(&process.ProcessEvent{Event: &process.ProcessEvent_Start{Start: &process.ProcessEvent_StartEvent{Pid: 42}}}); err != nil {
		return err
	}
	for i := range server.chunks {
		chunk := bytes.Repeat([]byte{byte(i % 251)}, server.size)
		outputs := []*process.ProcessEvent_DataEvent{
			{Output: &process.ProcessEvent_DataEvent_Stdout{Stdout: chunk}},
			{Output: &process.ProcessEvent_DataEvent_Stderr{Stderr: chunk}},
		}
		if server.pty {
			outputs = []*process.ProcessEvent_DataEvent{{Output: &process.ProcessEvent_DataEvent_Pty{Pty: chunk}}}
		}
		for _, data := range outputs {
			if err := send(&process.ProcessEvent{Event: &process.ProcessEvent_Data{Data: data}}); err != nil {
				return err
			}
		}
	}
	if server.failure {
		return connect.NewError(connect.CodeInternal, errors.New("secret output"))
	}
	message := "secret output"
	return send(&process.ProcessEvent{Event: &process.ProcessEvent_End{End: &process.ProcessEvent_EndEvent{Exited: true, ExitCode: server.exitCode, Status: "exited", Error: &message}}})
}
func (server *streamingProcessServer) Start(_ context.Context, _ *connect.Request[process.StartRequest], stream *connect.ServerStream[process.StartResponse]) error {
	return server.events(func(event *process.ProcessEvent) error { return stream.Send(&process.StartResponse{Event: event}) })
}
func (server *streamingProcessServer) Connect(_ context.Context, _ *connect.Request[process.ConnectRequest], stream *connect.ServerStream[process.ConnectResponse]) error {
	return server.events(func(event *process.ProcessEvent) error { return stream.Send(&process.ConnectResponse{Event: event}) })
}
func (server *streamingProcessServer) SendSignal(context.Context, *connect.Request[process.SendSignalRequest]) (*connect.Response[process.SendSignalResponse], error) {
	server.signals.Add(1)
	return connect.NewResponse(&process.SendSignalResponse{}), nil
}
func (server *streamingProcessServer) CloseStdin(context.Context, *connect.Request[process.CloseStdinRequest]) (*connect.Response[process.CloseStdinResponse], error) {
	server.eof.Add(1)
	return connect.NewResponse(&process.CloseStdinResponse{}), nil
}
func newStreamingTestSandbox(t *testing.T, fixture *streamingProcessServer) (*Sandbox, *closeTrackingTransport) {
	t.Helper()
	path, handler := processconnect.NewProcessHandler(fixture, connect.WithCodec(tolerantJSONCodec{}))
	mux := http.NewServeMux()
	mux.Handle(path, handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	transport := &closeTrackingTransport{base: http.DefaultTransport.(*http.Transport).Clone()}
	t.Cleanup(func() { transport.base.(*http.Transport).CloseIdleConnections() })
	client, err := NewClient(WithAPIURL(server.URL), WithSandboxURL(server.URL), WithHTTPClient(&http.Client{Transport: transport}))
	if err != nil {
		t.Fatal(err)
	}
	return client.sandboxFromAPI(api.Sandbox{SandboxID: "sbx", TemplateID: "base"}), transport
}
func waitStreamingTest(t *testing.T, handle *CommandHandle) (CommandResult, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	return handle.Wait(ctx)
}
func assertNoRetainedOutput(t *testing.T, handle *CommandHandle) {
	t.Helper()
	<-handle.Done
	if handle.result.Stdout != nil || handle.result.Stderr != nil {
		t.Fatal("streaming mode captured output")
	}
	for _, stream := range []*outputStream{handle.stdout, handle.stderr, handle.pty} {
		stream.mu.Lock()
		retained := len(stream.queue)
		stream.mu.Unlock()
		if retained != 0 || !stream.direct || cap(stream.output) != 0 {
			t.Fatalf("hidden output queue: %d", retained)
		}
	}
}

func TestCommandStreamingCallbacks(t *testing.T) {
	fixture := &streamingProcessServer{chunks: 2048, size: 1024}
	sandbox, transport := newStreamingTestSandbox(t, fixture)
	for _, attach := range []string{"start", "pid", "tag"} {
		t.Run(attach, func(t *testing.T) {
			var counts [2]int
			callback := func(index int) func([]byte) {
				return func(chunk []byte) {
					want := bytes.Repeat([]byte{byte(counts[index] % 251)}, fixture.size)
					if !bytes.Equal(chunk, want) {
						t.Errorf("unexpected callback chunk %d", counts[index])
					}
					counts[index]++
				}
			}
			policy := &CommandStreamingOptions{}
			var handle *CommandHandle
			var err error
			switch attach {
			case "start":
				handle, err = sandbox.Commands.Start(t.Context(), "flood", &CommandOptions{Streaming: policy, OnStdout: callback(0), OnStderr: callback(1)})
			case "pid":
				handle, err = sandbox.Commands.ConnectWithOptions(t.Context(), 42, "", &CommandConnectOptions{Streaming: policy, OnStdout: callback(0), OnStderr: callback(1)})
			case "tag":
				handle, err = sandbox.Commands.ConnectWithOptions(t.Context(), 0, "rpc", &CommandConnectOptions{Streaming: policy, OnStdout: callback(0), OnStderr: callback(1)})
			}
			if err != nil {
				t.Fatal(err)
			}
			defer handle.Close()
			result, err := waitStreamingTest(t, handle)
			if err != nil || result.PID != 42 || result.Status != "exited" {
				t.Fatalf("result: %+v %v", result, err)
			}
			if counts != [2]int{fixture.chunks, fixture.chunks} {
				t.Fatalf("callbacks: %v", counts)
			}
			assertNoRetainedOutput(t, handle)
			for _, channel := range []<-chan []byte{handle.Stdout, handle.Stderr, handle.PTY} {
				if _, ok := <-channel; ok {
					t.Fatal("unused channel is open")
				}
			}
		})
	}
	if transport.closed.Load() != 3 {
		t.Fatal("response bodies not released")
	}
}

func TestCommandStreamingChannels(t *testing.T) {
	fixture := &streamingProcessServer{chunks: 32, size: 100001}
	sandbox, _ := newStreamingTestSandbox(t, fixture)
	for _, stderr := range []bool{false, true} {
		t.Run(fmt.Sprintf("stderr=%t", stderr), func(t *testing.T) {
			handle, err := sandbox.Commands.Start(t.Context(), "flood", &CommandOptions{Streaming: &CommandStreamingOptions{StdoutChannel: true, StderrChannel: stderr}})
			if err != nil {
				t.Fatal(err)
			}
			defer handle.Close()
			// One receiver drains both enabled channels without assuming event boundaries.
			counts := [2]int{}
			stdout, stderrChannel := handle.Stdout, handle.Stderr
			for stdout != nil || stderrChannel != nil {
				var chunk []byte
				var ok bool
				index := 0
				select {
				case chunk, ok = <-stdout:
					if !ok {
						stdout = nil
						continue
					}
				case chunk, ok = <-stderrChannel:
					index = 1
					if !ok {
						stderrChannel = nil
						continue
					}
				case <-time.After(5 * time.Second):
					t.Fatal("channel stalled")
				}
				if len(chunk) > commandStreamChunkBytes {
					t.Fatal("unbounded channel slice")
				}
				for _, value := range chunk {
					if value != byte((counts[index]/fixture.size)%251) {
						t.Fatal("bytes lost or reordered")
					}
					counts[index]++
				}
				// Retaining and mutating a delivered slice must not affect subsequent data.
				clear(chunk)
			}
			result, err := waitStreamingTest(t, handle)
			if err != nil || result.PID != 42 {
				t.Fatalf("wait: %+v %v", result, err)
			}
			want := [2]int{fixture.chunks * fixture.size, 0}
			if stderr {
				want[1] = want[0]
			}
			if counts != want {
				t.Fatalf("bytes: %v want %v", counts, want)
			}
			assertNoRetainedOutput(t, handle)
		})
	}
}

func TestCommandStreamingBlockedDeliveryCancellation(t *testing.T) {
	fixture := &streamingProcessServer{chunks: 100, size: 1024}
	sandbox, transport := newStreamingTestSandbox(t, fixture)
	for i := range 30 {
		ctx, cancel := context.WithCancel(t.Context())
		handle, err := sandbox.Commands.ConnectWithOptions(ctx, 42, "", &CommandConnectOptions{Streaming: &CommandStreamingOptions{StdoutChannel: true}})
		if err != nil {
			cancel()
			t.Fatal(err)
		}
		// Receiving one chunk ensures delivery has started; the next unread chunk
		// must prevent the end event from being observed even if the server has exited.
		<-handle.Stdout
		select {
		case <-handle.Done:
			t.Fatal("unread channel silently discarded output")
		default:
		}
		if i%2 == 0 {
			cancel()
		} else {
			handle.Close()
		}
		result, err := waitStreamingTest(t, handle)
		cancel()
		if !errors.Is(err, context.Canceled) || result.PID != 42 {
			t.Fatalf("cancel: %+v %v", result, err)
		}
		assertNoRetainedOutput(t, handle)
		if _, ok := <-handle.Stdout; ok {
			t.Fatal("channel retained data after cancellation")
		}
		handle.Close()
	}
	if transport.closed.Load() != 30 || fixture.signals.Load() != 0 || fixture.eof.Load() != 0 {
		t.Fatal("detach leaked a response or signaled the process")
	}
}

func TestCommandStreamingFailures(t *testing.T) {
	for _, kind := range []string{"exit", "transport", "oversized"} {
		t.Run(kind, func(t *testing.T) {
			fixture := &streamingProcessServer{chunks: 1, size: 1024}
			switch kind {
			case "exit":
				fixture.exitCode = 9
			case "transport":
				fixture.failure = true
			case "oversized":
				fixture.size = commandStreamMessageBytes
			}
			sandbox, transport := newStreamingTestSandbox(t, fixture)
			handle, err := sandbox.Commands.Start(t.Context(), "fail", &CommandOptions{Streaming: &CommandStreamingOptions{}})
			if err != nil {
				t.Fatal(err)
			}
			defer handle.Close()
			result, err := waitStreamingTest(t, handle)
			if err == nil || result.PID != 42 || strings.Contains(fmt.Sprintf("%+v", err), "secret") {
				t.Fatalf("failure: %+v %v", result, err)
			}
			if kind == "exit" {
				var exit *CommandExitError
				if !errors.As(err, &exit) || exit.Result.ExitCode != 9 || exit.Result.Stdout != nil || exit.Result.Stderr != nil || exit.Message != "" {
					t.Fatalf("exit: %#v", err)
				}
			} else {
				want := connect.CodeInternal
				if kind == "oversized" {
					want = connect.CodeResourceExhausted
				}
				if connect.CodeOf(err) != want {
					t.Fatalf("code: %s want %s", connect.CodeOf(err), want)
				}

			}
			assertNoRetainedOutput(t, handle)
			if transport.closed.Load() != 1 {
				t.Fatal("unclosed response")
			}
		})
	}
}

func TestCommandStreamingValidationAndPTY(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	callback := func([]byte) {}
	for _, options := range []*CommandConnectOptions{
		{Streaming: &CommandStreamingOptions{StdoutChannel: true}, OnStdout: callback},
		{Streaming: &CommandStreamingOptions{StderrChannel: true}, OnStderr: callback},
		{Streaming: &CommandStreamingOptions{PTYChannel: true}, OnPTY: callback},
	} {
		if _, err := sandbox.Commands.ConnectWithOptions(t.Context(), 7, "", options); err == nil {
			t.Fatal("accepted fan-out")
		}
	}
	if _, err := sandbox.Commands.Start(t.Context(), "echo", &CommandOptions{Streaming: &CommandStreamingOptions{StdoutChannel: true}, OnStdout: callback}); err == nil {
		t.Fatal("accepted start fan-out")
	}
	if _, err := sandbox.PTY.Create(t.Context(), "sh", &PTYOptions{Streaming: &CommandStreamingOptions{PTYChannel: true}, OnPTY: callback}); err == nil {
		t.Fatal("accepted PTY fan-out")
	}
	for _, useCallback := range []bool{false, true} {
		var output bytes.Buffer
		options := &PTYOptions{Streaming: &CommandStreamingOptions{PTYChannel: !useCallback}}
		if useCallback {
			options.OnPTY = func(chunk []byte) { output.Write(chunk) }
		}
		handle, err := sandbox.PTY.Create(t.Context(), "sh", options)
		if err != nil {
			t.Fatal(err)
		}
		for chunk := range handle.PTY {
			output.Write(chunk)
		}
		if _, err := waitStreamingTest(t, handle); err != nil {
			t.Fatal(err)
		}
		if output.String() != "pty\n" {
			t.Fatalf("PTY: %q", output.String())
		}
		assertNoRetainedOutput(t, handle)
	}
	result, err := sandbox.Commands.Run(t.Context(), "flood", &CommandOptions{Streaming: &CommandStreamingOptions{StdoutChannel: true}})
	if err != nil || result.Stdout != nil || result.PID != 7 {
		t.Fatalf("streaming Run: %+v %v", result, err)
	}
	handle, err := sandbox.Commands.ConnectWithOptions(t.Context(), 7, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err = waitStreamingTest(t, handle)
	if err != nil || result.PID != 7 || string(result.Stderr) != "connected" {
		t.Fatalf("collecting Connect: %+v %v", result, err)
	}
	handle.Close()
}

func TestCommandCancellationDuringCallback(t *testing.T) {
	sandbox, transport := newStreamingTestSandbox(t, &streamingProcessServer{chunks: 2, size: 16})
	entered, release := make(chan struct{}), make(chan struct{})
	defer close(release)
	handle, err := sandbox.Commands.Start(t.Context(), "callback", &CommandOptions{Streaming: &CommandStreamingOptions{}, OnStdout: func([]byte) { close(entered); <-release }})
	if err != nil {
		t.Fatal(err)
	}
	<-entered
	handle.Close()
	deadline := time.Now().Add(5 * time.Second)
	for transport.closed.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if transport.closed.Load() != 1 {
		t.Fatal("callback prevented response closure")
	}
	select {
	case <-handle.Done:
		t.Fatal("callback was abandoned")
	default:
	}
	release <- struct{}{}
	if _, err := waitStreamingTest(t, handle); !errors.Is(err, context.Canceled) {
		t.Fatalf("callback cancellation: %v", err)
	}
	assertNoRetainedOutput(t, handle)
}

func BenchmarkCommandOutputRetention(b *testing.B) {
	for _, streaming := range []bool{false, true} {
		for _, chunks := range []int{1024, 16384} {
			b.Run(fmt.Sprintf("streaming=%t/chunks=%d", streaming, chunks), func(b *testing.B) {
				var policy *CommandStreamingOptions
				if streaming {
					policy = &CommandStreamingOptions{}
				}
				event := &process.ProcessEvent{Event: &process.ProcessEvent_Data{Data: &process.ProcessEvent_DataEvent{Output: &process.ProcessEvent_DataEvent_Stdout{Stdout: make([]byte, 1024)}}}}
				b.ReportAllocs()
				for b.Loop() {
					handle := newCommandHandle(nil, "", policy)
					n := 0
					handle.receive(b.Context(), func() (*process.ProcessEvent, bool) { n++; return event, n <= chunks }, func() error { return nil }, func() error { return nil }, outputCallbacks{stdout: func([]byte) {}})
					retained := cap(handle.result.Stdout)
					handle.stdout.mu.Lock()
					for _, chunk := range handle.stdout.queue {
						retained += cap(chunk)
					}
					handle.stdout.mu.Unlock()
					b.ReportMetric(float64(retained), "retained-B")
					handle.Close()
					runtime.KeepAlive(handle)
				}
			})
		}
	}
}

func TestCommandStreamingLiveRetention(t *testing.T) {
	fixture := &streamingProcessServer{chunks: 4096, size: 1024}
	sandbox, _ := newStreamingTestSandbox(t, fixture)
	halfway, release := make(chan struct{}), make(chan struct{})
	count := 0
	handle, err := sandbox.Commands.Start(t.Context(), "long", &CommandOptions{Streaming: &CommandStreamingOptions{}, OnStdout: func([]byte) {
		count++
		if count == 2048 {
			close(halfway)
			<-release
		}
	}})
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	// The callback barrier keeps the handle live and stops all receiver writes.
	<-halfway
	if handle.result.Stdout != nil || handle.result.Stderr != nil {
		t.Error("live result accumulated output")
	}
	for _, stream := range []*outputStream{handle.stdout, handle.stderr, handle.pty} {
		stream.mu.Lock()
		if len(stream.queue) != 0 {
			t.Error("live callback accumulated queued output")
		}
		stream.mu.Unlock()
	}
	close(release)
	if _, err := waitStreamingTest(t, handle); err != nil {
		t.Fatal(err)
	}
	if count != fixture.chunks {
		t.Fatalf("chunks: %d", count)
	}
}

func TestCommandCollectingCloseReleasesTail(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	handle, err := sandbox.Commands.Start(t.Context(), "flood", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := waitStreamingTest(t, handle); err != nil {
		t.Fatal(err)
	}
	handle.Close()
	for _, stream := range []*outputStream{handle.stdout, handle.stderr, handle.pty} {
		stream.mu.Lock()
		retained := len(stream.queue)
		stream.mu.Unlock()
		if retained != 0 {
			t.Fatal("Close retained collecting tail")
		}
		// A slice already in flight may win its send concurrently with abort.
		for range stream.output {
		}
	}
}

func TestPTYStreamingDetach(t *testing.T) {
	sandbox, closeServer := newEnvdTestSandbox(t)
	defer closeServer()
	for _, attach := range []bool{false, true} {
		var handle *CommandHandle
		var err error
		if attach {
			handle, err = sandbox.PTY.ConnectWithOptions(t.Context(), 7, "", &CommandConnectOptions{Streaming: &CommandStreamingOptions{StderrChannel: true}})
		} else {
			handle, err = sandbox.PTY.Create(t.Context(), "sh", &PTYOptions{Streaming: &CommandStreamingOptions{PTYChannel: true}})
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := handle.PID(t.Context()); err != nil {
			t.Fatal(err)
		}
		handle.Close()
		if _, err := waitStreamingTest(t, handle); !errors.Is(err, context.Canceled) {
			t.Fatalf("PTY detach: %v", err)
		}
		assertNoRetainedOutput(t, handle)
	}
}

func TestPTYStreamingConnectOutput(t *testing.T) {
	fixture := &streamingProcessServer{chunks: 40, size: 128, pty: true}
	sandbox, _ := newStreamingTestSandbox(t, fixture)
	for _, callback := range []bool{false, true} {
		var got bytes.Buffer
		options := &CommandConnectOptions{Streaming: &CommandStreamingOptions{PTYChannel: !callback}}
		if callback {
			options.OnPTY = func(chunk []byte) { got.Write(chunk) }
		}
		handle, err := sandbox.PTY.ConnectWithOptions(t.Context(), 42, "", options)
		if err != nil {
			t.Fatal(err)
		}
		for chunk := range handle.PTY {
			got.Write(chunk)
		}
		if _, err := waitStreamingTest(t, handle); err != nil {
			t.Fatal(err)
		}
		var want []byte
		for i := range fixture.chunks {
			want = append(want, bytes.Repeat([]byte{byte(i)}, fixture.size)...)
		}
		if !bytes.Equal(got.Bytes(), want) {
			t.Fatal("PTY attach lost output")
		}
		assertNoRetainedOutput(t, handle)
		handle.Close()
	}
}
