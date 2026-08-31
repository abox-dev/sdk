#!/usr/bin/env python3
"""Build the public documentation contract from checked-in SDK snapshots."""

from __future__ import annotations

import hashlib
import json
import os
import re
import shutil
import subprocess
import tempfile
from pathlib import Path

import yaml

from public_openapi import filter_public_openapi

ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "reference"
HTTP_METHODS = {"get", "post", "put", "patch", "delete"}


def sha256(path: Path) -> str:
    return hashlib.sha256(path.read_bytes()).hexdigest()


def operation_key(method: str, path: str) -> str:
    return f"{method.upper()} {path}"


def strip_extension(value):
    if isinstance(value, dict):
        return {
            key: strip_extension(item)
            for key, item in value.items()
            if not key.startswith("x-")
            and not (isinstance(item, dict) and item.get("x-not-implemented") is True)
        }
    if isinstance(value, list):
        return [strip_extension(item) for item in value]
    return value


def resolve_ref(document: dict, value: dict) -> tuple[dict, str | None]:
    current = value
    reference_name = None
    visited = set()
    while isinstance(current, dict) and "$ref" in current:
        reference = current["$ref"]
        if not reference.startswith("#/") or reference in visited:
            break
        visited.add(reference)
        reference_name = reference.rsplit("/", 1)[-1]
        resolved = document
        for part in reference[2:].split("/"):
            resolved = resolved[part.replace("~1", "/").replace("~0", "~")]
        current = resolved
    return current, reference_name


def markdown_text(value: object) -> str:
    return " ".join(str(value or "").split())


def schema_type(document: dict, schema: dict) -> str:
    resolved, reference_name = resolve_ref(document, schema)
    if reference_name:
        return reference_name
    if "oneOf" in resolved:
        return " | ".join(schema_type(document, item) for item in resolved["oneOf"])
    if "anyOf" in resolved:
        return " | ".join(schema_type(document, item) for item in resolved["anyOf"])
    if resolved.get("type") == "array":
        return f"array<{schema_type(document, resolved.get('items', {}))}>"
    if resolved.get("enum"):
        return " | ".join(str(item) for item in resolved["enum"])
    return resolved.get("type", "object")


def render_named_item(
    name: str, type_name: str, qualifiers: list[str], description: object
) -> list[str]:
    details = " · ".join([f"`{type_name}`", *qualifiers])
    lines = [f"- **`{name}`** · {details}", ""]
    normalized_description = markdown_text(description)
    if normalized_description:
        lines.extend([f"  {normalized_description}", ""])
    return lines


def schema_fields(document: dict, schema: dict) -> dict:
    """Return the object schema whose fields describe a response value."""
    resolved, _ = resolve_ref(document, schema)
    visited = set()
    while isinstance(resolved, dict) and resolved.get("type") == "array":
        marker = id(resolved)
        if marker in visited:
            return {}
        visited.add(marker)
        items = resolved.get("items")
        if not isinstance(items, dict):
            return {}
        resolved, _ = resolve_ref(document, items)
    return resolved if isinstance(resolved, dict) else {}


def render_schema(document: dict, schema: dict) -> list[str]:
    resolved, reference_name = resolve_ref(document, schema)
    lines = [f"Schema: `{reference_name or schema_type(document, resolved)}`", ""]
    fields = schema_fields(document, resolved)
    properties = fields.get("properties", {})
    if not properties:
        return lines
    required = set(fields.get("required", []))
    for name, property_schema in properties.items():
        property_value, _ = resolve_ref(document, property_schema)
        lines.extend(
            render_named_item(
                name,
                schema_type(document, property_schema),
                ["required" if name in required else "optional"],
                property_value.get("description", ""),
            )
        )
    return lines


def render_operation_markdown(document: dict, record: dict) -> str:
    operation = document["paths"][record["path"]][record["method"].lower()]
    lines = [f"# {record['method']} {record['path']}", ""]
    if operation.get("summary"):
        lines.extend([operation["summary"].strip(), ""])
    if operation.get("description"):
        lines.extend([operation["description"].strip(), ""])

    parameters = []
    for parameter in operation.get("parameters", []):
        resolved, _ = resolve_ref(document, parameter)
        parameters.append(resolved)
    if parameters:
        lines.extend(["## Parameters", ""])
        for parameter in parameters:
            lines.extend(
                render_named_item(
                    parameter.get("name", ""),
                    schema_type(document, parameter.get("schema", {})),
                    [
                        parameter.get("in", ""),
                        "required" if parameter.get("required") else "optional",
                    ],
                    parameter.get("description", ""),
                )
            )

    if operation.get("requestBody"):
        request_body, _ = resolve_ref(document, operation["requestBody"])
        lines.extend(
            [
                "## Request body",
                "",
                f"Required: {'yes' if request_body.get('required') else 'no'}",
                "",
            ]
        )
        if request_body.get("description"):
            lines.extend([request_body["description"].strip(), ""])
        for content_type, media in request_body.get("content", {}).items():
            lines.extend([f"### {content_type}", ""])
            if media.get("schema"):
                lines.extend(render_schema(document, media["schema"]))

    lines.extend(["## Responses", ""])
    for status, response in operation.get("responses", {}).items():
        resolved, _ = resolve_ref(document, response)
        lines.extend([f"### {status}", ""])
        if resolved.get("description"):
            lines.extend([resolved["description"].strip(), ""])
        for content_type, media in resolved.get("content", {}).items():
            lines.extend([f"Content-Type: `{content_type}`", ""])
            if media.get("schema"):
                lines.extend(render_schema(document, media["schema"]))
    return "\n".join(lines).rstrip() + "\n"


