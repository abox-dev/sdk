[@abox-dev/sdk](../README.md) / EntryInfo

# Interface: EntryInfo

Sandbox filesystem object information.

## Extends

- [`WriteInfo`](WriteInfo.md)

## Properties

### group

> **group**: `string`

Group owner of the filesystem object.

***

### metadata?

> `optional` **metadata?**: `Record`\<`string`, `string`\>

User-defined metadata stored on the file as `user.agentbox.*` extended
attributes. On writes this reflects the metadata supplied on upload; on
reads (`getInfo`, `list`, `rename`) it reflects any `user.agentbox.*` xattr on
the file, including ones set out-of-band. `undefined` when none is set.

#### Inherited from

[`WriteInfo`](WriteInfo.md).[`metadata`](WriteInfo.md#metadata)

***

### mode

> **mode**: `number`

File mode and permission bits.

***

### modifiedTime?

> `optional` **modifiedTime?**: `Date`

Last modification time of the filesystem object.

***

### name

> **name**: `string`

Name of the filesystem object.

#### Inherited from

[`WriteInfo`](WriteInfo.md).[`name`](WriteInfo.md#name)

***

### owner

> **owner**: `string`

Owner of the filesystem object.

***

### path

> **path**: `string`

Path to the filesystem object.

#### Inherited from

[`WriteInfo`](WriteInfo.md).[`path`](WriteInfo.md#path)

***

### permissions

> **permissions**: `string`

String representation of file permissions (e.g. 'rwxr-xr-x').

***

### size

> **size**: `number`

Size of the filesystem object in bytes.

***

### symlinkTarget?

> `optional` **symlinkTarget?**: `string`

If the filesystem object is a symlink, this is the target of the symlink.

***

### type?

> `optional` **type?**: [`FileType`](../enumerations/FileType.md)

Type of the filesystem object.

#### Inherited from

[`WriteInfo`](WriteInfo.md).[`type`](WriteInfo.md#type)
