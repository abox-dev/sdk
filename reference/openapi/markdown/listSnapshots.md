# GET /snapshots

List snapshots

List all snapshots for the team

## Parameters

- **`sandboxID`** · `string` · query · optional

- **`name`** · `string` · query · optional

  Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-team/my-snapshot" or "my-snapshot:v1").

- **`limit`** · `integer` · query · optional

  Maximum number of items to return per page

  Format: `int32`

  Default: `100`

  Minimum: `1`

  Maximum: `100`

- **`nextToken`** · `string` · query · optional

  Cursor to start the list from

## Responses

### 200

Successfully returned snapshots

#### Response headers

- **`X-Next-Token`** · `string` · response header

  Cursor to fetch the next page of results, if more exist

Content-Type: `application/json`

Schema: `array<SnapshotInfo>`

- **`snapshotID`** · `string` · required

  Identifier of the snapshot template including the tag. Uses namespace/alias when a name was provided (e.g. team-slug/my-snapshot:default), otherwise falls back to the raw template ID (e.g. abc123:default).

- **`names`** · `array<string>` · required

  Full names of the snapshot template including team namespace and tag (e.g. team-slug/my-snapshot:v2)

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
