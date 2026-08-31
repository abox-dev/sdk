[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxApi

# Class: SandboxApi

## Extends

- `ClientFactory`

## Extended by

- [`Sandbox`](../../index/classes/Sandbox.md)

## Constructors

### Constructor

> `protected` **new SandboxApi**(): `SandboxApi`

#### Returns

`SandboxApi`

#### Overrides

`ClientFactory.constructor`

## Methods

### connectSandbox()

> `protected` `static` **connectSandbox**(`sandboxId`, `opts?`): `Promise`\<\{ `envdAccessToken`: `string` \| `undefined`; `envdVersion`: `string`; `sandboxDomain`: `string` \| `undefined`; `sandboxId`: `string`; `trafficAccessToken`: `string` \| `undefined`; \}\>

#### Parameters

##### sandboxId

`string`

##### opts?

[`SandboxConnectOpts`](../type-aliases/SandboxConnectOpts.md)

#### Returns

`Promise`\<\{ `envdAccessToken`: `string` \| `undefined`; `envdVersion`: `string`; `sandboxDomain`: `string` \| `undefined`; `sandboxId`: `string`; `trafficAccessToken`: `string` \| `undefined`; \}\>

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

***

### listSnapshots()

> `static` **listSnapshots**(`opts?`): [`SnapshotPaginator`](SnapshotPaginator.md)

List all snapshots.

#### Parameters

##### opts?

[`SnapshotListOpts`](../interfaces/SnapshotListOpts.md)

list options including filters and pagination.

#### Returns

[`SnapshotPaginator`](SnapshotPaginator.md)

paginator for listing snapshots.

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

***

### setTimeout()

> `static` **setTimeout**(`sandboxId`, `timeoutMs`, `opts?`): `Promise`\<`void`\>

Set the timeout of the specified sandbox.
After the timeout expires the sandbox will be automatically killed.

This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to Sandbox.setTimeout.

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
