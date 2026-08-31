[@abox-dev/code-interpreter](../../README.md) / [messaging](../README.md) / Result

# Class: Result

Represents the data to be displayed as a result of executing a cell in a Jupyter notebook.
The result is similar to the structure returned by ipython kernel: https://ipython.readthedocs.io/en/stable/development/execution.html#execution-semantics

The result can contain multiple types of data, such as text, images, plots, etc. Each type of data is represented
as a string, and the result can contain multiple types of data. The display calls don't have to have text representation,
for the actual result the representation is always present for the result, the other representations are always optional.

## Constructors

### Constructor

> **new Result**(`rawData`, `isMainResult`): `Result`

#### Parameters

##### rawData

[`RawData`](../type-aliases/RawData.md)

##### isMainResult

`boolean`

#### Returns

`Result`

## Properties

### chart?

> `readonly` `optional` **chart?**: [`ChartTypes`](../../charts/type-aliases/ChartTypes.md)

Contains the chart data.

***

### data?

> `readonly` `optional` **data?**: `Record`\<`string`, `unknown`\>

Contains the data from DataFrame.

***

### extra?

> `readonly` `optional` **extra?**: `any`

Extra data that can be included. Not part of the standard types.

***

### html?

> `readonly` `optional` **html?**: `string`

HTML representation of the data.

***

### isMainResult

> `readonly` **isMainResult**: `boolean`

***

### javascript?

> `readonly` `optional` **javascript?**: `string`

JavaScript representation of the data.

***

### jpeg?

> `readonly` `optional` **jpeg?**: `string`

JPEG representation of the data.

***

### json?

> `readonly` `optional` **json?**: `string`

JSON representation of the data.

***

### latex?

> `readonly` `optional` **latex?**: `string`

LaTeX representation of the data.

***

### markdown?

> `readonly` `optional` **markdown?**: `string`

Markdown representation of the data.

***

### pdf?

> `readonly` `optional` **pdf?**: `string`

PDF representation of the data.

***

### png?

> `readonly` `optional` **png?**: `string`

PNG representation of the data.

***

### raw

> `readonly` **raw**: [`RawData`](../type-aliases/RawData.md)

***

### svg?

> `readonly` `optional` **svg?**: `string`

SVG representation of the data.

***

### text?

> `readonly` `optional` **text?**: `string`

Text representation of the result.

## Methods

### formats()

> **formats**(): `string`[]

Returns all the formats available for the result.

#### Returns

`string`[]

Array of strings representing the formats available for the result.

***

### toJSON()

> **toJSON**(): `object`

Returns the serializable representation of the result.

#### Returns

`object`

##### extra?

> `optional` **extra?**: `any`

##### html

> **html**: `string` \| `undefined`

##### javascript

> **javascript**: `string` \| `undefined`

##### jpeg

> **jpeg**: `string` \| `undefined`

##### json

> **json**: `string` \| `undefined`

##### latex

> **latex**: `string` \| `undefined`

##### markdown

> **markdown**: `string` \| `undefined`

##### pdf

> **pdf**: `string` \| `undefined`

##### png

> **png**: `string` \| `undefined`

##### svg

> **svg**: `string` \| `undefined`

##### text

> **text**: `string` \| `undefined`
