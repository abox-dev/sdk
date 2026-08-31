# PUT /sandboxes/{sandboxID}/network

Update sandbox network

Update the network configuration for a running sandbox. Replaces the current egress rules with the provided configuration. Omitting field clears it.

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |

## Request body

Required: yes

### application/json

Schema: `SandboxNetworkUpdateConfig`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `allowOut` | `array<string>` | no | List of allowed destinations for egress traffic. Each entry can be a CIDR block (e.g. "8.8.8.8/32"), a bare IP address (e.g. "8.8.8.8"), or a domain name (e.g. "example.com", "*.example.com"). Allowed entries always take precedence over denied entries. |
| `denyOut` | `array<string>` | no | List of denied CIDR blocks or IP addresses for egress traffic. Domain names are not supported for deny rules. |
| `rules` | `object` | no | Per-domain transform rules. Replaces all existing rules when provided. |
| `allow_internet_access` | `boolean` | no | Allow sandbox to access the internet. When set to false, it behaves the same as specifying denyOut to 0.0.0.0/0 in the network config. |

## Responses

### 204

Successfully updated the sandbox network configuration

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 404

Not found

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 409

Conflict

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
