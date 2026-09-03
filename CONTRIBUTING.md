# Contributing

Open an issue or pull request at [abox-dev/sdk](https://github.com/abox-dev/sdk). Include tests for behavior changes and keep JavaScript, Python sync/async, and Go APIs aligned where applicable.

For Go changes, run `make go-check`. It includes the exported GoDoc gate. Generated clients under `packages/go-sdk/internal/gen` and Go reference pages under `reference/sdk/go` must be regenerated with `make generate` and must not be edited manually. `packages/go-sdk/go.mod` defines the minimum supported Go version; CI also tests every newer supported minor listed in `RELEASING.md`.

Use the development and generation commands documented in the root README. By contributing, you agree that your contribution is licensed under the license applicable to the package you modify.
