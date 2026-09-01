# GET /metrics

Service stats

## Parameters

- **`Agentbox-Sandbox-Id`** · `string` · header · required

  Identifier of the sandbox that receives the request.

- **`Agentbox-Sandbox-Port`** · `integer` · header · required

  Internal envd HTTP port exposed through the sandbox proxy.

  Default: `49983`

## Responses

### 200

The resource usage metrics of the service

Content-Type: `application/json`

Schema: `Metrics`

Resource usage metrics

- **`ts`** · `integer` · optional

  Unix timestamp in UTC for current sandbox time

  Format: `int64`

- **`cpu_count`** · `integer` · optional

  Number of CPU cores

- **`cpu_used_pct`** · `number` · optional

  CPU usage percentage

  Format: `float`

- **`mem_total`** · `integer` · optional

  Total virtual memory in bytes

- **`mem_used`** · `integer` · optional

  Used virtual memory in bytes

- **`mem_cache`** · `integer` · optional

  Cached memory (page cache) in bytes

- **`mem_total_mib`** · `integer` · optional

  Total virtual memory in MiB

- **`mem_used_mib`** · `integer` · optional

  Used virtual memory in MiB

- **`disk_used`** · `integer` · optional

  Used disk space in bytes

- **`disk_total`** · `integer` · optional

  Total disk space in bytes
