# POST /templates/tags

Assign template tags

Assign tag(s) to a template build

## Request body

Required: yes

### application/json

Schema: `AssignTemplateTagsRequest`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `target` | `string` | yes | Target template in "name:tag" format |
| `tags` | `array<string>` | yes | Tags to assign to the template |

## Responses

### 201

Tag assigned successfully

Content-Type: `application/json`

Schema: `AssignedTemplateTags`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `tags` | `array<string>` | yes | Assigned tags of the template |
| `buildID` | `string` | yes | Identifier of the build associated with these tags |

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
