package consumer_test

import (
	"testing"

	"github.com/abox-dev/sdk/packages/go-sdk"
	"github.com/abox-dev/sdk/packages/go-sdk/codeinterpreter"
)

func TestPublicPackagesCompile(t *testing.T) {
	client, err := agentbox.NewClient(agentbox.WithAPIKey("test"))
	if err != nil {
		t.Fatal(err)
	}
	if client.Sandboxes == nil || client.Templates == nil || codeinterpreter.DefaultTemplate == "" {
		t.Fatal("public services are unavailable")
	}
	_ = agentbox.SandboxInfo{State: agentbox.SandboxRunning}
	_ = agentbox.SnapshotListOptions{Limit: 10}
	_ = agentbox.TemplateInfoOptions{Limit: 10}
	_ = agentbox.ForkResult{}
	_ = agentbox.SandboxRequestOptions{}
}
