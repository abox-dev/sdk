[agentbox-sdk-monorepo](../README.md) / WriteInfo

# Interface: WriteInfo

Sandbox filesystem object information.

## Extended by

- [`EntryInfo`](EntryInfo.md)

## Properties

### metadata?

> `optional` **metadata?**: `Record`\<`string`, `string`\>

User-defined metadata stored on the file as `user.agentbox.*` extended
attributes. On writes this reflects the metadata supplied on upload; on
reads (`getInfo`, `list`, `rename`) it reflects any `user.agentbox.*` xattr on
the file, including ones set out-of-band. `undefined` when none is set.

***

### name

> **name**: `string`

Name of the filesystem object.

***

### path

> **path**: `string`

Path to the filesystem object.

***

### type?

> `optional` **type?**: [`FileType`](../enumerations/FileType.md)

Type of the filesystem object.
