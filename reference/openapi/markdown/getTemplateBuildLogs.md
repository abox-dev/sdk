# GET /templates/{templateID}/builds/{buildID}/logs

Template build logs

Get template build logs

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `templateID` | path | yes | `string` |  |
| `buildID` | path | yes | `string` |  |
| `pageCursor` | query | no | `string` | Opaque continuation cursor returned as nextCursor by the previous page |
| `cursor` | query | no | `integer` | Starting timestamp of the logs that should be returned in milliseconds |
| `limit` | query | no | `integer` | Maximum number of logs that should be returned |
| `direction` | query | no | `LogsDirection` |  |
| `level` | query | no | `LogLevel` |  |
| `source` | query | no | `LogsSource` | Source of the logs that should be returned from |

## Responses

### 200

Successfully returned the template build logs

Content-Type: `application/json`

Schema: `TemplateBuildLogsResponse`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `logs` | `array<BuildLogEntry>` | yes | Build logs structured |
| `nextCursor` | `string` | no | Opaque continuation cursor for the next page |
| `source` | `LogsSource` | no | Source of the logs that should be returned |

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
