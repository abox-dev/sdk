[@abox-dev/sdk](../README.md) / WatchOpts

# Interface: WatchOpts

Options for watching a directory.

## Extends

- [`FilesystemRequestOpts`](FilesystemRequestOpts.md)

## Properties

### allowNetworkMounts?

> `optional` **allowNetworkMounts?**: `boolean`

Allow watching paths on network filesystem mounts (NFS, CIFS, SMB, FUSE),
which are rejected by default. Events on network mounts may be unreliable
or not delivered at all.

Requires envd 0.6.4 or later. Watching with this option against an older sandbox
throws a `TemplateError`.

***

### includeEntry?

> `optional` **includeEntry?**: `boolean`

Include the [EntryInfo](EntryInfo.md) of the affected entry in each FilesystemEvent.

The entry is populated best-effort and may be `undefined` for events where the
entry no longer exists at the path (e.g. remove or rename-away events).

Requires envd 0.6.3 or later. Watching with this option against an older sandbox
throws a `TemplateError`.

***

### onExit?

> `optional` **onExit?**: (`err?`) => `void` \| `Promise`\<`void`\>

Callback to call when the watch operation stops.

#### Parameters

##### err?

`Error`

#### Returns

`void` \| `Promise`\<`void`\>

***

### recursive?

> `optional` **recursive?**: `boolean`

Watch the directory recursively

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

### timeoutMs?

> `optional` **timeoutMs?**: `number`

Timeout for the watch operation in **milliseconds**.
You can pass `0` to disable the timeout.

#### Default

```ts
60_000 // 60 seconds
```

***

### user?

> `optional` **user?**: `string`

User to use for the operation in the sandbox.
This affects the resolution of relative paths and ownership of the created filesystem objects.

#### Inherited from

[`FilesystemRequestOpts`](FilesystemRequestOpts.md).[`user`](FilesystemRequestOpts.md#user)
