package agentbox

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
)

// CreateSandboxOptions configures a new sandbox.
type CreateSandboxOptions struct {
	Template            string
	Timeout             time.Duration
	AutoPause           *bool
	AutoPauseMemory     *bool
	AutoResume          *bool
	Secure              *bool
	AllowInternetAccess *bool
	Env                 map[string]string
	Metadata            map[string]string
	Network             *SandboxNetworkConfig
	IAM                 *SandboxIAM
}

// ConnectSandboxOptions configures connecting or resuming a sandbox.
type ConnectSandboxOptions struct{ Timeout time.Duration }

// ListSandboxOptions filters one sandbox list page.
type ListSandboxOptions struct {
	Metadata  map[string]string
	States    []SandboxState
	NextToken string
	Limit     int
}

// MetricsOptions selects a metrics interval.
type MetricsOptions struct{ Start, End time.Time }

// SandboxLogOptions filters one page of sandbox logs.
type SandboxLogOptions struct {
	Cursor    string
	Timestamp int64
	Limit     int
	Direction string
	Level     string
	Search    string
}

// PauseOptions configures snapshot behavior while pausing.
type PauseOptions struct{ Memory *bool }

// ForkOptions configures sandbox forks.
type ForkOptions struct {
	Count   int
	Timeout time.Duration
}

// ForkResult is one ordered fork outcome. Exactly one of Sandbox or Err is set.
type ForkResult struct {
	Sandbox *Sandbox
	Err     error
}

// SnapshotListOptions filters and paginates snapshots.
type SnapshotListOptions struct {
	SandboxID string
	Name      string
	NextToken string
	Limit     int
}

// SandboxRequestOptions configures a request to a service inside a sandbox.
type SandboxRequestOptions struct {
	Direct      bool
	Headers     http.Header
	ContentType string
}

// SandboxService manages sandboxes owned by a client.
type SandboxService struct{ client *Client }

// Sandbox is a connected AgentBox sandbox.
type Sandbox struct {
	client             *Client
	ID                 string
	TemplateID         string
	Alias              string
	Domain             string
	EnvdVersion        string
	envdAccessToken    string
	trafficAccessToken string

	Commands *CommandService
	PTY      *PTYService
	Files    *FileService
}

// RequestTimeout returns the default unary request timeout configured on the client.
func (sandbox *Sandbox) RequestTimeout() time.Duration { return sandbox.client.config.requestTimeout }

func (sandbox *Sandbox) unaryContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return withRequestTimeout(ctx, sandbox.RequestTimeout())
}

// String returns a credential-free sandbox description.
func (sandbox Sandbox) String() string {
	return fmt.Sprintf("Sandbox{ID:%q, TemplateID:%q, Alias:%q, Domain:%q, EnvdVersion:%q}", sandbox.ID, sandbox.TemplateID, sandbox.Alias, sandbox.Domain, sandbox.EnvdVersion)
}

// GoString returns a credential-free sandbox description for %#v formatting.
func (sandbox Sandbox) GoString() string { return sandbox.String() }

