# POST /v3/templates

Create template (v3)

Create a new template

## Request body

Required: yes

### application/json

Schema: `TemplateBuildRequestV3`

- **`name`** · `string` · optional

  Name of the template. Can include a tag with colon separator (e.g. "my-template" or "my-template:v1"). If tag is included, it will be treated as if the tag was provided in the tags array.

  Maximum length: `128`

- **`tags`** · `array<string>` · optional

  Tags to assign to the template build

- **`cpuCount`** · `CPUCount` · optional

  CPU cores for the sandbox

  Format: `int32`

  Minimum: `1`

- **`memoryMB`** · `MemoryMB` · optional

  Memory for the sandbox in MiB

  Format: `int32`

  Minimum: `128`

## Responses

### 202

The build was requested successfully

Content-Type: `application/json`

Schema: `TemplateRequestResponseV3`

- **`templateID`** · `string` · required

  Identifier of the template

- **`buildID`** · `string` · required

  Identifier of the last successful build for given template

- **`public`** · `boolean` · required

  Whether the template is public or only accessible by the team

- **`names`** · `array<string>` · required

  Names of the template

- **`tags`** · `array<string>` · required

  Tags assigned to the template build

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 403

Forbidden

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

  Format: `int32`

- **`message`** · `string` · required

  Error
