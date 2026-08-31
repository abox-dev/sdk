[agentbox-sdk-monorepo](../README.md) / waitForTimeout

# Function: waitForTimeout()

> **waitForTimeout**(`timeout`): [`ReadyCmd`](../classes/ReadyCmd.md)

Wait for a specified timeout before considering the sandbox ready.
Uses `sleep` command to wait for a fixed duration.

## Parameters

### timeout

`number`

Time to wait in milliseconds (minimum: 1000ms / 1 second)

## Returns

[`ReadyCmd`](../classes/ReadyCmd.md)

ReadyCmd that waits for the specified duration

## Example

```ts
import { Template, waitForTimeout } from '@abox-dev/sdk'

const template = Template()
  .fromNodeImage()
  .setStartCmd('npm start', waitForTimeout(5000)) // Wait 5 seconds
```
