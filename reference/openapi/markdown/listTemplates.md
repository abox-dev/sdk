# GET /v2/templates

List templates (v2)

List all templates

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `teamID` | query | no | `string` |  |
| `nextToken` | query | no | `string` | Cursor to start the list from |
| `limit` | query | no | `integer` | Maximum number of items to return per page |

## Responses

### 200

Successfully returned all templates

Content-Type: `application/json`

Schema: `array<Template>`

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

### 403

Forbidden

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
