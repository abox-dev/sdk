#!/usr/bin/env python3
"""Fast assertions for the generated public documentation contract."""

import hashlib
import json
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[1]
REFERENCE = ROOT / "reference"
BANNED_PATH_PARTS = ("/admin/", "/access-tokens", "/api-keys", "/volumes")
BANNED_ENVD_PATHS = {
    "/init",
    "/freeze",
    "/unfreeze",
    "/collapse",
    "/fsfreeze",
    "/fsthaw",
}
FORBIDDEN_REFERENCE_PROPERTIES = {
    "clientID",
    "envdAccessToken",
    "trafficAccessToken",
    "volumeMounts",
}
FORBIDDEN_SECURITY_SCHEMES = {
    "AccessTokenAuth",
    "AdminApiKeyAuth",
    "AdminTeamAuth",
    "AuthProviderBearerAuth",
    "AuthProviderTeamAuth",
}


def resolve_local_ref(document: dict, value: dict) -> dict:
    current = value
    visited = set()
    while isinstance(current, dict) and "$ref" in current:
        reference = current["$ref"]
        if not reference.startswith("#/") or reference in visited:
            break
        visited.add(reference)
        resolved = document
        for part in reference[2:].split("/"):
            resolved = resolved[part.replace("~1", "/").replace("~0", "~")]
        current = resolved
    return current


def array_item_properties(document: dict, schema: dict) -> dict:
    resolved = resolve_local_ref(document, schema)
    if not isinstance(resolved, dict) or resolved.get("type") != "array":
        return {}
    visited = set()
    while isinstance(resolved, dict) and resolved.get("type") == "array":
        marker = id(resolved)
        if marker in visited:
            return {}
        visited.add(marker)
        items = resolved.get("items")
        if not isinstance(items, dict):
            return {}
        resolved = resolve_local_ref(document, items)
    return resolved.get("properties", {}) if isinstance(resolved, dict) else {}


def assert_local_refs_resolve(document: dict) -> None:
    def visit(value):
        if isinstance(value, dict):
            reference = value.get("$ref")
            if isinstance(reference, str) and reference.startswith("#/"):
                resolved = document
                for part in reference[2:].split("/"):
                    resolved = resolved[part.replace("~1", "/").replace("~0", "~")]
            for child in value.values():
                visit(child)
        elif isinstance(value, list):
            for child in value:
                visit(child)

    visit(document)


def assert_public_schema(value) -> None:
    if isinstance(value, dict):
        properties = value.get("properties")
        if isinstance(properties, dict):
            assert not (set(properties) & FORBIDDEN_REFERENCE_PROPERTIES)
            assert not any(
                isinstance(schema, dict) and schema.get("deprecated") is True
                for schema in properties.values()
            )
        for child in value.values():
            assert_public_schema(child)
    elif isinstance(value, list):
        for child in value:
            assert_public_schema(child)


