[@abox-dev/sdk](../README.md) / FilesystemListOpts

# Interface: FilesystemListOpts

Options for the sandbox filesystem operations.

## Extends

- [`FilesystemRequestOpts`](FilesystemRequestOpts.md)

## Properties

### depth?

> `optional` **depth?**: `number`

Depth of the directory to list.

***

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for requests to the API in **milliseconds**.

#### Default

```ts
60_000 // 60 seconds
```

#### Inherited from

[`FilesystemRequestOpts`](FilesystemRequestOpts.md).[`requestTimeoutMs`](FilesystemRequestOpts.md#requesttimeoutms)

***

### signal?

> `optional` **signal?**: `AbortSignal`

An optional `AbortSignal` that can be used to cancel the in-flight request.
When the signal is aborted, the underlying `fetch` is aborted and the
returned promise rejects with an `AbortError`.

#### Inherited from

[`FilesystemRequestOpts`](FilesystemRequestOpts.md).[`signal`](FilesystemRequestOpts.md#signal)

***

### user?

> `optional` **user?**: `string`

User to use for the operation in the sandbox.
This affects the resolution of relative paths and ownership of the created filesystem objects.

#### Inherited from

[`FilesystemRequestOpts`](FilesystemRequestOpts.md).[`user`](FilesystemRequestOpts.md#user)
