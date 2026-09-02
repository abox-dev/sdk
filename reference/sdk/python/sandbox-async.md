## `agentbox.sandbox_async.main`

**Classes:**

- [**AsyncSandbox**](#agentbox.sandbox_async.main.AsyncSandbox) – AgentBox cloud sandbox is a secure and isolated cloud environment.

**Attributes:**

- [**logger**](#agentbox.sandbox_async.main.logger) – 

### `agentbox.sandbox_async.main.AsyncSandbox`

```python
AsyncSandbox(**opts)
```

Bases: <code>[SandboxApi](#agentbox.sandbox_async.sandbox_api.SandboxApi)</code>

AgentBox cloud sandbox is a secure and isolated cloud environment.

The sandbox allows you to:
- Access Linux OS
- Create, list, and delete files and directories
- Run commands
- Run isolated code
- Access the internet

See the sandbox documentation at https://docs.agentbox.ru/en/sdk/sandboxes/.

Use the `AsyncSandbox.create()` to create a new sandbox.

Example:
```python
from agentbox import AsyncSandbox

sandbox = await AsyncSandbox.create()
```

**Functions:**

- [**connect**](#agentbox.sandbox_async.main.AsyncSandbox.connect) – Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
- [**create**](#agentbox.sandbox_async.main.AsyncSandbox.create) – Create a new sandbox.
- [**create_snapshot**](#agentbox.sandbox_async.main.AsyncSandbox.create_snapshot) – Create a snapshot of the sandbox's current state.
- [**delete_snapshot**](#agentbox.sandbox_async.main.AsyncSandbox.delete_snapshot) – Delete a snapshot.
- [**fork**](#agentbox.sandbox_async.main.AsyncSandbox.fork) – Fork the sandbox.
- [**get_info**](#agentbox.sandbox_async.main.AsyncSandbox.get_info) – Get sandbox information like sandbox ID, template, metadata, started at/end at date.
- [**get_metrics**](#agentbox.sandbox_async.main.AsyncSandbox.get_metrics) – Get the metrics of the current sandbox.
- [**is_running**](#agentbox.sandbox_async.main.AsyncSandbox.is_running) – Check if the sandbox is running.
- [**kill**](#agentbox.sandbox_async.main.AsyncSandbox.kill) – Kill the sandbox specified by sandbox ID.
- [**list**](#agentbox.sandbox_async.main.AsyncSandbox.list) – List sandboxes.
- [**list_snapshots**](#agentbox.sandbox_async.main.AsyncSandbox.list_snapshots) – List snapshots for this sandbox.
- [**pause**](#agentbox.sandbox_async.main.AsyncSandbox.pause) – Pause the sandbox.
- [**set_timeout**](#agentbox.sandbox_async.main.AsyncSandbox.set_timeout) – Set the timeout of the specified sandbox.
- [**update_network**](#agentbox.sandbox_async.main.AsyncSandbox.update_network) – Update the network configuration of the sandbox.

**Attributes:**

- [**commands**](#agentbox.sandbox_async.main.AsyncSandbox.commands) (<code>[Commands](#agentbox.sandbox_async.commands.command.Commands)</code>) – Module for running commands in the sandbox.
- [**files**](#agentbox.sandbox_async.main.AsyncSandbox.files) (<code>[Filesystem](#agentbox.sandbox_async.filesystem.filesystem.Filesystem)</code>) – Module for interacting with the sandbox filesystem.
- [**pty**](#agentbox.sandbox_async.main.AsyncSandbox.pty) (<code>[Pty](#agentbox.sandbox_async.commands.pty.Pty)</code>) – Module for interacting with the sandbox pseudo-terminal.

Applications should use `AsyncSandbox.create()` or `AsyncSandbox.connect()`.

#### `agentbox.sandbox_async.main.AsyncSandbox.commands`

```python
commands: Commands
```

Module for running commands in the sandbox.

#### `agentbox.sandbox_async.main.AsyncSandbox.connect`

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

#### `agentbox.sandbox_async.main.AsyncSandbox.create`

```python
create(
    template=None,
    timeout=None,
    metadata=None,
    envs=None,
    secure=True,
    allow_internet_access=True,
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
:param network: Sandbox network configuration. ``allow_out``/``deny_out`` may also be a callable receiving a :class:`SandboxNetworkSelectorContext` (``ctx.all_traffic``, ``ctx.rules``) and returning a list of strings. Per-host transform rules are nested under ``network.rules``; a rule's ``transform`` may be a callable receiving a :class:`SandboxNetworkTransformContext` of placeholder strings (``ctx.iam.tokens[name]``).
:param iam: Sandbox workload identity configuration. Each token contains ``audience`` and ``token_type``. Registered tokens are exposed to ``network.rules`` ``transform`` callables as ``ctx.iam.tokens[name]`` placeholders, which the egress proxy resolves per request
:param lifecycle: Sandbox lifecycle configuration — ``on_timeout``: ``"kill"`` or ``"pause"`` (omitted from the request when unset, leaving the API's default, currently ``"kill"``, in effect), or an object ``{"action": "pause"|"kill", "keep_memory": bool}`` where ``keep_memory`` set to ``False`` makes a timeout auto-pause filesystem-only (cold-boots on resume; cannot be combined with ``auto_resume``); an omitted ``keep_memory`` leaves the snapshot kind to the API; ``auto_resume``: leave unset to let the API pick the behavior, set ``False`` to opt out explicitly, or ``True`` (only when ``on_timeout`` action is ``"pause"``). Example: ``{"on_timeout": {"action": "pause", "keep_memory": False}}``
:param logger: Logger used for request and response logging for this sandbox. Accepts any standard library `logging.Logger`. When omitted, no request/response logging is emitted.

:return: A Sandbox instance for the new sandbox

Use this method instead of using the constructor to create a new sandbox.

#### `agentbox.sandbox_async.main.AsyncSandbox.create_snapshot`

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

#### `agentbox.sandbox_async.main.AsyncSandbox.delete_snapshot`

```python
delete_snapshot(snapshot_id, **opts)
```

Delete a snapshot.

:param snapshot_id: Snapshot ID
:return: `True` if the snapshot was deleted, `False` if it was not found

#### `agentbox.sandbox_async.main.AsyncSandbox.files`

```python
files: Filesystem
```

Module for interacting with the sandbox filesystem.

#### `agentbox.sandbox_async.main.AsyncSandbox.fork`

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

#### `agentbox.sandbox_async.main.AsyncSandbox.get_info`

```python
get_info(**opts)
```

Get sandbox information like sandbox ID, template, metadata, started at/end at date.

:return: Sandbox info

#### `agentbox.sandbox_async.main.AsyncSandbox.get_metrics`

```python
get_metrics(start=None, end=None, **opts)
```

Get the metrics of the current sandbox.

:param start: Start time for the metrics, defaults to the start of the sandbox
:param end: End time for the metrics, defaults to the current time

:return: List of sandbox metrics containing CPU, memory and disk usage information

#### `agentbox.sandbox_async.main.AsyncSandbox.is_running`

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

#### `agentbox.sandbox_async.main.AsyncSandbox.kill`

```python
kill(**opts)
```

Kill the sandbox specified by sandbox ID.

:return: `True` if the sandbox was killed, `False` if the sandbox was not found

#### `agentbox.sandbox_async.main.AsyncSandbox.list`

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

#### `agentbox.sandbox_async.main.AsyncSandbox.list_snapshots`

```python
list_snapshots(limit=None, next_token=None, name=None, **opts)
```

List snapshots for this sandbox.

:param limit: Maximum number of snapshots to return per page
:param next_token: Token for pagination
:param name: Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-project/my-snapshot" or "my-snapshot:v1")

:return: Paginator for listing snapshots

#### `agentbox.sandbox_async.main.AsyncSandbox.pause`

```python
pause(keep_memory=True, **opts)
```

Pause the sandbox.

:param keep_memory: When `False`, the in-memory state is dropped and only the filesystem is persisted (no memory snapshot); resuming such a sandbox cold-boots (reboots) it from disk, losing running processes and open connections. Defaults to `True` (full memory snapshot).

:return: `True` if the sandbox got paused, `False` if the sandbox was already paused

#### `agentbox.sandbox_async.main.AsyncSandbox.pty`

```python
pty: Pty
```

Module for interacting with the sandbox pseudo-terminal.

#### `agentbox.sandbox_async.main.AsyncSandbox.set_timeout`

```python
set_timeout(timeout, **opts)
```

Set the timeout of the specified sandbox.
This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to `.set_timeout`.

The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.

:param timeout: Timeout for the sandbox in **seconds**

#### `agentbox.sandbox_async.main.AsyncSandbox.update_network`

```python
update_network(network, **opts)
```

Update the network configuration of the sandbox.

Replaces the current egress configuration atomically — fields that are
omitted are cleared on the server.

:param network: New network configuration.

### `agentbox.sandbox_async.main.logger`

```python
logger = logging.getLogger(__name__)
```

