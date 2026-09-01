# POST /sandboxes/{sandboxID}/snapshots

Create snapshot

Create a persistent snapshot from the sandbox's current state. Snapshots can be used to create new sandboxes and persist beyond the original sandbox's lifetime.

## Parameters

- **`sandboxID`** · `string` · path · required

## Request body

Required: yes

### application/json

Schema: `SandboxSnapshotRequest`

- **`name`** · `string` · optional

  Optional name for the snapshot template. If a snapshot template with this name already exists, a new build will be assigned to the existing template instead of creating a new one.

## Responses

### 201

Snapshot created successfully

Content-Type: `application/json`

Schema: `SnapshotInfo`

- **`snapshotID`** · `string` · required

  Identifier of the snapshot template including the tag. Uses namespace/alias when a name was provided (e.g. team-slug/my-snapshot:default), otherwise falls back to the raw template ID (e.g. abc123:default).

- **`names`** · `array<string>` · required

  Full names of the snapshot template including team namespace and tag (e.g. team-slug/my-snapshot:v2)

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
