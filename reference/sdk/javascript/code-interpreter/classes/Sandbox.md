[agentbox-sdk-monorepo](../README.md) / Sandbox

# Class: Sandbox

AgentBox cloud sandbox is a secure and isolated cloud environment.

The sandbox allows you to:
- Access Linux OS
- Create, list, and delete files and directories
- Run commands
- Run isolated code
- Access the internet

See the [Code Interpreter guide](https://docs.agentbox.ru/en/sdk/code-interpreter/).

Use [Sandbox.create](#create) to create a new sandbox.

## Example

```ts
import { Sandbox } from '@abox-dev/code-interpreter'

const sandbox = await Sandbox.create()
```

## Extends

- `Sandbox`

## Properties

### commands

> `readonly` **commands**: `Commands`

Module for running commands in the sandbox

#### Inherited from

`BaseSandbox.commands`

***

### files

> `readonly` **files**: `Filesystem`

Module for interacting with the sandbox filesystem

#### Inherited from

`BaseSandbox.files`

***

### pty

> `readonly` **pty**: `Pty`

Module for interacting with the sandbox pseudo-terminals

#### Inherited from

`BaseSandbox.pty`

***

### sandboxDomain

> `readonly` **sandboxDomain**: `string`

Domain where the sandbox is hosted.

#### Inherited from

`BaseSandbox.sandboxDomain`

***

### sandboxId

> `readonly` **sandboxId**: `string`

Unique identifier of the sandbox.

#### Inherited from

`BaseSandbox.sandboxId`

***

### trafficAccessToken?

> `readonly` `optional` **trafficAccessToken?**: `string`

Traffic access token for accessing sandbox services with restricted public traffic.

#### Inherited from

`BaseSandbox.trafficAccessToken`

## Methods

### connect()

> **connect**(`opts?`): `Promise`\<`Sandbox`\>

Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
Sandbox must be either running or be paused.

With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

#### Parameters

##### opts?

`SandboxConnectOpts`

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

#### Inherited from

`BaseSandbox.connect`

***

### createCodeContext()

> **createCodeContext**(`opts?`): `Promise`\<[`Context`](../type-aliases/Context.md)\>

Creates a new context to run code in.

#### Parameters

##### opts?

[`CreateCodeContextOpts`](../interfaces/CreateCodeContextOpts.md)

options for creating the context.

#### Returns

`Promise`\<[`Context`](../type-aliases/Context.md)\>

context object.

***

### createSnapshot()

> **createSnapshot**(`opts?`): `Promise`\<`SnapshotInfo`\>

Create a snapshot of the sandbox's current state.

The sandbox will be paused while the snapshot is being created.
The snapshot can be used to create new sandboxes with the same filesystem and state.
Snapshots are persistent and survive sandbox deletion.

Use the returned `snapshotId` with `Sandbox.create(snapshotId)` to create a new sandbox from the snapshot.

#### Parameters

##### opts?

`CreateSnapshotOpts`

snapshot creation options including optional name and connection options.

#### Returns

`Promise`\<`SnapshotInfo`\>

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

#### Inherited from

`BaseSandbox.createSnapshot`

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

#### Inherited from

`BaseSandbox.downloadUrl`

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

`SandboxForkOpts`

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

#### Inherited from

`BaseSandbox.fork`

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

#### Inherited from

`BaseSandbox.getHost`

***

### getInfo()

> **getInfo**(`opts?`): `Promise`\<`SandboxInfo`\>

Get sandbox information like sandbox ID, template, metadata, started at/end at date.

#### Parameters

##### opts?

`Pick`\<`SandboxOpts`, `"signal"` \| `"requestTimeoutMs"`\>

connection options.

#### Returns

`Promise`\<`SandboxInfo`\>

information about the sandbox

#### Inherited from

`BaseSandbox.getInfo`

***

### getMetrics()

> **getMetrics**(`opts?`): `Promise`\<`SandboxMetrics`[]\>

Get the metrics of the sandbox.

#### Parameters

##### opts?

`SandboxMetricsOpts`

connection options.

#### Returns

`Promise`\<`SandboxMetrics`[]\>

List of sandbox metrics containing CPU, memory and disk usage information.

#### Inherited from

`BaseSandbox.getMetrics`

***

### isRunning()

> **isRunning**(`opts?`): `Promise`\<`boolean`\>

Check if the sandbox is running.

#### Parameters

##### opts?

`Pick`\<`ConnectionOpts`, `"signal"` \| `"requestTimeoutMs"`\>

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

#### Inherited from

`BaseSandbox.isRunning`

***

### kill()

> **kill**(`opts?`): `Promise`\<`boolean`\>

Kill the sandbox.

#### Parameters

##### opts?

`Pick`\<`SandboxOpts`, `"signal"` \| `"requestTimeoutMs"`\>

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox was killed, `false` if the sandbox was not found.

#### Inherited from

`BaseSandbox.kill`

***

### listCodeContexts()

> **listCodeContexts**(): `Promise`\<[`Context`](../type-aliases/Context.md)[]\>

List all contexts.

#### Returns

`Promise`\<[`Context`](../type-aliases/Context.md)[]\>

list of contexts.

***

### listSnapshots()

> **listSnapshots**(`opts?`): `SnapshotPaginator`

List all snapshots created from this sandbox.

#### Parameters

##### opts?

`Omit`\<`SnapshotListOpts`, `"sandboxId"`\>

list options.

#### Returns

`SnapshotPaginator`

paginator for listing snapshots from this sandbox.

#### Inherited from

`BaseSandbox.listSnapshots`

***

### pause()

> **pause**(`opts?`): `Promise`\<`boolean`\>

Pause a sandbox by its ID.

#### Parameters

##### opts?

`SandboxPauseOpts`

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

#### Inherited from

`BaseSandbox.pause`

***

### removeCodeContext()

> **removeCodeContext**(`context`): `Promise`\<`void`\>

Removes a context.

#### Parameters

##### context

`string` \| [`Context`](../type-aliases/Context.md)

context to remove.

#### Returns

`Promise`\<`void`\>

void.

***

### restartCodeContext()

> **restartCodeContext**(`context`): `Promise`\<`void`\>

Restart a context.

#### Parameters

##### context

`string` \| [`Context`](../type-aliases/Context.md)

context to restart.

#### Returns

`Promise`\<`void`\>

void.

***

### runCode()

#### Call Signature

> **runCode**(`code`, `opts?`): `Promise`\<[`Execution`](../interfaces/Execution.md)\>

Run the code for the specified language.

Specify the `language` or `context` option to run the code as a different language or in a different `Context`.
If no language is specified, Python is used.

You can reference previously defined variables, imports, and functions in the code.

##### Parameters

###### code

`string`

code to execute.

###### opts?

[`RunCodeOpts`](../interfaces/RunCodeOpts.md) & `object`

options for executing the code.

##### Returns

`Promise`\<[`Execution`](../interfaces/Execution.md)\>

`Execution` result object.

#### Call Signature

> **runCode**(`code`, `opts?`): `Promise`\<[`Execution`](../interfaces/Execution.md)\>

Runs the code in the specified context, if not specified, the default context is used.

Specify the `language` or `context` option to run the code as a different language or in a different `Context`.

You can reference previously defined variables, imports, and functions in the code.

##### Parameters

###### code

`string`

code to execute.

###### opts?

[`RunCodeOpts`](../interfaces/RunCodeOpts.md) & `object`

options for executing the code

##### Returns

`Promise`\<[`Execution`](../interfaces/Execution.md)\>

`Execution` result object

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

`Pick`\<`SandboxOpts`, `"signal"` \| `"requestTimeoutMs"`\>

connection options.

#### Returns

`Promise`\<`void`\>

#### Inherited from

`BaseSandbox.setTimeout`

***

### updateNetwork()

> **updateNetwork**(`network`, `opts?`): `Promise`\<`void`\>

Update the network configuration of the sandbox.

Replaces the current egress configuration atomically — fields that are
omitted are cleared on the server.

#### Parameters

##### network

`SandboxNetworkUpdate`

new network configuration.

##### opts?

`Pick`\<`SandboxOpts`, `"signal"` \| `"requestTimeoutMs"`\>

connection options.

#### Returns

`Promise`\<`void`\>

#### Inherited from

`BaseSandbox.updateNetwork`

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

#### Inherited from

`BaseSandbox.uploadUrl`

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

`SandboxConnectOpts`

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

#### Inherited from

`BaseSandbox.connect`

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

`SandboxOpts`

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

##### Inherited from

`BaseSandbox.create`

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

`SandboxOpts`

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

##### Inherited from

`BaseSandbox.create`

***

### createSnapshot()

> `static` **createSnapshot**(`sandboxId`, `opts?`): `Promise`\<`SnapshotInfo`\>

Create a snapshot from a sandbox.

The sandbox will be paused while the snapshot is being created.
The snapshot can be used to create new sandboxes with the same state.
The snapshot is a persistent image that survives sandbox deletion.

#### Parameters

##### sandboxId

`string`

sandbox ID to create snapshot from.

##### opts?

`CreateSnapshotOpts`

snapshot creation options including optional name and connection options.

#### Returns

`Promise`\<`SnapshotInfo`\>

snapshot information including the snapshot name that can be used with Sandbox.create().

#### Inherited from

`BaseSandbox.createSnapshot`

***

### deleteSnapshot()

> `static` **deleteSnapshot**(`snapshotId`, `opts?`): `Promise`\<`boolean`\>

Delete a snapshot.

#### Parameters

##### snapshotId

`string`

snapshot ID.

##### opts?

`SandboxApiOpts`

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the snapshot was deleted, `false` if it was not found.

#### Inherited from

`BaseSandbox.deleteSnapshot`

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

`SandboxForkOpts`

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

#### Inherited from

`BaseSandbox.fork`

***

### getInfo()

> `static` **getInfo**(`sandboxId`, `opts?`): `Promise`\<`SandboxInfo`\>

Get sandbox information like sandbox ID, template, metadata, started at/end at date.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

`SandboxApiOpts`

connection options.

#### Returns

`Promise`\<`SandboxInfo`\>

sandbox information.

#### Inherited from

`BaseSandbox.getInfo`

***

### getMetrics()

> `static` **getMetrics**(`sandboxId`, `opts?`): `Promise`\<`SandboxMetrics`[]\>

Get the metrics of the sandbox.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

`SandboxMetricsOpts`

sandbox metrics options.

#### Returns

`Promise`\<`SandboxMetrics`[]\>

List of sandbox metrics containing CPU, memory and disk usage information.

#### Inherited from

`BaseSandbox.getMetrics`

***

### kill()

> `static` **kill**(`sandboxId`, `opts?`): `Promise`\<`boolean`\>

Kill the sandbox specified by sandbox ID.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

`SandboxApiOpts`

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox was found and killed, `false` otherwise.

#### Inherited from

`BaseSandbox.kill`

***

### list()

> `static` **list**(`opts?`): `SandboxPaginator`

List sandboxes.

By default (no `query.state` set in `opts`), returns sandboxes in both
`running` and `paused` states. To filter by state, pass
`opts.query.state = [...]`.

#### Parameters

##### opts?

`SandboxListOpts`

connection options, plus optional `query` to filter by
  metadata or state, and `limit` / `nextToken` for pagination.

#### Returns

`SandboxPaginator`

a SandboxPaginator that yields pages of sandboxes
  (running and paused by default). Iterate pages via
  `await paginator.nextItems()` while `paginator.hasNext` is `true`.

#### Inherited from

`BaseSandbox.list`

***

### listSnapshots()

> `static` **listSnapshots**(`opts?`): `SnapshotPaginator`

List all snapshots.

#### Parameters

##### opts?

`SnapshotListOpts`

list options including filters and pagination.

#### Returns

`SnapshotPaginator`

paginator for listing snapshots.

#### Inherited from

`BaseSandbox.listSnapshots`

***

### pause()

> `static` **pause**(`sandboxId`, `opts?`): `Promise`\<`boolean`\>

Pause the sandbox specified by sandbox ID.

#### Parameters

##### sandboxId

`string`

sandbox ID.

##### opts?

`SandboxPauseOpts`

pause options, including `keepMemory` and connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the sandbox got paused, `false` if the sandbox was already paused.

#### Inherited from

`BaseSandbox.pause`

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

`SandboxApiOpts`

connection options.

#### Returns

`Promise`\<`void`\>

#### Inherited from

`BaseSandbox.setTimeout`

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

`SandboxNetworkUpdate`

new network configuration.

##### opts?

`SandboxApiOpts`

connection options.

#### Returns

`Promise`\<`void`\>

#### Inherited from

`BaseSandbox.updateNetwork`
