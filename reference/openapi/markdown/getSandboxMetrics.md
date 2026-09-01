# GET /sandboxes/{sandboxID}/metrics

Sandbox metrics

Get sandbox metrics

## Parameters

- **`sandboxID`** · `string` · path · required

- **`start`** · `integer` · query · optional

  Unix timestamp for the start of the interval, in seconds, for which the metrics

  Format: `int64`

  Minimum: `0`

- **`end`** · `integer` · query · optional

  Format: `int64`

  Minimum: `0`

## Responses

### 200

Successfully returned the sandbox metrics

Content-Type: `application/json`

Schema: `array<SandboxMetric>`

- **`timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

  Format: `int64`

- **`cpuCount`** · `integer` · required

  Number of CPU cores

  Format: `int32`

- **`cpuUsedPct`** · `number` · required

  CPU usage percentage

  Format: `float`

- **`memUsed`** · `integer` · required

  Memory used in bytes

  Format: `int64`

- **`memTotal`** · `integer` · required

  Total memory in bytes

  Format: `int64`

- **`memCache`** · `integer` · required

  Cached memory (page cache) in bytes

  Format: `int64`

- **`diskUsed`** · `integer` · required

  Disk used in bytes

  Format: `int64`

- **`diskTotal`** · `integer` · required

  Total disk space in bytes

  Format: `int64`

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
