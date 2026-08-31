# GET /v2/sandboxes/{sandboxID}/logs

Sandbox logs (v2)

Get sandbox logs

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |
| `pageCursor` | query | no | `string` | Opaque continuation cursor returned as nextCursor by the previous page |
| `cursor` | query | no | `integer` | Starting timestamp of the logs that should be returned in milliseconds |
| `limit` | query | no | `integer` | Maximum number of logs that should be returned |
| `direction` | query | no | `LogsDirection` | Direction of the logs that should be returned |
| `level` | query | no | `LogLevel` | Minimum log level to return. Logs below this level are excluded |
| `search` | query | no | `string` | Case-sensitive substring match on log message content |

## Responses

### 200

Successfully returned the sandbox logs

Content-Type: `application/json`

Schema: `SandboxLogsV2Response`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `logs` | `array<SandboxLogEntry>` | yes | Sandbox logs structured |
| `nextCursor` | `string` | no | Opaque continuation cursor for the next page |

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 404

Not found

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |
