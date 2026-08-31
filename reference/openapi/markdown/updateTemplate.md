# PATCH /v2/templates/{templateID}

Update template (v2)

Update template

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `templateID` | path | yes | `string` |  |

## Request body

Required: yes

### application/json

Schema: `TemplateUpdateRequest`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `public` | `boolean` | no | Whether the template is public or only accessible by the team |

## Responses

### 200

The template was updated successfully

Content-Type: `application/json`

Schema: `TemplateUpdateResponse`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `names` | `array<string>` | yes | Names of the template (namespace/alias format when namespaced) |

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

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |
