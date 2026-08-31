[@abox-dev/sdk](../README.md) / TemplateBuildStatusResponse

# Type Alias: TemplateBuildStatusResponse

> **TemplateBuildStatusResponse** = `object`

Response from getting build status.

## Properties

### buildID

> **buildID**: `string`

Build identifier.

***

### logEntries

> **logEntries**: `LogEntry`[]

Build log entries.

***

### reason?

> `optional` **reason?**: [`BuildStatusReason`](BuildStatusReason.md)

Reason for the current status (typically for errors).

***

### status

> **status**: [`TemplateBuildStatus`](TemplateBuildStatus.md)

Current status of the build.

***

### templateID

> **templateID**: `string`

Template identifier.
