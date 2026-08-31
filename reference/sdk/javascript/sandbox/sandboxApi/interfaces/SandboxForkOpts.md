[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxForkOpts

# Interface: SandboxForkOpts

Options for forking a sandbox.

## Extends

- `ConnectionOpts`

## Properties

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

### apiUrl?

> `optional` **apiUrl?**: `string`

**`Internal`**

API Url to use for the API.

#### Default

AGENTBOX_API_URL // environment variable or `https://api.${domain}`

#### Inherited from

`ConnectionOpts.apiUrl`

***

### count?

> `optional` **count?**: `number`

Number of forked sandboxes to create.

All forks boot from the same snapshot — the snapshot is captured once
regardless of count. Each fork succeeds or fails independently; the
outcome of each is reported in its entry of the returned array.

#### Default

```ts
1
```

***

### debug?

> `optional` **debug?**: `boolean`

**`Internal`**

If true the SDK starts in the debug mode and connects to the local envd API server.

#### Default

AGENTBOX_DEBUG // environment variable or `false`

#### Inherited from

`ConnectionOpts.debug`

***

### domain?

> `optional` **domain?**: `string`

Domain to use for the API.

#### Default

AGENTBOX_DOMAIN // environment variable or `agentbox-runtime.ru`

#### Inherited from

`ConnectionOpts.domain`

***

### logger?

> `optional` **logger?**: `Logger`

Logger to use for logging messages. It can accept any object that implements `Logger` interface—for example, console.

#### Inherited from

`ConnectionOpts.logger`

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

**`Internal`**

Sandbox Url to use for the API.

#### Default

AGENTBOX_SANDBOX_URL // environment variable, `https://sandbox.${domain}`

#### Inherited from

`ConnectionOpts.sandboxUrl`

***

### signal?

> `optional` **signal?**: `AbortSignal`

An optional `AbortSignal` that can be used to cancel the in-flight request.
When the signal is aborted, the underlying `fetch` is aborted and the
returned promise rejects with an `AbortError`.

#### Inherited from

`ConnectionOpts.signal`

***

### timeoutMs?

> `optional` **timeoutMs?**: `number`

Timeout for the forked sandboxes in **milliseconds**.
Maximum time a sandbox can be kept alive is 24 hours (86_400_000 milliseconds) for Pro users and 1 hour (3_600_000 milliseconds) for Hobby users.

#### Default

```ts
300_000 // 5 minutes
```