def normalize_operation_auth(
    document: dict, operation: dict, public_scheme_names: set[str]
) -> dict | None:
    """Return one documentation-safe auth mechanism for an operation."""
    security = operation.get("security", document.get("security"))
    if not security:
        return None
    if not isinstance(security, list) or not all(
        isinstance(requirement, dict) for requirement in security
    ):
        raise SystemExit("Operation security must be an OpenAPI security array")

    anonymous_allowed = any(not requirement for requirement in security)
    public_requirements = [
        requirement
        for requirement in security
        if requirement and set(requirement).issubset(public_scheme_names)
    ]
    public_schemes = {
        name for requirement in public_requirements for name in requirement
    }
    if not public_schemes:
        if anonymous_allowed:
            return None
        raise SystemExit(
            "Authenticated public operation has no allowed public auth scheme"
        )
    if len(public_schemes) != 1 or any(
        len(requirement) != 1 for requirement in public_requirements
    ):
        raise SystemExit(
            "Public reference supports exactly one normalized auth mechanism"
        )

    scheme_name = next(iter(public_schemes))
    scheme = document.get("components", {}).get("securitySchemes", {}).get(scheme_name)
    if not isinstance(scheme, dict):
        raise SystemExit(f"Public auth scheme is missing: {scheme_name}")
    if scheme.get("type") != "apiKey" or scheme.get("in") != "header":
        raise SystemExit(f"Public auth scheme {scheme_name} must be an apiKey header")
    header = scheme.get("name")
    if not isinstance(header, str) or not header:
        raise SystemExit(f"Public auth scheme {scheme_name} has no header name")
    return {
        "type": "apiKey",
        "header": header,
        "required": not anonymous_allowed,
    }


def apply_public_operation_security(
    source: dict,
    public: dict,
    records: list[dict],
    public_scheme_names: set[str],
) -> None:
    """Restore only the normalized public auth contract after filtering."""
    source_schemes = source.get("components", {}).get("securitySchemes", {})
    published_schemes = {}
    scheme_by_header = {}
    for name in sorted(public_scheme_names):
        scheme = source_schemes.get(name)
        if not isinstance(scheme, dict):
            raise SystemExit(f"Public auth scheme is missing: {name}")
        if scheme.get("type") != "apiKey" or scheme.get("in") != "header":
            raise SystemExit(f"Public auth scheme {name} must be an apiKey header")
        header = scheme.get("name")
        if not isinstance(header, str) or not header:
            raise SystemExit(f"Public auth scheme {name} has no header name")
        if header in scheme_by_header:
            raise SystemExit(f"Multiple public auth schemes use header {header}")
        published_schemes[name] = strip_extension(scheme)
        scheme_by_header[header] = name

    public.setdefault("components", {})["securitySchemes"] = published_schemes
    for record in records:
        operation = public["paths"][record["path"]][record["method"].lower()]
        auth = record["auth"]
        if auth is None:
            operation["security"] = []
            continue
        scheme_name = scheme_by_header.get(auth["header"])
        if scheme_name is None:
            raise SystemExit(
                f"No public auth scheme for {record['operationId']} header {auth['header']}"
            )
        requirement = {scheme_name: []}
        operation["security"] = [requirement] if auth["required"] else [{}, requirement]


