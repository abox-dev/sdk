[@abox-dev/sdk](../README.md) / waitForProcess

# Function: waitForProcess()

> **waitForProcess**(`processName`): [`ReadyCmd`](../classes/ReadyCmd.md)

Wait for a process with a specific name to be running.
Uses `pgrep` to check if a process exists.

## Parameters

### processName

`string`

Name of the process to wait for

## Returns

[`ReadyCmd`](../classes/ReadyCmd.md)

ReadyCmd that checks for the process

## Example

```ts
import { Template, waitForProcess } from '@abox-dev/sdk'

const template = Template()
  .fromBaseImage()
  .setStartCmd('./my-daemon', waitForProcess('my-daemon'))
```
