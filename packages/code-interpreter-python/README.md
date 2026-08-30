<p align="center"><img src="https://raw.githubusercontent.com/abox-dev/sdk/main/readme-assets/agentbox-logo-email.png" alt="AgentBox" width="240"></p>

# `abox-code-interpreter`

Stateful Python, JavaScript, and TypeScript execution in an AgentBox sandbox. The distribution is named `abox-code-interpreter`; the Python package is `agentbox_code_interpreter`.

```bash
pip install abox-code-interpreter
```

Create an API key using the [API key guide](https://docs.agentbox.ru/en/quickstart/api-key/) and set `AGENTBOX_API_KEY`.

```python
from agentbox_code_interpreter import Sandbox

with Sandbox.create() as sandbox:
    execution = sandbox.run_code("x = 1; x += 1; x")
    print(execution.text)
```

`AsyncSandbox` provides the equivalent async API. See the [Code Interpreter guide](https://docs.agentbox.ru/en/sdk/code-interpreter/) and [API reference](https://docs.agentbox.ru/en/sdk/api-reference/python/code-interpreter/).
