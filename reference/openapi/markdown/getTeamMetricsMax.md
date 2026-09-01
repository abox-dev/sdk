# GET /teams/{teamID}/metrics/max

Maximum team metrics

Get the maximum metrics for the team in the given interval

## Parameters

- **`teamID`** · `string` · path · required

- **`start`** · `integer` · query · optional

  Unix timestamp for the start of the interval, in seconds, for which the metrics

  Format: `int64`

  Minimum: `0`

- **`end`** · `integer` · query · optional

  Format: `int64`

  Minimum: `0`

- **`metric`** · `concurrent_sandboxes | sandbox_start_rate` · query · required

  Metric to retrieve the maximum value for

  Allowed values: `concurrent_sandboxes` | `sandbox_start_rate`

## Responses

### 200

Successfully returned the team metrics

Content-Type: `application/json`

Schema: `MaxTeamMetric`

Team metric with timestamp

- **`timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

  Format: `int64`

- **`value`** · `number` · required

  The maximum value of the requested metric in the given interval

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

### 403

Forbidden

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
