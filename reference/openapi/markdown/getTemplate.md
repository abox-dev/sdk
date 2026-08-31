# GET /templates/{templateID}

List template builds

List all builds for a template

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `templateID` | path | yes | `string` |  |
| `nextToken` | query | no | `string` | Cursor to start the list from |
| `limit` | query | no | `integer` | Maximum number of items to return per page |

## Responses

### 200

Successfully returned the template with its builds

Content-Type: `application/json`

Schema: `TemplateWithBuilds`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `templateID` | `string` | yes | Identifier of the template |
| `public` | `boolean` | yes | Whether the template is public or only accessible by the team |
| `aliases` | `array<string>` | yes | Aliases of the template |
| `names` | `array<string>` | yes | Names of the template (namespace/alias format when namespaced) |
| `createdAt` | `string` | yes | Time when the template was created |
| `updatedAt` | `string` | yes | Time when the template was last updated |
| `lastSpawnedAt` | `string` | yes | Time when the template was last used |
| `spawnCount` | `integer` | yes | Number of times the template was used |
| `builds` | `array<TemplateBuild>` | yes | List of builds for the template |

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
