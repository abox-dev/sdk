package codeinterpreter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/abox-dev/sdk/packages/go-sdk"
)

type rewriteTransport struct{ target *url.URL }

func (transport rewriteTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	copied := *request.URL
	copied.Scheme, copied.Host = transport.target.Scheme, transport.target.Host
	clone.URL = &copied
	return http.DefaultTransport.RoundTrip(clone)
}

func TestCodeInterpreter(t *testing.T) {
	var failContexts atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		if failContexts.Load() && request.URL.Path == "/contexts" {
			http.Error(writer, "failure", http.StatusInternalServerError)
			return
		}
		switch {
		case request.URL.Path == "/sandboxes":
			writer.WriteHeader(http.StatusCreated)
			fmt.Fprint(writer, `{"sandboxID":"sbx","templateID":"code-interpreter","envdVersion":"1","domain":"example.test","envdAccessToken":"envd","trafficAccessToken":"traffic"}`)
		case request.URL.Path == "/sandboxes/sbx/connect":
			fmt.Fprint(writer, `{"sandboxID":"sbx","templateID":"code-interpreter","envdVersion":"1","domain":"example.test"}`)
		case request.URL.Path == "/execute":
			var input struct {
				Code string `json:"code"`
			}
			json.NewDecoder(request.Body).Decode(&input)
			if input.Code == "request-slow" {
				time.Sleep(100 * time.Millisecond)
				return
			}
			if input.Code == "http-error" {
				http.Error(writer, "failure", http.StatusInternalServerError)
				return
			}
			if input.Code == "slow" {
				writer.(http.Flusher).Flush()
				time.Sleep(100 * time.Millisecond)
				return
			}
			fmt.Fprintln(writer, `{"type":"stdout","text":"hello"}`)
			fmt.Fprintln(writer, `{"type":"stderr","text":"warning"}`)
			fmt.Fprintln(writer, `{"type":"result","text":"42","html":"<b>42</b>","is_main_result":true,"unknown":{"value":1},"chart":{"type":"future","title":"chart","elements":[],"future":true}}`)
			fmt.Fprintln(writer, `{"type":"error","name":"ValueError","value":"bad","traceback":"trace"}`)
			fmt.Fprintln(writer, `{"type":"number_of_executions","execution_count":3}`)
		case request.URL.Path == "/contexts" && request.Method == http.MethodPost:
			var input CreateContextOptions
			json.NewDecoder(request.Body).Decode(&input)
			if input.Cwd == "badjson" {
				fmt.Fprint(writer, `{`)
				return
			}
			if input.Language == TypeScript {
				http.Error(writer, "failure", http.StatusInternalServerError)
				return
			}
			fmt.Fprint(writer, `{"id":"ctx","language":"python","cwd":"/home/user"}`)
		case request.URL.Path == "/contexts" && request.Method == http.MethodGet:
			fmt.Fprint(writer, `[{"id":"ctx","language":"python","cwd":"/home/user"}]`)
		case strings.HasPrefix(request.URL.Path, "/contexts/"):
			writer.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(writer, request)
		}
	}))
	defer server.Close()
	target, _ := url.Parse(server.URL)
	httpClient := &http.Client{Transport: rewriteTransport{target: target}}
	client, err := NewClient(agentbox.WithAPIURL(server.URL), agentbox.WithSandboxURL(server.URL), agentbox.WithHTTPClient(httpClient))
	if err != nil {
		t.Fatal(err)
	}
	sandbox, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if sandbox.TemplateID != DefaultTemplate {
		t.Fatalf("template: %s", sandbox.TemplateID)
	}
	stdout, stderr, results, callbackErrors := 0, 0, 0, 0
	execution, err := sandbox.RunCode(context.Background(), "print(42)", &RunCodeOptions{Language: Python, Env: map[string]string{"A": "B"}, OnStdout: func(OutputMessage) { stdout++ }, OnStderr: func(OutputMessage) { stderr++ }, OnResult: func(Result) { results++ }, OnError: func(ExecutionError) { callbackErrors++ }})
	if err != nil {
		t.Fatal(err)
	}
	if stdout != 1 || stderr != 1 || results != 1 || callbackErrors != 1 || execution.Text() != "42" || execution.ExecutionCount != 3 || execution.Results[0].Chart.Type != ChartUnknown || len(execution.Results[0].Extra) != 1 {
		t.Fatalf("execution: %#v", execution)
	}
	if _, err := sandbox.RunCode(context.Background(), "x", &RunCodeOptions{Language: Python, Context: &Context{ID: "ctx"}}); err == nil {
		t.Fatal("expected language/context validation")
	}
	if _, err := sandbox.RunCode(context.Background(), "slow", &RunCodeOptions{ExecutionTimeout: 10 * time.Millisecond}); err == nil {
		t.Fatal("expected execution timeout")
	} else {
		var timeout *agentbox.TimeoutError
		if !errors.As(err, &timeout) {
			t.Fatalf("timeout type: %T", err)
		}
	}
	if _, err := sandbox.RunCode(context.Background(), "request-slow", &RunCodeOptions{RequestTimeout: 10 * time.Millisecond}); err == nil {
		t.Fatal("expected request timeout")
	}
	if _, err := sandbox.RunCode(context.Background(), "http-error", nil); err == nil {
		t.Fatal("expected HTTP error")
	}
	if _, err := sandbox.RunCode(context.Background(), "", nil); err == nil {
		t.Fatal("expected empty-code validation")
	}
	if _, err := sandbox.RunCode(context.Background(), "x", &RunCodeOptions{RequestTimeout: -time.Second}); err == nil {
		t.Fatal("expected timeout validation")
	}
	created, err := sandbox.CreateContext(context.Background(), &CreateContextOptions{Language: JavaScript, Cwd: "/tmp"})
	if err != nil || created.ID != "ctx" {
		t.Fatalf("create context: %#v %v", created, err)
	}
	contexts, err := sandbox.ListContexts(context.Background())
	if err != nil || len(contexts) != 1 {
		t.Fatalf("contexts: %#v %v", contexts, err)
	}
	if err := sandbox.RestartContext(context.Background(), "ctx"); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.RemoveContext(context.Background(), "ctx"); err != nil {
		t.Fatal(err)
	}
	if err := sandbox.RemoveContext(context.Background(), ""); err == nil {
		t.Fatal("expected context ID validation")
	}
	if err := sandbox.RestartContext(context.Background(), ""); err == nil {
		t.Fatal("expected context ID validation")
	}
	if _, err := sandbox.CreateContext(context.Background(), &CreateContextOptions{Cwd: "badjson"}); err == nil {
		t.Fatal("expected invalid context response")
	}
	if _, err := sandbox.CreateContext(context.Background(), &CreateContextOptions{Language: TypeScript}); err == nil {
		t.Fatal("expected context HTTP error")
	}
	if _, err := client.Connect(context.Background(), "sbx", nil); err != nil {
		t.Fatal(err)
	}
	noTimeoutClient, err := NewClient(agentbox.WithAPIURL(server.URL), agentbox.WithSandboxURL(server.URL), agentbox.WithHTTPClient(httpClient), agentbox.WithRequestTimeout(0))
	if err != nil {
		t.Fatal(err)
	}
	noTimeoutSandbox, err := noTimeoutClient.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := noTimeoutSandbox.ListContexts(context.Background()); err != nil {
		t.Fatal(err)
	}
	failContexts.Store(true)
	if _, err := sandbox.ListContexts(context.Background()); err == nil {
		t.Fatal("expected list contexts error")
	}
}

