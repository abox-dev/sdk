// Package codeinterpreter runs Python, JavaScript, and TypeScript in AgentBox.
package codeinterpreter

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/abox-dev/sdk/packages/go-sdk"
)

const (
	DefaultTemplate         = "code-interpreter"
	JupyterPort             = 49999
	defaultExecutionTimeout = time.Minute
)

// Client wraps the core client with Code Interpreter creation helpers.
type Client struct{ Core *agentbox.Client }

func NewClient(options ...agentbox.ClientOption) (*Client, error) {
	client, err := agentbox.NewClient(options...)
	if err != nil {
		return nil, err
	}
	return &Client{Core: client}, nil
}

// Sandbox is a core sandbox with notebook-kernel APIs.
type Sandbox struct{ *agentbox.Sandbox }

func (client *Client) Create(ctx context.Context, options *agentbox.CreateSandboxOptions) (*Sandbox, error) {
	if options == nil {
		options = &agentbox.CreateSandboxOptions{}
	} else {
		cloned := *options
		options = &cloned
	}
	if options.Template == "" {
		options.Template = DefaultTemplate
	}
	sandbox, err := client.Core.Sandboxes.Create(ctx, options)
	if err != nil {
		return nil, err
	}
	return &Sandbox{Sandbox: sandbox}, nil
}
func (client *Client) Connect(ctx context.Context, id string, options *agentbox.ConnectSandboxOptions) (*Sandbox, error) {
	sandbox, err := client.Core.Sandboxes.Connect(ctx, id, options)
	if err != nil {
		return nil, err
	}
	return &Sandbox{Sandbox: sandbox}, nil
}

// Language is a Code Interpreter runtime language.
type Language string

const (
	Python     Language = "python"
	JavaScript Language = "javascript"
	TypeScript Language = "typescript"
)

// Context identifies a persistent kernel context.
type Context struct {
	ID       string `json:"id"`
	Language string `json:"language"`
	Cwd      string `json:"cwd"`
}
type CreateContextOptions struct {
	Language       Language      `json:"language,omitempty"`
	Cwd            string        `json:"cwd,omitempty"`
	RequestTimeout time.Duration `json:"-"`
}

// RunCodeOptions configures code execution and streaming callbacks.
type RunCodeOptions struct {
	Language         Language
	Context          *Context
	Env              map[string]string
	RequestTimeout   time.Duration
	ExecutionTimeout time.Duration
	OnStdout         func(OutputMessage)
	OnStderr         func(OutputMessage)
	OnResult         func(Result)
	OnError          func(ExecutionError)
}

