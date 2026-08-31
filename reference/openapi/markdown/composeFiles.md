# POST /files/compose

Compose multiple files into a single file using zero-copy concatenation. Source files are deleted after successful composition.

## Request body

Required: yes

### application/json

Schema: `ComposeRequest`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `source_paths` | `array<string>` | yes | Ordered list of source file paths to concatenate |
| `destination` | `string` | yes | Destination file path for the composed file |
| `username` | `string` | no | User for setting ownership and resolving relative paths |

## Responses

### 200

Files composed successfully

Content-Type: `application/json`

Schema: `EntryInfo`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `path` | `string` | yes | Path to the file |
| `name` | `string` | yes | Name of the file |
| `type` | `file` | yes | Type of the file |
| `metadata` | `object` | no | User-defined metadata stored as extended attributes on the file. |

### 400

Invalid path

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `message` | `string` | yes | Error message |
| `code` | `integer` | yes | Error code |

### 401

Invalid user

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `message` | `string` | yes | Error message |
| `code` | `integer` | yes | Error code |

### 404

File not found

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `message` | `string` | yes | Error message |
| `code` | `integer` | yes | Error code |

### 500

Internal server error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `message` | `string` | yes | Error message |
| `code` | `integer` | yes | Error code |

### 507

Not enough disk space

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `message` | `string` | yes | Error message |
| `code` | `integer` | yes | Error code |
