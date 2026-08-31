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


if __name__ == "__main__":
    main()
