[agentbox-sdk-monorepo](../README.md) / CreateCodeContextOpts

# Interface: CreateCodeContextOpts

Options for creating a code context.

## Properties

### cwd?

> `optional` **cwd?**: `string`

Working directory for the context.

#### Default

```ts
/home/user
```

***

### language?

> `optional` **language?**: [`RunCodeLanguage`](../type-aliases/RunCodeLanguage.md)

Language for the context.

#### Default

```ts
python
```

***

### requestTimeoutMs?

> `optional` **requestTimeoutMs?**: `number`

Timeout for the request in **milliseconds**.

#### Default

```ts
30_000 // 30 seconds
```
