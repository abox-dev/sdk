# GET /templates/{templateID}/builds/{buildID}/status

Template build status

Get template build info

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `templateID` | path | yes | `string` |  |
| `buildID` | path | yes | `string` |  |
| `logsOffset` | query | no | `integer` | Index of the starting build log that should be returned with the template |
| `limit` | query | no | `integer` | Maximum number of logs that should be returned |
| `level` | query | no | `LogLevel` |  |

## Responses

### 200

Successfully returned the template

Content-Type: `application/json`

Schema: `TemplateBuildInfo`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `logs` | `array<string>` | yes | Build logs |
| `logEntries` | `array<BuildLogEntry>` | yes | Build logs structured |
| `templateID` | `string` | yes | Identifier of the template |
| `buildID` | `string` | yes | Identifier of the build |
| `status` | `TemplateBuildStatus` | yes | Status of the template build |
| `reason` | `BuildStatusReason` | no |  |

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
