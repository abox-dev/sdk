package agentbox

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"connectrpc.com/connect"
	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/filesystem/filesystemconnect"
	process "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process"
	"github.com/abox-dev/sdk/packages/go-sdk/internal/gen/envd/process/processconnect"
)

type userProcessServer struct{ testProcessServer }

func (userProcessServer) Start(_ context.Context, request *connect.Request[process.StartRequest], stream *connect.ServerStream[process.StartResponse]) error {
	user, _, _ := (&http.Request{Header: request.Header()}).BasicAuth()
	if err := stream.Send(&process.StartResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_Start{Start: &process.ProcessEvent_StartEvent{Pid: 7}}}}); err != nil {
		return err
	}
	data := &process.ProcessEvent_DataEvent{Output: &process.ProcessEvent_DataEvent_Stdout{Stdout: []byte(user)}}
	if request.Msg.Pty != nil {
		data.Output = &process.ProcessEvent_DataEvent_Pty{Pty: []byte(user)}
	}
	if err := stream.Send(&process.StartResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_Data{Data: data}}}); err != nil {
		return err
	}
	return stream.Send(&process.StartResponse{Event: &process.ProcessEvent{Event: &process.ProcessEvent_End{End: &process.ProcessEvent_EndEvent{Exited: true}}}})
}

