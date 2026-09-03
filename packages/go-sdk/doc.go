// Package agentbox provides the official Go client for AgentBox sandboxes.
//
// Create a client, start a sandbox, and run a command:
//
//	client, err := agentbox.NewClient()
//	if err != nil {
//		log.Fatal(err)
//	}
//	sandbox, err := client.Sandboxes.Create(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer sandbox.Kill(context.Background())
//	result, err := sandbox.Commands.Run(ctx, "echo", &agentbox.CommandOptions{
//		Args: []string{"Hello from AgentBox"},
//	})
package agentbox
