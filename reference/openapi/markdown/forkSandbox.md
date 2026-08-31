# POST /sandboxes/{sandboxID}/fork

Fork sandbox

Fork the sandbox: checkpoint the running sandbox in place (it is briefly paused, snapshotted with its full memory state, and resumed on its node, keeping its ID and expiration untouched) and create count new sandboxes from that snapshot. Returns one result per requested fork, each carrying either the created sandbox or the error that prevented it from starting. A non-201 status means the request failed before any fork was attempted.

## Parameters

- **`sandboxID`** · `string` · path · required

## Request body

Required: no

### application/json

Schema: `SandboxForkRequest`

- **`timeout`** · `integer` · optional

  Time to live for the new forked sandboxes in seconds.

- **`count`** · `integer` · optional

  Number of forked sandboxes to create. All forks boot from the same snapshot, so the snapshot is captured once regardless of count. Each fork succeeds or fails independently; the outcome of each is reported in its entry of the response list.

## Responses

### 201

The sandbox was snapshotted and the forks were attempted; each entry reports one fork's outcome

Content-Type: `application/json`

Schema: `array<SandboxForkResult>`

- **`sandbox`** · `Sandbox` · optional

- **`error`** · `Error` · optional

### 409

Conflict

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

### 401

Authentication error

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
