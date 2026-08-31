# GET /templates/{templateID}/builds/{buildID}/logs

Template build logs

Get template build logs

## Parameters

- **`templateID`** · `string` · path · required

- **`buildID`** · `string` · path · required

- **`pageCursor`** · `string` · query · optional

  Opaque continuation cursor returned as nextCursor by the previous page

- **`cursor`** · `integer` · query · optional

  Starting timestamp of the logs that should be returned in milliseconds

- **`limit`** · `integer` · query · optional

  Maximum number of logs that should be returned

- **`direction`** · `LogsDirection` · query · optional

- **`level`** · `LogLevel` · query · optional

- **`source`** · `LogsSource` · query · optional

  Source of the logs that should be returned from

## Responses

### 200

Successfully returned the template build logs

Content-Type: `application/json`

Schema: `TemplateBuildLogsResponse`

- **`logs`** · `array<BuildLogEntry>` · required

  Build logs structured

- **`nextCursor`** · `string` · optional

  Opaque continuation cursor for the next page

- **`source`** · `LogsSource` · optional

  Source of the logs that should be returned

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 404

Not found

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error