// Create starts and connects to a sandbox.
func (service *SandboxService) Create(ctx context.Context, options *CreateSandboxOptions) (*Sandbox, error) {
	if options == nil {
		options = &CreateSandboxOptions{}
	}
	if options.IAM != nil && options.IAM.Tokens != nil {
		for name := range *options.IAM.Tokens {
			if err := ValidateIAMTokenName(name); err != nil {
				return nil, err
			}
		}
	}
	template := options.Template
	if template == "" {
		template = "base"
	}
	timeout, err := durationSeconds(options.Timeout, defaultSandboxTimeout)
	if err != nil {
		return nil, err
	}
	var network *api.SandboxNetworkConfig
	if options.Network != nil {
		converted, err := convertModel[api.SandboxNetworkConfig](*options.Network)
		if err != nil {
			return nil, err
		}
		network = &converted
	}
	var iam *api.SandboxIam
	if options.IAM != nil {
		converted, err := convertModel[api.SandboxIam](*options.IAM)
		if err != nil {
			return nil, err
		}
		iam = &converted
	}
	body := api.NewSandbox{
		TemplateID: template, Timeout: &timeout, AutoPause: options.AutoPause,
		AutoPauseMemory: options.AutoPauseMemory, Secure: options.Secure,
		AllowInternetAccess: options.AllowInternetAccess, Network: network,
		Iam: iam,
	}
	if options.AutoResume != nil {
		body.AutoResume = &api.SandboxAutoResumeConfig{Enabled: *options.AutoResume}
	}
	if options.Env != nil {
		env := api.EnvVars(options.Env)
		body.EnvVars = &env
	}
	if options.Metadata != nil {
		metadata := api.SandboxMetadata(options.Metadata)
		body.Metadata = &metadata
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.PostSandboxesWithResponse(requestCtx, body)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON201 == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return service.client.sandboxFromAPI(*response.JSON201), nil
}

// Connect connects to or resumes a sandbox.
func (service *SandboxService) Connect(ctx context.Context, id string, options *ConnectSandboxOptions) (*Sandbox, error) {
	if strings.TrimSpace(id) == "" {
		return nil, &InvalidArgumentError{Message: "sandbox ID cannot be empty"}
	}
	timeout := defaultSandboxTimeout
	if options != nil && options.Timeout != 0 {
		timeout = options.Timeout
	}
	seconds, err := durationSeconds(timeout, defaultSandboxTimeout)
	if err != nil {
		return nil, err
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.PostSandboxesSandboxIDConnectWithResponse(requestCtx, id, api.ConnectSandbox{Timeout: seconds})
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	value := response.JSON200
	if value == nil {
		value = response.JSON201
	}
	if value == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return service.client.sandboxFromAPI(*value), nil
}

// List returns one list page. The continuation token is read from X-Next-Token.
func (service *SandboxService) List(ctx context.Context, options *ListSandboxOptions) (Page[ListedSandbox], error) {
	parameters := &api.GetV2SandboxesParams{}
	if options != nil {
		if len(options.Metadata) > 0 {
			values := url.Values{}
			for key, value := range options.Metadata {
				values.Set(key, value)
			}
			encoded := values.Encode()
			parameters.Metadata = &encoded
		}
		if len(options.States) > 0 {
			states := make([]api.SandboxState, len(options.States))
			for index, state := range options.States {
				states[index] = api.SandboxState(state)
			}
			parameters.State = &states
		}
		if options.NextToken != "" {
			parameters.NextToken = &options.NextToken
		}
		if options.Limit > 0 {
			limit := int32(options.Limit)
			parameters.Limit = &limit
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetV2SandboxesWithResponse(requestCtx, parameters)
	if err != nil {
		return Page[ListedSandbox]{}, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return Page[ListedSandbox]{}, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	items, err := convertModel[[]ListedSandbox](*response.JSON200)
	return Page[ListedSandbox]{Items: items, NextToken: response.HTTPResponse.Header.Get("X-Next-Token")}, err
}

// Info returns current sandbox state and configuration.
func (service *SandboxService) Info(ctx context.Context, id string) (*SandboxInfo, error) {
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetSandboxesSandboxIDWithResponse(requestCtx, id)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	result, err := convertModel[SandboxInfo](*response.JSON200)
	return &result, err
}

// Kill permanently stops a sandbox. It returns false when it did not exist.
func (service *SandboxService) Kill(ctx context.Context, id string) (bool, error) {
	if service.client.config.debug {
		return true, nil
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.DeleteSandboxesSandboxID(requestCtx, id)
	if err != nil {
		return false, normalizeRequestError(err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if response.StatusCode != http.StatusNoContent {
		return false, decodeHTTPError(response)
	}
	return true, nil
}

// Snapshots returns one page of snapshots.
func (service *SandboxService) Snapshots(ctx context.Context, options *SnapshotListOptions) (Page[SnapshotInfo], error) {
	parameters := &api.GetSnapshotsParams{}
	if options != nil {
		if options.SandboxID != "" {
			parameters.SandboxID = &options.SandboxID
		}
		if options.Name != "" {
			parameters.Name = &options.Name
		}
		if options.NextToken != "" {
			parameters.NextToken = &options.NextToken
		}
		if options.Limit > 0 {
			value := int32(options.Limit)
			parameters.Limit = &value
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetSnapshotsWithResponse(requestCtx, parameters)
	if err != nil {
		return Page[SnapshotInfo]{}, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return Page[SnapshotInfo]{}, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	items, err := convertModel[[]SnapshotInfo](*response.JSON200)
	return Page[SnapshotInfo]{Items: items, NextToken: response.HTTPResponse.Header.Get("X-Next-Token")}, err
}

// DeleteSnapshot deletes a snapshot. It returns false when it did not exist.
func (service *SandboxService) DeleteSnapshot(ctx context.Context, snapshotID string) (bool, error) {
	if strings.TrimSpace(snapshotID) == "" {
		return false, &InvalidArgumentError{Message: "snapshot ID cannot be empty"}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.DeleteTemplatesTemplateIDWithResponse(requestCtx, snapshotID)
	if err != nil {
		return false, normalizeRequestError(err)
	}
	if response.StatusCode() == http.StatusNotFound {
		return false, nil
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return false, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return true, nil
}

// Metrics returns the latest metrics for the requested sandbox IDs.
func (service *SandboxService) Metrics(ctx context.Context, sandboxIDs ...string) (map[string]SandboxMetric, error) {
	if len(sandboxIDs) == 0 {
		return nil, &InvalidArgumentError{Message: "at least one sandbox ID is required"}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetSandboxesMetricsWithResponse(requestCtx, &api.GetSandboxesMetricsParams{SandboxIds: sandboxIDs})
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return convertModel[map[string]SandboxMetric](response.JSON200.Sandboxes)
}

// Logs returns one page of structured sandbox logs.
func (service *SandboxService) Logs(ctx context.Context, id string, options *SandboxLogOptions) (Page[SandboxLogEntry], error) {
	params := &api.GetV2SandboxesSandboxIDLogsParams{}
	if options != nil {
		if options.Cursor != "" {
			params.PageCursor = &options.Cursor
		}
		if options.Timestamp != 0 {
			params.Cursor = &options.Timestamp
		}
		if options.Limit > 0 {
			limit := int32(options.Limit)
			params.Limit = &limit
		}
		if options.Direction != "" {
			value := api.LogsDirection(options.Direction)
			params.Direction = &value
		}
		if options.Level != "" {
			value := api.LogLevel(options.Level)
			params.Level = &value
		}
		if options.Search != "" {
			params.Search = &options.Search
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, service.client.config.requestTimeout)
	defer cancel()
	response, err := service.client.api.GetV2SandboxesSandboxIDLogsWithResponse(requestCtx, id, params)
	if err != nil {
		return Page[SandboxLogEntry]{}, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return Page[SandboxLogEntry]{}, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	items, err := convertModel[[]SandboxLogEntry](response.JSON200.Logs)
	if err != nil {
		return Page[SandboxLogEntry]{}, err
	}
	page := Page[SandboxLogEntry]{Items: items}
	if response.JSON200.NextCursor != nil {
		page.NextToken = *response.JSON200.NextCursor
	}
	return page, nil
}

// Info returns information about this sandbox.
func (sandbox *Sandbox) Info(ctx context.Context) (*SandboxInfo, error) {
	return sandbox.client.Sandboxes.Info(ctx, sandbox.ID)
}

// Logs returns one page of sandbox logs.
func (sandbox *Sandbox) Logs(ctx context.Context, options *SandboxLogOptions) (Page[SandboxLogEntry], error) {
	return sandbox.client.Sandboxes.Logs(ctx, sandbox.ID, options)
}

// Kill permanently stops this sandbox. It returns false when it was not found.
func (sandbox *Sandbox) Kill(ctx context.Context) (bool, error) {
	return sandbox.client.Sandboxes.Kill(ctx, sandbox.ID)
}

// IsRunning reports whether envd is reachable. A 502 response means the
// sandbox is no longer running.
func (sandbox *Sandbox) IsRunning(ctx context.Context) (bool, error) {
	requestCtx, cancel := sandbox.unaryContext(ctx)
	defer cancel()
	response, err := sandbox.Request(requestCtx, envdPort, http.MethodGet, "/health", nil, false)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadGateway {
		return false, nil
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return false, decodeHTTPError(response)
	}
	return true, nil
}

// SetTimeout changes the sandbox expiration timeout from now.
func (sandbox *Sandbox) SetTimeout(ctx context.Context, timeout time.Duration) error {
	seconds, err := durationSeconds(timeout, 0)
	if err != nil {
		return err
	}
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.PostSandboxesSandboxIDTimeoutWithResponse(requestCtx, sandbox.ID, api.SandboxTimeoutRequest{Timeout: seconds})
	if err != nil {
		return normalizeRequestError(err)
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return nil
}

// KeepAlive extends the sandbox lifetime. A zero duration uses the server default.
func (sandbox *Sandbox) KeepAlive(ctx context.Context, duration time.Duration) error {
	body := api.SandboxRefreshRequest{}
	if duration != 0 {
		seconds, err := durationSecondsInt(duration)
		if err != nil {
			return err
		}
		body.Duration = &seconds
	}
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.PostSandboxesSandboxIDRefreshesWithResponse(requestCtx, sandbox.ID, body)
	if err != nil {
		return normalizeRequestError(err)
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return nil
}

// Pause pauses the sandbox, optionally retaining memory.
func (sandbox *Sandbox) Pause(ctx context.Context, options *PauseOptions) error {
	body := api.SandboxPauseRequest{}
	if options != nil {
		body.Memory = options.Memory
	}
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.PostSandboxesSandboxIDPauseWithResponse(requestCtx, sandbox.ID, body)
	if err != nil {
		return normalizeRequestError(err)
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return nil
}

// UpdateNetwork atomically replaces sandbox egress rules.
func (sandbox *Sandbox) UpdateNetwork(ctx context.Context, config SandboxNetworkConfig, allowInternetAccess *bool) error {
	body, err := convertModel[api.SandboxNetworkUpdateConfig](config)
	if err != nil {
		return err
	}
	body.AllowInternetAccess = allowInternetAccess
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.PutSandboxesSandboxIDNetworkWithResponse(requestCtx, sandbox.ID, body)
	if err != nil {
		return normalizeRequestError(err)
	}
	if response.StatusCode() < 200 || response.StatusCode() >= 300 {
		return decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return nil
}

// Metrics returns sandbox metrics for the requested interval.
func (sandbox *Sandbox) Metrics(ctx context.Context, options *MetricsOptions) ([]SandboxMetric, error) {
	parameters := &api.GetSandboxesSandboxIDMetricsParams{}
	if options != nil {
		if !options.Start.IsZero() {
			v := options.Start.Unix()
			parameters.Start = &v
		}
		if !options.End.IsZero() {
			v := options.End.Unix()
			parameters.End = &v
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.GetSandboxesSandboxIDMetricsWithResponse(requestCtx, sandbox.ID, parameters)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON200 == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	return convertModel[[]SandboxMetric](*response.JSON200)
}

// Fork creates one or more sandboxes from this sandbox's current state.
func (sandbox *Sandbox) Fork(ctx context.Context, options *ForkOptions) ([]ForkResult, error) {
	body := api.SandboxForkRequest{}
	if options != nil {
		if options.Count > 0 {
			count := int32(options.Count)
			body.Count = &count
		}
		if options.Timeout != 0 {
			seconds, err := durationSeconds(options.Timeout, 0)
			if err != nil {
				return nil, err
			}
			body.Timeout = &seconds
		}
	}
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.PostSandboxesSandboxIDForkWithResponse(requestCtx, sandbox.ID, body)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON201 == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	result := make([]ForkResult, 0, len(*response.JSON201))
	for _, item := range *response.JSON201 {
		if item.Sandbox != nil {
			result = append(result, ForkResult{Sandbox: sandbox.client.sandboxFromAPI(*item.Sandbox)})
			continue
		}
		if item.Error != nil {
			result = append(result, ForkResult{Err: decodeStatusError(int(item.Error.Code), strconv.Itoa(int(item.Error.Code)), []byte(item.Error.Message))})
			continue
		}
		result = append(result, ForkResult{Err: &SandboxError{APIError: APIError{Message: "fork result contained neither sandbox nor error"}}})
	}
	return result, nil
}

// CreateSnapshot stores this sandbox as a template snapshot.
func (sandbox *Sandbox) CreateSnapshot(ctx context.Context, name string) (*SnapshotInfo, error) {
	body := api.SandboxSnapshotRequest{}
	if name != "" {
		body.Name = &name
	}
	requestCtx, cancel := withRequestTimeout(ctx, sandbox.client.config.requestTimeout)
	defer cancel()
	response, err := sandbox.client.api.PostSandboxesSandboxIDSnapshotsWithResponse(requestCtx, sandbox.ID, body)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	if response.JSON201 == nil {
		return nil, decodeStatusError(response.StatusCode(), response.Status(), response.Body)
	}
	result, err := convertModel[SnapshotInfo](*response.JSON201)
	return &result, err
}

// Host returns the public hostname for a port exposed by the sandbox.
func (sandbox *Sandbox) Host(port int) string {
	return fmt.Sprintf("%d-%s.%s", port, sandbox.ID, sandbox.Domain)
}

// Request performs an authenticated request against a service inside the sandbox.
// Direct bypasses the stable proxy hostname while retaining AgentBox routing headers.
func (sandbox *Sandbox) Request(ctx context.Context, port int, method, path string, body io.Reader, direct bool) (*http.Response, error) {
	return sandbox.RequestWithOptions(ctx, port, method, path, body, &SandboxRequestOptions{Direct: direct})
}

// RequestWithOptions performs an authenticated request with custom headers.
func (sandbox *Sandbox) RequestWithOptions(ctx context.Context, port int, method, path string, body io.Reader, options *SandboxRequestOptions) (*http.Response, error) {
	if options == nil {
		options = &SandboxRequestOptions{}
	}
	endpoint := strings.TrimRight(sandbox.envdURL(port, options.Direct), "/") + "/" + strings.TrimLeft(path, "/")
	request, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("agentbox: create sandbox request: %w", err)
	}
	request.Header = sandbox.envdHeaders(port)
	for key, values := range options.Headers {
		request.Header[key] = slices.Clone(values)
	}
	if options.ContentType != "" {
		request.Header.Set("Content-Type", options.ContentType)
	} else if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := sandbox.client.httpClient.Do(request)
	if err != nil {
		return nil, normalizeRequestError(err)
	}
	return response, nil
}

func (client *Client) sandboxFromAPI(value api.Sandbox) *Sandbox {
	domain := client.config.domain
	if value.Domain != nil && *value.Domain != "" {
		domain = *value.Domain
	}
	sandbox := &Sandbox{client: client, ID: value.SandboxID, TemplateID: value.TemplateID, EnvdVersion: value.EnvdVersion, Domain: domain}
	if value.Alias != nil {
		sandbox.Alias = *value.Alias
	}
	if value.EnvdAccessToken != nil {
		sandbox.envdAccessToken = *value.EnvdAccessToken
	}
	if value.TrafficAccessToken != nil {
		sandbox.trafficAccessToken = *value.TrafficAccessToken
	}
	sandbox.Commands = newCommandService(sandbox)
	sandbox.PTY = &PTYService{commands: sandbox.Commands}
	sandbox.Files = newFileService(sandbox)
	return sandbox
}

func (sandbox *Sandbox) envdURL(port int, direct bool) string {
	if sandbox.client.config.debug {
		return fmt.Sprintf("http://localhost:%d", port)
	}
	if sandbox.client.config.sandboxURL != "" {
		parsed, _ := url.Parse(sandbox.client.config.sandboxURL)
		if direct {
			return parsed.Scheme + "://" + fmt.Sprintf("%d-%s.%s", port, sandbox.ID, parsed.Host)
		}
		return sandbox.client.config.sandboxURL
	}
	if !direct && sandbox.Domain == defaultDomain {
		return "https://sandbox." + sandbox.Domain
	}
	return "https://" + sandbox.Host(port)
}

func (sandbox *Sandbox) envdHeaders(port int) http.Header {
	headers := make(http.Header)
	headers.Set("Agentbox-Sandbox-Id", sandbox.ID)
	headers.Set("Agentbox-Sandbox-Port", strconv.Itoa(port))
	if sandbox.envdAccessToken != "" {
		headers.Set("X-Access-Token", sandbox.envdAccessToken)
	}
	if sandbox.trafficAccessToken != "" {
		headers.Set("Agentbox-Traffic-Access-Token", sandbox.trafficAccessToken)
	}
	return headers
}

func (sandbox *Sandbox) resolveUser(user string) string {
	if user == "" && !envdAtLeast(sandbox.EnvdVersion, 0, 4, 0) {
		return "user"
	}
	return user
}

func envdAtLeast(version string, major, minor, patch int) bool {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if prerelease, _, found := strings.Cut(version, "-"); found {
		version = prerelease
	}
	parts := strings.Split(version, ".")
	if len(parts) == 0 || len(parts) > 3 {
		return true
	}
	actual := [3]int{}
	for index, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil || value < 0 {
			return true
		}
		actual[index] = value
	}
	required := [3]int{major, minor, patch}
	for index := range actual {
		if actual[index] != required[index] {
			return actual[index] > required[index]
		}
	}
	return true
}

func durationSeconds(value, fallback time.Duration) (int32, error) {
	if value == 0 {
		value = fallback
	}
	if value <= 0 || value > time.Duration(^uint32(0)>>1)*time.Second {
		return 0, &InvalidArgumentError{Message: "duration must be positive and fit in seconds"}
	}
	return int32(value / time.Second), nil
}
func durationSecondsInt(value time.Duration) (int, error) {
	seconds, err := durationSeconds(value, 0)
	return int(seconds), err
}
func normalizeRequestError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return &TimeoutError{APIError: APIError{Message: "request timed out", Cause: err}}
	}
	if isConnectionError(err) {
		return &SandboxError{APIError: APIError{Message: "connection failed", Cause: err}}
	}
	return err
}

func fileSignature(path, operation, user, token string, expiration time.Time) (string, int64, error) {
	if token == "" {
		return "", 0, &AuthenticationError{APIError: APIError{Message: "envd access token is required for signed URLs"}}
	}
	raw := path + ":" + operation + ":" + user + ":" + token
	unix := int64(0)
	if !expiration.IsZero() {
		unix = expiration.Unix()
		raw += ":" + strconv.FormatInt(unix, 10)
	}
	digest := sha256.Sum256([]byte(raw))
	return "v1_" + strings.TrimRight(base64.StdEncoding.EncodeToString(digest[:]), "="), unix, nil
}
