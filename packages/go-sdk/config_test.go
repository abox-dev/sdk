package agentbox

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
)

func TestClientConfigurationAndTransport(t *testing.T) {
	t.Setenv("AGENTBOX_API_KEY", "environment-key")
	t.Setenv("AGENTBOX_DOMAIN", "environment.test")
	t.Setenv("AGENTBOX_DEBUG", "false")
	var received http.Header
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		received = request.Header.Clone()
		writer.Header().Set("Content-Type", "application/json")
		io.WriteString(writer, `[]`)
	}))
	defer server.Close()
	client, err := NewClient(
		WithAPIKey("option-key"), WithAPIURL(server.URL), WithDomain("option.test"),
		WithSandboxURL(server.URL), WithDebug(true), WithRequestTimeout(time.Second),
		WithHeaders(http.Header{"X-Custom": {"value"}}),
		WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))),
	)
	if err != nil {
		t.Fatal(err)
	}
	if client.config.apiKey != "option-key" || client.config.domain != "option.test" || !client.config.debug {
		t.Fatalf("unexpected config: %#v", client.config)
	}
	if _, err := client.Sandboxes.List(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if received.Get("X-API-KEY") != "option-key" || received.Get("X-Custom") != "value" || received.Get("User-Agent") != "agentbox-go-sdk/"+Version {
		t.Fatalf("unexpected headers: %v", received)
	}
}

func TestTransportScopesCredentialsAndConfiguresH2C(t *testing.T) {
	var sandboxHeaders http.Header
	sandboxServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		sandboxHeaders = request.Header.Clone()
		writer.WriteHeader(http.StatusNoContent)
	}))
	defer sandboxServer.Close()
	apiServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		io.WriteString(writer, `[]`)
	}))
	defer apiServer.Close()
	client, err := NewClient(WithAPIURL(apiServer.URL), WithAPIKey("secret"), WithHeaders(http.Header{"X-Control": {"only"}}))
	if err != nil {
		t.Fatal(err)
	}
	request, _ := http.NewRequest(http.MethodGet, sandboxServer.URL+"/path?token=private", nil)
	response, err := client.httpClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if sandboxHeaders.Get("X-API-KEY") != "" || sandboxHeaders.Get("X-Control") != "" || sandboxHeaders.Get("User-Agent") == "" {
		t.Fatalf("credentials escaped control plane: %v", sandboxHeaders)
	}

	debugClient, err := NewClient(WithAPIURL(apiServer.URL), WithDebug(true))
	if err != nil {
		t.Fatal(err)
	}
	transport := debugClient.envdClient.Transport.(*sdkTransport).base.(*http.Transport)
	if transport.Protocols == nil || !transport.Protocols.UnencryptedHTTP2() || transport.Protocols.HTTP1() {
		t.Fatalf("expected prior-knowledge h2c transport, got %v", transport.Protocols)
	}
	if transportAPIHost("::bad") != "" {
		t.Fatal("invalid API URL unexpectedly has a host")
	}
	custom := &http.Client{Timeout: 37 * time.Second}
	customClient, err := NewClient(WithAPIURL(apiServer.URL), WithHTTPClient(custom))
	if err != nil {
		t.Fatal(err)
	}
	if customClient.httpClient.Timeout != custom.Timeout {
		t.Fatalf("custom HTTP timeout = %s, want %s", customClient.httpClient.Timeout, custom.Timeout)
	}
}

func TestInvalidClientOptions(t *testing.T) {
	t.Setenv("AGENTBOX_DEBUG", "invalid")
	if _, err := NewClient(); err == nil {
		t.Fatal("expected invalid environment error")
	}
	t.Setenv("AGENTBOX_DEBUG", "false")
	tests := []ClientOption{WithDomain(" "), WithAPIURL("ftp://bad"), WithSandboxURL("bad"), WithProxy("bad"), WithRequestTimeout(-1), WithHTTPClient(nil), nil}
	for _, option := range tests {
		if _, err := NewClient(option); err == nil {
			t.Fatalf("expected failure for option %#v", option)
		}
	}
	if _, err := NewClient(WithHTTPClient(http.DefaultClient), WithProxy("http://localhost:8080")); err == nil {
		t.Fatal("expected conflicting transport options")
	}
	t.Setenv("AGENTBOX_MAX_CONNECTIONS", "17")
	t.Setenv("AGENTBOX_MAX_KEEPALIVE_CONNECTIONS", "3")
	t.Setenv("AGENTBOX_KEEPALIVE_EXPIRY", "2.5")
	client, err := NewClient(WithProxy("http://localhost:8080"))
	if err != nil {
		t.Fatal(err)
	}
	transport := client.httpClient.Transport.(*sdkTransport).base.(*http.Transport)
	if transport.MaxIdleConns != 17 || transport.MaxIdleConnsPerHost != 3 || transport.IdleConnTimeout != 2500*time.Millisecond {
		t.Fatalf("pool config: %#v", transport)
	}
	if roundTripper(nil) == nil {
		t.Fatal("nil round tripper")
	}
}

