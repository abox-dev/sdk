# DELETE /templates/{templateID}

Delete template

Delete a template

## Parameters

- **`templateID`** · `string` · path · required

## Responses

### 204

The template was deleted successfully

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
