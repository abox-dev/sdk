package agentbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const maxErrorBody = 1 << 20

type sdkTransport struct {
	base    http.RoundTripper
	headers http.Header
	apiKey  string
	apiHost string
	logger  *slog.Logger
}

func (transport *sdkTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	clone := request.Clone(request.Context())
	clone.Header = request.Header.Clone()
	if clone.URL.Host == transport.apiHost {
		for key, values := range transport.headers {
			if _, exists := clone.Header[key]; !exists {
				clone.Header[key] = slices.Clone(values)
			}
		}
		if transport.apiKey != "" && clone.Header.Get("X-API-KEY") == "" {
			clone.Header.Set("X-API-KEY", transport.apiKey)
		}
	}
	if clone.Header.Get("User-Agent") == "" {
		clone.Header.Set("User-Agent", "agentbox-go-sdk/"+Version)
	}
	started := time.Now()
	response, err := transport.base.RoundTrip(clone)
	if transport.logger != nil {
		safeURL := clone.URL.Scheme + "://" + clone.URL.Host + clone.URL.EscapedPath()
		attributes := []any{"method", clone.Method, "url", safeURL, "duration", time.Since(started)}
		if response != nil {
			attributes = append(attributes, "status", response.StatusCode)
		}
		if err != nil {
			attributes = append(attributes, "error", err)
		}
		transport.logger.DebugContext(request.Context(), "AgentBox HTTP request", attributes...)
	}
	return response, err
}

func newHTTPClient(config clientConfig) *http.Client {
	apiURL, _ := url.Parse(config.apiURL)
	apiHost := ""
	if apiURL != nil {
		apiHost = apiURL.Host
	}
	if config.httpClient != nil {
		client := *config.httpClient
		client.Transport = &sdkTransport{base: roundTripper(client.Transport), headers: config.headers.Clone(), apiKey: config.apiKey, apiHost: apiHost, logger: config.logger}
		return &client
	}
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.ForceAttemptHTTP2 = true
	base.MaxIdleConns = envPositiveInt("AGENTBOX_MAX_CONNECTIONS", 200)
	base.MaxIdleConnsPerHost = envPositiveInt("AGENTBOX_MAX_KEEPALIVE_CONNECTIONS", 20)
	base.IdleConnTimeout = envDurationSeconds("AGENTBOX_KEEPALIVE_EXPIRY", 300*time.Second)
	if config.proxyURL != nil {
		base.Proxy = http.ProxyURL(config.proxyURL)
	}
	return &http.Client{Transport: &sdkTransport{base: base, headers: config.headers.Clone(), apiKey: config.apiKey, apiHost: apiHost, logger: config.logger}}
}

// newEnvdHTTPClient enables prior-knowledge h2c for local debug envd. Production
// envd uses HTTPS and negotiates HTTP/2 normally. A caller-supplied client or
// proxy remains authoritative because it may provide its own routing transport.
func newEnvdHTTPClient(config clientConfig, standard *http.Client) *http.Client {
	if !config.debug || config.httpClient != nil || config.proxyURL != nil {
		return standard
	}
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.MaxIdleConns = envPositiveInt("AGENTBOX_MAX_CONNECTIONS", 200)
	base.MaxIdleConnsPerHost = envPositiveInt("AGENTBOX_MAX_KEEPALIVE_CONNECTIONS", 20)
	base.IdleConnTimeout = envDurationSeconds("AGENTBOX_KEEPALIVE_EXPIRY", 300*time.Second)
	base.Protocols = new(http.Protocols)
	base.Protocols.SetUnencryptedHTTP2(true)
	return &http.Client{Transport: &sdkTransport{base: base, apiHost: transportAPIHost(config.apiURL), logger: config.logger}}
}

func transportAPIHost(value string) string {
	parsed, _ := url.Parse(value)
	if parsed == nil {
		return ""
	}
	return parsed.Host
}

func roundTripper(value http.RoundTripper) http.RoundTripper {
	if value == nil {
		return http.DefaultTransport
	}
	return value
}

func envPositiveInt(name string, fallback int) int {
	value, err := strconv.Atoi(osEnv(name))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}

func envDurationSeconds(name string, fallback time.Duration) time.Duration {
	value, err := strconv.ParseFloat(osEnv(name), 64)
	if err != nil || value <= 0 {
		return fallback
	}
	return time.Duration(value * float64(time.Second))
}

var osEnv = func(name string) string { return strings.TrimSpace(getenv(name)) }
var getenv = os.Getenv

func withRequestTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		return context.WithCancel(ctx)
	}
	return context.WithTimeout(ctx, timeout)
}

func decodeHTTPError(response *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(response.Body, maxErrorBody))
	return decodeStatusError(response.StatusCode, response.Status, body)
}

func decodeFileHTTPError(response *http.Response) error {
	err := decodeHTTPError(response)
	if response.StatusCode != http.StatusNotFound {
		return err
	}
	var missing *SandboxNotFoundError
	if errors.As(err, &missing) {
		return &FileNotFoundError{APIError: missing.APIError}
	}
	return err
}

func decodeStatusError(statusCode int, status string, body []byte) error {
	message := strings.TrimSpace(string(body))
	var payload struct {
		Message string `json:"message"`
		Code    any    `json:"code"`
	}
	if json.Unmarshal(body, &payload) == nil && payload.Message != "" {
		message = payload.Message
	}
	if message == "" {
		message = status
	}
	apiError := APIError{StatusCode: statusCode, Message: message}
	if payload.Code != nil {
		apiError.Code = fmt.Sprint(payload.Code)
	}
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return &AuthenticationError{APIError: apiError}
	case http.StatusNotFound:
		return &SandboxNotFoundError{APIError: apiError}
	case http.StatusRequestEntityTooLarge:
		return &NotEnoughSpaceError{APIError: apiError}
	case http.StatusTooManyRequests:
		return &RateLimitError{APIError: apiError}
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return &TimeoutError{APIError: apiError}
	default:
		return &SandboxError{APIError: apiError}
	}
}

func isConnectionError(err error) bool {
	var networkError net.Error
	return errors.As(err, &networkError)
}
