# POST /v3/templates

Create template (v3)

Create a new template

## Request body

Required: yes

### application/json

Schema: `TemplateBuildRequestV3`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | `string` | no | Name of the template. Can include a tag with colon separator (e.g. "my-template" or "my-template:v1"). If tag is included, it will be treated as if the tag was provided in the tags array. |
| `tags` | `array<string>` | no | Tags to assign to the template build |
| `cpuCount` | `CPUCount` | no | CPU cores for the sandbox |
| `memoryMB` | `MemoryMB` | no | Memory for the sandbox in MiB |

## Responses

### 202

The build was requested successfully

Content-Type: `application/json`

Schema: `TemplateRequestResponseV3`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `templateID` | `string` | yes | Identifier of the template |
| `buildID` | `string` | yes | Identifier of the last successful build for given template |
| `public` | `boolean` | yes | Whether the template is public or only accessible by the team |
| `names` | `array<string>` | yes | Names of the template |
| `tags` | `array<string>` | yes | Tags assigned to the template build |

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
