[agentbox-sdk-monorepo](../README.md) / waitForPort

# Function: waitForPort()

> **waitForPort**(`port`): [`ReadyCmd`](../classes/ReadyCmd.md)

Wait for a port to be listening.
Uses `ss` command to check if a port is open and listening.

## Parameters

### port

`number`

Port number to wait for

## Returns

[`ReadyCmd`](../classes/ReadyCmd.md)

ReadyCmd that checks for the port

## Example

```ts
import { Template, waitForPort } from '@abox-dev/sdk'

const template = Template()
  .fromPythonImage()
  .setStartCmd('python -m http.server 8000', waitForPort(8000))
```
