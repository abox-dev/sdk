<p align="center">
  <a href="https://agentbox.ru">
    <img src="readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240">
  </a>
</p>

# AgentBox SDK

Official JavaScript, Python, and CLI clients for running isolated AgentBox sandboxes and code interpreters.

| Package                     | Install                                  | Import                       |
| --------------------------- | ---------------------------------------- | ---------------------------- |
| JavaScript SDK              | `npm install @abox-dev/sdk`              | `@abox-dev/sdk`              |
| Python SDK                  | `pip install abox-sdk`                   | `agentbox`                   |
| JavaScript Code Interpreter | `npm install @abox-dev/code-interpreter` | `@abox-dev/code-interpreter` |
| Python Code Interpreter     | `pip install abox-code-interpreter`      | `agentbox_code_interpreter`  |
| CLI                         | `npm install --global @abox-dev/cli`     | `agentbox`                   |

## Quick start

Create an API key as described in the [API key guide](https://docs.agentbox.ru/en/quickstart/api-key/) and export it:

```bash
export AGENTBOX_API_KEY=ab_...
```

JavaScript:

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

Python:

```python
from agentbox import Sandbox

with Sandbox.create() as sandbox:
    result = sandbox.commands.run('echo "Hello from AgentBox"')
    print(result.stdout)
```

## Configuration

The SDKs use `AGENTBOX_API_KEY` and optionally `AGENTBOX_PROJECT_ID`, `AGENTBOX_DOMAIN`, `AGENTBOX_API_URL`, and `AGENTBOX_SANDBOX_URL`. The production defaults are `agentbox-runtime.ru`, `api.agentbox-runtime.ru`, and `sandbox.agentbox-runtime.ru`.

The CLI stores local configuration in `~/.agentbox/config.json`; environment variables take precedence. See [CLI configuration](https://docs.agentbox.ru/en/cli/configuration/).

## Documentation

- [Core SDK](https://docs.agentbox.ru/en/sdk/)
- [Sandboxes](https://docs.agentbox.ru/en/sdk/sandboxes/)
- [Commands](https://docs.agentbox.ru/en/sdk/commands/)
- [Files](https://docs.agentbox.ru/en/sdk/files/)
- [Templates and builds](https://docs.agentbox.ru/en/sdk/templates/)
- [Code Interpreter](https://docs.agentbox.ru/en/sdk/code-interpreter/)
- [CLI](https://docs.agentbox.ru/en/cli/)
- [API reference](https://docs.agentbox.ru/en/sdk/api-reference/)
- [Examples](https://docs.agentbox.ru/en/examples/)

## Development

Contributor setup, spec synchronization, verification, KVM testing, versioning,
and publication are documented in [RELEASING.md](RELEASING.md).

AgentBox SDK is derived from upstream work described in [UPSTREAM.md](UPSTREAM.md). Licensing notices are in [LICENSE](LICENSE) and [NOTICE](NOTICE).
