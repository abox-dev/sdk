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
