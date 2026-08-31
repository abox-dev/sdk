"""Shared filtering policies for generated public OpenAPI artifacts."""

from __future__ import annotations

from typing import Any

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

CLIENT_REMOVED_PROPERTIES = {"volumeMounts"}
REFERENCE_REMOVED_PROPERTIES = CLIENT_REMOVED_PROPERTIES | {
    "envdAccessToken",
    "trafficAccessToken",
}


def remove_properties(value: Any, names: set[str]) -> None:
    if isinstance(value, dict):
        properties = value.get("properties")
        if isinstance(properties, dict):
            removed = set(properties) & names
            for name in removed:
                del properties[name]
            required = value.get("required")
            if isinstance(required, list):
                value["required"] = [name for name in required if name not in removed]
        for child in value.values():
            remove_properties(child, names)
    elif isinstance(value, list):
        for child in value:
            remove_properties(child, names)


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
            raise RuntimeError(
                f"public schema still references removed component: {ref}"
            )
        for child in value.values():
            assert_no_removed_refs(child)
    elif isinstance(value, list):
        for child in value:
            assert_no_removed_refs(child)


def component_reference(reference: str) -> tuple[str, str] | None:
    prefix = "#/components/"
    if not reference.startswith(prefix):
        return None
    parts = reference[len(prefix) :].split("/")
    if len(parts) < 2:
        return None

    def decode(value: str) -> str:
        return value.replace("~1", "/").replace("~0", "~")

    return decode(parts[0]), decode(parts[1])


def collect_component_references(value: Any) -> set[tuple[str, str]]:
    references = set()
    if isinstance(value, dict):
        reference = value.get("$ref")
        if isinstance(reference, str):
            component = component_reference(reference)
            if component is not None:
                references.add(component)
        for child in value.values():
            references.update(collect_component_references(child))
    elif isinstance(value, list):
        for child in value:
            references.update(collect_component_references(child))
    return references


def prune_unreferenced_components(document: dict) -> None:
    """Keep only components reachable from the published OpenAPI surface."""
    components = document.get("components")
    if not isinstance(components, dict):
        return

    roots = {key: value for key, value in document.items() if key != "components"}
    pending = list(collect_component_references(roots))
    reachable = set()
    while pending:
        component = pending.pop()
        if component in reachable:
            continue
        section, name = component
        definition = components.get(section, {}).get(name)
        if definition is None:
            raise RuntimeError(
                f"public OpenAPI references missing component: {section}/{name}"
            )
        reachable.add(component)
        pending.extend(collect_component_references(definition) - reachable)

    for section in list(components):
        definitions = components[section]
        if not isinstance(definitions, dict):
            continue
        components[section] = {
            name: definition
            for name, definition in definitions.items()
            if (section, name) in reachable
        }
        if not components[section]:
            del components[section]


def filter_public_openapi(document: dict, *, for_reference: bool = False) -> dict:
    """Mutate and return a client-safe or stricter documentation-safe contract."""
    remove_properties(
        document,
        REFERENCE_REMOVED_PROPERTIES if for_reference else CLIENT_REMOVED_PROPERTIES,
    )
    remove_deprecated_contract(document)
    document.pop("security", None)
    document.get("components", {}).pop("securitySchemes", None)

    schemas = document.get("components", {}).get("schemas", {})
    for name in REMOVED_SCHEMAS:
        schemas.pop(name, None)

    assert_no_removed_refs(document)
    if for_reference:
        prune_unreferenced_components(document)
    return document
