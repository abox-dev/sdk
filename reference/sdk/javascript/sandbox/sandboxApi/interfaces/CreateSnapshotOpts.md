[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / CreateSnapshotOpts

# Interface: CreateSnapshotOpts

Options for creating a snapshot.

## Extends

- [`SandboxApiOpts`](SandboxApiOpts.md)

## Properties

### apiHeaders?

> `optional` **apiHeaders?**: `Record`\<`string`, `string`\>

Additional headers to send with AgentBox API requests.

#### Inherited from

[`SandboxForkOpts`](SandboxForkOpts.md).[`apiHeaders`](SandboxForkOpts.md#apiheaders)

***

### apiKey?

> `optional` **apiKey?**: `string`

AgentBox API key to use for authentication.

#### Default

```ts
AGENTBOX_API_KEY // environment variable
```

#### Inherited from

[`SandboxForkOpts`](SandboxForkOpts.md).[`apiKey`](SandboxForkOpts.md#apikey)

***

### debug?

> `optional` **debug?**: `boolean`

**`Internal`**

If true the SDK starts in the debug mode and connects to the local envd API server.

#### Default

AGENTBOX_DEBUG // environment variable or `false`

#### Inherited from

[`SandboxForkOpts`](SandboxForkOpts.md).[`debug`](SandboxForkOpts.md#debug)

***

### domain?

> `optional` **domain?**: `string`

Domain to use for the API.

#### Default

AGENTBOX_DOMAIN // environment variable or `agentbox-runtime.ru`

#### Inherited from

[`SandboxForkOpts`](SandboxForkOpts.md).[`domain`](SandboxForkOpts.md#domain)

***

### name?

> `optional` **name?**: `string`

Optional name for the snapshot template.
If a snapshot template with this name already exists, a new build will be assigned
to the existing template instead of creating a new one.

***

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for requests to the API in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```

#### Inherited from

[`SandboxForkOpts`](SandboxForkOpts.md).[`requestTimeoutMs`](SandboxForkOpts.md#requesttimeoutms)

***

### signal?

> `optional` **signal?**: `AbortSignal`

An optional `AbortSignal` that can be used to cancel the in-flight request.
When the signal is aborted, the underlying `fetch` is aborted and the
returned promise rejects with an `AbortError`.

#### Inherited from

[`SandboxForkOpts`](SandboxForkOpts.md).[`signal`](SandboxForkOpts.md#signal)
