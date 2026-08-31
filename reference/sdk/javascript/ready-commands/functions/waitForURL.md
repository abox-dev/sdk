[@abox-dev/sdk](../README.md) / waitForURL

# Function: waitForURL()

> **waitForURL**(`url`, `statusCode?`): [`ReadyCmd`](../classes/ReadyCmd.md)

Wait for a URL to return a specific HTTP status code.
Uses `curl` to make HTTP requests and check the response status.

## Parameters

### url

`string`

URL to check (e.g., 'http://localhost:3000/health')

### statusCode?

`number` = `200`

Expected HTTP status code (default: 200)

## Returns

[`ReadyCmd`](../classes/ReadyCmd.md)

ReadyCmd that checks the URL

## Example

```ts
import { Template, waitForURL } from '@abox-dev/sdk'

const template = Template()
  .fromNodeImage()
  .setStartCmd('npm start', waitForURL('http://localhost:3000/health'))
```
