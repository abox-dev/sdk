# POST /sandboxes/{sandboxID}/refreshes

Refresh sandbox

Refresh the sandbox extending its time to live

## Parameters

- **`sandboxID`** · `string` · path · required

## Request body

Required: no

### application/json

Schema: `SandboxRefreshRequest`

- **`duration`** · `integer` · optional

  Duration for which the sandbox should be kept alive in seconds

  Minimum: `0`

  Maximum: `3600`

## Responses

### 204

Successfully refreshed the sandbox

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
