[@abox-dev/sdk](../README.md) / LogEntryEnd

# Class: LogEntryEnd

Special log entry indicating the end of a build process.

## Extends

- [`LogEntry`](LogEntry.md)

## Constructors

### Constructor

> **new LogEntryEnd**(`timestamp`, `message`): `LogEntryEnd`

#### Parameters

##### timestamp

`Date`

##### message

`string`

#### Returns

`LogEntryEnd`

#### Overrides

[`LogEntry`](LogEntry.md).[`constructor`](LogEntry.md#constructor)

## Properties

### level

> `readonly` **level**: [`LogEntryLevel`](../type-aliases/LogEntryLevel.md)

#### Inherited from

[`LogEntry`](LogEntry.md).[`level`](LogEntry.md#level)

***

### message

> `readonly` **message**: `string`

#### Inherited from

[`LogEntry`](LogEntry.md).[`message`](LogEntry.md#message)

***

### timestamp

> `readonly` **timestamp**: `Date`

#### Inherited from

[`LogEntry`](LogEntry.md).[`timestamp`](LogEntry.md#timestamp)

## Methods

### toString()

> **toString**(): `string`

#### Returns

`string`

#### Inherited from

[`LogEntry`](LogEntry.md).[`toString`](LogEntry.md#tostring)
