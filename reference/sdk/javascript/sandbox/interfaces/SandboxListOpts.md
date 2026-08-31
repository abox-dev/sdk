[agentbox-sdk-monorepo](../README.md) / SandboxListOpts

# Interface: SandboxListOpts

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

Number of sandboxes to return per page.

#### Default

```ts
100
```

***

### nextToken?

> `optional` **nextToken?**: `string`

Token to the next page.

***

### query?

> `optional` **query?**: `object`

Filter the list of sandboxes, e.g. by metadata `metadata:{"key": "value"}`, if there are multiple filters they are combined with AND.

#### metadata?

> `optional` **metadata?**: `Record`\<`string`, `string`\>

#### state?

> `optional` **state?**: [`SandboxState`](../type-aliases/SandboxState.md)[]

Filter the list of sandboxes by state.

##### Default

```ts
['running', 'paused']
```

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