def main() -> None:
    manifest = json.loads((REFERENCE / "manifest.json").read_text())
    assert manifest["schemaVersion"] == 1
    for relative, expected in manifest["files"].items():
        actual = hashlib.sha256((REFERENCE / relative).read_bytes()).hexdigest()
        assert actual == expected, relative

    operations = json.loads((REFERENCE / "openapi/operations.json").read_text())
    ids = [operation["operationId"] for operation in operations]
    routes = {
        (operation["spec"], operation["method"], operation["path"])
        for operation in operations
    }
    assert len(ids) == len(set(ids)), "operationId values must be unique"
    for operation in operations:
        assert isinstance(operation.get("summary"), str), operation["operationId"]
        assert operation["summary"].strip() == operation["summary"]
        assert operation["summary"], operation["operationId"]
        assert "auth" in operation, operation["operationId"]
        auth = operation["auth"]
        if operation["spec"] == "control-plane":
            assert auth == {
                "type": "apiKey",
                "header": "X-API-Key",
                "required": True,
            }, operation["operationId"]
        elif operation["path"] == "/health":
            assert auth is None, operation["operationId"]
        else:
            assert auth == {
                "type": "apiKey",
                "header": "X-Access-Token",
                "required": False,
            }, operation["operationId"]
        markdown = REFERENCE / operation["markdown"]
        assert markdown.is_file(), operation["operationId"]
        rendered = markdown.read_text()
        assert "## Responses" in rendered, operation["operationId"]
    assert not any(
        part in operation["path"]
        for operation in operations
        for part in BANNED_PATH_PARTS
    )
    assert not any(operation["path"] in BANNED_ENVD_PATHS for operation in operations)

    config = yaml.safe_load((ROOT / "reference-config/operations.yaml").read_text())
    control = yaml.safe_load((ROOT / config["controlPlane"]["source"]).read_text())
    public_tags = set(config["controlPlane"]["publicTags"])
    expected = set()
    for path, path_item in control["paths"].items():
        for method, operation in path_item.items():
            if method not in {"get", "post", "put", "patch", "delete"}:
                continue
            if operation.get("deprecated") or operation.get("x-not-implemented"):
                continue
            if set(operation.get("tags", [])) & public_tags:
                expected.add(("control-plane", method.upper(), path))
    expected.update(
        {
            ("control-plane", "GET", "/teams/{teamID}/metrics"),
            ("control-plane", "GET", "/teams/{teamID}/metrics/max"),
        }
    )
    assert expected == {route for route in routes if route[0] == "control-plane"}

    rendered_specs = "\n".join(
        (REFERENCE / name).read_text()
        for name in ("openapi/control-plane.yml", "openapi/envd.yml")
    )
    assert "x-not-implemented" not in rendered_specs
    assert "SandboxEgressProxyConfig" not in rendered_specs
    for forbidden in FORBIDDEN_REFERENCE_PROPERTIES | FORBIDDEN_SECURITY_SCHEMES:
        assert forbidden not in rendered_specs

    public_documents = {}
    for name in ("openapi/control-plane.yml", "openapi/envd.yml"):
        document = yaml.safe_load((REFERENCE / name).read_text())
        public_documents[name.removeprefix("openapi/").removesuffix(".yml")] = document
        assert "securitySchemes" not in document.get("components", {})
        for path_item in document["paths"].values():
            for method, operation in path_item.items():
                if method.lower() in {"get", "post", "put", "patch", "delete"}:
                    assert "security" not in operation
        assert_public_schema(document)
        assert_local_refs_resolve(document)

    rendered_markdown = "\n".join(
        path.read_text() for path in (REFERENCE / "openapi/markdown").glob("*.md")
    )
    for forbidden in FORBIDDEN_REFERENCE_PROPERTIES:
        assert forbidden not in rendered_markdown
    assert "| Field | Type | Required | Description |" not in rendered_markdown
    assert "| Name | In | Required | Type | Description |" not in rendered_markdown

    for operation in operations:
        document = public_documents[operation["spec"]]
        operation_schema = document["paths"][operation["path"]][
            operation["method"].lower()
        ]
        markdown = (REFERENCE / operation["markdown"]).read_text()
        for response in operation_schema.get("responses", {}).values():
            resolved_response = resolve_local_ref(document, response)
            for media in resolved_response.get("content", {}).values():
                for field in array_item_properties(document, media.get("schema", {})):
                    assert f"- **`{field}`** ·" in markdown, (
                        operation["operationId"],
                        field,
                    )

    cli_markdown = "\n".join(
        path.read_text() for path in (REFERENCE / "sdk/cli").glob("*.md")
    )
    assert "[object Object]" not in cli_markdown

    javascript_files = {
        str(path.relative_to(REFERENCE))
        for path in (REFERENCE / "sdk/javascript").rglob("*.md")
    }
    assert not any("mapEntryInfo" in path for path in javascript_files)
    assert not any("formatSandboxTimeoutError" in path for path in javascript_files)
    assert not any("deserializeChart" in path for path in javascript_files)
    assert not any("parseOutput" in path for path in javascript_files)
    assert not any("extractError" in path for path in javascript_files)

    list_sandboxes = (REFERENCE / "openapi/markdown/listSandboxes.md").read_text()
    assert "Metadata query used to filter the sandboxes" in list_sandboxes
    assert "### 200" in list_sandboxes
    assert "Schema: `array<ListedSandbox>`" in list_sandboxes
    assert "- **`sandboxID`** · `string` · required" in list_sandboxes
    assert "- **`metadata`** · `string` · query · optional" in list_sandboxes
    create_sandbox = (REFERENCE / "openapi/markdown/createSandbox.md").read_text()
    assert "Schema: `NewSandbox`" in create_sandbox


if __name__ == "__main__":
    main()
