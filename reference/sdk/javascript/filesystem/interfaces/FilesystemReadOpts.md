[@abox-dev/sdk](../README.md) / FilesystemReadOpts

# Interface: FilesystemReadOpts

Options for reading files from the sandbox filesystem.

## Extends

- [`FilesystemRequestOpts`](FilesystemRequestOpts.md)

## Properties

### gzip?

> `optional` **gzip?**: `boolean`

When true, the download will request gzip-encoded responses.

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

### streamIdleTimeoutMs?

> `optional` **streamIdleTimeoutMs?**: `number`

Idle timeout for a streamed read (`format: 'stream'`) in **milliseconds**:
abort if no chunk arrives from the server within this window *while
reading*. It bounds only the wire — a slow or paused consumer never trips
it (a consumer that holds the stream but stops reading is reclaimed
server-side). Defaults to the request timeout (60s); pass `0` to disable.

***

### user?

> `optional` **user?**: `string`

User to use for the operation in the sandbox.
This affects the resolution of relative paths and ownership of the created filesystem objects.

#### Inherited from

[`FilesystemRequestOpts`](FilesystemRequestOpts.md).[`user`](FilesystemRequestOpts.md#user)
