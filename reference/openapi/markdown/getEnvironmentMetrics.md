# GET /metrics

Service stats

## Responses

### 200

The resource usage metrics of the service

Content-Type: `application/json`

Schema: `Metrics`

- **`ts`** · `integer` · optional

  Unix timestamp in UTC for current sandbox time

- **`cpu_count`** · `integer` · optional

  Number of CPU cores

- **`cpu_used_pct`** · `number` · optional

  CPU usage percentage

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
