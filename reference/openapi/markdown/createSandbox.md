# POST /sandboxes

Create sandbox

Create a sandbox from the template

## Request body

Required: yes

### application/json

Schema: `NewSandbox`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `templateID` | `string` | yes | Identifier of the required template |
| `timeout` | `integer` | no | Time to live for the sandbox in seconds. |
| `autoPause` | `boolean` | no | Automatically pauses the sandbox after the timeout |
| `autoPauseMemory` | `boolean` | no | Controls the snapshot kind taken when the sandbox auto-pauses on timeout (only relevant when autoPause is true). When false, the auto-pause drops the in-memory state and persists only the filesystem (a filesystem-only snapshot); resuming it cold-boots (reboots) the sandbox from disk. Such a snapshot cannot be auto-resumed by traffic and must be resumed explicitly, so it cannot be combined with autoResume. Defaults to true (full memory snapshot). |
| `autoResume` | `SandboxAutoResumeConfig` | no | Auto-resume configuration for paused sandboxes. |
| `secure` | `boolean` | no | Secure all system communication with sandbox |
| `allow_internet_access` | `boolean` | no | Allow sandbox to access the internet. When set to false, it behaves the same as specifying denyOut to 0.0.0.0/0 in the network config. |
| `network` | `SandboxNetworkConfig` | no |  |
| `metadata` | `SandboxMetadata` | no |  |
| `envVars` | `EnvVars` | no |  |
| `mcp` | `Mcp` | no | MCP configuration for the sandbox |
| `iam` | `SandboxIam` | no | Sandbox workload identity configuration. A non-empty, valid tokens map enables workload identity for the sandbox. |
| `volumeMounts` | `array<SandboxVolumeMount>` | no |  |

## Responses

### 201

The sandbox was created successfully

Content-Type: `application/json`

Schema: `Sandbox`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `templateID` | `string` | yes | Identifier of the template from which is the sandbox created |
| `sandboxID` | `string` | yes | Identifier of the sandbox |
| `alias` | `string` | no | Alias of the template |
| `clientID` | `string` | yes | Identifier of the client |
| `envdVersion` | `EnvdVersion` | yes | Version of the envd running in the sandbox |
| `envdAccessToken` | `string` | no | Access token used for envd communication |
| `trafficAccessToken` | `string` | no | Token required for accessing sandbox via proxy. |
| `domain` | `string` | no | Base domain where the sandbox traffic is accessible |

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 400

Bad request

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
