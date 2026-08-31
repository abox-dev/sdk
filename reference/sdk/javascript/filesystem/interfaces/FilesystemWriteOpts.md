[@abox-dev/sdk](../README.md) / FilesystemWriteOpts

# Interface: FilesystemWriteOpts

Options for writing files to the sandbox filesystem.

## Extends

- [`FilesystemRequestOpts`](FilesystemRequestOpts.md)

## Properties

### gzip?

> `optional` **gzip?**: `boolean`

When true, the upload will be gzip-compressed. Implies the
`application/octet-stream` upload.

Requires envd 0.5.7 or later — when not supported by the sandbox's envd
version, the upload falls back to uncompressed `multipart/form-data`.

***

### metadata?

> `optional` **metadata?**: `Record`\<`string`, `string`\>

User-defined metadata to persist on the uploaded file(s) as extended
attributes. Keys are lowercased by the sandbox, so they may differ in case
when read back. Invalid keys or values throw an `InvalidArgumentError`.
The same metadata is applied to every file in a multi-file upload.
Requires envd 0.6.2 or later.

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

### useOctetStream?

> `optional` **useOctetStream?**: `boolean`

When true, the upload uses `application/octet-stream` instead of `multipart/form-data`.
Outside the browser, `ReadableStream` data is then streamed to the sandbox
instead of being buffered in memory.

Defaults to `undefined`, which uses octet-stream when any entry is a
`ReadableStream` (so streamed uploads aren't buffered) and
`multipart/form-data` otherwise; browsers always use `multipart/form-data`
since they can't stream request bodies. Requires envd 0.5.7 or later — when
not supported by the sandbox's envd version, the upload falls back to
`multipart/form-data`.

***

### user?

> `optional` **user?**: `string`

User to use for the operation in the sandbox.
This affects the resolution of relative paths and ownership of the created filesystem objects.

#### Inherited from

[`FilesystemRequestOpts`](FilesystemRequestOpts.md).[`user`](FilesystemRequestOpts.md#user)
