[@abox-dev/sdk](../README.md) / FilesystemRequestOpts

# Interface: FilesystemRequestOpts

Options for the sandbox filesystem operations.

## Extends

- `Partial`\<`Pick`\<`ConnectionOpts`, `"requestTimeoutMs"` \| `"signal"`\>\>

## Extended by

- [`FilesystemWriteOpts`](FilesystemWriteOpts.md)
- [`FilesystemReadOpts`](FilesystemReadOpts.md)
- [`FilesystemListOpts`](FilesystemListOpts.md)
- [`WatchOpts`](WatchOpts.md)

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

***

### user?

> `optional` **user?**: `string`

User to use for the operation in the sandbox.
This affects the resolution of relative paths and ownership of the created filesystem objects.
