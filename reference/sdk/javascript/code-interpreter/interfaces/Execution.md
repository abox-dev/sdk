[agentbox-sdk-monorepo](../README.md) / Execution

# Interface: Execution

Represents the result of a cell execution.

## Properties

### error?

> `optional` **error?**: [`ExecutionError`](ExecutionError.md)

An Error object if an error occurred, null otherwise.

***

### executionCount?

> `optional` **executionCount?**: `number`

Execution count of the cell.

***

### logs

> **logs**: [`Logs`](../type-aliases/Logs.md)

Logs printed to stdout and stderr during execution.

***

### results

> **results**: [`Result`](Result.md)[] = `[]`

List of result of the cell (interactively interpreted last line), display calls (e.g. matplotlib plots).

## Accessors

### text

#### Get Signature

> **get** **text**(): `string` \| `undefined`

Returns the text representation of the main result of the cell.

##### Returns

`string` \| `undefined`

## Methods

### toJSON()

> **toJSON**(): `object`

Returns the serializable representation of the execution result.

#### Returns

`object`

##### error

> **error**: [`ExecutionError`](ExecutionError.md) \| `undefined`

##### logs

> **logs**: [`Logs`](../type-aliases/Logs.md)

##### results

> **results**: [`Result`](Result.md)[]
