# POST /files/compose

Compose multiple files into a single file using zero-copy concatenation. Source files are deleted after successful composition.

## Parameters

- **`Agentbox-Sandbox-Id`** · `string` · header · required

  Identifier of the sandbox that receives the request.

- **`Agentbox-Sandbox-Port`** · `integer` · header · required

  Internal envd HTTP port exposed through the sandbox proxy.

  Default: `49983`

## Request body

Required: yes

### application/json

Schema: `ComposeRequest`

- **`source_paths`** · `array<string>` · required

  Ordered list of source file paths to concatenate

- **`destination`** · `string` · required

  Destination file path for the composed file

- **`username`** · `string` · optional

  User for setting ownership and resolving relative paths

## Responses

### 200

Files composed successfully

Content-Type: `application/json`

Schema: `EntryInfo`

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

### 404

File not found

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
