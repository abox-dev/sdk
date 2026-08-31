# POST /sandboxes/{sandboxID}/connect

Connect sandbox

Returns sandbox details. If the sandbox is paused, it will be resumed. TTL is only extended.

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |

## Request body

Required: yes

### application/json

Schema: `ConnectSandbox`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `timeout` | `integer` | yes | Timeout in seconds from the current time after which the sandbox should expire |

## Responses

### 200

The sandbox was already running

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

### 201

The sandbox was resumed successfully

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

### 400

Bad request

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

### 404

Not found

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
