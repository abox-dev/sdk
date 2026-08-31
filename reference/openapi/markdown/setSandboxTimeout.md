# POST /sandboxes/{sandboxID}/timeout

Set sandbox timeout

Set the timeout for the sandbox. The sandbox will expire x seconds from the time of the request. Calling this method multiple times overwrites the TTL, each time using the current timestamp as the starting point to measure the timeout duration.

## Parameters

- **`sandboxID`** · `string` · path · required

## Request body

Required: no

### application/json

Schema: `SandboxTimeoutRequest`

- **`timeout`** · `integer` · required

  Timeout in seconds from the current time after which the sandbox should expire

## Responses

### 204

Successfully set the sandbox timeout

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 404

Not found

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
