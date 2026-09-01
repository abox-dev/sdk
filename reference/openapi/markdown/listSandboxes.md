# GET /v2/sandboxes

List sandboxes (v2)

List all sandboxes

## Parameters

- **`metadata`** · `string` · query · optional

  Metadata query used to filter the sandboxes (e.g. "user=abc&app=prod"). Each key and values must be URL encoded.

- **`state`** · `array<SandboxState>` · query · optional

  Filter sandboxes by one or more states

  Allowed values for `SandboxState`: `running` | `paused`

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

Successfully returned all running sandboxes

#### Response headers

- **`X-Next-Token`** · `string` · response header

  Cursor to fetch the next page of results, if more exist

- **`X-Total-Running`** · `integer` · response header

  Number of running sandboxes matching the filters, before pagination is applied. Only present when running sandboxes were requested.

  Format: `int32`

Content-Type: `application/json`

Schema: `array<ListedSandbox>`

- **`templateID`** · `string` · required

  Identifier of the template from which is the sandbox created

- **`alias`** · `string` · optional

  Alias of the template

- **`sandboxID`** · `string` · required

  Identifier of the sandbox

- **`startedAt`** · `string` · required

  Time when the sandbox was started

  Format: `date-time`

- **`endAt`** · `string` · required

  Time when the sandbox will expire

  Format: `date-time`

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

- **`metadata`** · `SandboxMetadata` · optional

- **`metadata.*`** · `string` · additional property

  Metadata of the sandbox

- **`state`** · `SandboxState` · required

  State of the sandbox

  Allowed values for `SandboxState`: `running` | `paused`

- **`envdVersion`** · `EnvdVersion` · required

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

### 400

Bad request

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
