[agentbox-sdk-monorepo](../README.md) / Sandbox

# Class: Sandbox

AgentBox cloud sandbox is a secure and isolated cloud environment.

The sandbox allows you to:
- Access Linux OS
- Create, list, and delete files and directories
- Run commands
- Run isolated code
- Access the internet

Check the [sandbox documentation](https://docs.agentbox.ru/en/sdk/sandboxes/).

Use [Sandbox.create](#create) to create a new sandbox.

## Example

```ts
import { Sandbox } from '@abox-dev/sdk'

const sandbox = await Sandbox.create()
```

## Extends

- `SandboxApi`

## Properties

### commands

> `readonly` **commands**: `Commands`

Module for running commands in the sandbox

***

### connectionConfig

> `protected` `readonly` **connectionConfig**: `ConnectionConfig`

***

### envdAccessToken?

> `protected` `readonly` `optional` **envdAccessToken?**: `string`

***

### envdPort

> `protected` `readonly` **envdPort**: `49983` = `49983`

***

### files

> `readonly` **files**: `Filesystem`

Module for interacting with the sandbox filesystem

***

### mcpPort

> `protected` `readonly` **mcpPort**: `50005` = `50005`

***

### pty

> `readonly` **pty**: `Pty`

Module for interacting with the sandbox pseudo-terminals

***

### sandboxDomain

> `readonly` **sandboxDomain**: `string`

Domain where the sandbox is hosted.

***

### sandboxId

> `readonly` **sandboxId**: `string`

Unique identifier of the sandbox.

***

### trafficAccessToken?

> `readonly` `optional` **trafficAccessToken?**: `string`

Traffic access token for accessing sandbox services with restricted public traffic.

***

### defaultMcpTemplate

> `protected` `readonly` `static` **defaultMcpTemplate**: `string` = `'mcp-gateway'`

***

### defaultSandboxTimeoutMs

> `protected` `readonly` `static` **defaultSandboxTimeoutMs**: `300000` = `DEFAULT_SANDBOX_TIMEOUT_MS`

***

### defaultTemplate

> `protected` `readonly` `static` **defaultTemplate**: `string` = `'base'`

## Methods

### connect()

> **connect**(`opts?`): `Promise`\<`Sandbox`\>

Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
Sandbox must be either running or be paused.

With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

#### Parameters

##### opts?

[`SandboxConnectOpts`](../type-aliases/SandboxConnectOpts.md)

connection options.

#### Returns

`Promise`\<`Sandbox`\>

A running sandbox instance

#### Example

```ts
const sandbox = await Sandbox.create()
await sandbox.pause()

// Connect to the same sandbox.
const sameSandbox = await sandbox.connect()
```

***

### createSnapshot()

> **createSnapshot**(`opts?`): `Promise`\<[`SnapshotInfo`](../interfaces/SnapshotInfo.md)\>

Create a snapshot of the sandbox's current state.

The sandbox will be paused while the snapshot is being created.
The snapshot can be used to create new sandboxes with the same filesystem and state.
Snapshots are persistent and survive sandbox deletion.

Use the returned `snapshotId` with `Sandbox.create(snapshotId)` to create a new sandbox from the snapshot.

#### Parameters

##### opts?

[`CreateSnapshotOpts`](../interfaces/CreateSnapshotOpts.md)

snapshot creation options including optional name and connection options.

#### Returns

`Promise`\<[`SnapshotInfo`](../interfaces/SnapshotInfo.md)\>

snapshot information including the snapshot ID.

#### Example

```ts
const sandbox = await Sandbox.create()
await sandbox.files.write('/app/state.json', '{"step": 1}')

// Create a snapshot
const snapshot = await sandbox.createSnapshot({ name: 'my-snapshot' })

// Create a new sandbox from the snapshot
const newSandbox = await Sandbox.create(snapshot.snapshotId)
```

***

### downloadUrl()

> **downloadUrl**(`path`, `opts?`): `Promise`\<`string`\>

Get the URL to download a file from the sandbox.

#### Parameters

##### path

`string`

path to the file in the sandbox.

##### opts?

`SandboxUrlOpts`

download url options.

#### Returns

`Promise`\<`string`\>

URL for downloading file.

***

### fork()

> **fork**(`opts?`): `Promise`\<(`Error` \| `Sandbox`)[]\>

Fork the sandbox.

The sandbox is checkpointed in place (briefly paused, snapshotted with its
full memory state, and resumed — its ID and expiration stay untouched) and
`count` new sandboxes are created from that snapshot. All forks boot from
the same snapshot, so the snapshot is captured once regardless of count.

Each fork succeeds or fails independently — the returned array contains
one entry per requested fork, either a running Sandbox instance or
an `Error` describing why that fork failed to start
(`Promise.allSettled`-style). Per-fork error codes map to the same error
classes as other API errors (e.g. 429 to `RateLimitError`).

#### Parameters

##### opts?

[`SandboxForkOpts`](../interfaces/SandboxForkOpts.md)

fork options — `count`, `timeoutMs` and connection options.

#### Returns

`Promise`\<(`Error` \| `Sandbox`)[]\>

array with one entry per requested fork — a sandbox instance or an error.

#### Example

```ts
const sandbox = await Sandbox.create()

const [fork1, fork2] = await sandbox.fork({ count: 2 })
if (fork1 instanceof Sandbox) {
  await fork1.commands.run('echo "hello from fork"')
}
```

***

### getHost()

> **getHost**(`port`): `string`

Get the host address for the specified sandbox port.
You can then use this address to connect to the sandbox port from outside the sandbox via HTTP or WebSocket.

#### Parameters

##### port

`number`

number of the port in the sandbox.

#### Returns

`string`

host address of the sandbox port.

#### Example

```ts
const sandbox = await Sandbox.create()
// Start an HTTP server
await sandbox.commands.run('python3 -m http.server 3000', { background: true })
// Get the hostname of the HTTP server
const serverURL = sandbox.getHost(3000)
```

***

### getInfo()

> **getInfo**(`opts?`): `Promise`\<[`SandboxInfo`](../interfaces/SandboxInfo.md)\>

Get sandbox information like sandbox ID, template, metadata, started at/end at date.

#### Parameters

##### opts?

`Pick`\<[`SandboxOpts`](../interfaces/SandboxOpts.md), `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<[`SandboxInfo`](../interfaces/SandboxInfo.md)\>

information about the sandbox

***

### getMcpToken()

> **getMcpToken**(): `Promise`\<`string` \| `undefined`\>

Get the MCP token for the sandbox.

#### Returns

`Promise`\<`string` \| `undefined`\>

MCP token for the sandbox, or undefined if MCP is not enabled.

***

### getMcpUrl()

> **getMcpUrl**(): `string`

Get the MCP URL for the sandbox.

#### Returns

`string`

MCP URL for the sandbox.

***

### getMetrics()

> **getMetrics**(`opts?`): `Promise`\<[`SandboxMetrics`](../interfaces/SandboxMetrics.md)[]\>

Get the metrics of the sandbox.

#### Parameters

##### opts?

[`SandboxMetricsOpts`](../interfaces/SandboxMetricsOpts.md)

connection options.

#### Returns

`Promise`\<[`SandboxMetrics`](../interfaces/SandboxMetrics.md)[]\>

List of sandbox metrics containing CPU, memory and disk usage information.

***

### isRunning()

> **isRunning**(`opts?`): `Promise`\<`boolean`\>

Check if the sandbox is running.

#### Parameters

##### opts?

`Pick`\<`ConnectionOpts`, `"requestTimeoutMs"` \| `"signal"`\>

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox is running, `false` otherwise.

#### Example

```ts
const sandbox = await Sandbox.create()
await sandbox.isRunning() // Returns true

await sandbox.kill()
await sandbox.isRunning() // Returns false
```

***

### kill()

> **kill**(`opts?`): `Promise`\<`boolean`\>

Kill the sandbox.

#### Parameters

##### opts?

`Pick`\<[`SandboxOpts`](../interfaces/SandboxOpts.md), `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox was killed, `false` if the sandbox was not found.

***

### listSnapshots()

> **listSnapshots**(`opts?`): [`SnapshotPaginator`](../interfaces/SnapshotPaginator.md)

List all snapshots created from this sandbox.

#### Parameters

##### opts?

`Omit`\<[`SnapshotListOpts`](../interfaces/SnapshotListOpts.md), `"sandboxId"`\>

list options.

#### Returns

[`SnapshotPaginator`](../interfaces/SnapshotPaginator.md)

paginator for listing snapshots from this sandbox.

***

### pause()

> **pause**(`opts?`): `Promise`\<`boolean`\>

Pause a sandbox by its ID.

#### Parameters

##### opts?

[`SandboxPauseOpts`](../interfaces/SandboxPauseOpts.md)

connection options, plus `keepMemory` to control the snapshot
kind. When `opts.keepMemory` is `false`, the in-memory state is dropped and
only the filesystem is persisted (a filesystem-only snapshot); resuming such
a sandbox cold-boots (reboots) it from disk, losing running processes and
open connections. Defaults to `true` (full memory snapshot).

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox got paused, `false` if the sandbox was already paused.

#### Example

```ts
const sandbox = await Sandbox.create()
await sandbox.pause()

// filesystem-only snapshot (resume reboots the sandbox)
await sandbox.pause({ keepMemory: false })
```

***

### setTimeout()

> **setTimeout**(`timeoutMs`, `opts?`): `Promise`\<`void`\>

Set the timeout of the sandbox.

This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to `.setTimeout`.
Maximum time a sandbox can be kept alive is 24 hours (86_400_000 milliseconds) for Pro users and 1 hour (3_600_000 milliseconds) for Hobby users.

#### Parameters

##### timeoutMs

`number`

timeout in **milliseconds**.

##### opts?

`Pick`\<[`SandboxOpts`](../interfaces/SandboxOpts.md), `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<`void`\>

***

### updateNetwork()

> **updateNetwork**(`network`, `opts?`): `Promise`\<`void`\>

Update the network configuration of the sandbox.

Replaces the current egress configuration atomically — fields that are
omitted are cleared on the server.

#### Parameters

##### network

[`SandboxNetworkUpdate`](../type-aliases/SandboxNetworkUpdate.md)

new network configuration.

##### opts?

`Pick`\<[`SandboxOpts`](../interfaces/SandboxOpts.md), `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<`void`\>

***

### uploadUrl()

> **uploadUrl**(`path?`, `opts?`): `Promise`\<`string`\>

Get the URL to upload a file to the sandbox.

You have to send a POST request to this URL with the file as multipart/form-data.

#### Parameters

##### path?

`string`

path to the file in the sandbox.

##### opts?

`SandboxUrlOpts`

download url options.

#### Returns

`Promise`\<`string`\>

URL for uploading file.

***

### connect()

> `static` **connect**\<`S`\>(`this`, `sandboxId`, `opts?`): `Promise`\<`InstanceType`\<`S`\>\>

Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
Sandbox must be either running or be paused.

With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

#### Type Parameters

##### S

`S` *extends* *typeof* `Sandbox`

#### Parameters

##### this

`S`

##### sandboxId

`string`

sandbox ID.

##### opts?

[`SandboxConnectOpts`](../type-aliases/SandboxConnectOpts.md)

connection options.

#### Returns

`Promise`\<`InstanceType`\<`S`\>\>

A running sandbox instance

#### Example

```ts
const sandbox = await Sandbox.create()
const sandboxId = sandbox.sandboxId

// Connect to the same sandbox.
const sameSandbox = await Sandbox.connect(sandboxId)
```

***

### connectSandbox()

> `protected` `static` **connectSandbox**(`sandboxId`, `opts?`): `Promise`\<\{ `envdAccessToken`: `string` \| `undefined`; `envdVersion`: `string`; `sandboxDomain`: `string` \| `undefined`; `sandboxId`: `string`; `trafficAccessToken`: `string` \| `undefined`; \}\>

#### Parameters

##### sandboxId

`string`

##### opts?

[`SandboxConnectOpts`](../type-aliases/SandboxConnectOpts.md)

#### Returns

`Promise`\<\{ `envdAccessToken`: `string` \| `undefined`; `envdVersion`: `string`; `sandboxDomain`: `string` \| `undefined`; `sandboxId`: `string`; `trafficAccessToken`: `string` \| `undefined`; \}\>

#### Inherited from

`SandboxApi.connectSandbox`

***

### create()

#### Call Signature

> `static` **create**\<`S`\>(`this`, `opts?`): `Promise`\<`InstanceType`\<`S`\>\>

Create a new sandbox from the default `base` sandbox template.

##### Type Parameters

###### S

`S` *extends* *typeof* `Sandbox`

##### Parameters

###### this

`S`

###### opts?

[`SandboxOpts`](../interfaces/SandboxOpts.md)

connection options.

##### Returns

`Promise`\<`InstanceType`\<`S`\>\>

sandbox instance for the new sandbox.

##### Example

```ts
const sandbox = await Sandbox.create()
```

##### Constructs

Sandbox

#### Call Signature

> `static` **create**\<`S`\>(`this`, `template`, `opts?`): `Promise`\<`InstanceType`\<`S`\>\>

Create a new sandbox from the specified sandbox template.

##### Type Parameters

###### S

`S` *extends* *typeof* `Sandbox`

##### Parameters

###### this

`S`

###### template

`string`

sandbox template name or ID.

###### opts?

[`SandboxOpts`](../interfaces/SandboxOpts.md)

connection options.

##### Returns

`Promise`\<`InstanceType`\<`S`\>\>

sandbox instance for the new sandbox.

##### Example

```ts
const sandbox = await Sandbox.create('<template-name-or-id>')
```

##### Constructs

Sandbox

***

### createSandbox()

> `protected` `static` **createSandbox**(`template`, `timeoutMs`, `opts?`): `Promise`\<\{ `envdAccessToken`: `string` \| `undefined`; `envdVersion`: `string`; `sandboxDomain`: `string` \| `undefined`; `sandboxId`: `string`; `trafficAccessToken`: `string` \| `undefined`; \}\>

#### Parameters

##### template

`string`

##### timeoutMs

`number`

##### opts?

[`SandboxOpts`](../interfaces/SandboxOpts.md)

#### Returns

`Promise`\<\{ `envdAccessToken`: `string` \| `undefined`; `envdVersion`: `string`; `sandboxDomain`: `string` \| `undefined`; `sandboxId`: `string`; `trafficAccessToken`: `string` \| `undefined`; \}\>

#### Inherited from

`SandboxApi.createSandbox`

***

### createSnapshot()

> `static` **createSnapshot**(`sandboxId`, `opts?`): `Promise`\<[`SnapshotInfo`](../interfaces/SnapshotInfo.md)\>

Create a snapshot from a sandbox.

The sandbox will be paused while the snapshot is being created.
The snapshot can be used to create new sandboxes with the same state.
The snapshot is a persistent image that survives sandbox deletion.

#### Parameters

##### sandboxId

`string`

sandbox ID to create snapshot from.

##### opts?

[`CreateSnapshotOpts`](../interfaces/CreateSnapshotOpts.md)

snapshot creation options including optional name and connection options.

#### Returns

`Promise`\<[`SnapshotInfo`](../interfaces/SnapshotInfo.md)\>

snapshot information including the snapshot name that can be used with Sandbox.create().

#### Inherited from

`SandboxApi.createSnapshot`

***

### deleteSnapshot()

> `static` **deleteSnapshot**(`snapshotId`, `opts?`): `Promise`\<`boolean`\>

Delete a snapshot.

#### Parameters

##### snapshotId

`string`

snapshot ID.

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the snapshot was deleted, `false` if it was not found.

#### Inherited from

`SandboxApi.deleteSnapshot`

***

### fork()

> `static` **fork**\<`S`\>(`this`, `sandboxId`, `opts?`): `Promise`\<(`Error` \| `InstanceType`\<`S`\>)[]\>

Fork a running sandbox specified by sandbox ID.

The sandbox is checkpointed in place (briefly paused, snapshotted with its
full memory state, and resumed — its ID and expiration stay untouched) and
`count` new sandboxes are created from that snapshot. All forks boot from
the same snapshot, so the snapshot is captured once regardless of count.

Each fork succeeds or fails independently — the returned array contains
one entry per requested fork, either a running Sandbox instance or
an `Error` describing why that fork failed to start
(`Promise.allSettled`-style). Per-fork error codes map to the same error
classes as other API errors (e.g. 429 to `RateLimitError`).

#### Type Parameters

##### S

`S` *extends* *typeof* `Sandbox`

#### Parameters

##### this

`S`

##### sandboxId

`string`

sandbox ID.

##### opts?

[`SandboxForkOpts`](../interfaces/SandboxForkOpts.md)

fork options — `count`, `timeoutMs` and connection options.

#### Returns

`Promise`\<(`Error` \| `InstanceType`\<`S`\>)[]\>

array with one entry per requested fork — a sandbox instance or an error.

#### Example

```ts
const sandbox = await Sandbox.create()

const [fork1, fork2] = await Sandbox.fork(sandbox.sandboxId, { count: 2 })
if (fork1 instanceof Sandbox) {
  await fork1.commands.run('echo "hello from fork"')
}
```

***

### forkSandbox()

> `protected` `static` **forkSandbox**(`sandboxId`, `timeoutMs`, `count`, `opts?`): `Promise`\<`SandboxForkResponse`[]\>

#### Parameters

##### sandboxId

`string`

##### timeoutMs

`number`

##### count

`number`

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

#### Returns

`Promise`\<`SandboxForkResponse`[]\>

#### Inherited from

`SandboxApi.forkSandbox`

***

### getInfo()

> `static` **getInfo**(`sandboxId`, `opts?`): `Promise`\<[`SandboxInfo`](../interfaces/SandboxInfo.md)\>

Get sandbox information like sandbox ID, template, metadata, started at/end at date.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

connection options.

#### Returns

`Promise`\<[`SandboxInfo`](../interfaces/SandboxInfo.md)\>

sandbox information.

#### Inherited from

`SandboxApi.getInfo`

***

### getMetrics()

> `static` **getMetrics**(`sandboxId`, `opts?`): `Promise`\<[`SandboxMetrics`](../interfaces/SandboxMetrics.md)[]\>

Get the metrics of the sandbox.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

[`SandboxMetricsOpts`](../interfaces/SandboxMetricsOpts.md)

sandbox metrics options.

#### Returns

`Promise`\<[`SandboxMetrics`](../interfaces/SandboxMetrics.md)[]\>

List of sandbox metrics containing CPU, memory and disk usage information.

#### Inherited from

`SandboxApi.getMetrics`

***

### kill()

> `static` **kill**(`sandboxId`, `opts?`): `Promise`\<`boolean`\>

Kill the sandbox specified by sandbox ID.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox was found and killed, `false` otherwise.

#### Inherited from

`SandboxApi.kill`

***

### list()

> `static` **list**(`opts?`): [`SandboxPaginator`](../interfaces/SandboxPaginator.md)

List sandboxes.

By default (no `query.state` set in `opts`), returns sandboxes in both
`running` and `paused` states. To filter by state, pass
`opts.query.state = [...]`.

#### Parameters

##### opts?

[`SandboxListOpts`](../interfaces/SandboxListOpts.md)

connection options, plus optional `query` to filter by
  metadata or state, and `limit` / `nextToken` for pagination.

#### Returns

[`SandboxPaginator`](../interfaces/SandboxPaginator.md)

a [SandboxPaginator](../interfaces/SandboxPaginator.md) that yields pages of sandboxes
  (running and paused by default). Iterate pages via
  `await paginator.nextItems()` while `paginator.hasNext` is `true`.

***

### listSnapshots()

> `static` **listSnapshots**(`opts?`): [`SnapshotPaginator`](../interfaces/SnapshotPaginator.md)

List all snapshots.

#### Parameters

##### opts?

[`SnapshotListOpts`](../interfaces/SnapshotListOpts.md)

list options including filters and pagination.

#### Returns

[`SnapshotPaginator`](../interfaces/SnapshotPaginator.md)

paginator for listing snapshots.

#### Inherited from

`SandboxApi.listSnapshots`

***

### pause()

> `static` **pause**(`sandboxId`, `opts?`): `Promise`\<`boolean`\>

Pause the sandbox specified by sandbox ID.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

[`SandboxPauseOpts`](../interfaces/SandboxPauseOpts.md)

pause options, including `keepMemory` and connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox got paused, `false` if the sandbox was already paused.

#### Inherited from

`SandboxApi.pause`

***

### setTimeout()

> `static` **setTimeout**(`sandboxId`, `timeoutMs`, `opts?`): `Promise`\<`void`\>

Set the timeout of the specified sandbox.
After the timeout expires the sandbox will be automatically killed.

This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to [Sandbox.setTimeout](#settimeout-1).

Maximum time a sandbox can be kept alive is 24 hours (86_400_000 milliseconds) for Pro users and 1 hour (3_600_000 milliseconds) for Hobby users.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### timeoutMs

`number`

timeout in **milliseconds**.

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

connection options.

#### Returns

`Promise`\<`void`\>

#### Inherited from

`SandboxApi.setTimeout`

***

### updateNetwork()

> `static` **updateNetwork**(`sandboxId`, `network`, `opts?`): `Promise`\<`void`\>

Update the network configuration of a running sandbox.

Replaces the current egress configuration atomically — fields that are
omitted are cleared on the server.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### network

[`SandboxNetworkUpdate`](../type-aliases/SandboxNetworkUpdate.md)

new network configuration.

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

connection options.

#### Returns

`Promise`\<`void`\>

#### Inherited from

`SandboxApi.updateNetwork`