def build_openapi(name: str, config: dict) -> list[dict]:
    source = ROOT / config["source"]
    document = yaml.safe_load(source.read_text())
    assigned = config["operations"]
    public_auth_schemes = set(config.get("publicAuthSchemes", []))
    seen = set()
    output_paths = {}

    if name == "controlPlane":
        public_tags = set(config["publicTags"])
        candidates = set()
        for path, path_item in document["paths"].items():
            for method, operation in path_item.items():
                if (
                    method not in HTTP_METHODS
                    or operation.get("deprecated")
                    or operation.get("x-not-implemented") is True
                ):
                    continue
                tags = set(operation.get("tags", []))
                if tags & public_tags:
                    candidates.add(operation_key(method, path))
        candidates.update(
            {
                "GET /teams/{teamID}/metrics",
                "GET /teams/{teamID}/metrics/max",
            }
        )
        missing = sorted(candidates - set(assigned))
        if missing:
            raise SystemExit(
                "Public operations require an id, group, slug and summary:\n  "
                + "\n  ".join(missing)
            )

    records = []
    for path, path_item in document["paths"].items():
        for method, operation in path_item.items():
            key = operation_key(method, path)
            if key not in assigned:
                continue
            metadata = assigned[key]
            required = {"id", "group", "slug", "summary"}
            if set(metadata) != required:
                raise SystemExit(f"{key} must define exactly {sorted(required)}")
            if (
                not isinstance(metadata["summary"], str)
                or not metadata["summary"].strip()
            ):
                raise SystemExit(f"{key} must define a non-empty summary")
            if metadata["id"] in seen:
                raise SystemExit(f"Duplicate operation id: {metadata['id']}")
            seen.add(metadata["id"])
            public_operation = strip_extension(operation)
            public_operation["operationId"] = metadata["id"]
            public_operation["tags"] = [metadata["group"]]
            output_paths.setdefault(path, {})[method] = public_operation
            records.append(
                {
                    "operationId": metadata["id"],
                    "method": method.upper(),
                    "path": path,
                    "group": metadata["group"],
                    "slug": metadata["slug"],
                    "summary": markdown_text(metadata["summary"]),
                    "spec": "control-plane" if name == "controlPlane" else "envd",
                    "markdown": f"openapi/markdown/{metadata['id']}.md",
                    "auth": normalize_operation_auth(
                        document, operation, public_auth_schemes
                    ),
                }
            )

    unknown = sorted(
        set(assigned)
        - {
            operation_key(method, path)
            for path, item in document["paths"].items()
            for method in item
            if method in HTTP_METHODS
        }
    )
    if unknown:
        raise SystemExit(
            "Configured operations are missing from the spec:\n  "
            + "\n  ".join(unknown)
        )

    public = strip_extension(document)
    public["paths"] = output_paths
    public["tags"] = [
        {"name": group} for group in sorted({record["group"] for record in records})
    ]
    filter_public_openapi(public, for_reference=True)
    apply_public_operation_security(document, public, records, public_auth_schemes)
    destination = (
        OUT
        / "openapi"
        / ("control-plane.yml" if name == "controlPlane" else "envd.yml")
    )
    destination.parent.mkdir(parents=True, exist_ok=True)
    destination.write_text(yaml.safe_dump(public, sort_keys=False, allow_unicode=True))
    for record in records:
        markdown = OUT / record["markdown"]
        markdown.parent.mkdir(parents=True, exist_ok=True)
        markdown.write_text(render_operation_markdown(public, record))
    return sorted(records, key=lambda item: (item["group"], item["slug"]))


def proto_reference(source: Path, destination: Path, title: str) -> None:
    text = source.read_text()
    service_match = re.search(r"service\s+(\w+)\s*\{(.*?)\n\}", text, re.S)
    if not service_match:
        raise SystemExit(f"No service found in {source}")
    service, body = service_match.groups()
    rpcs = re.findall(
        r"(?:\s*//\s*(.*?)\n)?\s*rpc\s+(\w+)\s*\(([^)]+)\)\s*returns\s*\(([^)]+)\)",
        body,
    )
    lines = [f"# {title}", "", f"Service: `{service}`", ""]
    for comment, name, request, response in rpcs:
        lines.extend(
            [
                f"## {name}",
                "",
                comment.strip() if comment else "Public RPC exposed by envd.",
                "",
                f"- Request: `{request.strip()}`",
                f"- Response: `{response.strip()}`",
                "",
            ]
        )
    lines.extend(["## Message types", ""])
    for match in re.finditer(
        r"(?:^|\n)(?:\s*//\s*(.*?)\n)?\s*message\s+(\w+)\s*\{", text
    ):
        comment, name = match.groups()
        depth = 1
        cursor = match.end()
        while depth and cursor < len(text):
            if text[cursor] == "{":
                depth += 1
            elif text[cursor] == "}":
                depth -= 1
            cursor += 1
        body = text[match.end() : cursor - 1]
        fields = re.findall(
            r"(?:\s*//\s*(.*?)\n)?\s*(repeated\s+|optional\s+)?([.\w]+)\s+(\w+)\s*=\s*(\d+)",
            body,
        )
        lines.extend([f"### {name}", ""])
        if comment:
            lines.extend([comment.strip(), ""])
        if fields:
            lines.extend(
                [
                    "| Field | Type | Number | Description |",
                    "| --- | --- | ---: | --- |",
                ]
            )
            for field_comment, modifier, field_type, field_name, number in fields:
                type_name = f"{modifier.strip()} {field_type}".strip()
                lines.append(
                    f"| `{field_name}` | `{type_name}` | {number} | {field_comment.strip() if field_comment else ''} |"
                )
            lines.append("")
    destination.parent.mkdir(parents=True, exist_ok=True)
    destination.write_text("\n".join(lines))


