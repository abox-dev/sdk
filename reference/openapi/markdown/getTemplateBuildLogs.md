# GET /templates/{templateID}/builds/{buildID}/logs

Template build logs

Get template build logs

## Parameters

- **`templateID`** · `string` · path · required

- **`buildID`** · `string` · path · required

- **`pageCursor`** · `string` · query · optional

  Opaque continuation cursor returned as nextCursor by the previous page

  Maximum length: `512`

- **`cursor`** · `integer` · query · optional

  Starting timestamp of the logs that should be returned in milliseconds

  Format: `int64`

  Minimum: `0`

- **`limit`** · `integer` · query · optional

  Maximum number of logs that should be returned

  Format: `int32`

  Default: `100`

  Minimum: `0`

  Maximum: `100`

- **`direction`** · `LogsDirection` · query · optional

  Allowed values for `LogsDirection`: `forward` | `backward`

- **`level`** · `LogLevel` · query · optional

  Allowed values for `LogLevel`: `debug` | `info` | `warn` | `error`

- **`source`** · `LogsSource` · query · optional

  Source of the logs that should be returned from

  Allowed values for `LogsSource`: `temporary` | `persistent`

## Responses

### 200

Successfully returned the template build logs

Content-Type: `application/json`

Schema: `TemplateBuildLogsResponse`

- **`logs`** · `array<BuildLogEntry>` · required

  Build logs structured

- **`logs.id`** · `string` · optional

  Stable identifier used to reconcile overlapping live log pages

- **`logs.timestamp`** · `string` · required

  Timestamp of the log entry

  Format: `date-time`

- **`logs.message`** · `string` · required

  Log message content

- **`logs.level`** · `LogLevel` · required

  State of the sandbox

  Allowed values for `LogLevel`: `debug` | `info` | `warn` | `error`

- **`logs.step`** · `string` · optional

  Step in the build process related to the log entry

- **`nextCursor`** · `string` · optional

  Opaque continuation cursor for the next page

- **`source`** · `LogsSource` · optional

  Source of the logs that should be returned

  Allowed values for `LogsSource`: `temporary` | `persistent`

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
