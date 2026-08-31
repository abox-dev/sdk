<p align="center"><a href="https://agentbox.ru"><img src="https://raw.githubusercontent.com/abox-dev/sdk/main/readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240"></a></p>

# `@abox-dev/code-interpreter`

Stateful Python, JavaScript, and TypeScript execution in an AgentBox sandbox.

```bash
npm install @abox-dev/code-interpreter
```

Create an API key using the [API key guide](https://docs.agentbox.ru/en/quickstart/api-key/) and set `AGENTBOX_API_KEY`.

```ts
import { Sandbox } from '@abox-dev/code-interpreter'

const sandbox = await Sandbox.create()
try {
  const execution = await sandbox.runCode('x = 1; x += 1; x')
  console.log(execution.text)
} finally {
  await sandbox.kill()
}
```

See the [Code Interpreter guide](https://docs.agentbox.ru/en/sdk/code-interpreter/) and [API reference](https://docs.agentbox.ru/en/sdk/api-reference/javascript/code-interpreter/).
