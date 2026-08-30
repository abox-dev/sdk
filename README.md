<p align="center">
  <img src="readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240">
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

Install JavaScript dependencies with `pnpm install --frozen-lockfile` and Python dependencies with `uv sync --frozen` in each Python package. `make generate` regenerates clients only from the checked-in snapshots under `spec/`; maintainers update them with `make sync-specs MONO_DIR=/path/to/mono`.

Run the complete local verification before publishing:

```bash
pnpm format:check
pnpm lint
pnpm typecheck
pnpm test
make generate
./scripts/build-release.sh
node scripts/check-public-artifacts.mjs release
./scripts/test-release-artifacts.sh
```

With the AgentBox KVM development contour running, test those exact release artifacts with:

```bash
AGENTBOX_API_KEY=ab_... \
AGENTBOX_API_URL=http://localhost:3000 \
AGENTBOX_SANDBOX_URL=http://localhost:3002 \
./scripts/test-release-runtime.sh
```

The runtime smoke covers the JavaScript and Python sync/async SDKs, both Code Interpreter packages, the CLI lifecycle, private traffic routing, and a real temporary template build. It removes the test sandboxes and template when it finishes.

After publishing, run the same checks from clean npm and PyPI installs with
`./scripts/test-published-runtime.sh`. It defaults to version `0.1.0` and does
not use workspace links or local package artifacts.

AgentBox SDK is derived from upstream work described in [UPSTREAM.md](UPSTREAM.md). Licensing notices are in [LICENSE](LICENSE) and [NOTICE](NOTICE).