func TestConnectErrorMappingsAndCodec(t *testing.T) {
	for code, target := range map[connect.Code]error{connect.CodeUnauthenticated: &AuthenticationError{}, connect.CodePermissionDenied: &AuthenticationError{}, connect.CodeNotFound: &SandboxError{}, connect.CodeResourceExhausted: &RateLimitError{}, connect.CodeCanceled: &TimeoutError{}, connect.CodeDeadlineExceeded: &TimeoutError{}, connect.CodeUnavailable: &TimeoutError{}, connect.CodeInvalidArgument: &InvalidArgumentError{}, connect.CodeInternal: &SandboxError{}} {
		mapped := connectError(connect.NewError(code, errors.New("failure")))
		if reflect.TypeOf(mapped) != reflect.TypeOf(target) {
			t.Fatalf("code %s mapped to %T", code, mapped)
		}
	}
	if connectError(nil) != nil {
		t.Fatal("nil error changed")
	}
	plain := errors.New("plain")
	if connectError(plain) != plain {
		t.Fatal("plain error changed")
	}
	codec := tolerantJSONCodec{}
	if codec.Name() != "json" {
		t.Fatal(codec.Name())
	}
	if _, err := codec.Marshal("bad"); err == nil {
		t.Fatal("expected marshal type error")
	}
	if err := codec.Unmarshal(nil, "bad"); err == nil {
		t.Fatal("expected unmarshal type error")
	}
	network := &net.DNSError{Err: "failure", Name: "example.test"}
	if !isConnectionError(network) {
		t.Fatal("network error not detected")
	}
	if _, ok := normalizeRequestError(network).(*SandboxError); !ok {
		t.Fatal("network error not normalized")
	}
	if timeout := normalizeRequestError(context.DeadlineExceeded); !errors.Is(timeout, context.DeadlineExceeded) {
		t.Fatal("deadline not preserved")
	} else {
		var typed *TimeoutError
		if !errors.As(timeout, &typed) {
			t.Fatalf("deadline mapped to %T", timeout)
		}
	}
	ctx, cancel := withRequestTimeout(context.Background(), 0)
	cancel()
	if ctx.Err() == nil {
		t.Fatal("zero-timeout context not cancellable")
	}
	client, err := NewClient(WithAPIURL("http://example.test"), WithHTTPClient(&http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) { return nil, network })}), WithLogger(slog.New(slog.NewTextHandler(io.Discard, nil))))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Sandboxes.List(context.Background(), nil); err == nil {
		t.Fatal("expected transport error")
	}
	if _, err := NewClient(WithHTTPClient(&http.Client{})); err != nil {
		t.Fatal(err)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (function roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return function(request)
}

func TestErrorsAndHelpers(t *testing.T) {
	cause := errors.New("cause")
	values := []error{
		&SandboxError{APIError: APIError{StatusCode: 500, Message: "sandbox", Cause: cause}},
		&AuthenticationError{APIError: APIError{Message: "auth", Cause: cause}},
		&RateLimitError{APIError: APIError{Message: "rate", Cause: cause}},
		&TimeoutError{APIError: APIError{Message: "timeout", Cause: cause}},
		&FileNotFoundError{APIError: APIError{Message: "file", Cause: cause}},
		&NotEnoughSpaceError{APIError: APIError{Message: "disk", Cause: cause}},
		&SandboxNotFoundError{APIError: APIError{Message: "missing", Cause: cause}},
		&TemplateError{APIError: APIError{Message: "template", Cause: cause}},
		&BuildError{APIError: APIError{Message: "build", Cause: cause}},
		&FileUploadError{APIError: APIError{Message: "upload", Cause: cause}},
		&InvalidArgumentError{Message: "argument", Cause: cause},
	}
	for _, value := range values {
		if !strings.Contains(value.Error(), "agentbox:") || !errors.Is(value, cause) {
			t.Fatalf("bad error %T: %v", value, value)
		}
	}
	for status, target := range map[int]error{401: &AuthenticationError{}, 403: &AuthenticationError{}, 404: &SandboxNotFoundError{}, 413: &NotEnoughSpaceError{}, 429: &RateLimitError{}, 502: &TimeoutError{}, 504: &TimeoutError{}, 500: &SandboxError{}} {
		err := decodeStatusError(status, http.StatusText(status), []byte(`{"message":"failure"}`))
		if reflect.TypeOf(err) != reflect.TypeOf(target) {
			t.Fatalf("status %d mapped to %T", status, err)
		}
	}
	if err := decodeStatusError(http.StatusBadRequest, "bad request", []byte("{\"message\":\"json-message\",\"code\":42}")); !strings.Contains(err.Error(), "json-message") {
		t.Fatal(err)
	} else if sandboxErr := err.(*SandboxError); sandboxErr.Code != "42" {
		t.Fatalf("error code: %q", sandboxErr.Code)
	}
	if err := decodeStatusError(http.StatusTeapot, "teapot", nil); !strings.Contains(err.Error(), "teapot") {
		t.Fatal(err)
	}
}
