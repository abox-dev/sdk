# GET /files

Download a file

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

## Responses

### 200

Entire file downloaded successfully.

Content-Type: `application/octet-stream`

Schema: `string`

The raw file content

Format: `binary`

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

### 404

File not found

Content-Type: `application/json`

Schema: `Error`

- **`message`** · `string` · required

  Error message

- **`code`** · `integer` · required

  Error code

### 406

Requested encoding is not supported

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
