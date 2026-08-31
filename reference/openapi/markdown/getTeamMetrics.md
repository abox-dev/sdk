# GET /teams/{teamID}/metrics

Team metrics

Get metrics for the team

## Parameters

- **`teamID`** · `string` · path · required

- **`start`** · `integer` · query · optional

  Unix timestamp for the start of the interval, in seconds, for which the metrics

- **`end`** · `integer` · query · optional

## Responses

### 200

Successfully returned the team metrics

Content-Type: `application/json`

Schema: `array<TeamMetric>`

- **`timestampUnix`** · `integer` · required

  Timestamp of the metric entry in Unix time (seconds since epoch)

- **`concurrentSandboxes`** · `integer` · required

  The number of concurrent sandboxes for the team

- **`sandboxStartRate`** · `number` · required

  Number of sandboxes started per second

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
