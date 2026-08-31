[agentbox-sdk-monorepo](../README.md) / Pty

# Interface: Pty

Module for interacting with PTYs (pseudo-terminals) in the sandbox.

## Methods

### connect()

> **connect**(`pid`, `opts?`): `Promise`\<[`CommandHandle`](CommandHandle.md)\>

Connect to a running PTY.

#### Parameters

##### pid

`number`

process ID of the PTY to connect to. You can get the list of running PTYs using [Commands.list](Commands.md#list).

##### opts?

`PtyConnectOpts`

connection options.

#### Returns

`Promise`\<[`CommandHandle`](CommandHandle.md)\>

handle to interact with the PTY.

***

### create()

> **create**(`opts`): `Promise`\<[`CommandHandle`](CommandHandle.md)\>

Create a new PTY (pseudo-terminal).

#### Parameters

##### opts

`PtyCreateOpts`

options for creating the PTY.

#### Returns

`Promise`\<[`CommandHandle`](CommandHandle.md)\>

handle to interact with the PTY.

***

### kill()

> **kill**(`pid`, `opts?`): `Promise`\<`boolean`\>

Kill a running PTY specified by process ID.
It uses `SIGKILL` signal to kill the PTY.

#### Parameters

##### pid

`number`

process ID of the PTY.

##### opts?

`Pick`\<`ConnectionOpts`, `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the PTY was killed, `false` if the PTY was not found.

***

### resize()

> **resize**(`pid`, `size`, `opts?`): `Promise`\<`void`\>

Resize PTY.
Call this when the terminal window is resized and the number of columns and rows has changed.

#### Parameters

##### pid

`number`

process ID of the PTY.

##### size

new size of the PTY.

###### cols

`number`

###### rows

`number`

##### opts?

`Pick`\<`ConnectionOpts`, `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<`void`\>

***

### sendInput()

> **sendInput**(`pid`, `data`, `opts?`): `Promise`\<`void`\>

Send input to a PTY.

#### Parameters

##### pid

`number`

process ID of the PTY.

##### data

`Uint8Array`

input data to send to the PTY.

##### opts?

`Pick`\<`ConnectionOpts`, `"requestTimeoutMs"` \| `"signal"`\>

connection options.

#### Returns

`Promise`\<`void`\>
