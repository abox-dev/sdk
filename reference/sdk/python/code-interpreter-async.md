## `agentbox_code_interpreter.code_interpreter_async`

**Classes:**

- [**AsyncSandbox**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox) – AgentBox cloud sandbox is a secure and isolated cloud environment.

**Attributes:**

- [**logger**](#agentbox_code_interpreter.code_interpreter_async.logger) – 

### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox`

Bases: <code>[AsyncSandbox](#agentbox.AsyncSandbox)</code>

AgentBox cloud sandbox is a secure and isolated cloud environment.

The sandbox allows you to:
- Access Linux OS
- Create, list, and delete files and directories
- Run commands
- Run isolated code
- Access the internet

See the [Code Interpreter guide](https://docs.agentbox.ru/en/sdk/code-interpreter/).

Use the `AsyncSandbox.create()` to create a new sandbox.

Example:
```python
from agentbox_code_interpreter import AsyncSandbox
sandbox = await AsyncSandbox.create()
```

**Functions:**

- [**connect**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.connect) – Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
- [**create**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.create) – Create a new sandbox.
- [**create_code_context**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.create_code_context) – Creates a new context to run code in.
- [**create_snapshot**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.create_snapshot) – Create a snapshot of the sandbox's current state.
- [**delete_snapshot**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.delete_snapshot) – Delete a snapshot.
- [**fork**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.fork) – Fork the sandbox.
- [**get_info**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.get_info) – Get sandbox information like sandbox ID, template, metadata, started at/end at date.
- [**get_mcp_token**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.get_mcp_token) – Get the MCP token for the sandbox.
- [**get_metrics**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.get_metrics) – Get the metrics of the current sandbox.
- [**is_running**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.is_running) – Check if the sandbox is running.
- [**kill**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.kill) – Kill the sandbox specified by sandbox ID.
- [**list**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.list) – List sandboxes.
- [**list_code_contexts**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.list_code_contexts) – List all contexts.
- [**list_snapshots**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.list_snapshots) – List snapshots for this sandbox.
- [**pause**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.pause) – Pause the sandbox.
- [**remove_code_context**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.remove_code_context) – Removes a context.
- [**restart_code_context**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.restart_code_context) – Restart a context.
- [**run_code**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.run_code) – 
- [**set_timeout**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.set_timeout) – Set the timeout of the specified sandbox.
- [**update_network**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.update_network) – Update the network configuration of the sandbox.

**Attributes:**

- [**commands**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.commands) (<code>[Commands](#agentbox.sandbox_async.commands.command.Commands)</code>) – Module for running commands in the sandbox.
- [**default_template**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.default_template) – 
- [**files**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.files) (<code>[Filesystem](#agentbox.sandbox_async.filesystem.filesystem.Filesystem)</code>) – Module for interacting with the sandbox filesystem.
- [**pty**](#agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.pty) (<code>[Pty](#agentbox.sandbox_async.commands.pty.Pty)</code>) – Module for interacting with the sandbox pseudo-terminal.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.commands`

```python
commands: Commands
```

Module for running commands in the sandbox.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.connect`

```python
connect(timeout=None, **opts)
```

Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
Sandbox must be either running or be paused.

With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

:param timeout: Timeout for the sandbox in **seconds**
    For running sandboxes, the timeout will update only if the new timeout is longer than the existing one.
:return: A running sandbox instance

@example
```python
sandbox = await AsyncSandbox.create()
await sandbox.pause()

# Another code block
same_sandbox = await sandbox.connect()
```

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.create`

```python
create(
    template=None,
    timeout=None,
    metadata=None,
    envs=None,
    secure=True,
    allow_internet_access=True,
    mcp=None,
    network=None,
    iam=None,
    lifecycle=None,
    logger=None,
    **opts
)
```

Create a new sandbox.

By default, the sandbox is created from the default `base` sandbox template.

:param template: Sandbox template name or ID
:param timeout: Timeout for the sandbox in **seconds**, default to 300 seconds. The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.
:param metadata: Custom metadata for the sandbox
:param envs: Custom environment variables for the sandbox
:param secure: Envd is secured with access token and cannot be used without it, defaults to `True`.
:param allow_internet_access: Allow sandbox to access the internet, defaults to `True`. If set to `False`, it works the same as setting network `deny_out` to `[0.0.0.0/0]`.
:param mcp: MCP server to enable in the sandbox
:param network: Sandbox network configuration. ``allow_out``/``deny_out`` may also be a callable receiving a :class:`SandboxNetworkSelectorContext` (``ctx.all_traffic``, ``ctx.rules``) and returning a list of strings. Per-host transform rules are nested under ``network.rules``; a rule's ``transform`` may be a callable receiving a :class:`SandboxNetworkTransformContext` of placeholder strings (``ctx.iam.tokens[name]``).
:param iam: Sandbox workload identity configuration. Each token contains ``audience`` and ``token_type``. Registered tokens are exposed to ``network.rules`` ``transform`` callables as ``ctx.iam.tokens[name]`` placeholders, which the egress proxy resolves per request
:param lifecycle: Sandbox lifecycle configuration — ``on_timeout``: ``"kill"`` or ``"pause"`` (omitted from the request when unset, leaving the API's default, currently ``"kill"``, in effect), or an object ``{"action": "pause"|"kill", "keep_memory": bool}`` where ``keep_memory`` set to ``False`` makes a timeout auto-pause filesystem-only (cold-boots on resume; cannot be combined with ``auto_resume``); an omitted ``keep_memory`` leaves the snapshot kind to the API; ``auto_resume``: leave unset to let the API pick the behavior, set ``False`` to opt out explicitly, or ``True`` (only when ``on_timeout`` action is ``"pause"``). Example: ``{"on_timeout": {"action": "pause", "keep_memory": False}}``
:param logger: Logger used for request and response logging for this sandbox. Accepts any standard library `logging.Logger`. When omitted, no request/response logging is emitted.

:return: A Sandbox instance for the new sandbox

Use this method instead of using the constructor to create a new sandbox.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.create_code_context`

```python
create_code_context(cwd=None, language=None, request_timeout=None)
```

Creates a new context to run code in.

:param cwd: Set the current working directory for the context, defaults to `/home/user`
:param language: Language of the context. If not specified, defaults to Python
:param request_timeout: Timeout for the request in **milliseconds**

:return: Context object

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.create_snapshot`

```python
create_snapshot(name=None, **opts)
```

Create a snapshot of the sandbox's current state.

The sandbox will be paused while the snapshot is being created.
The snapshot can be used to create new sandboxes with the same filesystem and state.
Snapshots are persistent and survive sandbox deletion.

Use the returned `snapshot_id` with `AsyncSandbox.create(snapshot_id)` to create a new sandbox from the snapshot.

:param name: Optional name for the snapshot template. If a snapshot template with this name already exists, a new build will be assigned to the existing template instead of creating a new one.

:return: Snapshot information including the snapshot ID and names

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.default_template`

```python
default_template = DEFAULT_TEMPLATE
```

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.delete_snapshot`

```python
delete_snapshot(snapshot_id, **opts)
```

Delete a snapshot.

:param snapshot_id: Snapshot ID
:return: `True` if the snapshot was deleted, `False` if it was not found

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.files`

```python
files: Filesystem
```

Module for interacting with the sandbox filesystem.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.fork`

```python
fork(timeout=None, count=None, **opts)
```

Fork the sandbox.

The sandbox is checkpointed in place (briefly paused, snapshotted with
its full memory state, and resumed — its ID and expiration stay
untouched) and `count` new sandboxes are created from that snapshot.
All forks boot from the same snapshot, so the snapshot is captured once
regardless of count.

Each fork succeeds or fails independently — the returned list contains
one entry per requested fork, either a running `AsyncSandbox` instance
or an exception describing why that fork failed to start. Per-fork
error codes map to the same exception classes as other API errors
(e.g. 429 to `RateLimitException`).

:param timeout: Timeout for the forked sandboxes in **seconds**, defaults to 300 seconds
:param count: Number of forked sandboxes to create, defaults to 1

:return: List with one entry per requested fork — a sandbox instance or an exception

@example
```python
sandbox = await AsyncSandbox.create()

fork1, fork2 = await sandbox.fork(count=2)
```

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.get_info`

```python
get_info(**opts)
```

Get sandbox information like sandbox ID, template, metadata, started at/end at date.

:return: Sandbox info

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.get_mcp_token`

```python
get_mcp_token()
```

Get the MCP token for the sandbox.

:return: MCP token for the sandbox, or None if MCP is not enabled.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.get_metrics`

```python
get_metrics(start=None, end=None, **opts)
```

Get the metrics of the current sandbox.

:param start: Start time for the metrics, defaults to the start of the sandbox
:param end: End time for the metrics, defaults to the current time

:return: List of sandbox metrics containing CPU, memory and disk usage information

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.is_running`

```python
is_running(request_timeout=None)
```

Check if the sandbox is running.

:param request_timeout: Timeout for the request in **seconds**

:return: `True` if the sandbox is running, `False` otherwise

Example
```python
sandbox = await AsyncSandbox.create()
await sandbox.is_running() # Returns True

await sandbox.kill()
await sandbox.is_running() # Returns False
```

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.kill`

```python
kill(**opts)
```

Kill the sandbox specified by sandbox ID.

:return: `True` if the sandbox was killed, `False` if the sandbox was not found

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.list`

```python
list(query=None, limit=None, next_token=None, **opts)
```

List sandboxes.

By default (no `query.state` set), returns sandboxes in both `running`
and `paused` states. To filter by state, pass `query=SandboxQuery(state=[...])`.

:param query: Filter the list of sandboxes by metadata or state, e.g. `SandboxQuery(metadata={"key": "value"})` or `SandboxQuery(state=[SandboxState.RUNNING])`
:param limit: Maximum number of sandboxes to return per page
:param next_token: Token for pagination

:return: An `AsyncSandboxPaginator` that yields pages of sandboxes (running and paused by default). Iterate pages via `await paginator.next_items()` while `paginator.has_next` is True.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.list_code_contexts`

```python
list_code_contexts()
```

List all contexts.

:return: List of contexts.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.list_snapshots`

```python
list_snapshots(limit=None, next_token=None, name=None, **opts)
```

List snapshots for this sandbox.

:param limit: Maximum number of snapshots to return per page
:param next_token: Token for pagination
:param name: Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-project/my-snapshot" or "my-snapshot:v1")

:return: Paginator for listing snapshots

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.pause`

```python
pause(keep_memory=True, **opts)
```

Pause the sandbox.

:param keep_memory: When `False`, the in-memory state is dropped and only the filesystem is persisted (no memory snapshot); resuming such a sandbox cold-boots (reboots) it from disk, losing running processes and open connections. Defaults to `True` (full memory snapshot).

:return: `True` if the sandbox got paused, `False` if the sandbox was already paused

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.pty`

```python
pty: Pty
```

Module for interacting with the sandbox pseudo-terminal.

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.remove_code_context`

```python
remove_code_context(context)
```

Removes a context.

:param context: Context to remove. Can be a Context object or a context ID string.

:return: None

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.restart_code_context`

```python
restart_code_context(context)
```

Restart a context.

:param context: Context to restart. Can be a Context object or a context ID string.

:return: None

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.run_code`

```python
run_code(
    code,
    language=None,
    context=None,
    on_stdout=None,
    on_stderr=None,
    on_result=None,
    on_error=None,
    envs=None,
    timeout=None,
    request_timeout=None,
)
```

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.set_timeout`

```python
set_timeout(timeout, **opts)
```

Set the timeout of the specified sandbox.
This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to `.set_timeout`.

The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.

:param timeout: Timeout for the sandbox in **seconds**

#### `agentbox_code_interpreter.code_interpreter_async.AsyncSandbox.update_network`

```python
update_network(network, **opts)
```

Update the network configuration of the sandbox.

Replaces the current egress configuration atomically — fields that are
omitted are cleared on the server.

:param network: New network configuration.

### `agentbox_code_interpreter.code_interpreter_async.logger`

```python
logger = logging.getLogger(__name__)
```