func TestModelsAndHTTPError(t *testing.T) {
	message := OutputMessage{Line: "line"}
	if message.String() != "line" {
		t.Fatal(message)
	}
	executionError := ExecutionError{Name: "Error", Value: "bad"}
	if executionError.Error() != "Error: bad" {
		t.Fatal(executionError)
	}
	var chart Chart
	if err := json.Unmarshal([]byte(`{"type":"line","title":"x","elements":[],"extra":1}`), &chart); err != nil || chart.Type != ChartLine || len(chart.Extra) != 1 {
		t.Fatalf("chart: %#v %v", chart, err)
	}
	result := resultFromRaw(map[string]json.RawMessage{"type": json.RawMessage(`"result"`), "text": json.RawMessage(`"ok"`), "is_main_result": json.RawMessage(`true`), "other": json.RawMessage(`1`)})
	if result.Text != "ok" || len(result.Formats()) != 2 {
		t.Fatalf("result: %#v", result)
	}
	for status, target := range map[int]error{http.StatusBadGateway: &agentbox.TimeoutError{}, http.StatusGatewayTimeout: &agentbox.TimeoutError{}, http.StatusNotFound: &agentbox.SandboxError{}, http.StatusInternalServerError: &agentbox.SandboxError{}} {
		response := &http.Response{StatusCode: status, Status: http.StatusText(status), Body: io.NopCloser(strings.NewReader("failure"))}
		err := interpreterHTTPError(response)
		if fmt.Sprintf("%T", err) != fmt.Sprintf("%T", target) {
			t.Fatalf("status %d: %T", status, err)
		}
	}
	if (Execution{}).Text() != "" {
		t.Fatal("empty execution text")
	}
	if err := parseOutput([]byte(`not-json`), &Execution{}, &RunCodeOptions{}); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestClientFailures(t *testing.T) {
	t.Setenv("AGENTBOX_DEBUG", "invalid")
	if _, err := NewClient(); err == nil {
		t.Fatal("expected client config error")
	}
	t.Setenv("AGENTBOX_DEBUG", "false")
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, `{"message":"failure"}`)
	}))
	defer server.Close()
	target, _ := url.Parse(server.URL)
	client, err := NewClient(agentbox.WithAPIURL(server.URL), agentbox.WithSandboxURL(server.URL), agentbox.WithHTTPClient(&http.Client{Transport: rewriteTransport{target: target}}))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Create(context.Background(), nil); err == nil {
		t.Fatal("expected create error")
	}
	if _, err := client.Connect(context.Background(), "id", nil); err == nil {
		t.Fatal("expected connect error")
	}
}
