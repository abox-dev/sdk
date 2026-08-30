CODEGEN_IMAGE ?= agentbox-sdk-codegen

.PHONY: generate generate-in-container codegen-image sync-specs

# Generate exclusively from the checked-in snapshots under spec/.
generate: codegen-image
	docker run --rm -v "$$PWD:/workspace" $(CODEGEN_IMAGE) make generate-in-container

codegen-image:
	docker build -q -t $(CODEGEN_IMAGE) -f codegen.Dockerfile .

generate-in-container:
	cd packages/js-sdk && pnpm generate
	cd packages/python-sdk && make generate

# Maintainer-only update from a local mono checkout.
sync-specs:
	@test -n "$(MONO_DIR)" || (echo "Usage: make sync-specs MONO_DIR=/path/to/mono" >&2; exit 2)
	./scripts/sync-specs.sh "$(MONO_DIR)"
