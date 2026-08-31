[@abox-dev/code-interpreter](../../README.md) / [sandbox](../README.md) / RunCodeOpts

# Interface: RunCodeOpts

Options for running code.

## Properties

### envs?

> `optional` **envs?**: `Record`\<`string`, `string`\>

Custom environment variables for code execution.

#### Default

```ts
{}
```

***

### onError?

> `optional` **onError?**: (`error`) => `any`

Callback for handling the `ExecutionError` object.

#### Parameters

##### error

[`ExecutionError`](../../messaging/classes/ExecutionError.md)

#### Returns

`any`

***

### onResult?

> `optional` **onResult?**: (`data`) => `any`

Callback for handling the final execution result.

#### Parameters

##### data

[`Result`](../../messaging/classes/Result.md)

#### Returns

`any`

***

### onStderr?

> `optional` **onStderr?**: (`output`) => `any`

Callback for handling stderr messages.

#### Parameters

##### output

[`OutputMessage`](../../messaging/classes/OutputMessage.md)

#### Returns

`any`

***

### onStdout?

> `optional` **onStdout?**: (`output`) => `any`

Callback for handling stdout messages.

#### Parameters

##### output

[`OutputMessage`](../../messaging/classes/OutputMessage.md)

#### Returns

`any`

***

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for the request in **milliseconds**.

#### Default

```ts
30_000 // 30 seconds
```

***

### timeoutMs?

> `optional` **timeoutMs?**: `number`

Timeout for the code execution in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```
