# GET /v2/sandboxes

List sandboxes (v2)

List all sandboxes

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `metadata` | query | no | `string` | Metadata query used to filter the sandboxes (e.g. "user=abc&app=prod"). Each key and values must be URL encoded. |
| `state` | query | no | `array<SandboxState>` | Filter sandboxes by one or more states |
| `nextToken` | query | no | `string` | Cursor to start the list from |
| `limit` | query | no | `integer` | Maximum number of items to return per page |

## Responses

### 200

Successfully returned all running sandboxes

Content-Type: `application/json`

Schema: `array<ListedSandbox>`

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
