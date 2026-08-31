# GET /v2/templates

List templates (v2)

List all templates

## Parameters

- **`teamID`** · `string` · query · optional

- **`nextToken`** · `string` · query · optional

  Cursor to start the list from

- **`limit`** · `integer` · query · optional

  Maximum number of items to return per page

## Responses

### 200

Successfully returned all templates

Content-Type: `application/json`

Schema: `array<Template>`

- **`templateID`** · `string` · required

  Identifier of the template

- **`buildID`** · `string` · required

  Identifier of the last successful build for given template

- **`cpuCount`** · `CPUCount` · required

  CPU cores for the sandbox

- **`memoryMB`** · `MemoryMB` · required

  Memory for the sandbox in MiB

- **`diskSizeMB`** · `DiskSizeMB` · required

  Disk size for the sandbox in MiB

- **`public`** · `boolean` · required

  Whether the template is public or only accessible by the team

- **`names`** · `array<string>` · required

  Names of the template (namespace/alias format when namespaced)

- **`createdAt`** · `string` · required

  Time when the template was created

- **`updatedAt`** · `string` · required

  Time when the template was last updated

- **`createdBy`** · `object` · required

- **`lastSpawnedAt`** · `string` · required

  Time when the template was last used

- **`spawnCount`** · `integer` · required

  Number of times the template was used

- **`buildCount`** · `integer` · required

  Number of times the template was built

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`buildStatus`** · `TemplateBuildStatus` · required

  Status of the template build

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 403

Forbidden

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error
