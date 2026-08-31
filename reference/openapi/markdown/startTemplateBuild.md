# POST /v2/templates/{templateID}/builds/{buildID}

Start template build (v2)

Start the build

## Parameters

| Name | In | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `templateID` | path | yes | `string` |  |
| `buildID` | path | yes | `string` |  |

## Request body

Required: yes

### application/json

Schema: `TemplateBuildStartV2`

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `fromImage` | `string` | no | Image to use as a base for the template build |
| `fromTemplate` | `string` | no | Template to use as a base for the template build |
| `fromImageRegistry` | `FromImageRegistry` | no |  |
| `force` | `boolean` | no | Whether the whole build should be forced to run regardless of the cache |
| `steps` | `array<TemplateStep>` | no | List of steps to execute in the template build |
| `startCmd` | `string` | no | Start command to execute in the template after the build |
| `readyCmd` | `string` | no | Ready check command to execute in the template after the build |

## Responses

### 202

The build has started

### 401

Authentication error

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
