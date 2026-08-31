# GET /files

Download a file

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `path` | query | no | `string` | Path to the file, URL encoded. Can be relative to the user's home directory (e.g. "file.txt" resolves to ~/file.txt). |
| `username` | query | no | `string` | User for setting file ownership and resolving relative paths. Defaults to the sandbox's default user. |
| `signature` | query | no | `string` | Signature used for file access permission verification. |
| `signature_expiration` | query | no | `integer` | Unix timestamp (seconds) after which the signature expires. Only used with the signature parameter. |

## Responses

### 200

Entire file downloaded successfully.

Content-Type: `application/octet-stream`

Schema: `string`

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

### 406

Requested encoding is not supported

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
