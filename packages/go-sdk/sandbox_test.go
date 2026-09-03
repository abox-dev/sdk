package agentbox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
)

func TestSandboxLifecycleAPI(t *testing.T) {
	requests := make(chan *http.Request, 32)
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests <- request.Clone(request.Context())
		writer.Header().Set("Content-Type", "application/json")
		switch {
		case request.Method == http.MethodPost && request.URL.Path == "/sandboxes":
			writer.WriteHeader(201)
			fmt.Fprint(writer, `{"sandboxID":"sbx","templateID":"base","envdVersion":"1.0","domain":"example.test","envdAccessToken":"envd","trafficAccessToken":"traffic"}`)
		case request.URL.Path == "/v2/sandboxes":
			fmt.Fprint(writer, `[{"sandboxID":"sbx","templateID":"base","envdVersion":"1","cpuCount":2,"memoryMB":1024,"diskSizeMB":1024,"startedAt":"2024-01-01T00:00:00Z","endAt":"2024-01-01T01:00:00Z","state":"running"}]`)
		case request.URL.Path == "/sandboxes/sbx" && request.Method == http.MethodGet:
			fmt.Fprint(writer, `{"sandboxID":"sbx","templateID":"base","envdVersion":"1","cpuCount":2,"memoryMB":1024,"diskSizeMB":1024,"startedAt":"2024-01-01T00:00:00Z","endAt":"2024-01-01T01:00:00Z","state":"running"}`)
		case request.URL.Path == "/sandboxes/sbx/connect":
			writer.WriteHeader(200)
			fmt.Fprint(writer, `{"sandboxID":"sbx","templateID":"base","envdVersion":"1"}`)
		case request.URL.Path == "/sandboxes/sbx/metrics":
			fmt.Fprint(writer, `[{"timestampUnix":1,"cpuCount":2,"cpuUsedPct":1,"memUsed":1,"memTotal":2,"memCache":0,"diskUsed":1,"diskTotal":2}]`)
		case request.URL.Path == "/sandboxes/metrics":
			fmt.Fprint(writer, `{"sandboxes":{"sbx":{"timestampUnix":1,"cpuCount":2,"cpuUsedPct":1,"memUsed":1,"memTotal":2,"memCache":0,"diskUsed":1,"diskTotal":2}}}`)
		case request.URL.Path == "/v2/sandboxes/sbx/logs":
			fmt.Fprint(writer, `{"logs":[{"timestamp":"2024-01-01T00:00:00Z","level":"info","message":"ready","fields":{}}],"nextCursor":"next"}`)
		case request.URL.Path == "/sandboxes/sbx/fork":
			writer.WriteHeader(201)
			fmt.Fprint(writer, `[{"sandbox":{"sandboxID":"fork","templateID":"base","envdVersion":"1"}},{"error":{"code":409,"message":"failed"}}]`)
		case request.URL.Path == "/sandboxes/sbx/snapshots":
			writer.WriteHeader(201)
			fmt.Fprint(writer, `{"snapshotID":"snap:default","names":["snap:default"]}`)
		case request.URL.Path == "/snapshots":
			fmt.Fprint(writer, `[{"snapshotID":"snap:default","names":["snap:default"]}]`)
		case request.Method == http.MethodDelete && request.URL.Path == "/templates/missing":
			writer.WriteHeader(http.StatusNotFound)
		case request.Method == http.MethodDelete && request.URL.Path == "/sandboxes/missing":
			writer.WriteHeader(http.StatusNotFound)
		case request.Method == http.MethodDelete:
			writer.WriteHeader(http.StatusNoContent)
		default:
			writer.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()
	client, err := NewClient(WithAPIURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	sandbox, err := client.Sandboxes.Create(ctx, &CreateSandboxOptions{Template: "base", Timeout: time.Minute, Env: map[string]string{"A": "B"}, Metadata: map[string]string{"m": "v"}})
	if err != nil {
		t.Fatal(err)
	}
	if sandbox.ID != "sbx" || sandbox.Host(8080) != "8080-sbx.example.test" || sandbox.Commands == nil || sandbox.Files == nil || sandbox.PTY == nil {
		t.Fatalf("bad sandbox: %#v", sandbox)
	}
	if page, err := client.Sandboxes.List(ctx, &ListSandboxOptions{Metadata: map[string]string{"app": "test"}, States: []SandboxState{SandboxRunning}, NextToken: "token", Limit: 10}); err != nil || len(page.Items) != 1 {
		t.Fatalf("list: %v %#v", err, page)
	}
	if _, err := client.Sandboxes.Connect(ctx, "sbx", &ConnectSandboxOptions{Timeout: time.Minute}); err != nil {
		t.Fatal(err)
	}
	if _, err := sandbox.Info(ctx); err != nil {
		t.Fatal(err)
	}
	if metrics, err := sandbox.Metrics(ctx, &MetricsOptions{Start: time.Unix(1, 0), End: time.Unix(2, 0)}); err != nil || len(metrics) != 1 {
		t.Fatalf("metrics: %v", err)
	}
	if metrics, err := client.Sandboxes.Metrics(ctx, "sbx"); err != nil || len(metrics) != 1 {
		t.Fatalf("all metrics: %v", err)
	}
	if logs, err := sandbox.Logs(ctx, &SandboxLogOptions{Cursor: "cursor", Timestamp: 1, Limit: 2, Direction: "forward", Level: "info", Search: "ready"}); err != nil || logs.NextToken != "next" {
		t.Fatalf("logs: %v %#v", err, logs)
	}
	if err := sandbox.SetTimeout(ctx, time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.KeepAlive(ctx, time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.Pause(ctx, &PauseOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.UpdateNetwork(ctx, SandboxNetworkConfig{}, nil); err != nil {
		t.Fatal(err)
	}
	if forks, err := sandbox.Fork(ctx, &ForkOptions{Count: 2, Timeout: time.Minute}); err != nil || len(forks) != 2 || forks[0].Sandbox == nil || forks[1].Err == nil {
		t.Fatalf("fork: %v %#v", err, forks)
	}
	if snapshot, err := sandbox.CreateSnapshot(ctx, "snap"); err != nil || snapshot.SnapshotID == "" {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshots, err := client.Sandboxes.Snapshots(ctx, &SnapshotListOptions{SandboxID: "sbx", Name: "snap", Limit: 10}); err != nil || len(snapshots.Items) != 1 {
		t.Fatalf("snapshots: %v", err)
	}
	if deleted, err := client.Sandboxes.DeleteSnapshot(ctx, "snap"); err != nil || !deleted {
		t.Fatalf("delete snapshot: %v %v", deleted, err)
	}
	if deleted, err := client.Sandboxes.DeleteSnapshot(ctx, "missing"); err != nil || deleted {
		t.Fatalf("missing snapshot: %v %v", deleted, err)
	}
	if killed, err := sandbox.Kill(ctx); err != nil || !killed {
		t.Fatal(err)
	}
	if killed, err := client.Sandboxes.Kill(ctx, "missing"); err != nil || killed {
		t.Fatalf("missing sandbox kill: %v %v", killed, err)
	}
	close(requests)
	for request := range requests {
		if request.Header.Get("User-Agent") == "" {
			t.Fatal("missing user agent")
		}
	}
}

func TestSandboxValidationAndSigning(t *testing.T) {
	client, _ := NewClient(WithAPIURL("http://localhost"))
	debugClient, _ := NewClient(WithAPIURL("http://localhost"), WithDebug(true))
	if killed, err := debugClient.Sandboxes.Kill(t.Context(), "debug"); err != nil || !killed {
		t.Fatalf("debug kill: %v %v", killed, err)
	}
	if _, err := client.Sandboxes.Connect(context.Background(), "", nil); err == nil {
		t.Fatal("expected validation")
	}
	if _, err := client.Sandboxes.Metrics(context.Background()); err == nil {
		t.Fatal("expected metrics validation")
	}
	if _, err := durationSeconds(-time.Second, 0); err == nil {
		t.Fatal("expected duration validation")
	}
	if _, err := client.Sandboxes.DeleteSnapshot(context.Background(), ""); err == nil {
		t.Fatal("expected snapshot validation")
	}
	placeholders, err := IAMTokenPlaceholders("aws", "gcp")
	if err != nil || placeholders["aws"] != "${agentbox.identity.tokens.aws}" {
		t.Fatalf("IAM placeholders: %#v %v", placeholders, err)
	}
	for _, name := range []string{"", "bad}", "bad\nname"} {
		if _, err := IAMTokenPlaceholder(name); err == nil {
			t.Fatalf("expected invalid IAM name %q", name)
		}
	}
	tokens := SandboxIAMTokens{"bad{": {Audience: "aud", TokenType: "jwt"}}
	if _, err := client.Sandboxes.Create(context.Background(), &CreateSandboxOptions{IAM: &SandboxIAM{Tokens: &tokens}}); err == nil {
		t.Fatal("expected invalid IAM config")
	}
	sandbox := client.sandboxFromAPI(api.Sandbox{SandboxID: "id", TemplateID: "base", EnvdVersion: "1"})
	if _, err := sandbox.Files.SignedReadURL("/x", "", time.Time{}); err == nil {
		t.Fatal("expected missing token")
	}
	sandbox.envdAccessToken = "secret"
	sandbox.trafficAccessToken = "traffic-secret"
	if rendered := fmt.Sprintf("%+v %#v", sandbox, sandbox); strings.Contains(rendered, "secret") {
		t.Fatalf("sandbox formatting leaked credentials: %s", rendered)
	}
	read, err := sandbox.Files.SignedReadURL("/x", "user", time.Unix(100, 0))
	if err != nil || !strings.Contains(read, "signature=v1_") || !strings.Contains(read, "signature_expiration=100") {
		t.Fatalf("signed URL: %s %v", read, err)
	}
	write, err := sandbox.Files.SignedWriteURL("/x", "", time.Time{})
	if err != nil || write == read {
		t.Fatalf("signed write URL: %v", err)
	}
}

func TestSandboxHealthAndModelConversionFailures(t *testing.T) {
	status := http.StatusBadGateway
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(status)
	}))
	defer server.Close()
	client, err := NewClient(WithAPIURL(server.URL), WithSandboxURL(server.URL))
	if err != nil {
		t.Fatal(err)
	}
	sandbox := client.sandboxFromAPI(api.Sandbox{SandboxID: "id", TemplateID: "base", EnvdVersion: "1"})
	if running, err := sandbox.IsRunning(t.Context()); err != nil || running {
		t.Fatalf("bad gateway health: %v %v", running, err)
	}
	status = http.StatusInternalServerError
	if _, err := sandbox.IsRunning(t.Context()); err == nil {
		t.Fatal("expected health error")
	}
	if _, err := convertModel[any](make(chan int)); err == nil {
		t.Fatal("expected model encode error")
	}
	if _, err := convertModel[chan int]("invalid"); err == nil {
		t.Fatal("expected model decode error")
	}
}

func decodeRequest(t *testing.T, request *http.Request, value any) {
	t.Helper()
	if err := json.NewDecoder(request.Body).Decode(value); err != nil {
		t.Fatal(err)
	}
}

func TestControlPlaneFailures(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, `{"message":"failure"}`)
	}))
	defer server.Close()
	client, _ := NewClient(WithAPIURL(server.URL))
	ctx := context.Background()
	sandbox := client.sandboxFromAPI(api.Sandbox{SandboxID: "id", TemplateID: "base", EnvdVersion: "1"})
	calls := []func() error{
		func() error { _, err := client.Sandboxes.Create(ctx, nil); return err },
		func() error { _, err := client.Sandboxes.Connect(ctx, "id", nil); return err },
		func() error { _, err := client.Sandboxes.List(ctx, nil); return err },
		func() error { _, err := client.Sandboxes.Info(ctx, "id"); return err },
		func() error { _, err := client.Sandboxes.Kill(ctx, "id"); return err },
		func() error { _, err := client.Sandboxes.Snapshots(ctx, nil); return err },
		func() error { _, err := client.Sandboxes.DeleteSnapshot(ctx, "id"); return err },
		func() error { _, err := client.Sandboxes.Metrics(ctx, "id"); return err },
		func() error { _, err := client.Sandboxes.Logs(ctx, "id", nil); return err },
		func() error { return sandbox.SetTimeout(ctx, time.Second) },
		func() error { return sandbox.KeepAlive(ctx, 0) },
		func() error { return sandbox.Pause(ctx, nil) },
		func() error { return sandbox.UpdateNetwork(ctx, SandboxNetworkConfig{}, nil) },
		func() error { _, err := sandbox.Metrics(ctx, nil); return err },
		func() error { _, err := sandbox.Fork(ctx, nil); return err },
		func() error { _, err := sandbox.CreateSnapshot(ctx, ""); return err },
	}
	for index, call := range calls {
		if err := call(); err == nil {
			t.Fatalf("call %d unexpectedly succeeded", index)
		}
	}
}
