[@abox-dev/sdk](../README.md) / CommandStartOpts

# Interface: CommandStartOpts

Options for starting a new command.

## Extends

- [`CommandRequestOpts`](CommandRequestOpts.md)

## Properties

### background?

> `optional` **background?**: `boolean`

If true, starts command in the background and the method returns immediately.
You can use CommandHandle.wait to wait for the command to finish.

***

### cwd?

> `optional` **cwd?**: `string`

Working directory for the command.

#### Default

```ts
// home directory of the user used to start the command
```

***

### envs?

> `optional` **envs?**: `Record`\<`string`, `string`\>

Environment variables used for the command.

This overrides the default environment variables from `Sandbox` constructor.

#### Default

`{}`

***

### onStderr?

> `optional` **onStderr?**: (`data`) => `void` \| `Promise`\<`void`\>

Callback for command stderr output.

#### Parameters

##### data

`string`

#### Returns

`void` \| `Promise`\<`void`\>

***

### onStdout?

> `optional` **onStdout?**: (`data`) => `void` \| `Promise`\<`void`\>

Callback for command stdout output.

#### Parameters

##### data

`string`

#### Returns

`void` \| `Promise`\<`void`\>

***

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for requests to the API in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```

#### Inherited from

[`CommandRequestOpts`](CommandRequestOpts.md).[`requestTimeoutMs`](CommandRequestOpts.md#requesttimeoutms)

***

### signal?

> `optional` **signal?**: `AbortSignal`

An optional `AbortSignal` that can be used to cancel the in-flight request.
When the signal is aborted, the underlying `fetch` is aborted and the
returned promise rejects with an `AbortError`.

#### Inherited from

[`CommandRequestOpts`](CommandRequestOpts.md).[`signal`](CommandRequestOpts.md#signal)

***

### stdin?

> `optional` **stdin?**: `boolean`

If true, command stdin is kept open and you can send data to it using [Commands.sendStdin](../classes/Commands.md#sendstdin) or CommandHandle.sendStdin.

#### Default

```ts
false
```

***

### timeoutMs?

> `optional` **timeoutMs?**: `number`

Timeout for the command in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```

***

### user?

> `optional` **user?**: `string`

User to run the command as.

#### Default

`default Sandbox user (as specified in the template)`