def run_sdk_generators() -> None:
    groups = json.loads((ROOT / "reference-config/sdk-groups.json").read_text())
    for name, settings in groups["javascript"].items():
        if set(settings) != {"entryPoint", "tsconfig"}:
            raise SystemExit(
                f"JavaScript group {name} must define entryPoint and tsconfig"
            )
        entry_point = ROOT / settings["entryPoint"]
        tsconfig = ROOT / settings["tsconfig"]
        if not entry_point.is_file() or not tsconfig.is_file():
            raise SystemExit(f"JavaScript group {name} has an invalid configuration")
        destination = OUT / "sdk/javascript" / name
        command = [
            "pnpm",
            "exec",
            "typedoc",
            "--plugin",
            "typedoc-plugin-markdown",
            "--entryPointStrategy",
            "resolve",
            "--tsconfig",
            str(tsconfig),
            "--readme",
            "none",
            "--disableSources",
            "--excludeInternal",
            "--excludePrivate",
            "--excludeProtected",
            "--hidePageHeader",
            "--out",
            str(destination),
            str(entry_point),
        ]
        subprocess.run(command, cwd=ROOT, check=True)

    with tempfile.TemporaryDirectory(prefix="agentbox-reference-") as temporary:
        python_root = Path(temporary)
        shutil.copytree(ROOT / "packages/python-sdk/agentbox", python_root / "agentbox")
        shutil.copytree(
            ROOT / "packages/code-interpreter-python/agentbox_code_interpreter",
            python_root / "agentbox_code_interpreter",
        )
        for namespace in (
            "sandbox_sync",
            "sandbox_async",
            "template",
            "template_sync",
            "template_async",
        ):
            (python_root / "agentbox" / namespace / "__init__.py").touch()

        for name, entries in groups["python"].items():
            destination = OUT / "sdk/python" / f"{name}.md"
            destination.parent.mkdir(parents=True, exist_ok=True)
            command = ["griffe2md", *entries, "-o", str(destination)]
            environment = os.environ.copy()
            environment["PYTHONPATH"] = os.pathsep.join(
                [str(python_root), environment.get("PYTHONPATH", "")]
            )
            subprocess.run(command, cwd=ROOT, env=environment, check=True)

    subprocess.run(
        [
            "tsx",
            "--tsconfig",
            "../../reference-config/cli-reference.json",
            "src/reference.ts",
            str(OUT / "sdk/cli"),
        ],
        cwd=ROOT / "packages/cli",
        check=True,
    )


def package_versions() -> dict:
    packages = {
        "@abox-dev/sdk": ROOT / "packages/js-sdk/package.json",
        "@abox-dev/code-interpreter": ROOT
        / "packages/code-interpreter-js/package.json",
        "@abox-dev/cli": ROOT / "packages/cli/package.json",
    }
    versions = {
        name: json.loads(path.read_text())["version"] for name, path in packages.items()
    }
    for name, path in {
        "abox-sdk": ROOT / "packages/python-sdk/pyproject.toml",
        "abox-code-interpreter": ROOT
        / "packages/code-interpreter-python/pyproject.toml",
    }.items():
        versions[name] = re.search(
            r'^version = "([^"]+)"', path.read_text(), re.M
        ).group(1)
    return versions


def main() -> None:
    if OUT.exists():
        shutil.rmtree(OUT)
    config = yaml.safe_load((ROOT / "reference-config/operations.yaml").read_text())
    operations = []
    operations.extend(build_openapi("controlPlane", config["controlPlane"]))
    operations.extend(build_openapi("envd", config["envd"]))
    (OUT / "openapi/operations.json").write_text(
        json.dumps(operations, indent=2) + "\n"
    )
    proto_reference(
        ROOT / "spec/envd/process/process.proto",
        OUT / "connect/process.md",
        "Process API",
    )
    proto_reference(
        ROOT / "spec/envd/filesystem/filesystem.proto",
        OUT / "connect/filesystem.md",
        "Filesystem API",
    )
    run_sdk_generators()

    source = json.loads((ROOT / "spec/source.json").read_text())
    files = {
        str(path.relative_to(OUT)): sha256(path)
        for path in sorted(OUT.rglob("*"))
        if path.is_file()
    }
    manifest = {
        "schemaVersion": config["schemaVersion"],
        "packages": package_versions(),
        "monoRevision": source["revision"],
        "files": files,
    }
    (OUT / "manifest.json").write_text(
        json.dumps(manifest, indent=2, sort_keys=True) + "\n"
    )


if __name__ == "__main__":
    main()
