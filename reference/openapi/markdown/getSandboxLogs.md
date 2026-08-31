# GET /v2/sandboxes/{sandboxID}/logs

Sandbox logs (v2)

Get sandbox logs

## Parameters

- **`sandboxID`** · `string` · path · required

- **`pageCursor`** · `string` · query · optional

  Opaque continuation cursor returned as nextCursor by the previous page

- **`cursor`** · `integer` · query · optional

  Starting timestamp of the logs that should be returned in milliseconds

- **`limit`** · `integer` · query · optional

  Maximum number of logs that should be returned

- **`direction`** · `LogsDirection` · query · optional

  Direction of the logs that should be returned

- **`level`** · `LogLevel` · query · optional

  Minimum log level to return. Logs below this level are excluded

- **`search`** · `string` · query · optional

  Case-sensitive substring match on log message content

## Responses

### 200

Successfully returned the sandbox logs

Content-Type: `application/json`

Schema: `SandboxLogsV2Response`

- **`logs`** · `array<SandboxLogEntry>` · required

  Sandbox logs structured

- **`logs.id`** · `string` · optional

  Stable identifier used to reconcile overlapping live log pages

- **`logs.timestamp`** · `string` · required

  Timestamp of the log entry

- **`logs.message`** · `string` · required

  Log message content

- **`logs.level`** · `LogLevel` · required

  State of the sandbox

- **`logs.fields`** · `object` · required

- **`logs.fields.*`** · `string` · additional property

- **`nextCursor`** · `string` · optional

  Opaque continuation cursor for the next page

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
