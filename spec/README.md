# API snapshots

The OpenAPI and envd files in this directory are checked-in snapshots of `mono/infra`. Their source revision and SHA-256 digests are recorded in `source.json`.

Do not edit synchronized files manually. Maintainers update them from a local mono checkout:

```bash
make sync-specs MONO_DIR=/path/to/mono
make generate
```

`make generate` is deterministic and uses only these checked-in files. It performs no network fetches. `mcp-server.json` and the `buf-*.gen.yaml` files are maintained in this repository.
