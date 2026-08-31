# PATCH /v2/templates/{templateID}

Update template (v2)

Update template

## Parameters

- **`templateID`** · `string` · path · required

## Request body

Required: yes

### application/json

Schema: `TemplateUpdateRequest`

- **`public`** · `boolean` · optional

  Whether the template is public or only accessible by the team

## Responses

### 200

The template was updated successfully

Content-Type: `application/json`

Schema: `TemplateUpdateResponse`

- **`names`** · `array<string>` · required

  Names of the template (namespace/alias format when namespaced)

### 400

Bad request

Content-Type: `application/json`

Schema: `Error`

- **`code`** · `integer` · required

  Error code

- **`message`** · `string` · required

  Error

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
