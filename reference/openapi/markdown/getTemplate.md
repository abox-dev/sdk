# GET /templates/{templateID}

List template builds

List all builds for a template

## Parameters

- **`templateID`** · `string` · path · required

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

Successfully returned the template with its builds

#### Response headers

- **`X-Next-Token`** · `string` · response header

  Cursor to fetch the next page of results, if more exist

Content-Type: `application/json`

Schema: `TemplateWithBuilds`

- **`templateID`** · `string` · required

  Identifier of the template

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

- **`lastSpawnedAt`** · `string | null` · required

  Time when the template was last used

  Format: `date-time`

- **`spawnCount`** · `integer` · required

  Number of times the template was used

  Format: `int64`

- **`builds`** · `array<TemplateBuild>` · required

  List of builds for the template

- **`builds.buildID`** · `string` · required

  Identifier of the build

  Format: `uuid`

- **`builds.status`** · `TemplateBuildStatus` · required

  Status of the template build

  Allowed values for `TemplateBuildStatus`: `building` | `waiting` | `ready` | `error`

- **`builds.createdAt`** · `string` · required

  Time when the build was created

  Format: `date-time`

- **`builds.updatedAt`** · `string` · required

  Time when the build was last updated

  Format: `date-time`

- **`builds.finishedAt`** · `string` · optional

  Time when the build was finished

  Format: `date-time`

- **`builds.cpuCount`** · `CPUCount` · required

  CPU cores for the sandbox

  Format: `int32`

  Minimum: `1`

- **`builds.memoryMB`** · `MemoryMB` · required

  Memory for the sandbox in MiB

  Format: `int32`

  Minimum: `128`

- **`builds.diskSizeMB`** · `DiskSizeMB` · optional

  Disk size for the sandbox in MiB

  Format: `int32`

  Minimum: `0`

- **`builds.envdVersion`** · `EnvdVersion` · optional

  Version of the envd running in the sandbox

### 401

Authentication error

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
