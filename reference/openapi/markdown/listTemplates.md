# GET /v2/templates

List templates (v2)

List all templates

## Parameters

- **`teamID`** · `string` · query · optional

- **`nextToken`** · `string` · query · optional

  Cursor to start the list from

- **`limit`** · `integer` · query · optional

  Maximum number of items to return per page

  Format: `int32`

  Default: `100`

  Minimum: `1`

  Maximum: `100`

## Responses

### 200

Successfully returned all templates

#### Response headers

- **`X-Next-Token`** · `string` · response header

  Cursor to fetch the next page of results, if more exist

Content-Type: `application/json`

Schema: `array<Template>`

- **`templateID`** · `string` · required

  Identifier of the template

- **`buildID`** · `string` · required

  Identifier of the last successful build for given template

- **`cpuCount`** · `CPUCount` · required

  CPU cores for the sandbox

  Format: `int32`

  Minimum: `1`

- **`memoryMB`** · `MemoryMB` · required

  Memory for the sandbox in MiB

  Format: `int32`

  Minimum: `128`

- **`diskSizeMB`** · `DiskSizeMB` · required

  Disk size for the sandbox in MiB

  Format: `int32`

  Minimum: `0`

- **`public`** · `boolean` · required

  Whether the template is public or only accessible by the team

- **`names`** · `array<string>` · required

  Names of the template (namespace/alias format when namespaced)

- **`createdAt`** · `string` · required

  Time when the template was created

  Format: `date-time`

- **`updatedAt`** · `string` · required

  Time when the template was last updated

  Format: `date-time`

- **`createdBy`** · `object` · required

- **`lastSpawnedAt`** · `string` · required

  Time when the template was last used

  Format: `date-time`

- **`spawnCount`** · `integer` · required

  Number of times the template was used

  Format: `int64`

- **`buildCount`** · `integer` · required

  Number of times the template was built

  Format: `int32`

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`buildStatus`** · `TemplateBuildStatus` · required

  Status of the template build

  Allowed values for `TemplateBuildStatus`: `building` | `waiting` | `ready` | `error`

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 403

Forbidden

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error
