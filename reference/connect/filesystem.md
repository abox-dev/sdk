# Filesystem API

Service: `Filesystem`

## Stat

Public RPC exposed by envd.

- Request: `StatRequest`
- Response: `StatResponse`

## MakeDir

Public RPC exposed by envd.

- Request: `MakeDirRequest`
- Response: `MakeDirResponse`

## Move

Public RPC exposed by envd.

- Request: `MoveRequest`
- Response: `MoveResponse`

## ListDir

Public RPC exposed by envd.

- Request: `ListDirRequest`
- Response: `ListDirResponse`

## Remove

Public RPC exposed by envd.

- Request: `RemoveRequest`
- Response: `RemoveResponse`

## WatchDir

Public RPC exposed by envd.

- Request: `WatchDirRequest`
- Response: `stream WatchDirResponse`

## CreateWatcher

Non-streaming versions of WatchDir

- Request: `CreateWatcherRequest`
- Response: `CreateWatcherResponse`

## GetWatcherEvents

Public RPC exposed by envd.

- Request: `GetWatcherEventsRequest`
- Response: `GetWatcherEventsResponse`

## RemoveWatcher

Public RPC exposed by envd.

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
| `allow_network_mounts` | `bool` | 4 | Events on network mounts may be unreliable or not delivered at all. |

### FilesystemEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `name` | `string` | 1 |  |
| `type` | `EventType` | 2 |  |
| `entry` | `optional EntryInfo` | 3 | events, where the entry no longer exists at this path). |

### WatchDirResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `start` | `StartEvent` | 1 |  |
| `filesystem` | `FilesystemEvent` | 2 |  |
| `keepalive` | `KeepAlive` | 3 |  |

### StartEvent

### KeepAlive

### CreateWatcherRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `path` | `string` | 1 |  |
| `recursive` | `bool` | 2 |  |
| `include_entry` | `bool` | 3 | If true, each FilesystemEvent includes the EntryInfo of the affected entry, when available. |
| `allow_network_mounts` | `bool` | 4 | Events on network mounts may be unreliable or not delivered at all. |

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
