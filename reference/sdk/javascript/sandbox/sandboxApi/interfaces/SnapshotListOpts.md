[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SnapshotListOpts

# Interface: SnapshotListOpts

Options for listing snapshots.

## Extends

- `Omit`\<[`SandboxApiOpts`](SandboxApiOpts.md), `"signal"`\>

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

### limit?

> `optional` **limit?**: `number`

Number of snapshots to return per page.

#### Default

```ts
100
```

***

### name?

> `optional` **name?**: `string`

Filter snapshots by name or ID, optionally tag-qualified
(e.g. "my-snapshot", "my-project/my-snapshot" or "my-snapshot:v1").

***

### nextToken?

> `optional` **nextToken?**: `string`

Token to the next page.

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

### sandboxId?

> `optional` **sandboxId?**: `string`

Filter snapshots by source sandbox ID.
