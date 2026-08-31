# POST /sandboxes

Create sandbox

Create a sandbox from the template

## Request body

Required: yes

### application/json

Schema: `NewSandbox`

- **`templateID`** · `string` · required

  Identifier of the required template

- **`timeout`** · `integer` · optional

  Time to live for the sandbox in seconds.

- **`autoPause`** · `boolean` · optional

  Automatically pauses the sandbox after the timeout

- **`autoPauseMemory`** · `boolean` · optional

  Controls the snapshot kind taken when the sandbox auto-pauses on timeout (only relevant when autoPause is true). When false, the auto-pause drops the in-memory state and persists only the filesystem (a filesystem-only snapshot); resuming it cold-boots (reboots) the sandbox from disk. Such a snapshot cannot be auto-resumed by traffic and must be resumed explicitly, so it cannot be combined with autoResume. Defaults to true (full memory snapshot).

- **`autoResume`** · `SandboxAutoResumeConfig` · optional

  Auto-resume configuration for paused sandboxes.

- **`secure`** · `boolean` · optional

  Secure all system communication with sandbox

- **`allow_internet_access`** · `boolean` · optional

  Allow sandbox to access the internet. When set to false, it behaves the same as specifying denyOut to 0.0.0.0/0 in the network config.

- **`network`** · `SandboxNetworkConfig` · optional

- **`metadata`** · `SandboxMetadata` · optional

- **`envVars`** · `EnvVars` · optional

- **`mcp`** · `Mcp` · optional

  MCP configuration for the sandbox

- **`iam`** · `SandboxIam` · optional

  Sandbox workload identity configuration. A non-empty, valid tokens map enables workload identity for the sandbox.

## Responses

### 201

The sandbox was created successfully

Content-Type: `application/json`

Schema: `Sandbox`

- **`templateID`** · `string` · required

  Identifier of the template from which is the sandbox created

- **`sandboxID`** · `string` · required

  Identifier of the sandbox

- **`alias`** · `string` · optional

  Alias of the template

- **`envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`domain`** · `string` · optional

  Base domain where the sandbox traffic is accessible

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
