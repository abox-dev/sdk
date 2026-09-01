[agentbox-sdk-monorepo](../README.md) / SandboxOpts

# Interface: SandboxOpts

Options for creating a new Sandbox.

## Extends

- `ConnectionOpts`

## Properties

### allowInternetAccess?

> `optional` **allowInternetAccess?**: `boolean`

Allow sandbox to access the internet. If set to `False`, it works the same as setting network `denyOut` to `[0.0.0.0/0]`.

#### Default

```ts
true
```

***

### apiHeaders?

> `optional` **apiHeaders?**: `Record`\<`string`, `string`\>

Additional headers to send with AgentBox API requests.

#### Inherited from

`ConnectionOpts.apiHeaders`

***

### apiKey?

> `optional` **apiKey?**: `string`

AgentBox API key to use for authentication.

#### Default

```ts
AGENTBOX_API_KEY // environment variable
```

#### Inherited from

`ConnectionOpts.apiKey`

***

### domain?

> `optional` **domain?**: `string`

Domain to use for the API.

#### Default

AGENTBOX_DOMAIN // environment variable or `agentbox-runtime.ru`

#### Inherited from

`ConnectionOpts.domain`

***

### envs?

> `optional` **envs?**: `Record`\<`string`, `string`\>

Custom environment variables for the sandbox.

Used when executing commands and code in the sandbox.
Can be overridden with the `envs` argument when executing commands or code.

#### Default

```ts
{}
```

***

### iam?

> `optional` **iam?**: [`SandboxIamOpts`](SandboxIamOpts.md)

Sandbox workload identity configuration. Providing a non-empty
`tokens` map enables workload identity for the sandbox.

Registered tokens are exposed to [SandboxNetworkOpts.rules](../type-aliases/SandboxNetworkOpts.md#rules)
`transform` callbacks as `iam.tokens.<name>` placeholders, which the egress
proxy resolves per request.

#### Example

```ts
const sandbox = await Sandbox.create({
  iam: {
    tokens: {
      aws: {
        audience: 'sts.amazonaws.com',
        tokenType: 'JWT-SVID',
      },
    },
  },
})
```

***

### lifecycle?

> `optional` **lifecycle?**: [`SandboxLifecycle`](../type-aliases/SandboxLifecycle.md)

Sandbox lifecycle configuration.

***

### logger?

> `optional` **logger?**: `Logger`

Logger to use for logging messages. It can accept any object that implements `Logger` interface—for example, console.

#### Inherited from

`ConnectionOpts.logger`

***

### mcp?

> `optional` **mcp?**: `McpServer`

MCP server to enable in the sandbox

#### Default

```ts
undefined
```

***

### metadata?

> `optional` **metadata?**: `Record`\<`string`, `string`\>

Custom metadata for the sandbox.

#### Default

```ts
{}
```

***

### network?

> `optional` **network?**: [`SandboxNetworkOpts`](../type-aliases/SandboxNetworkOpts.md)

Sandbox network configuration

***

### proxy?

> `optional` **proxy?**: `string`

Proxy URL to use for requests. In case of a sandbox it applies to all
requests made to the returned sandbox.

#### Example

```ts
'http://user:pass@127.0.0.1:8080'
```

#### Inherited from

`ConnectionOpts.proxy`

***

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for requests to the API in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```

#### Inherited from

`ConnectionOpts.requestTimeoutMs`

***

### sandboxUrl?

> `optional` **sandboxUrl?**: `string`

Sandbox URL. Used for local development

#### Overrides

`ConnectionOpts.sandboxUrl`

***

### secure?

> `optional` **secure?**: `boolean`

Secure all traffic coming to the sandbox controller with auth token

#### Default

```ts
true
```

***

### signal?

> `optional` **signal?**: `AbortSignal`

An optional `AbortSignal` that can be used to cancel the in-flight request.
When the signal is aborted, the underlying `fetch` is aborted and the
returned promise rejects with an `AbortError`.

#### Inherited from

`ConnectionOpts.signal`

***

### template?

> `optional` **template?**: `string`

Sandbox template name or ID.

#### Default

'base' (or 'mcp-gateway' when `mcp` option is set)

***

### timeoutMs?

> `optional` **timeoutMs?**: `number`

Timeout for the sandbox in **milliseconds**.
Maximum time a sandbox can be kept alive is 24 hours (86_400_000 milliseconds) for Pro users and 1 hour (3_600_000 milliseconds) for Hobby users.

#### Default

```ts
300_000 // 5 minutes
```
