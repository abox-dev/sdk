package agentbox

import (
	"fmt"
	"net/http"

	api "github.com/abox-dev/sdk/packages/go-sdk/internal/gen/api"
)

// Client is an AgentBox control-plane client.
type Client struct {
	config     clientConfig
	httpClient *http.Client
	envdClient *http.Client
	api        *api.ClientWithResponses

	Sandboxes *SandboxService
	Templates *TemplateService
}

// NewClient creates a client. Configuration defaults to AGENTBOX_* environment
// variables and can be overridden with options.
func NewClient(options ...ClientOption) (*Client, error) {
	config, err := applyOptions(options)
	if err != nil {
		return nil, err
	}
	httpClient := newHTTPClient(config)
	envdClient := newEnvdHTTPClient(config, httpClient)
	generated, err := api.NewClientWithResponses(config.apiURL, api.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("agentbox: initialize API client: %w", err)
	}
	client := &Client{config: config, httpClient: httpClient, envdClient: envdClient, api: generated}
	client.Sandboxes = &SandboxService{client: client}
	client.Templates = &TemplateService{client: client}
	return client, nil
}
