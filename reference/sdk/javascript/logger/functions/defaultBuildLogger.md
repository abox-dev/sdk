[agentbox-sdk-monorepo](../README.md) / defaultBuildLogger

# Function: defaultBuildLogger()

> **defaultBuildLogger**(`options?`): (`logEntry`) => `void`

Create a default build logger with animated timer display.

## Parameters

### options?

Logger configuration options

#### minLevel?

[`LogEntryLevel`](../type-aliases/LogEntryLevel.md)

Minimum log level to display (default: 'info')

## Returns

Logger function that accepts LogEntry instances

(`logEntry`) => `void`

## Example

```ts
import { Template, defaultBuildLogger } from '@abox-dev/sdk'

const template = Template().fromPythonImage()

await Template.build(template, {
  alias: 'my-template',
  onBuildLogs: defaultBuildLogger({ minLevel: 'debug' })
})
```
