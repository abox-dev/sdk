# GET /templates/{templateID}/builds/{buildID}/status

Template build status

Get template build info

## Parameters

- **`templateID`** · `string` · path · required

- **`buildID`** · `string` · path · required

- **`logsOffset`** · `integer` · query · optional

  Index of the starting build log that should be returned with the template

  Format: `int32`

  Default: `0`

  Minimum: `0`

- **`limit`** · `integer` · query · optional

  Maximum number of logs that should be returned

  Format: `int32`

  Default: `100`

  Minimum: `0`

  Maximum: `100`

- **`level`** · `LogLevel` · query · optional

  Allowed values for `LogLevel`: `debug` | `info` | `warn` | `error`

## Responses

### 200

Successfully returned the template

Content-Type: `application/json`

Schema: `TemplateBuildInfo`

- **`logs`** · `array<string>` · required

  Build logs

- **`logEntries`** · `array<BuildLogEntry>` · required

  Build logs structured

- **`logEntries.id`** · `string` · optional

  Stable identifier used to reconcile overlapping live log pages

- **`logEntries.timestamp`** · `string` · required

  Timestamp of the log entry

  Format: `date-time`

- **`logEntries.message`** · `string` · required

  Log message content

- **`logEntries.level`** · `LogLevel` · required

  State of the sandbox

  Allowed values for `LogLevel`: `debug` | `info` | `warn` | `error`

- **`logEntries.step`** · `string` · optional

  Step in the build process related to the log entry

- **`templateID`** · `string` · required

  Identifier of the template

- **`buildID`** · `string` · required

  Identifier of the build

- **`status`** · `TemplateBuildStatus` · required

  Status of the template build

  Allowed values for `TemplateBuildStatus`: `building` | `waiting` | `ready` | `error`

- **`reason`** · `BuildStatusReason` · optional

- **`reason.message`** · `string` · required

  Message with the status reason, currently reporting only for error status

- **`reason.step`** · `string` · optional

  Step that failed

- **`reason.logEntries`** · `array<BuildLogEntry>` · optional

  Log entries related to the status reason

- **`reason.logEntries.id`** · `string` · optional

  Stable identifier used to reconcile overlapping live log pages

- **`reason.logEntries.timestamp`** · `string` · required

  Timestamp of the log entry

  Format: `date-time`

- **`reason.logEntries.message`** · `string` · required

  Log message content

- **`reason.logEntries.level`** · `LogLevel` · required

  State of the sandbox

  Allowed values for `LogLevel`: `debug` | `info` | `warn` | `error`

- **`reason.logEntries.step`** · `string` · optional

  Step in the build process related to the log entry

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
