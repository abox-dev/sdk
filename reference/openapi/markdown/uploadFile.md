# POST /files

Upload a file and ensure the parent directories exist. If the file exists, it will be overwritten.

Any request header of the form `X-Metadata-<key>: <value>` is persisted
as a user-defined extended attribute on the uploaded file. The
`X-Metadata-` prefix is stripped and the remaining header name is
lowercased to form the metadata key; the resulting map is returned on
`EntryInfo` lookups (e.g. `Stat`, `ListDir`).

Each upload replaces the file's metadata with the keys provided in
that request: keys previously stored but absent from the new request
are removed, and an upload that sends no `X-Metadata-*` header clears
all existing metadata.

Both keys and values must be printable US-ASCII (bytes `0x20`-`0x7E`)
and are rejected with HTTP 400 otherwise. Each key is capped at 246
bytes (the Linux VFS xattr-name limit minus the namespace prefix), and
the combined size of all metadata on a file (keys plus values, with the
namespace prefix counted per key) is capped at 4096 bytes to stay within
the filesystem's per-inode xattr budget. Multiple files in a single
multipart upload receive the same metadata. If the same
`X-Metadata-<key>` header is sent more than once, only the first
value is used.

## Parameters

- **`Agentbox-Sandbox-Id`** · `string` · header · required

  Identifier of the sandbox that receives the request.

- **`Agentbox-Sandbox-Port`** · `integer` · header · required

  Internal envd HTTP port exposed through the sandbox proxy.

  Default: `49983`

- **`path`** · `string` · query · optional

  Path to the file, URL encoded. Can be relative to the user's home directory (e.g. "file.txt" resolves to ~/file.txt).

- **`username`** · `string` · query · optional

  User for setting file ownership and resolving relative paths. Defaults to the sandbox's default user.

- **`signature`** · `string` · query · optional

  Signature used for file access permission verification.

- **`signature_expiration`** · `integer` · query · optional

  Unix timestamp (seconds) after which the signature expires. Only used with the signature parameter.

## Request body

Required: yes

### multipart/form-data

Schema: `object`

- **`file`** · `string` · optional

  Format: `binary`

### application/octet-stream

Schema: `string`

Raw file content. The 'path' query parameter is required when using this content type.

Format: `binary`

## Responses

### 200

The file was uploaded successfully.

Content-Type: `application/json`

Schema: `array<EntryInfo>`

- **`path`** · `string` · required

  Path to the file

- **`name`** · `string` · required

  Name of the file

- **`type`** · `file` · required

  Type of the file

  Allowed values: `file`

- **`metadata`** · `object` · optional

  User-defined metadata stored as extended attributes on the file.

- **`metadata.*`** · `string` · additional property

### 400

Invalid path

Content-Type: `application/json`

Schema: `Error`

- **`message`** · `string` · required

  Error message

- **`code`** · `integer` · required

  Error code

### 401

Invalid user

Content-Type: `application/json`

Schema: `Error`

- **`message`** · `string` · required

  Error message

- **`code`** · `integer` · required

  Error code

### 500

Internal server error

Content-Type: `application/json`

Schema: `Error`

- **`message`** · `string` · required

  Error message

- **`code`** · `integer` · required

  Error code

### 507

Not enough disk space

Content-Type: `application/json`

Schema: `Error`

- **`message`** · `string` · required

  Error message

- **`code`** · `integer` · required

  Error code
