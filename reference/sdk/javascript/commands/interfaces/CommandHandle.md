[agentbox-sdk-monorepo](../README.md) / CommandHandle

# Interface: CommandHandle

Command execution handle.

It provides methods for waiting for the command to finish, retrieving stdout/stderr, and killing the command.

## Implements

- `Omit`\<[`CommandResult`](CommandResult.md), `"exitCode"` \| `"error"`\>
- `Partial`\<`Pick`\<[`CommandResult`](CommandResult.md), `"exitCode"` \| `"error"`\>\>

## Properties

### pid

> `readonly` **pid**: `number`

process ID of the command.

## Accessors

### error

#### Get Signature

> **get** **error**(): `string` \| `undefined`

Error message from command execution.

##### Returns

`string` \| `undefined`

#### Implementation of

[`CommandResult`](CommandResult.md).[`error`](CommandResult.md#error)

***

### exitCode

#### Get Signature

> **get** **exitCode**(): `number` \| `undefined`

Command execution exit code.
`0` if the command finished successfully.

It is `undefined` if the command is still running.

##### Returns

`number` \| `undefined`

#### Implementation of

[`CommandResult`](CommandResult.md).[`exitCode`](CommandResult.md#exitcode)

***

### stderr

#### Get Signature

> **get** **stderr**(): `string`

Command execution stderr output.

##### Returns

`string`

#### Implementation of

[`CommandResult`](CommandResult.md).[`stderr`](CommandResult.md#stderr)

***

### stdout

#### Get Signature

> **get** **stdout**(): `string`

Command execution stdout output.

##### Returns

`string`

#### Implementation of

[`CommandResult`](CommandResult.md).[`stdout`](CommandResult.md#stdout)

## Methods

### closeStdin()

> **closeStdin**(`opts?`): `Promise`\<`void`\>

Close the command stdin.

This signals EOF to the command. The command must have been started with
`stdin: true`.

#### Parameters

##### opts?

[`CommandRequestOpts`](CommandRequestOpts.md)

connection options.

#### Returns

`Promise`\<`void`\>

***

### disconnect()

> **disconnect**(): `Promise`\<`void`\>

Disconnect from the command.

The command is not killed, but SDK stops receiving events from the command.
You can reconnect to the command using [Commands.connect](Commands.md#connect).

Once it returns, the `onStdout`/`onStderr`/`onPty` callbacks are guaranteed
not to fire for output produced after this call. It does not wait for the
event handler to drain, so it returns promptly even for an idle command
whose stream produces no further output.

#### Returns

`Promise`\<`void`\>

***

### kill()

> **kill**(): `Promise`\<`boolean`\>

Kill the command.
It uses `SIGKILL` signal to kill the command.

#### Returns

`Promise`\<`boolean`\>

`true` if the command was killed successfully, `false` if the command was not found.

***

### sendStdin()

> **sendStdin**(`data`, `opts?`): `Promise`\<`void`\>

Send data to the command stdin.

The command must have been started with `stdin: true`.

#### Parameters

##### data

`string` \| `Uint8Array`\<`ArrayBufferLike`\>

data to send to the command.

##### opts?

[`CommandRequestOpts`](CommandRequestOpts.md)

connection options.

#### Returns

`Promise`\<`void`\>

***

### wait()

> **wait**(): `Promise`\<[`CommandResult`](CommandResult.md)\>

Wait for the command to finish and return the result.
If the command exits with a non-zero exit code, it throws a `CommandExitError`.

#### Returns

`Promise`\<[`CommandResult`](CommandResult.md)\>

`CommandResult` result of command execution.
