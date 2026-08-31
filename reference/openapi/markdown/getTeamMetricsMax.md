# GET /teams/{teamID}/metrics/max

Maximum team metrics

Get the maximum metrics for the team in the given interval

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `teamID` | path | yes | `string` |  |
| `start` | query | no | `integer` | Unix timestamp for the start of the interval, in seconds, for which the metrics |
| `end` | query | no | `integer` |  |
| `metric` | query | yes | `concurrent_sandboxes \| sandbox_start_rate` | Metric to retrieve the maximum value for |

## Responses

### 200

Successfully returned the team metrics

Content-Type: `application/json`

Schema: `MaxTeamMetric`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `timestampUnix` | `integer` | yes | Timestamp of the metric entry in Unix time (seconds since epoch) |
| `value` | `number` | yes | The maximum value of the requested metric in the given interval |

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 403

Forbidden

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `code` | `integer` | yes | Error code |
| `message` | `string` | yes | Error |
