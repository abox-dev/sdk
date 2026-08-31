# GET /templates/{templateID}

List template builds

List all builds for a template

## Parameters

- **`templateID`** · `string` · path · required

- **`nextToken`** · `string` · query · optional

  Cursor to start the list from

- **`limit`** · `integer` · query · optional

  Maximum number of items to return per page

## Responses

### 200

Successfully returned the template with its builds

Content-Type: `application/json`

Schema: `TemplateWithBuilds`

- **`templateID`** · `string` · required

  Identifier of the template

- **`public`** · `boolean` · required

  Whether the template is public or only accessible by the team

- **`names`** · `array<string>` · required

  Names of the template (namespace/alias format when namespaced)

- **`createdAt`** · `string` · required

  Time when the template was created

- **`updatedAt`** · `string` · required

  Time when the template was last updated

- **`lastSpawnedAt`** · `string` · required

  Time when the template was last used

- **`spawnCount`** · `integer` · required

  Number of times the template was used

- **`builds`** · `array<TemplateBuild>` · required

  List of builds for the template

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
