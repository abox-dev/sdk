# GET /templates/{templateID}/files/{hash}

Template build file upload URL

Get an upload link for a tar file containing build layer files

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `templateID` | path | yes | `string` |  |
| `hash` | path | yes | `string` |  |

## Responses

### 201

The upload link where to upload the tar file

Content-Type: `application/json`

Schema: `TemplateBuildFileUpload`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `present` | `boolean` | yes | Whether the file is already present in the cache |
| `url` | `string` | no | Url where the file should be uploaded to |

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

### 404

Not found

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
