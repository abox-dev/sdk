ROOT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)

generate-api:
	redocly bundle python-sdk --config $(ROOT_DIR)/redocly.yaml -o $(ROOT_DIR)/spec/openapi_generated.python-sdk.yml
	python $(ROOT_DIR)/scripts/filter-public-openapi.py $(ROOT_DIR)/spec/openapi_generated.python-sdk.yml
	openapi-python-client generate --output-path $(ROOT_DIR)/packages/python-sdk/agentbox/api/api --overwrite --path $(ROOT_DIR)/spec/openapi_generated.python-sdk.yml
	rm -rf agentbox/api/client
	mv agentbox/api/api/agent_box_api_client agentbox/api/client
	rm -rf agentbox/api/api
	ruff format .

# Requires `buf` and the codegen dependency group (`make init`), which
# provides the `protoc-gen-connectrpc` plugin.
generate-envd:
	cd $(ROOT_DIR)/spec/envd && pwd && buf generate --template buf-python.gen.yaml
	./scripts/fix-python-pb.sh

	ruff format .

generate: generate-api generate-envd

typecheck:
	ty check

lint:
	ruff check .
	ruff format --check .

format:
	ruff format .

.PHONY: test
test:
	uv run pytest --verbose --numprocesses=4
