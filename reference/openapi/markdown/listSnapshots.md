# GET /snapshots

List snapshots

List all snapshots for the team

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `sandboxID` | query | no | `string` |  |
| `name` | query | no | `string` | Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-team/my-snapshot" or "my-snapshot:v1"). |
| `limit` | query | no | `integer` | Maximum number of items to return per page |
| `nextToken` | query | no | `string` | Cursor to start the list from |

## Responses

### 200

Successfully returned snapshots

Content-Type: `application/json`

Schema: `array<SnapshotInfo>`

### 401

Authentication error

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
