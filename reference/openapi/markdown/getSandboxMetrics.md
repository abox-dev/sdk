# GET /sandboxes/{sandboxID}/metrics

Sandbox metrics

Get sandbox metrics

## Parameters

- **`sandboxID`** · `string` · path · required

- **`start`** · `integer` · query · optional

  Unix timestamp for the start of the interval, in seconds, for which the metrics

- **`end`** · `integer` · query · optional

## Responses

### 200

Successfully returned the sandbox metrics

Content-Type: `application/json`

Schema: `array<SandboxMetric>`

- **`timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

- **`cpuCount`** · `integer` · required

  Number of CPU cores

- **`cpuUsedPct`** · `number` · required

  CPU usage percentage

- **`memUsed`** · `integer` · required

  Memory used in bytes

- **`memTotal`** · `integer` · required

  Total memory in bytes

- **`memCache`** · `integer` · required

  Cached memory (page cache) in bytes

- **`diskUsed`** · `integer` · required

  Disk used in bytes

- **`diskTotal`** · `integer` · required

  Total disk space in bytes

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

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
