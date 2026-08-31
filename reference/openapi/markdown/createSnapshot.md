# POST /sandboxes/{sandboxID}/snapshots

Create snapshot

Create a persistent snapshot from the sandbox's current state. Snapshots can be used to create new sandboxes and persist beyond the original sandbox's lifetime.

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | path | yes | `string` |  |

## Request body

Required: yes

### application/json

Schema: `SandboxSnapshotRequest`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | `string` | no | Optional name for the snapshot template. If a snapshot template with this name already exists, a new build will be assigned to the existing template instead of creating a new one. |

## Responses

### 201

Snapshot created successfully

Content-Type: `application/json`

Schema: `SnapshotInfo`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `snapshotID` | `string` | yes | Identifier of the snapshot template including the tag. Uses namespace/alias when a name was provided (e.g. team-slug/my-snapshot:default), otherwise falls back to the raw template ID (e.g. abc123:default). |
| `names` | `array<string>` | yes | Full names of the snapshot template including team namespace and tag (e.g. team-slug/my-snapshot:v2) |

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 404

Not found

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |
