# GET /templates/aliases/{alias}

Check template alias

Check if template with given alias exists

## Parameters

- **`alias`** · `string` · path · required

## Responses

### 200

Successfully queried template by alias

Content-Type: `application/json`

Schema: `TemplateAliasResponse`

- **`templateID`** · `string` · required

  Identifier of the template

- **`public`** · `boolean` · required

  Whether the template is public or only accessible by the team

### 400

Bad request

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

### 404

Not found

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
