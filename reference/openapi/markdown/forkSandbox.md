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

  Format: `int32`

  Default: `15`

  Minimum: `0`

- **`count`** · `integer` · optional

  Number of forked sandboxes to create. All forks boot from the same snapshot, so the snapshot is captured once regardless of count. Each fork succeeds or fails independently; the outcome of each is reported in its entry of the response list.

  Format: `int32`

  Default: `1`

  Minimum: `1`

  Maximum: `100`

## Responses

### 201

The sandbox was snapshotted and the forks were attempted; each entry reports one fork's outcome

Content-Type: `application/json`

Schema: `array<SandboxForkResult>`

- **`sandbox`** · `Sandbox` · optional

- **`sandbox.templateID`** · `string` · required

  Identifier of the template from which is the sandbox created

- **`sandbox.sandboxID`** · `string` · required

  Identifier of the sandbox

- **`sandbox.alias`** · `string` · optional

  Alias of the template

- **`sandbox.envdVersion`** · `EnvdVersion` · required

  Version of the envd running in the sandbox

- **`sandbox.domain`** · `string` · optional

  Base domain where the sandbox traffic is accessible

- **`error`** · `Error` · optional

- **`error.code`** · `integer` · required

  Error code

  Format: `int32`

- **`error.message`** · `string` · required

  Error

### 409

Conflict

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

### 401

Authentication error

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