func newUserSelectionSandbox(t *testing.T, check func(*http.Request)) *Sandbox {
	t.Helper()
	mux := http.NewServeMux()
	path, handler := processconnect.NewProcessHandler(userProcessServer{}, connect.WithCodec(tolerantJSONCodec{}))
	mux.Handle(path, handler)
	path, handler = filesystemconnect.NewFilesystemHandler(testFilesystemServer{}, connect.WithCodec(tolerantJSONCodec{}))
	mux.Handle(path, handler)
	mux.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, "hello")
			return
		}
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Error(err)
			http.Error(w, "invalid multipart upload", 400)
			return
		}
		defer r.MultipartForm.RemoveAll()
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"name":"file.txt","path":"/file.txt","type":"file"}]`)
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if check != nil {
			check(r)
		}
		user, _, _ := r.BasicAuth()
		if user == "unknown" {
			_ = connect.NewErrorWriter().Write(w, r, connect.NewError(connect.CodeUnauthenticated, errors.New("unknown user")))
			return
		}
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)
	client, err := NewClient(WithAPIURL(server.URL), WithSandboxURL(server.URL), WithAPIKey("test"))
	if err != nil {
		t.Fatal(err)
	}
	return client.sandboxFromAPI(api.Sandbox{SandboxID: "sbx", TemplateID: "base", EnvdVersion: "0.6.4", EnvdAccessToken: ptr("envd-token"), TrafficAccessToken: ptr("traffic-token")})
}

// These operations exercise public entry points, including the HTTP and RPC
// transports, without depending on the internal header helper.
func userSelectionOperations(ctx context.Context, sandbox *Sandbox, user string, nilOptions bool) map[string]func() error {
	file := &FileOptions{User: user}
	list := &ListFilesOptions{User: user, Depth: 2}
	write := &WriteFileOptions{User: user}
	watch := &WatchOptions{User: user}
	command := &CommandOptions{User: user}
	pty := &PTYOptions{User: user}
	if nilOptions {
		file, list, write, watch, command, pty = nil, nil, nil, nil, nil, nil
	}
	return map[string]func() error{
		"run": func() error { _, err := sandbox.Commands.Run(ctx, "id", command); return err },
		"start": func() error {
			h, err := sandbox.Commands.Start(ctx, "id", command)
			if err != nil {
				return err
			}
			defer h.Close()
			_, err = h.Wait(ctx)
			return err
		},
		"pty": func() error {
			h, err := sandbox.PTY.Create(ctx, "sh", pty)
			if err != nil {
				return err
			}
			defer h.Close()
			_, err = h.Wait(ctx)
			return err
		},
		"watch": func() error {
			h, err := sandbox.Files.Watch(ctx, "~", watch)
			if err != nil {
				return err
			}
			defer h.Close()
			for range h.Events {
			}
			return h.Close()
		},
		"stat":   func() error { _, err := sandbox.Files.Stat(ctx, "~", file); return err },
		"exists": func() error { _, err := sandbox.Files.Exists(ctx, "~", file); return err },
		"list":   func() error { _, err := sandbox.Files.List(ctx, "~", list); return err },
		"mkdir":  func() error { _, err := sandbox.Files.MakeDir(ctx, "~/dir", file); return err },
		"rename": func() error { _, err := sandbox.Files.Rename(ctx, "~/a", "~/b", file); return err },
		"remove": func() error { return sandbox.Files.Remove(ctx, "~/dir", file) },
		"read":   func() error { _, err := sandbox.Files.ReadText(ctx, "~/file", file); return err },
		"readTo": func() error { _, err := sandbox.Files.ReadTo(ctx, "~/file", io.Discard, file); return err },
		"write":  func() error { _, err := sandbox.Files.WriteText(ctx, "~/file", "hello", write); return err },
		"batch": func() error {
			_, err := sandbox.Files.WriteBatch(ctx, []WriteFile{{Path: "~/file", Data: strings.NewReader("hello")}}, write)
			return err
		},
	}
}

func TestEnvdUserSelection(t *testing.T) {
	for _, tc := range []struct {
		name, version, user, want string
		nilOptions                bool
	}{
		{"default", "0.6.4", "", "", true},
		{"empty", "0.6.4", "", "", false},
		{"boundary", "0.4.0", "", "", true},
		{"legacy-default", "0.3.9", "", "user", true},
		{"legacy-empty", "0.3.9", "", "user", false},
		{"explicit", "0.6.4", "root", "root", false},
		{"legacy-explicit", "0.3.9", "root", "root", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sandbox := newUserSelectionSandbox(t, func(r *http.Request) {
				user, password, ok := r.BasicAuth()
				if r.URL.Path == "/files" {
					user = r.URL.Query().Get("username")
					if r.Header.Get("Authorization") != "" {
						t.Error("HTTP files unexpectedly carry RPC authentication")
					}
				} else if ok != (tc.want != "") || password != "" {
					t.Errorf("invalid Basic auth presence or password: %s", r.URL.Path)
				}
				if user != tc.want {
					t.Errorf("%s: user = %q, want %q", r.URL.Path, user, tc.want)
				}
				for key, want := range map[string]string{"X-Access-Token": "envd-token", "Agentbox-Traffic-Access-Token": "traffic-token", "Agentbox-Sandbox-Id": "sbx", "Agentbox-Sandbox-Port": "49983"} {
					if r.Header.Get(key) != want {
						t.Errorf("missing %s", key)
					}
				}
			})
			sandbox.EnvdVersion = tc.version
			for name, call := range userSelectionOperations(t.Context(), sandbox, tc.user, tc.nilOptions) {
				t.Run(name, func(t *testing.T) {
					if err := call(); err != nil {
						t.Fatal(err)
					}
				})
			}
			for _, sign := range []func(string, *FileURLOptions) (string, error){sandbox.Files.SignedReadURL, sandbox.Files.SignedWriteURL} {
				opts := &FileURLOptions{User: tc.user, Expiration: time.Unix(100, 0)}
				if tc.nilOptions {
					opts = nil
				}
				raw, err := sign("~/file", opts)
				if err != nil {
					t.Fatal(err)
				}
				u, err := url.Parse(raw)
				if err != nil {
					t.Fatal(err)
				}
				if u.Query().Get("username") != tc.want {
					t.Fatalf("signed URL user: %s", raw)
				}
			}
		})
	}
}

func TestEnvdInvalidUser(t *testing.T) {
	var requests atomic.Int64
	sandbox := newUserSelectionSandbox(t, func(*http.Request) { requests.Add(1) })
	for _, user := range []string{"root:ignored", "root\n", "root\x00", "root\x7f"} {
		for name, call := range userSelectionOperations(t.Context(), sandbox, user, false) {
			t.Run(user+"/"+name, func(t *testing.T) {
				var invalid *InvalidArgumentError
				if err := call(); !errors.As(err, &invalid) {
					t.Fatalf("expected argument error, got %v", err)
				}
			})
		}
		if _, err := sandbox.Files.SignedReadURL("~/file", &FileURLOptions{User: user}); err == nil {
			t.Fatal("accepted invalid signed URL user")
		}
	}
	if requests.Load() != 0 {
		t.Fatal("invalid user reached envd")
	}
	for name, call := range userSelectionOperations(t.Context(), sandbox, "unknown", false) {
		// Unknown HTTP users are handled by the existing file endpoint error contract.
		if name == "read" || name == "readTo" || name == "write" || name == "batch" {
			continue
		}
		t.Run("unknown/"+name, func(t *testing.T) {
			before := requests.Load()
			var auth *AuthenticationError
			if err := call(); !errors.As(err, &auth) {
				t.Fatalf("expected authentication error, got %v", err)
			}
			if requests.Load()-before != 1 {
				t.Fatal("unexpected retry or fallback")
			}
		})
	}
	_, err := sandbox.Commands.Run(t.Context(), "id", &CommandOptions{User: "unknown", Streaming: &CommandStreamingOptions{}})
	var auth *AuthenticationError
	if !errors.As(err, &auth) {
		t.Fatalf("streaming authentication error: %T %v", err, err)
	}

}

func TestEnvdConcurrentUsersAndStreaming(t *testing.T) {
	sandbox := newUserSelectionSandbox(t, nil)
	var wg sync.WaitGroup
	for i := range 24 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			user := []string{"root", "user"}[i%2]
			var output strings.Builder
			opts := &CommandOptions{User: user}
			if i%3 != 0 {
				opts.Streaming = &CommandStreamingOptions{}
				opts.OnStdout = func(p []byte) { output.Write(p) }
			}
			result, err := sandbox.Commands.Run(t.Context(), "id", opts)
			if err != nil {
				t.Error(err)
				return
			}
			if opts.Streaming == nil {
				output.Write(result.Stdout)
			}
			if output.String() != user {
				t.Errorf("identity crossed calls: got %q, want %q", output.String(), user)
			}
		}()
	}
	wg.Wait()
	var terminal strings.Builder
	h, err := sandbox.PTY.Create(t.Context(), "sh", &PTYOptions{User: "root", Streaming: &CommandStreamingOptions{}, OnPTY: func(p []byte) { terminal.Write(p) }})
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	if _, err := h.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	if terminal.String() != "root" {
		t.Fatalf("PTY user: %q", terminal.String())
	}
}

func TestEnvdUserHeaderScope(t *testing.T) {
	sandbox := newUserSelectionSandbox(t, func(r *http.Request) {
		if r.URL.Path == processconnect.ProcessStartProcedure {
			return
		}
		if r.Header.Get("Authorization") != "" {
			t.Errorf("user leaked to %s", r.URL.Path)
		}
	})
	if _, err := sandbox.Commands.Run(t.Context(), "id", &CommandOptions{User: "root"}); err != nil {
		t.Fatal(err)
	}
	// Existing-process operations do not select or change the process identity,
	// including on old envd versions.
	sandbox.EnvdVersion = "0.3.9"
	h, err := sandbox.Commands.ConnectWithOptions(t.Context(), 7, "", &CommandConnectOptions{Streaming: &CommandStreamingOptions{}})
	if err != nil {
		t.Fatal(err)
	}
	defer h.Close()
	if _, err := h.Wait(t.Context()); err != nil {
		t.Fatal(err)
	}
	for _, call := range []func() error{
		func() error { _, err := sandbox.Commands.List(t.Context()); return err },
		func() error { return sandbox.Commands.Kill(t.Context(), 7, "") },
		func() error { _, err := h.Write(t.Context(), []byte("input")); return err },
		func() error { return h.CloseStdin(t.Context()) },
		func() error { return sandbox.PTY.Resize(t.Context(), h, 80, 24) },
	} {
		if err := call(); err != nil {
			t.Fatal(err)
		}
	}
	_, _ = sandbox.client.Sandboxes.List(t.Context(), nil)
	request, err := http.NewRequestWithContext(t.Context(), http.MethodGet, sandbox.client.config.apiURL+"/unrelated", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := sandbox.client.httpClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
}

func TestWriteBatchUserAndMetadata(t *testing.T) {
	var calls atomic.Int64
	sandbox := newUserSelectionSandbox(t, func(r *http.Request) {
		index := calls.Add(1)
		if r.URL.Query().Get("username") != "root" {
			t.Error("batch lost user")
		}
		want := "common"
		if index == 1 {
			want = "file"
		}
		if r.Header.Get("X-Metadata-Kind") != want || r.Header.Get("X-Metadata-Shared") != "yes" {
			t.Errorf("batch metadata: %v", r.Header)
		}
	})
	opts := &WriteFileOptions{User: "root", Metadata: map[string]string{"kind": "common", "shared": "yes"}, RequestTimeout: time.Second}
	files := []WriteFile{
		{Path: "~/one", Data: strings.NewReader("one"), Metadata: map[string]string{"kind": "file"}},
		{Path: "~/two", Data: strings.NewReader("two")},
	}
	if entries, err := sandbox.Files.WriteBatch(t.Context(), files, opts); err != nil || len(entries) != 2 {
		t.Fatalf("batch: %#v %v", entries, err)
	}
	if opts.Metadata["kind"] != "common" || len(files[0].Metadata) != 1 {
		t.Fatal("batch modified caller options")
	}
}
