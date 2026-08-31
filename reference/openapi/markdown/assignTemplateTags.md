# POST /templates/tags

Assign template tags

Assign tag(s) to a template build

## Request body

Required: yes

### application/json

Schema: `AssignTemplateTagsRequest`

- **`target`** · `string` · required

  Target template in "name:tag" format

- **`tags`** · `array<string>` · required

  Tags to assign to the template

## Responses

### 201

Tag assigned successfully

Content-Type: `application/json`

Schema: `AssignedTemplateTags`

- **`tags`** · `array<string>` · required

  Assigned tags of the template

- **`buildID`** · `string` · required

  Identifier of the build associated with these tags

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

### 404

Not found

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
