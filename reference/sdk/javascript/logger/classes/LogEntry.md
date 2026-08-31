[agentbox-sdk-monorepo](../README.md) / LogEntry

# Class: LogEntry

Represents a single log entry from the template build process.

## Extended by

- [`LogEntryEnd`](LogEntryEnd.md)
- [`LogEntryStart`](LogEntryStart.md)

## Constructors

### Constructor

> **new LogEntry**(`timestamp`, `level`, `message`): `LogEntry`

#### Parameters

##### timestamp

`Date`

##### level

[`LogEntryLevel`](../type-aliases/LogEntryLevel.md)

##### message

`string`

#### Returns

`LogEntry`

## Properties

### level

> `readonly` **level**: [`LogEntryLevel`](../type-aliases/LogEntryLevel.md)

***

### message

> `readonly` **message**: `string`

***

### timestamp

> `readonly` **timestamp**: `Date`

## Methods

### toString()

> **toString**(): `string`

#### Returns

`string`
