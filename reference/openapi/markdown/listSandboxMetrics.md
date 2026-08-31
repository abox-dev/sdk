# GET /sandboxes/metrics

List sandbox metrics

List metrics for given sandboxes

## Parameters

- **`sandbox_ids`** · `array<string>` · query · required

  Comma-separated list of sandbox IDs to get metrics for

## Responses

### 200

Successfully returned all running sandboxes with metrics

Content-Type: `application/json`

Schema: `SandboxesWithMetrics`

- **`sandboxes`** · `object` · required

- **`sandboxes.*`** · `SandboxMetric` · additional property

  Metric entry with timestamp and line

- **`sandboxes.*.timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

- **`sandboxes.*.cpuCount`** · `integer` · required

  Number of CPU cores

- **`sandboxes.*.cpuUsedPct`** · `number` · required

  CPU usage percentage

- **`sandboxes.*.memUsed`** · `integer` · required

  Memory used in bytes

- **`sandboxes.*.memTotal`** · `integer` · required

  Total memory in bytes

- **`sandboxes.*.memCache`** · `integer` · required

  Cached memory (page cache) in bytes

- **`sandboxes.*.diskUsed`** · `integer` · required

  Disk used in bytes

- **`sandboxes.*.diskTotal`** · `integer` · required

  Total disk space in bytes

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 400

Bad request

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
