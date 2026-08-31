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

    for name in ("openapi/control-plane.yml", "openapi/envd.yml"):
        document = yaml.safe_load((REFERENCE / name).read_text())
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
    create_sandbox = (REFERENCE / "openapi/markdown/createSandbox.md").read_text()
    assert "Schema: `NewSandbox`" in create_sandbox


if __name__ == "__main__":
    main()
