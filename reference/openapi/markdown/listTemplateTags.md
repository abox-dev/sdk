# GET /templates/{templateID}/tags

List template tags

List all tags for a template

## Parameters

- **`templateID`** · `string` · path · required

## Responses

### 200

Successfully returned the template tags

Content-Type: `application/json`

Schema: `array<TemplateTag>`

- **`tag`** · `string` · required

  The tag name

- **`buildID`** · `string` · required

  Identifier of the build associated with this tag

  Format: `uuid`

- **`createdAt`** · `string` · required

  Time when the tag was assigned

  Format: `date-time`

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
