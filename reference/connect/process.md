# Process API

Service: `Process`

## List

Public RPC exposed by envd.

- Request: `ListRequest`
- Response: `ListResponse`

## Connect

Public RPC exposed by envd.

- Request: `ConnectRequest`
- Response: `stream ConnectResponse`

## Start

Public RPC exposed by envd.

- Request: `StartRequest`
- Response: `stream StartResponse`

## Update

Public RPC exposed by envd.

- Request: `UpdateRequest`
- Response: `UpdateResponse`

## StreamInput

Client input stream ensures ordering of messages

- Request: `stream StreamInputRequest`
- Response: `StreamInputResponse`

## SendInput

Public RPC exposed by envd.

- Request: `SendInputRequest`
- Response: `SendInputResponse`

## SendSignal

Public RPC exposed by envd.

- Request: `SendSignalRequest`
- Response: `SendSignalResponse`

## CloseStdin

Only works for non-PTY processes. For PTY, send Ctrl+D (0x04) instead.

- Request: `CloseStdinRequest`
- Response: `CloseStdinResponse`

## Message types

### PTY

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `size` | `Size` | 1 |  |
| `cols` | `uint32` | 1 |  |
| `rows` | `uint32` | 2 |  |

### Size

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `cols` | `uint32` | 1 |  |
| `rows` | `uint32` | 2 |  |

### ProcessConfig

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `cmd` | `string` | 1 |  |
| `args` | `repeated string` | 2 |  |
| `cwd` | `optional string` | 4 |  |

### ListRequest

### ProcessInfo

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `config` | `ProcessConfig` | 1 |  |
| `pid` | `uint32` | 2 |  |
| `tag` | `optional string` | 3 |  |

### ListResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `processes` | `repeated ProcessInfo` | 1 |  |

### StartRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessConfig` | 1 |  |
| `pty` | `optional PTY` | 2 |  |
| `tag` | `optional string` | 3 |  |
| `stdin` | `optional bool` | 4 | We default to true. New SDK versions will set this to false by default. |

### UpdateRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |
| `pty` | `optional PTY` | 2 |  |

### UpdateResponse

### ProcessEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `start` | `StartEvent` | 1 |  |
| `data` | `DataEvent` | 2 |  |
| `end` | `EndEvent` | 3 |  |
| `keepalive` | `KeepAlive` | 4 |  |
| `pid` | `uint32` | 1 |  |
| `stdout` | `bytes` | 1 |  |
| `stderr` | `bytes` | 2 |  |
| `pty` | `bytes` | 3 |  |
| `exit_code` | `sint32` | 1 |  |
| `exited` | `bool` | 2 |  |
| `status` | `string` | 3 |  |
| `error` | `optional string` | 4 |  |

### StartEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `pid` | `uint32` | 1 |  |

### DataEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `stdout` | `bytes` | 1 |  |
| `stderr` | `bytes` | 2 |  |
| `pty` | `bytes` | 3 |  |

### EndEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `exit_code` | `sint32` | 1 |  |
| `exited` | `bool` | 2 |  |
| `status` | `string` | 3 |  |
| `error` | `optional string` | 4 |  |

### KeepAlive

### StartResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `event` | `ProcessEvent` | 1 |  |

### ConnectResponse

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `event` | `ProcessEvent` | 1 |  |

### SendInputRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |
| `input` | `ProcessInput` | 2 |  |

### SendInputResponse

### ProcessInput

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `stdin` | `bytes` | 1 |  |
| `pty` | `bytes` | 2 |  |

### StreamInputRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `start` | `StartEvent` | 1 |  |
| `data` | `DataEvent` | 2 |  |
| `keepalive` | `KeepAlive` | 3 |  |
| `process` | `ProcessSelector` | 1 |  |
| `input` | `ProcessInput` | 2 |  |

### StartEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |

### DataEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `input` | `ProcessInput` | 2 |  |

### KeepAlive

### StreamInputResponse

### SendSignalRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |
| `signal` | `Signal` | 2 |  |

### SendSignalResponse

### CloseStdinRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |

### CloseStdinResponse

### ConnectRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |

### ProcessSelector

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `pid` | `uint32` | 1 |  |
| `tag` | `string` | 2 |  |
