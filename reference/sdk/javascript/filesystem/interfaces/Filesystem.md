[agentbox-sdk-monorepo](../README.md) / Filesystem

# Interface: Filesystem

Module for interacting with the sandbox filesystem.

## Methods

### exists()

> **exists**(`path`, `opts?`): `Promise`\<`boolean`\>

Check if a file or a directory exists.

#### Parameters

##### path

`string`

path to a file or a directory

##### opts?

`FilesystemRequestOpts`

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the file or directory exists, `false` otherwise

***

### getInfo()

> **getInfo**(`path`, `opts?`): `Promise`\<[`EntryInfo`](EntryInfo.md)\>

Get information about a file or directory.

#### Parameters

##### path

`string`

path to a file or directory.

##### opts?

`FilesystemRequestOpts`

connection options.

#### Returns

`Promise`\<[`EntryInfo`](EntryInfo.md)\>

information about the file or directory like name, type, and path.

***

### list()

> **list**(`path`, `opts?`): `Promise`\<[`EntryInfo`](EntryInfo.md)[]\>

List entries in a directory.

#### Parameters

##### path

`string`

path to the directory.

##### opts?

`FilesystemListOpts`

connection options.

#### Returns

`Promise`\<[`EntryInfo`](EntryInfo.md)[]\>

list of entries in the sandbox filesystem directory.

***

### makeDir()

> **makeDir**(`path`, `opts?`): `Promise`\<`boolean`\>

Create a new directory and all directories along the way if needed on the specified path.

#### Parameters

##### path

`string`

path to a new directory. For example '/dirA/dirB' when creating 'dirB'.

##### opts?

`FilesystemRequestOpts`

connection options.

#### Returns

`Promise`\<`boolean`\>

`true` if the directory was created, `false` if it already exists.

***

### read()

#### Call Signature

> **read**(`path`, `opts?`): `Promise`\<`string`\>

Read file content as a `string`.

You can pass `text`, `bytes`, `blob`, or `stream` to `opts.format` to change the return type.

##### Parameters

###### path

`string`

path to the file.

###### opts?

[`FilesystemReadOpts`](FilesystemReadOpts.md) & `object`

connection options.

##### Returns

`Promise`\<`string`\>

file content as string

#### Call Signature

> **read**(`path`, `opts?`): `Promise`\<`Uint8Array`\<`ArrayBufferLike`\>\>

Read file content as a `Uint8Array`.

You can pass `text`, `bytes`, `blob`, or `stream` to `opts.format` to change the return type.

##### Parameters

###### path

`string`

path to the file.

###### opts?

[`FilesystemReadOpts`](FilesystemReadOpts.md) & `object`

connection options.

##### Returns

`Promise`\<`Uint8Array`\<`ArrayBufferLike`\>\>

file content as `Uint8Array`

#### Call Signature

> **read**(`path`, `opts?`): `Promise`\<`Blob`\>

Read file content as a `Blob`.

You can pass `text`, `bytes`, `blob`, or `stream` to `opts.format` to change the return type.

##### Parameters

###### path

`string`

path to the file.

###### opts?

[`FilesystemReadOpts`](FilesystemReadOpts.md) & `object`

connection options.

##### Returns

`Promise`\<`Blob`\>

file content as `Blob`

#### Call Signature

> **read**(`path`, `opts?`): `Promise`\<`ReadableStream`\<`Uint8Array`\<`ArrayBufferLike`\>\>\>

Read file content as a `ReadableStream`.

You can pass `text`, `bytes`, `blob`, or `stream` to `opts.format` to change the return type.

The request timeout bounds only the initial handshake. The returned stream
holds a pooled connection until it is fully read, cancelled, errors, or the
idle timeout (`opts.streamIdleTimeoutMs`) fires—so consume it to the end or
cancel it (`opts.signal`).

##### Parameters

###### path

`string`

path to the file.

###### opts?

[`FilesystemReadOpts`](FilesystemReadOpts.md) & `object`

connection options.

##### Returns

`Promise`\<`ReadableStream`\<`Uint8Array`\<`ArrayBufferLike`\>\>\>

file content as `ReadableStream`

***

### remove()

> **remove**(`path`, `opts?`): `Promise`\<`void`\>

Remove a file or directory.

#### Parameters

##### path

`string`

path to a file or directory.

##### opts?

`FilesystemRequestOpts`

connection options.

#### Returns

`Promise`\<`void`\>

***

### rename()

> **rename**(`oldPath`, `newPath`, `opts?`): `Promise`\<[`EntryInfo`](EntryInfo.md)\>

Rename a file or directory.

#### Parameters

##### oldPath

`string`

path to the file or directory to rename.

##### newPath

`string`

new path for the file or directory.

##### opts?

`FilesystemRequestOpts`

connection options.

#### Returns

`Promise`\<[`EntryInfo`](EntryInfo.md)\>

information about renamed file or directory.

***

### watchDir()

> **watchDir**(`path`, `onEvent`, `opts?`): `Promise`\<[`WatchHandle`](WatchHandle.md)\>

Start watching a directory for filesystem events.

#### Parameters

##### path

`string`

path to directory to watch.

##### onEvent

(`event`) => `void` \| `Promise`\<`void`\>

callback to call when an event in the directory occurs.

##### opts?

`WatchOpts` & `object`

connection options.

#### Returns

`Promise`\<[`WatchHandle`](WatchHandle.md)\>

`WatchHandle` object for stopping watching directory.

***

### write()

#### Call Signature

> **write**(`path`, `data`, `opts?`): `Promise`\<[`WriteInfo`](WriteInfo.md)\>

Write content to a file.

Writing to a file that doesn't exist creates the file.

Writing to a file that already exists overwrites the file.

Writing to a file at path that doesn't exist creates the necessary directories.

##### Parameters

###### path

`string`

path to file.

###### data

`string` \| `Blob` \| `ReadableStream`\<`any`\> \| `ArrayBuffer`

data to write to the file. Data can be a string, `ArrayBuffer`, `Blob`, or `ReadableStream`.

###### opts?

[`FilesystemWriteOpts`](FilesystemWriteOpts.md)

connection options.

##### Returns

`Promise`\<[`WriteInfo`](WriteInfo.md)\>

information about the written file

#### Call Signature

> **write**(`files`, `opts?`): `Promise`\<[`WriteInfo`](WriteInfo.md)[]\>

Write content to a file.

Writing to a file that doesn't exist creates the file.

Writing to a file that already exists overwrites the file.

Writing to a file at path that doesn't exist creates the necessary directories.

##### Parameters

###### files

`WriteEntry`[]

###### opts?

[`FilesystemWriteOpts`](FilesystemWriteOpts.md)

connection options.

##### Returns

`Promise`\<[`WriteInfo`](WriteInfo.md)[]\>

information about the written file

***

### writeFiles()

> **writeFiles**(`files`, `opts?`): `Promise`\<[`WriteInfo`](WriteInfo.md)[]\>

Write multiple files.

Writing to a file that doesn't exist creates the file.

Writing to a file that already exists overwrites the file.

Writing to a file at path that doesn't exist creates the necessary directories.

#### Parameters

##### files

`WriteEntry`[]

list of files to write as `WriteEntry` objects, each containing `path` and `data`.

##### opts?

[`FilesystemWriteOpts`](FilesystemWriteOpts.md)

connection options.

#### Returns

`Promise`\<[`WriteInfo`](WriteInfo.md)[]\>

information about the written files
