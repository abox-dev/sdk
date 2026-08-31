[@abox-dev/code-interpreter](../../README.md) / [messaging](../README.md) / OutputMessage

# Class: OutputMessage

Represents an output message from the sandbox code execution.

## Constructors

### Constructor

> **new OutputMessage**(`line`, `timestamp`, `error`): `OutputMessage`

#### Parameters

##### line

`string`

The output line.

##### timestamp

`number`

Unix epoch in nanoseconds.

##### error

`boolean`

Whether the output is an error.

#### Returns

`OutputMessage`

## Properties

### error

> `readonly` **error**: `boolean`

Whether the output is an error.

***

### line

> `readonly` **line**: `string`

The output line.

***

### timestamp

> `readonly` **timestamp**: `number`

Unix epoch in nanoseconds.

## Methods

### toString()

> **toString**(): `string`

#### Returns

`string`
