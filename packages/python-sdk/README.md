<p align="center"><a href="https://agentbox.ru"><img src="https://raw.githubusercontent.com/abox-dev/sdk/main/readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240"></a></p>

# `abox-sdk`

The official Python SDK for AgentBox sandboxes and templates. The distribution is named `abox-sdk`; the Python package is `agentbox`.

```bash
pip install abox-sdk
```

Create an API key using the [API key guide](https://docs.agentbox.ru/en/quickstart/api-key/) and set `AGENTBOX_API_KEY`.

```python
from agentbox import Sandbox

with Sandbox.create() as sandbox:
    result = sandbox.commands.run('echo "Hello from AgentBox"')
    print(result.stdout)
```

Async and explicitly configured clients are supported:

```python
from agentbox import AgentBox

client = AgentBox(api_key="ab_...")
sandbox = await client.AsyncSandbox.create()
```

Documentation: [core SDK](https://docs.agentbox.ru/en/sdk/), [sandboxes](https://docs.agentbox.ru/en/sdk/sandboxes/), [commands](https://docs.agentbox.ru/en/sdk/commands/), [files](https://docs.agentbox.ru/en/sdk/files/), [templates](https://docs.agentbox.ru/en/sdk/templates/), and [API reference](https://docs.agentbox.ru/en/sdk-reference/python/).
