# POST /sandboxes/{sandboxID}/refreshes

Refresh sandbox

Refresh the sandbox extending its time to live

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |

## Request body

Required: no

### application/json

Schema: `SandboxRefreshRequest`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `duration` | `integer` | no | Duration for which the sandbox should be kept alive in seconds |

## Responses

### 204

Successfully refreshed the sandbox

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
