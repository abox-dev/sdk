package agentbox

import (
	"cmp"
	"errors"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultDomain         = "agentbox-runtime.ru"
	defaultRequestTimeout = 60 * time.Second
	defaultSandboxTimeout = 5 * time.Minute
	envdPort              = 49983
)

// ClientOption configures a Client.
type ClientOption func(*clientConfig) error

type clientConfig struct {
	apiKey         string
	domain         string
	apiURL         string
	sandboxURL     string
	debug          bool
	requestTimeout time.Duration
	httpClient     *http.Client
	proxyURL       *url.URL
	headers        http.Header
	logger         *slog.Logger
}

func defaultClientConfig() (clientConfig, error) {
	domain := cmp.Or(os.Getenv("AGENTBOX_DOMAIN"), defaultDomain)
	debugMode, err := strconv.ParseBool(cmp.Or(os.Getenv("AGENTBOX_DEBUG"), "false"))
	if err != nil {
		return clientConfig{}, &InvalidArgumentError{Message: "AGENTBOX_DEBUG must be true or false", Cause: err}
	}
	apiURL := os.Getenv("AGENTBOX_API_URL")
	if apiURL == "" {
		if debugMode {
			apiURL = "http://localhost:3000"
		} else {
			apiURL = "https://api." + domain
		}
	}
	return clientConfig{
		apiKey:         os.Getenv("AGENTBOX_API_KEY"),
		domain:         domain,
		apiURL:         apiURL,
		sandboxURL:     os.Getenv("AGENTBOX_SANDBOX_URL"),
		debug:          debugMode,
		requestTimeout: defaultRequestTimeout,
		headers:        make(http.Header),
	}, nil
}

// WithAPIKey sets the AgentBox API key. By default AGENTBOX_API_KEY is used.
func WithAPIKey(apiKey string) ClientOption {
	return func(config *clientConfig) error {
		config.apiKey = apiKey
		return nil
	}
}

// WithDomain sets the AgentBox runtime domain.
func WithDomain(domain string) ClientOption {
	return func(config *clientConfig) error {
		if strings.TrimSpace(domain) == "" {
			return &InvalidArgumentError{Message: "domain cannot be empty"}
		}
		config.domain = domain
		return nil
	}
}

// WithAPIURL overrides the control-plane API URL.
func WithAPIURL(value string) ClientOption {
	return func(config *clientConfig) error {
		if _, err := parseHTTPURL("API URL", value); err != nil {
			return err
		}
		config.apiURL = strings.TrimRight(value, "/")
		return nil
	}
}

// WithSandboxURL overrides the sandbox proxy URL.
func WithSandboxURL(value string) ClientOption {
	return func(config *clientConfig) error {
		if _, err := parseHTTPURL("sandbox URL", value); err != nil {
			return err
		}
		config.sandboxURL = strings.TrimRight(value, "/")
		return nil
	}
}

// WithDebug enables local envd routing.
func WithDebug(enabled bool) ClientOption {
	return func(config *clientConfig) error {
		config.debug = enabled
		return nil
	}
}

// WithRequestTimeout sets the default unary request timeout. Zero disables it.
func WithRequestTimeout(timeout time.Duration) ClientOption {
	return func(config *clientConfig) error {
		if timeout < 0 {
			return &InvalidArgumentError{Message: "request timeout cannot be negative"}
		}
		config.requestTimeout = timeout
		return nil
	}
}

// WithHTTPClient supplies the HTTP client used for every request. Its Timeout
// applies to complete streaming requests as well as unary requests; leave it at
// zero and use contexts or operation options for long-lived streams.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(config *clientConfig) error {
		if client == nil {
			return &InvalidArgumentError{Message: "HTTP client cannot be nil"}
		}
		config.httpClient = client
		return nil
	}
}

// WithProxy sets an HTTP or HTTPS proxy for SDK requests.
func WithProxy(value string) ClientOption {
	return func(config *clientConfig) error {
		proxyURL, err := parseHTTPURL("proxy URL", value)
		if err != nil {
			return err
		}
		config.proxyURL = proxyURL
		return nil
	}
}

// WithHeaders adds headers to control-plane requests.
func WithHeaders(headers http.Header) ClientOption {
	return func(config *clientConfig) error {
		maps.Copy(config.headers, headers)
		return nil
	}
}

// WithLogger enables structured request and lifecycle logging.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(config *clientConfig) error {
		config.logger = logger
		return nil
	}
}

func parseHTTPURL(name, value string) (*url.URL, error) {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return nil, &InvalidArgumentError{Message: "invalid " + name + ": " + value, Cause: err}
	}
	return parsed, nil
}

func applyOptions(options []ClientOption) (clientConfig, error) {
	config, err := defaultClientConfig()
	if err != nil {
		return clientConfig{}, err
	}
	for _, option := range options {
		if option == nil {
			return clientConfig{}, errors.New("agentbox: nil client option")
		}
		if err := option(&config); err != nil {
			return clientConfig{}, err
		}
	}
	if config.apiURL == "https://api."+defaultDomain && config.domain != defaultDomain {
		config.apiURL = "https://api." + config.domain
	}
	if config.httpClient != nil && config.proxyURL != nil {
		return clientConfig{}, &InvalidArgumentError{Message: "WithHTTPClient and WithProxy cannot be combined"}
	}
	return config, nil
}
