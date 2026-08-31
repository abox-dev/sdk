# GET /sandboxes/metrics

List sandbox metrics

List metrics for given sandboxes

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandbox_ids` | query | yes | `array<string>` | Comma-separated list of sandbox IDs to get metrics for |

## Responses

### 200

Successfully returned all running sandboxes with metrics

Content-Type: `application/json`

Schema: `SandboxesWithMetrics`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `sandboxes` | `object` | yes |  |

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 400

Bad request

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
