# POST /sandboxes/{sandboxID}/connect

Connect sandbox

Returns sandbox details. If the sandbox is paused, it will be resumed. TTL is only extended.

## Parameters

- **`sandboxID`** · `string` · path · required

## Request body

Required: yes

### application/json

Schema: `ConnectSandbox`

- **`timeout`** · `integer` · required

  Timeout in seconds from the current time after which the sandbox should expire

  Format: `int32`

  Minimum: `0`

## Responses

### 200

The sandbox was already running

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

- **`domain`** · `string | null` · optional

  Base domain where the sandbox traffic is accessible

### 201

The sandbox was resumed successfully

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

- **`domain`** · `string | null` · optional

  Base domain where the sandbox traffic is accessible

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 404

Not found

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
