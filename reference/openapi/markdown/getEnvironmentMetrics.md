# GET /metrics

Service stats

## Responses

### 200

The resource usage metrics of the service

Content-Type: `application/json`

Schema: `Metrics`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `ts` | `integer` | no | Unix timestamp in UTC for current sandbox time |
| `cpu_count` | `integer` | no | Number of CPU cores |
| `cpu_used_pct` | `number` | no | CPU usage percentage |
| `mem_total` | `integer` | no | Total virtual memory in bytes |
| `mem_used` | `integer` | no | Used virtual memory in bytes |
| `mem_cache` | `integer` | no | Cached memory (page cache) in bytes |
| `mem_total_mib` | `integer` | no | Total virtual memory in MiB |
| `mem_used_mib` | `integer` | no | Used virtual memory in MiB |
| `disk_used` | `integer` | no | Used disk space in bytes |
| `disk_total` | `integer` | no | Total disk space in bytes |
