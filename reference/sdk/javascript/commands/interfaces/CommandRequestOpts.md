[@abox-dev/sdk](../README.md) / CommandRequestOpts

# Interface: CommandRequestOpts

Options for sending a command request.

## Extends

- `Partial`\<`Pick`\<`ConnectionOpts`, `"requestTimeoutMs"` \| `"signal"`\>\>

## Extended by

- [`CommandStartOpts`](CommandStartOpts.md)

## Properties

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for requests to the API in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```

#### Inherited from

`Partial.requestTimeoutMs`

***

### signal?

> `optional` **signal?**: `AbortSignal`

An optional `AbortSignal` that can be used to cancel the in-flight request.
When the signal is aborted, the underlying `fetch` is aborted and the
returned promise rejects with an `AbortError`.

#### Inherited from

`Partial.signal`
