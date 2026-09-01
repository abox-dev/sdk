# Filesystem API

## Transport

AgentBox exposes this service through the Connect protocol over HTTP.
The AgentBox SDK supplies the routing and authorization headers automatically.

- Production base URL: `https://sandbox.agentbox-runtime.ru`
- Fully qualified service: `filesystem.Filesystem`
- RPC URL pattern: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/{RPC}`

### Request headers

- **`Agentbox-Sandbox-Id`** · required — Sandbox identifier.
- **`Agentbox-Sandbox-Port`** · required — Envd port routed by the sandbox proxy. Default: `49983`.
- **`X-Access-Token`** · conditional — Sandbox-scoped envd access token, when one was issued.

## Stat

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/Stat`
- Request: `StatRequest`
- Response: `StatResponse`

## MakeDir

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/MakeDir`
- Request: `MakeDirRequest`
- Response: `MakeDirResponse`

## Move

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/Move`
- Request: `MoveRequest`
- Response: `MoveResponse`

## ListDir

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/ListDir`
- Request: `ListDirRequest`
- Response: `ListDirResponse`

## Remove

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/Remove`
- Request: `RemoveRequest`
- Response: `RemoveResponse`

## WatchDir

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/WatchDir`
- Request: `WatchDirRequest`
- Response: `stream WatchDirResponse`

## CreateWatcher

Non-streaming versions of WatchDir

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/CreateWatcher`
- Request: `CreateWatcherRequest`
- Response: `CreateWatcherResponse`

## GetWatcherEvents

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/GetWatcherEvents`
- Request: `GetWatcherEventsRequest`
- Response: `GetWatcherEventsResponse`

## RemoveWatcher

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/filesystem.Filesystem/RemoveWatcher`
- Request: `RemoveWatcherRequest`
- Response: `RemoveWatcherResponse`

## Message types

### MoveRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `source` | `string` | 1 |  |
| `destination` | `string` | 2 |  |

### MoveResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `entry` | `EntryInfo` | 1 |  |

### MakeDirRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |

### MakeDirResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `entry` | `EntryInfo` | 1 |  |

### RemoveRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |

### RemoveResponse

### StatRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |

### StatResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `entry` | `EntryInfo` | 1 |  |

### EntryInfo

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `name` | `string` | 1 |  |
| `type` | `FileType` | 2 |  |
| `path` | `string` | 3 |  |
| `size` | `int64` | 4 |  |
| `mode` | `uint32` | 5 |  |
| `permissions` | `string` | 6 |  |
| `owner` | `string` | 7 |  |
| `group` | `string` | 8 |  |
| `modified_time` | `google.protobuf.Timestamp` | 9 |  |
| `symlink_target` | `optional string` | 10 | If the entry is a symlink, this field contains the target of the symlink. |
| `metadata` | `map<string, string>` | 11 | User-defined metadata stored as extended attributes (xattrs) on the file. Keys live under the `user.agentbox.` xattr namespace; the prefix is stripped here. Plain `user.*` xattrs written by other tooling are not reflected. |

### ListDirRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |
| `depth` | `uint32` | 2 |  |

### ListDirResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `entries` | `repeated EntryInfo` | 1 |  |

### WatchDirRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |
| `recursive` | `bool` | 2 |  |
| `include_entry` | `bool` | 3 | If true, each FilesystemEvent includes the EntryInfo of the affected entry, when available. |
| `allow_network_mounts` | `bool` | 4 | If true, allows watching paths on network filesystem mounts (NFS, CIFS, SMB, FUSE). Events on network mounts may be unreliable or not delivered at all. |

### FilesystemEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `name` | `string` | 1 |  |
| `type` | `EventType` | 2 |  |
| `entry` | `optional EntryInfo` | 3 | Info of the entry that triggered the event. Only populated when include_entry was requested and the entry could be stat-ed (e.g. not set for remove/rename-away events, where the entry no longer exists at this path). |

### WatchDirResponse

| Field | Type | Number | Oneof | Description |
| --- | --- | ---: | --- | --- |
| `start` | `StartEvent` | 1 | `event` |  |
| `filesystem` | `FilesystemEvent` | 2 | `event` |  |
| `keepalive` | `KeepAlive` | 3 | `event` |  |

### WatchDirResponse.StartEvent

### WatchDirResponse.KeepAlive

### CreateWatcherRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |
| `recursive` | `bool` | 2 |  |
| `include_entry` | `bool` | 3 | If true, each FilesystemEvent includes the EntryInfo of the affected entry, when available. |
| `allow_network_mounts` | `bool` | 4 | If true, allows watching paths on network filesystem mounts (NFS, CIFS, SMB, FUSE). Events on network mounts may be unreliable or not delivered at all. |

### CreateWatcherResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `watcher_id` | `string` | 1 |  |

### GetWatcherEventsRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `watcher_id` | `string` | 1 |  |

### GetWatcherEventsResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `events` | `repeated FilesystemEvent` | 1 |  |

### RemoveWatcherRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `watcher_id` | `string` | 1 |  |

### RemoveWatcherResponse

## Enum types

### FileType

| Value | Number | Description |
| --- | ---: | --- |
| `FILE_TYPE_UNSPECIFIED` | 0 |  |
| `FILE_TYPE_FILE` | 1 |  |
| `FILE_TYPE_DIRECTORY` | 2 |  |
| `FILE_TYPE_SYMLINK` | 3 |  |

### EventType

| Value | Number | Description |
| --- | ---: | --- |
| `EVENT_TYPE_UNSPECIFIED` | 0 |  |
| `EVENT_TYPE_CREATE` | 1 |  |
| `EVENT_TYPE_WRITE` | 2 |  |
| `EVENT_TYPE_REMOVE` | 3 |  |
| `EVENT_TYPE_RENAME` | 4 |  |
| `EVENT_TYPE_CHMOD` | 5 |  |
