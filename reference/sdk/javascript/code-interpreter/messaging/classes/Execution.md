[@abox-dev/code-interpreter](../../README.md) / [messaging](../README.md) / Execution

# Class: Execution

Represents the result of a cell execution.

## Constructors

### Constructor

> **new Execution**(`results?`, `logs?`, `error?`, `executionCount?`): `Execution`

#### Parameters

##### results?

[`Result`](Result.md)[] = `[]`

List of result of the cell (interactively interpreted last line), display calls (e.g. matplotlib plots).

##### logs?

[`Logs`](../type-aliases/Logs.md) = `...`

Logs printed to stdout and stderr during execution.

##### error?

[`ExecutionError`](ExecutionError.md)

An Error object if an error occurred, null otherwise.

##### executionCount?

`number`

Execution count of the cell.

#### Returns

`Execution`

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
