# GET /sandboxes/{sandboxID}

Sandbox

Get a sandbox by id

## Parameters

- **`sandboxID`** · `string` · path · required

## Responses

### 200

Successfully returned the sandbox

Content-Type: `application/json`

Schema: `SandboxDetail`

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

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`allowInternetAccess`** · `boolean` · optional

  Whether internet access was explicitly enabled or disabled for the sandbox. Null means it was not explicitly set.

- **`domain`** · `string` · optional

  Base domain where the sandbox traffic is accessible

- **`cpuCount`** · `CPUCount` · required

  CPU cores for the sandbox

- **`memoryMB`** · `MemoryMB` · required

  Memory for the sandbox in MiB

- **`diskSizeMB`** · `DiskSizeMB` · required

  Disk size for the sandbox in MiB

- **`metadata`** · `SandboxMetadata` · optional

- **`state`** · `SandboxState` · required

  State of the sandbox

- **`network`** · `SandboxNetworkConfig` · optional

- **`lifecycle`** · `SandboxLifecycle` · optional

  Sandbox lifecycle policy returned by sandbox info.

### 404

Not found

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

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error