// RunCode executes source code and collects streamed NDJSON output.
func (sandbox *Sandbox) RunCode(ctx context.Context, code string, options *RunCodeOptions) (*Execution, error) {
	if code == "" {
		return nil, &agentbox.InvalidArgumentError{Message: "code cannot be empty"}
	}
	if options == nil {
		options = &RunCodeOptions{}
	}
	if options.Context != nil && options.Language != "" {
		return nil, &agentbox.InvalidArgumentError{Message: "provide context or language, not both"}
	}
	if options.RequestTimeout < 0 || options.ExecutionTimeout < 0 {
		return nil, &agentbox.InvalidArgumentError{Message: "timeouts cannot be negative"}
	}
	payload := struct {
		Code      string            `json:"code"`
		ContextID string            `json:"context_id,omitempty"`
		Language  Language          `json:"language,omitempty"`
		Env       map[string]string `json:"env_vars,omitempty"`
	}{Code: code, Language: options.Language, Env: options.Env}
	if options.Context != nil {
		payload.ContextID = options.Context.ID
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	requestCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	requestTimeout := options.RequestTimeout
	if requestTimeout == 0 {
		requestTimeout = sandbox.RequestTimeout()
	}
	var requestTimer *time.Timer
	if requestTimeout > 0 {
		requestTimer = time.AfterFunc(requestTimeout, cancel)
	}
	response, err := sandbox.Request(requestCtx, JupyterPort, http.MethodPost, "/execute", bytes.NewReader(body), true)
	if requestTimer != nil && !requestTimer.Stop() && err == nil {
		response.Body.Close()
		return nil, &agentbox.TimeoutError{APIError: agentbox.APIError{Message: "code execution request timed out", Cause: context.DeadlineExceeded}}
	}
	if err != nil {
		if errors.Is(requestCtx.Err(), context.Canceled) && ctx.Err() == nil {
			return nil, &agentbox.TimeoutError{APIError: agentbox.APIError{Message: "code execution request timed out", Cause: context.DeadlineExceeded}}
		}
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, interpreterHTTPError(response)
	}
	executionTimeout := options.ExecutionTimeout
	if executionTimeout == 0 {
		executionTimeout = defaultExecutionTimeout
	}
	executionTimer := time.AfterFunc(executionTimeout, cancel)
	defer executionTimer.Stop()
	execution := &Execution{Logs: Logs{Stdout: []string{}, Stderr: []string{}}}
	scanner := bufio.NewScanner(response.Body)
	scanner.Buffer(make([]byte, 64*1024), 16*1024*1024)
	for scanner.Scan() {
		if err := parseOutput(scanner.Bytes(), execution, options); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		if requestCtx.Err() != nil && ctx.Err() == nil {
			return nil, &agentbox.TimeoutError{APIError: agentbox.APIError{Message: "code execution timed out", Cause: context.DeadlineExceeded}}
		}
		return nil, fmt.Errorf("codeinterpreter: read output: %w", err)
	}
	return execution, nil
}

func (sandbox *Sandbox) CreateContext(ctx context.Context, options *CreateContextOptions) (*Context, error) {
	if options == nil {
		options = &CreateContextOptions{}
	}
	var result Context
	if err := sandbox.contextRequest(ctx, http.MethodPost, "/contexts", options, options.RequestTimeout, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
func (sandbox *Sandbox) ListContexts(ctx context.Context) ([]Context, error) {
	var result []Context
	if err := sandbox.contextRequest(ctx, http.MethodGet, "/contexts", nil, 0, &result); err != nil {
		return nil, err
	}
	return result, nil
}
func (sandbox *Sandbox) RemoveContext(ctx context.Context, contextID string) error {
	if strings.TrimSpace(contextID) == "" {
		return &agentbox.InvalidArgumentError{Message: "context ID cannot be empty"}
	}
	return sandbox.contextRequest(ctx, http.MethodDelete, "/contexts/"+url.PathEscape(contextID), nil, 0, nil)
}
func (sandbox *Sandbox) RestartContext(ctx context.Context, contextID string) error {
	if strings.TrimSpace(contextID) == "" {
		return &agentbox.InvalidArgumentError{Message: "context ID cannot be empty"}
	}
	return sandbox.contextRequest(ctx, http.MethodPost, "/contexts/"+url.PathEscape(contextID)+"/restart", nil, 0, nil)
}
func (sandbox *Sandbox) contextRequest(ctx context.Context, method, path string, body any, timeout time.Duration, output any) error {
	var reader io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(encoded)
	}
	if timeout == 0 {
		timeout = sandbox.RequestTimeout()
	}
	var requestCtx context.Context
	var cancel context.CancelFunc
	if timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		requestCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()
	response, err := sandbox.Request(requestCtx, JupyterPort, method, path, reader, true)
	if err != nil {
		if errors.Is(requestCtx.Err(), context.DeadlineExceeded) {
			return &agentbox.TimeoutError{APIError: agentbox.APIError{Message: "context request timed out", Cause: requestCtx.Err()}}
		}
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return interpreterHTTPError(response)
	}
	if output != nil {
		if err := json.NewDecoder(response.Body).Decode(output); err != nil {
			return fmt.Errorf("codeinterpreter: decode response: %w", err)
		}
	}
	return nil
}

func interpreterHTTPError(response *http.Response) error {
	data, _ := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	message := strings.TrimSpace(string(data))
	if message == "" {
		message = response.Status
	}
	apiError := agentbox.APIError{StatusCode: response.StatusCode, Message: message}
	switch response.StatusCode {
	case http.StatusBadGateway, http.StatusGatewayTimeout:
		return &agentbox.TimeoutError{APIError: apiError}
	default:
		return &agentbox.SandboxError{APIError: apiError}
	}
}
