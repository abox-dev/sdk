[agentbox-sdk-monorepo](../README.md) / CommandResult

# Interface: CommandResult

Command execution result.

## Properties

### error?

> `optional` **error?**: `string`

Error message from command execution if it failed.

***

### exitCode

> **exitCode**: `number`

Command execution exit code.
`0` if the command finished successfully.

***

### stderr

> **stderr**: `string`

Command stderr output.

***

### stdout

> **stdout**: `string`

Command stdout output.
