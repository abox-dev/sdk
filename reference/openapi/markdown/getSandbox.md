# GET /sandboxes/{sandboxID}

Sandbox

Get a sandbox by id

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |

## Responses

### 200

Successfully returned the sandbox

Content-Type: `application/json`

Schema: `SandboxDetail`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `templateID` | `string` | yes | Identifier of the template from which is the sandbox created |
| `alias` | `string` | no | Alias of the template |
| `sandboxID` | `string` | yes | Identifier of the sandbox |
| `startedAt` | `string` | yes | Time when the sandbox was started |
| `endAt` | `string` | yes | Time when the sandbox will expire |
| `envdVersion` | `EnvdVersion` | yes | Version of the envd running in the sandbox |
| `allowInternetAccess` | `boolean` | no | Whether internet access was explicitly enabled or disabled for the sandbox. Null means it was not explicitly set. |
| `domain` | `string` | no | Base domain where the sandbox traffic is accessible |
| `cpuCount` | `CPUCount` | yes | CPU cores for the sandbox |
| `memoryMB` | `MemoryMB` | yes | Memory for the sandbox in MiB |
| `diskSizeMB` | `DiskSizeMB` | yes | Disk size for the sandbox in MiB |
| `metadata` | `SandboxMetadata` | no |  |
| `state` | `SandboxState` | yes | State of the sandbox |
| `network` | `SandboxNetworkConfig` | no |  |
| `lifecycle` | `SandboxLifecycle` | no | Sandbox lifecycle policy returned by sandbox info. |

### 404

Not found

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

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |
