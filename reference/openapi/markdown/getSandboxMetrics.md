# GET /sandboxes/{sandboxID}/metrics

Sandbox metrics

Get sandbox metrics

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |
| `start` | query | no | `integer` | Unix timestamp for the start of the interval, in seconds, for which the metrics |
| `end` | query | no | `integer` |  |

## Responses

### 200

Successfully returned the sandbox metrics

Content-Type: `application/json`

Schema: `array<SandboxMetric>`

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

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
