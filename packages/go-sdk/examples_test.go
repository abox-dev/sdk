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

func ExampleCommandService_Start_streaming() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := agentbox.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	sandbox, err := client.Sandboxes.Create(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer sandbox.Kill(context.Background())
	handle, err := sandbox.Commands.Start(ctx, "cat", &agentbox.CommandOptions{
		Stdin:     true,
		Streaming: &agentbox.CommandStreamingOptions{},
		OnStdout:  func(chunk []byte) { log.Printf("output: %s", chunk) },
		// With no callback or channel selected, stderr is discarded.
	})
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	pid, err := handle.PID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = handle.Write(ctx, []byte("hello\n"))
	// Close only detaches. The process remains alive with stdin open.
	_ = handle.Close()
	_, _ = handle.Wait(ctx)
	attached, err := sandbox.Commands.ConnectWithOptions(ctx, pid, "", &agentbox.CommandConnectOptions{
		Streaming: &agentbox.CommandStreamingOptions{},
		OnStdout:  func(chunk []byte) { log.Printf("output: %s", chunk) },
	})
	if err != nil {
		log.Fatal(err)
	}
	defer attached.Close()
	_, _ = attached.Write(ctx, []byte("hello again\n"))
	_ = attached.CloseStdin(ctx)
	_, _ = attached.Wait(ctx)
}
