#!/usr/bin/env python3
"""Remove backend-only compatibility shapes from a bundled SDK OpenAPI file."""

from __future__ import annotations

import argparse
from pathlib import Path
from typing import Any

import yaml


REMOVED_SCHEMAS = {
    "CreatedAccessToken",
    "NewVolume",
    "NewAccessToken",
    "SandboxVolumeMount",
    "TemplateLegacy",
    "Volume",
    "VolumeAndToken",
    "VolumeToken",
}


def remove_volume_mounts(value: Any) -> None:
    if isinstance(value, dict):
        properties = value.get("properties")
        if isinstance(properties, dict):
            properties.pop("volumeMounts", None)
        for child in value.values():
            remove_volume_mounts(child)
    elif isinstance(value, list):
        for child in value:
            remove_volume_mounts(child)


def remove_deprecated_contract(value: Any) -> None:
    """Remove deprecated operations and fields from the generated public API."""
    if not isinstance(value, dict):
        return

    paths = value.get("paths")
    if isinstance(paths, dict):
        for path_item in paths.values():
            if not isinstance(path_item, dict):
                continue
            for method in list(path_item):
                operation = path_item[method]
                if isinstance(operation, dict) and operation.get("deprecated") is True:
                    del path_item[method]
                elif isinstance(operation, dict):
                    operation.pop("security", None)

    properties = value.get("properties")
    if isinstance(properties, dict):
        removed = {
            name
            for name, schema in properties.items()
            if isinstance(schema, dict) and schema.get("deprecated") is True
        }
        for name in removed:
            del properties[name]
        required = value.get("required")
        if isinstance(required, list):
            value["required"] = [name for name in required if name not in removed]

    for child in value.values():
        if isinstance(child, dict):
            remove_deprecated_contract(child)
        elif isinstance(child, list):
            for item in child:
                remove_deprecated_contract(item)


def assert_no_removed_refs(value: Any) -> None:
    if isinstance(value, dict):
        ref = value.get("$ref")
        if isinstance(ref, str) and ref.rsplit("/", 1)[-1] in REMOVED_SCHEMAS:
            raise RuntimeError(f"public schema still references removed component: {ref}")
        for child in value.values():
            assert_no_removed_refs(child)
    elif isinstance(value, list):
        for child in value:
            assert_no_removed_refs(child)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("path", type=Path)
    args = parser.parse_args()

    document = yaml.safe_load(args.path.read_text())
    remove_volume_mounts(document)
    remove_deprecated_contract(document)
    document.pop("security", None)
    document.get("components", {}).pop("securitySchemes", None)

    schemas = document.get("components", {}).get("schemas", {})
    for name in REMOVED_SCHEMAS:
        schemas.pop(name, None)

    assert_no_removed_refs(document)
    args.path.write_text(yaml.safe_dump(document, sort_keys=False))


if __name__ == "__main__":
    main()
