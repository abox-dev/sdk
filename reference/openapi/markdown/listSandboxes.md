# GET /v2/sandboxes

List sandboxes (v2)

List all sandboxes

## Parameters

- **`metadata`** · `string` · query · optional

  Metadata query used to filter the sandboxes (e.g. "user=abc&app=prod"). Each key and values must be URL encoded.

- **`state`** · `array<SandboxState>` · query · optional

  Filter sandboxes by one or more states

- **`nextToken`** · `string` · query · optional

  Cursor to start the list from

- **`limit`** · `integer` · query · optional

  Maximum number of items to return per page

## Responses

### 200

Successfully returned all running sandboxes

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

- **`endAt`** · `string` · required

  Time when the sandbox will expire

- **`cpuCount`** · `CPUCount` · required

  CPU cores for the sandbox

- **`memoryMB`** · `MemoryMB` · required

  Memory for the sandbox in MiB

- **`diskSizeMB`** · `DiskSizeMB` · required

  Disk size for the sandbox in MiB

- **`metadata`** · `SandboxMetadata` · optional

- **`metadata.*`** · `string` · additional property

  Metadata of the sandbox

- **`state`** · `SandboxState` · required

  State of the sandbox

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 400

Bad request

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
