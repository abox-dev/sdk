[agentbox-sdk-monorepo](../README.md) / SnapshotInfo

# Interface: SnapshotInfo

Information about a snapshot.

## Properties

### names

> **names**: `string`[]

Full names of the snapshot template including project slug and tag (e.g. project-slug/my-snapshot:v2).

***

### snapshotId

> **snapshotId**: `string`

Snapshot identifier — template ID with tag, or namespaced name with tag (e.g. my-snapshot:latest).
Can be used with Sandbox.create() to create a new sandbox from this snapshot.
