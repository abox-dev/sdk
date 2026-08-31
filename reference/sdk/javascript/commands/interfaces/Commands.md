[agentbox-sdk-monorepo](../README.md) / Commands

# Interface: Commands

Module for starting and interacting with commands in the sandbox.

## Properties

### rpc

> `protected` `readonly` **rpc**: `Client`\<*typeof* `ProcessService`\>

## Methods

### closeStdin()

> **closeStdin**(`pid`, `opts?`): `Promise`\<`void`\>

Close command stdin.

This signals EOF to the command. The command must have been started with `stdin: true`.

#### Parameters

##### pid

`number`

process ID of the command. You can get the list of running commands using [Commands.list](#list).

##### opts?

[`CommandRequestOpts`](CommandRequestOpts.md)

connection options.

#### Returns

`Promise`\<`void`\>

***

### connect()

> **connect**(`pid`, `opts?`): `Promise`\<[`CommandHandle`](CommandHandle.md)\>

Connect to a running command.
You can use [CommandHandle.wait](CommandHandle.md#wait) to wait for the command to finish and get execution results.

#### Parameters

##### pid

`number`

process ID of the command to connect to. You can get the list of running commands using [Commands.list](#list).

##### opts?

[`CommandConnectOpts`](../type-aliases/CommandConnectOpts.md)

connection options.

#### Returns

`Promise`\<[`CommandHandle`](CommandHandle.md)\>

`CommandHandle` handle to interact with the running command.

***

### kill()

> **kill**(`pid`, `opts?`): `Promise`\<`boolean`\>

Kill a running command specified by its process ID.
It uses `SIGKILL` signal to kill the command.

#### Parameters

##### pid

`number`

process ID of the command. You can get the list of running commands using [Commands.list](#list).

##### opts?

[`CommandRequestOpts`](CommandRequestOpts.md)

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the command was killed, `false` if the command was not found.

***

### list()

> **list**(`opts?`): `Promise`\<[`ProcessInfo`](ProcessInfo.md)[]\>

List all running commands and PTY sessions.

#### Parameters

##### opts?

[`CommandRequestOpts`](CommandRequestOpts.md)

connection options.

#### Returns

`Promise`\<[`ProcessInfo`](ProcessInfo.md)[]\>

list of running commands and PTY sessions.

***

### run()

#### Call Signature

> **run**(`cmd`, `opts?`): `Promise`\<[`CommandResult`](CommandResult.md)\>

Start a new command and wait until it finishes executing.

##### Parameters

###### cmd

`string`

command to execute.

###### opts?

[`CommandStartOpts`](CommandStartOpts.md) & `object`

options for starting the command.

##### Returns

`Promise`\<[`CommandResult`](CommandResult.md)\>

`CommandResult` result of the command execution.

#### Call Signature

> **run**(`cmd`, `opts`): `Promise`\<[`CommandHandle`](CommandHandle.md)\>

Start a new command in the background.
You can use [CommandHandle.wait](CommandHandle.md#wait) to wait for the command to finish and get its result.

##### Parameters

###### cmd

`string`

command to execute.

###### opts

[`CommandStartOpts`](CommandStartOpts.md) & `object`

options for starting the command

##### Returns

`Promise`\<[`CommandHandle`](CommandHandle.md)\>

`CommandHandle` handle to interact with the running command.

#### Call Signature

> **run**(`cmd`, `opts?`): `Promise`\<[`CommandResult`](CommandResult.md) \| [`CommandHandle`](CommandHandle.md)\>

Start a new command.

##### Parameters

###### cmd

`string`

command to execute.

###### opts?

[`CommandStartOpts`](CommandStartOpts.md) & `object`

options for starting the command.
  - `opts.background: true` - runs in background, returns `CommandHandle`
  - `opts.background: false | undefined` - waits for completion, returns `CommandResult`

##### Returns

`Promise`\<[`CommandResult`](CommandResult.md) \| [`CommandHandle`](CommandHandle.md)\>

Either a `CommandHandle` or a `CommandResult` (depending on `opts.background`).

***

### sendStdin()

> **sendStdin**(`pid`, `data`, `opts?`): `Promise`\<`void`\>

Send data to command stdin.

#### Parameters

##### pid

`number`

process ID of the command. You can get the list of running commands using [Commands.list](#list).

##### data

`string` \| `Uint8Array`\<`ArrayBufferLike`\>

data to send to the command.

##### opts?

[`CommandRequestOpts`](CommandRequestOpts.md)

connection options.

#### Returns

`Promise`\<`void`\>
