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

- **Variant `AWSRegistry`** · discriminator `type` = `aws`

- **`fromImageRegistry<AWSRegistry>.type`** · `aws` · required

  Type of registry authentication

  Allowed values: `aws`

- **`fromImageRegistry<AWSRegistry>.awsAccessKeyId`** · `string` · required

  AWS Access Key ID for ECR authentication

- **`fromImageRegistry<AWSRegistry>.awsSecretAccessKey`** · `string` · required

  AWS Secret Access Key for ECR authentication

- **`fromImageRegistry<AWSRegistry>.awsRegion`** · `string` · required

  AWS Region where the ECR registry is located

- **Variant `GCPRegistry`** · discriminator `type` = `gcp`

- **`fromImageRegistry<GCPRegistry>.type`** · `gcp` · required

  Type of registry authentication

  Allowed values: `gcp`

- **`fromImageRegistry<GCPRegistry>.serviceAccountJson`** · `string` · required

  Service Account JSON for GCP authentication

- **Variant `GeneralRegistry`** · discriminator `type` = `registry`

- **`fromImageRegistry<GeneralRegistry>.type`** · `registry` · required

  Type of registry authentication

  Allowed values: `registry`

- **`fromImageRegistry<GeneralRegistry>.username`** · `string` · required

  Username to use for the registry

- **`fromImageRegistry<GeneralRegistry>.password`** · `string` · required

  Password to use for the registry

- **`force`** · `boolean` · optional

  Whether the whole build should be forced to run regardless of the cache

  Default: `false`

- **`steps`** · `array<TemplateStep>` · optional

  List of steps to execute in the template build

  Default: `[]`

- **`steps.type`** · `string` · required

  Type of the step

- **`steps.args`** · `array<string>` · optional

  Arguments for the step

  Default: `[]`

- **`steps.filesHash`** · `string` · optional

  Hash of the files used in the step

- **`steps.force`** · `boolean` · optional

  Whether the step should be forced to run regardless of the cache

  Default: `false`

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
