# GET /templates/aliases/{alias}

Check template alias

Check if template with given alias exists

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `alias` | path | yes | `string` |  |

## Responses

### 200

Successfully queried template by alias

Content-Type: `application/json`

Schema: `TemplateAliasResponse`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `templateID` | `string` | yes | Identifier of the template |
| `public` | `boolean` | yes | Whether the template is public or only accessible by the team |

### 400

Bad request

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
