[@abox-dev/sdk](../README.md) / waitForFile

# Function: waitForFile()

> **waitForFile**(`filename`): [`ReadyCmd`](../classes/ReadyCmd.md)

Wait for a file to exist.
Uses shell test command to check file existence.

## Parameters

### filename

`string`

Path to the file to wait for

## Returns

[`ReadyCmd`](../classes/ReadyCmd.md)

ReadyCmd that checks for the file

## Example

```ts
import { Template, waitForFile } from '@abox-dev/sdk'

const template = Template()
  .fromBaseImage()
  .setStartCmd('./init.sh', waitForFile('/tmp/ready'))
```
