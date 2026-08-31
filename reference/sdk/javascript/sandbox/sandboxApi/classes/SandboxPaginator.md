[@abox-dev/sdk](../../README.md) / [sandboxApi](../README.md) / SandboxPaginator

# Class: SandboxPaginator

Paginator for listing sandboxes.

## Example

```ts
const paginator = Sandbox.list()
while (paginator.hasNext) {
  const sandboxes = await paginator.nextItems()
  console.log(sandboxes)
}
```

## Extends

- `Paginator`\<[`SandboxInfo`](../interfaces/SandboxInfo.md), [`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)\>

## Constructors

### Constructor

> **new SandboxPaginator**(`opts?`): `SandboxPaginator`

#### Parameters

##### opts?

[`SandboxListOpts`](../interfaces/SandboxListOpts.md)

#### Returns

`SandboxPaginator`

#### Overrides

`Paginator<SandboxInfo, SandboxApiOpts>.constructor`

## Properties

### limit?

> `protected` `readonly` `optional` **limit?**: `number`

#### Inherited from

`Paginator.limit`

***

### opts?

> `protected` `readonly` `optional` **opts?**: [`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

#### Inherited from

`Paginator.opts`

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

> **nextItems**(`opts?`): `Promise`\<[`SandboxInfo`](../interfaces/SandboxInfo.md)[]\>

Get the next page of items.

#### Parameters

##### opts?

[`SandboxApiOpts`](../interfaces/SandboxApiOpts.md)

per-call connection options. When provided, this call uses
these options (e.g. `apiKey`, `domain`, `headers`, `requestTimeoutMs`,
`signal`) instead of the ones the paginator was constructed with.
Aborting a page via `signal` does not affect subsequent Paginator.nextItems
calls — pass a fresh signal each call you want to be cancellable.

#### Returns

`Promise`\<[`SandboxInfo`](../interfaces/SandboxInfo.md)[]\>

List of items

#### Throws

Error if there are no more items to fetch. Call this method only if `hasNext` is `true`.

#### Overrides

`Paginator.nextItems`

***

### updatePagination()

> `protected` **updatePagination**(`response`): `void`

Update the pagination state from a response, reading the `x-next-token`
header. Concrete paginators call this from Paginator.nextItems
after fetching a page.

#### Parameters

##### response

`Response`

#### Returns

`void`

#### Inherited from

`Paginator.updatePagination`
