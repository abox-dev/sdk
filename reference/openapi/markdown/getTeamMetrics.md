# GET /teams/{teamID}/metrics

Team metrics

Get metrics for the team

## Parameters

- **`teamID`** · `string` · path · required

- **`start`** · `integer` · query · optional

  Unix timestamp for the start of the interval, in seconds, for which the metrics

  Format: `int64`

  Minimum: `0`

- **`end`** · `integer` · query · optional

  Unix timestamp for the end of the interval, in seconds, for which the metrics

  Format: `int64`

  Minimum: `0`

## Responses

### 200

Successfully returned the team metrics

Content-Type: `application/json`

Schema: `array<TeamMetric>`

- **`timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

  Format: `int64`

- **`concurrentSandboxes`** · `integer` · required

  The number of concurrent sandboxes for the team

  Format: `int32`

- **`sandboxStartRate`** · `number` · required

  Number of sandboxes started per second

  Format: `float`

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
