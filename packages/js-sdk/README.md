<p align="center"><a href="https://agentbox.ru"><img src="https://raw.githubusercontent.com/abox-dev/sdk/main/readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240"></a></p>

# `@abox-dev/sdk`

The official JavaScript and TypeScript SDK for AgentBox sandboxes and templates.

```bash
npm install @abox-dev/sdk
```

Create an API key using the [API key guide](https://docs.agentbox.ru/en/quickstart/api-key/) and set `AGENTBOX_API_KEY`.

```ts
import { Sandbox } from '@abox-dev/sdk'

const sandbox = await Sandbox.create()
try {
  const result = await sandbox.commands.run('echo "Hello from AgentBox"')
  console.log(result.stdout)
} finally {
  await sandbox.kill()
}
```

For explicit per-client configuration:

```ts
import { AgentBox } from '@abox-dev/sdk'

const client = new AgentBox({ apiKey: 'ab_...' })
const sandbox = await client.Sandbox.create()
```

Documentation: [core SDK](https://docs.agentbox.ru/en/sdk/), [sandboxes](https://docs.agentbox.ru/en/sdk/sandboxes/), [commands](https://docs.agentbox.ru/en/sdk/commands/), [files](https://docs.agentbox.ru/en/sdk/files/), [templates](https://docs.agentbox.ru/en/sdk/templates/), and [API reference](https://docs.agentbox.ru/en/sdk-reference/javascript/).
