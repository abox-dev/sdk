# POST /v2/templates/{templateID}/builds/{buildID}

Start template build (v2)

Start the build

## Parameters

- **`templateID`** · `string` · path · required

- **`buildID`** · `string` · path · required

## Request body

Required: yes

### application/json

Schema: `TemplateBuildStartV2`

- **`fromImage`** · `string` · optional

  Image to use as a base for the template build

- **`fromTemplate`** · `string` · optional

  Template to use as a base for the template build

- **`fromImageRegistry`** · `FromImageRegistry` · optional

- **`force`** · `boolean` · optional

  Whether the whole build should be forced to run regardless of the cache

- **`steps`** · `array<TemplateStep>` · optional

  List of steps to execute in the template build

- **`steps.type`** · `string` · required

  Type of the step

- **`steps.args`** · `array<string>` · optional

  Arguments for the step

- **`steps.filesHash`** · `string` · optional

  Hash of the files used in the step

- **`steps.force`** · `boolean` · optional

  Whether the step should be forced to run regardless of the cache

- **`startCmd`** · `string` · optional

  Start command to execute in the template after the build

- **`readyCmd`** · `string` · optional

  Ready check command to execute in the template after the build

## Responses

### 202

The build has started

### 401

Authentication error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

### 500

Server error

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error
