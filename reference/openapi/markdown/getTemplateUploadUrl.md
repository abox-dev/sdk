# GET /templates/{templateID}/files/{hash}

Template build file upload URL

Get an upload link for a tar file containing build layer files

## Parameters

- **`templateID`** · `string` · path · required

- **`hash`** · `string` · path · required

## Responses

### 201

The upload link where to upload the tar file

Content-Type: `application/json`

Schema: `TemplateBuildFileUpload`

- **`present`** · `boolean` · required

  Whether the file is already present in the cache

- **`url`** · `string` · optional

  Url where the file should be uploaded to

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
