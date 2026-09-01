[agentbox-sdk-monorepo](../README.md) / SnapshotPaginator

# Interface: SnapshotPaginator

Paginator for listing snapshots.

## Example

```ts
const paginator = Sandbox.listSnapshots()
while (paginator.hasNext) {
  const snapshots = await paginator.nextItems()
  console.log(snapshots)
}
```

## Extends

- `Paginator`\<[`SnapshotInfo`](SnapshotInfo.md), [`SandboxApiOpts`](SandboxApiOpts.md)\>

## Accessors

### hasNext

#### Get Signature

> **get** **hasNext**(): `boolean`

Returns true if there are more items to fetch.

##### Returns

`boolean`

#### Inherited from

`Paginator.hasNext`

***

### nextToken

#### Get Signature

> **get** **nextToken**(): `string` \| `undefined`

Returns the next token to use for pagination.

##### Returns

`string` \| `undefined`

#### Inherited from

`Paginator.nextToken`

## Methods

### nextItems()

> **nextItems**(`opts?`): `Promise`\<[`SnapshotInfo`](SnapshotInfo.md)[]\>

Get the next page of items.

#### Parameters

##### opts?

[`SandboxApiOpts`](SandboxApiOpts.md)

per-call connection options. When provided, this call uses
these options (e.g. `apiKey`, `domain`, `headers`, `requestTimeoutMs`,
`signal`) instead of the ones the paginator was constructed with.
Aborting a page via `signal` does not affect subsequent Paginator.nextItems
calls — pass a fresh signal each call you want to be cancellable.

#### Returns

`Promise`\<[`SnapshotInfo`](SnapshotInfo.md)[]\>

List of items

#### Throws

Error if there are no more items to fetch. Call this method only if `hasNext` is `true`.

#### Overrides

`Paginator.nextItems`
