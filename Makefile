CODEGEN_IMAGE ?= agentbox-sdk-codegen

.PHONY: generate generate-in-container generate-go codegen-image sync-specs \
	go-format-check go-vet go-build go-test go-race go-coverage \
	go-doc-check go-integration go-consumer-check go-check check-agent-instructions

# Generate exclusively from the checked-in snapshots under spec/.
generate: codegen-image
	docker run --rm -v "$$PWD:/workspace" $(CODEGEN_IMAGE) make generate-in-container

codegen-image:
	docker build -q -t $(CODEGEN_IMAGE) -f codegen.Dockerfile .

generate-in-container:
	cd packages/js-sdk && pnpm generate
	cd packages/python-sdk && make generate
	$(MAKE) generate-go
	python scripts/generate-reference.py
	python scripts/test-reference-contract.py

generate-go:
	redocly bundle go-sdk --config redocly.yaml -o spec/openapi_generated.go-sdk.yml
	python scripts/filter-public-openapi.py spec/openapi_generated.go-sdk.yml
	oapi-codegen --config packages/go-sdk/internal/gen/api/oapi-codegen.yaml spec/openapi_generated.go-sdk.yml
	redocly bundle envd --config redocly.yaml -o spec/openapi_generated.go-envd.yml
	python packages/go-sdk/scripts/filter-go-envd.py spec/openapi_generated.go-envd.yml
	oapi-codegen --config packages/go-sdk/internal/gen/envdapi/oapi-codegen.yaml spec/openapi_generated.go-envd.yml
	cd spec/envd && buf generate --template buf-go.gen.yaml
	gofmt -w packages/go-sdk/internal/gen

check-agent-instructions:
	cmp -s AGENTS.md CLAUDE.md

go-format-check:
	cd packages/go-sdk && ./scripts/check-go-format.sh
	test -z "$$(gofmt -l scripts/check-go-docs.go)"

go-doc-check:
	go run ./scripts/check-go-docs.go packages/go-sdk packages/go-sdk/codeinterpreter

go-vet:
	cd packages/go-sdk && go vet ./...

go-build:
	cd packages/go-sdk && go build ./...

go-test:
	cd packages/go-sdk && go test ./...

go-race:
	cd packages/go-sdk && go test -race ./...

go-coverage:
	cd packages/go-sdk && ./scripts/check-go-coverage.sh

go-integration:
	cd packages/go-sdk && go test -tags=integration ./...

go-consumer-check:
	cd packages/go-sdk && ./scripts/test-go-consumer.sh

go-check: check-agent-instructions go-format-check go-doc-check go-vet go-build go-test go-race go-coverage go-consumer-check

# Maintainer-only update from a local mono checkout.
sync-specs:
	@test -n "$(MONO_DIR)" || (echo "Usage: make sync-specs MONO_DIR=/path/to/mono" >&2; exit 2)
	./scripts/sync-specs.sh "$(MONO_DIR)"
