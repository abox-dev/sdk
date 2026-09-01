# Process API

## Transport

AgentBox exposes this service through the Connect protocol over HTTP.
The AgentBox SDK supplies the routing and authorization headers automatically.

- Production base URL: `https://sandbox.agentbox-runtime.ru`
- Fully qualified service: `process.Process`
- RPC URL pattern: `POST https://sandbox.agentbox-runtime.ru/process.Process/{RPC}`

### Request headers

- **`Agentbox-Sandbox-Id`** · required — Sandbox identifier.
- **`Agentbox-Sandbox-Port`** · required — Envd port routed by the sandbox proxy. Default: `49983`.
- **`X-Access-Token`** · conditional — Sandbox-scoped envd access token, when one was issued.

## List

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/List`
- Request: `ListRequest`
- Response: `ListResponse`

## Connect

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/Connect`
- Request: `ConnectRequest`
- Response: `stream ConnectResponse`

## Start

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/Start`
- Request: `StartRequest`
- Response: `stream StartResponse`

## Update

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/Update`
- Request: `UpdateRequest`
- Response: `UpdateResponse`

## StreamInput

Client input stream ensures ordering of messages

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/StreamInput`
- Request: `stream StreamInputRequest`
- Response: `StreamInputResponse`

## SendInput

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/SendInput`
- Request: `SendInputRequest`
- Response: `SendInputResponse`

## SendSignal

Public RPC exposed by envd.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/SendSignal`
- Request: `SendSignalRequest`
- Response: `SendSignalResponse`

## CloseStdin

Close stdin to signal EOF to the process. Only works for non-PTY processes. For PTY, send Ctrl+D (0x04) instead.

- Endpoint: `POST https://sandbox.agentbox-runtime.ru/process.Process/CloseStdin`
- Request: `CloseStdinRequest`
- Response: `CloseStdinResponse`

## Message types

### PTY

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `size` | `Size` | 1 |  |

### PTY.Size

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `cols` | `uint32` | 1 |  |
| `rows` | `uint32` | 2 |  |

### ProcessConfig

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `cmd` | `string` | 1 |  |
| `args` | `repeated string` | 2 |  |
| `envs` | `map<string, string>` | 3 |  |
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
| `stdin` | `optional bool` | 4 | This is optional for backwards compatibility. We default to true. New SDK versions will set this to false by default. |

### UpdateRequest

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |
| `pty` | `optional PTY` | 2 |  |

### UpdateResponse

### ProcessEvent

| Field | Type | Number | Oneof | Description |
| --- | --- | ---: | --- | --- |
| `start` | `StartEvent` | 1 | `event` |  |
| `data` | `DataEvent` | 2 | `event` |  |
| `end` | `EndEvent` | 3 | `event` |  |
| `keepalive` | `KeepAlive` | 4 | `event` |  |

### ProcessEvent.StartEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `pid` | `uint32` | 1 |  |

### ProcessEvent.DataEvent

| Field | Type | Number | Oneof | Description |
| --- | --- | ---: | --- | --- |
| `stdout` | `bytes` | 1 | `output` |  |
| `stderr` | `bytes` | 2 | `output` |  |
| `pty` | `bytes` | 3 | `output` |  |

### ProcessEvent.EndEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `exit_code` | `sint32` | 1 |  |
| `exited` | `bool` | 2 |  |
| `status` | `string` | 3 |  |
| `error` | `optional string` | 4 |  |

### ProcessEvent.KeepAlive

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

| Field | Type | Number | Oneof | Description |
| --- | --- | ---: | --- | --- |
| `stdin` | `bytes` | 1 | `input` |  |
| `pty` | `bytes` | 2 | `input` |  |

### StreamInputRequest

| Field | Type | Number | Oneof | Description |
| --- | --- | ---: | --- | --- |
| `start` | `StartEvent` | 1 | `event` |  |
| `data` | `DataEvent` | 2 | `event` |  |
| `keepalive` | `KeepAlive` | 3 | `event` |  |

### StreamInputRequest.StartEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `process` | `ProcessSelector` | 1 |  |

### StreamInputRequest.DataEvent

| Field | Type | Number | Description |
| --- | --- | ---: | --- |
| `input` | `ProcessInput` | 2 |  |

### StreamInputRequest.KeepAlive

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

| Field | Type | Number | Oneof | Description |
| --- | --- | ---: | --- | --- |
| `pid` | `uint32` | 1 | `selector` |  |
| `tag` | `string` | 2 | `selector` |  |

## Enum types

### Signal

| Value | Number | Description |
| --- | ---: | --- |
| `SIGNAL_UNSPECIFIED` | 0 |  |
| `SIGNAL_SIGTERM` | 15 |  |
| `SIGNAL_SIGKILL` | 9 |  |
