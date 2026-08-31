# GET /teams/{teamID}/metrics/max

Maximum team metrics

Get the maximum metrics for the team in the given interval

## Parameters

- **`teamID`** · `string` · path · required

- **`start`** · `integer` · query · optional

  Unix timestamp for the start of the interval, in seconds, for which the metrics

- **`end`** · `integer` · query · optional

- **`metric`** · `concurrent_sandboxes | sandbox_start_rate` · query · required

  Metric to retrieve the maximum value for

## Responses

### 200

Successfully returned the team metrics

Content-Type: `application/json`

Schema: `MaxTeamMetric`

- **`timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

- **`value`** · `number` · required

  The maximum value of the requested metric in the given interval

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

### 403

Forbidden

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
