[@abox-dev/code-interpreter](../../README.md) / [messaging](../README.md) / ExecutionError

# Class: ExecutionError

Represents an error that occurred during the execution of a cell.
The error contains the name of the error, the value of the error, and the traceback.

## Constructors

### Constructor

> **new ExecutionError**(`name`, `value`, `traceback`): `ExecutionError`

#### Parameters

##### name

`string`

Name of the error.

##### value

`string`

Value of the error.

##### traceback

`string`

The raw traceback of the error.

#### Returns

`ExecutionError`

## Properties

### name

> **name**: `string`

Name of the error.

***

### traceback

> **traceback**: `string`

The raw traceback of the error.

***

### value

> **value**: `string`

Value of the error.
