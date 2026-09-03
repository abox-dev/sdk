package agentbox_test

import (
	"context"
	"log"

	"github.com/abox-dev/sdk/packages/go-sdk"
)

func ExampleClient() {
	client, err := agentbox.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	sandbox, err := client.Sandboxes.Create(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer sandbox.Kill(context.Background())
	_, _ = sandbox.Commands.Run(context.Background(), "echo", &agentbox.CommandOptions{Args: []string{"hello"}})
}

func ExampleTemplateBuilder() {
	template := agentbox.NewTemplate(".").FromPython("3.13").Copy("requirements.txt", "/app/", nil).PipInstall().Workdir("/app")
	_, _ = template.JSON()
}
