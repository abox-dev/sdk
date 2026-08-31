[agentbox-sdk-monorepo](../README.md) / FilesystemEvent

# Interface: FilesystemEvent

Information about a filesystem event.

## Properties

### entry?

> `optional` **entry?**: [`EntryInfo`](EntryInfo.md)

Information about the entry that triggered the event.

Only populated when the watch was started with `includeEntry: true` and the
sandbox's envd version supports it. It may be `undefined` for events where the
entry no longer exists at the path (e.g. remove or rename-away events).

***

### name

> **name**: `string`

Relative path to the filesystem object.

***

### type

> **type**: [`FilesystemEventType`](../enumerations/FilesystemEventType.md)

Filesystem operation event type.
