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


def markdown_cell(value: object) -> str:
    return str(value or "").replace("\n", " ").replace("|", "\\|").strip()


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


def render_schema(document: dict, schema: dict) -> list[str]:
    resolved, reference_name = resolve_ref(document, schema)
    lines = [f"Schema: `{reference_name or schema_type(document, resolved)}`", ""]
    properties = resolved.get("properties", {})
    if not properties:
        return lines
    required = set(resolved.get("required", []))
    lines.extend(
        [
            "| Field | Type | Required | Description |",
            "| --- | --- | --- | --- |",
        ]
    )
    for name, property_schema in properties.items():
        property_value, _ = resolve_ref(document, property_schema)
        lines.append(
            f"| `{name}` | `{markdown_cell(schema_type(document, property_schema))}` "
            f"| {'yes' if name in required else 'no'} "
            f"| {markdown_cell(property_value.get('description', ''))} |"
        )
    lines.append("")
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
        lines.extend(
            [
                "## Parameters",
                "",
                "| Name | In | Required | Type | Description |",
                "| --- | --- | --- | --- | --- |",
            ]
        )
        for parameter in parameters:
            lines.append(
                f"| `{parameter.get('name', '')}` | {parameter.get('in', '')} "
                f"| {'yes' if parameter.get('required') else 'no'} "
                f"| `{markdown_cell(schema_type(document, parameter.get('schema', {})))}` "
                f"| {markdown_cell(parameter.get('description', ''))} |"
            )
        lines.append("")

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


def build_openapi(name: str, config: dict) -> list[dict]:
    source = ROOT / config["source"]
    document = yaml.safe_load(source.read_text())
    assigned = config["operations"]
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
                "Public operations require an id, group and slug:\n  "
                + "\n  ".join(missing)
            )

    records = []
    for path, path_item in document["paths"].items():
        for method, operation in path_item.items():
            key = operation_key(method, path)
            if key not in assigned:
                continue
            metadata = assigned[key]
            required = {"id", "group", "slug"}
            if set(metadata) != required:
                raise SystemExit(f"{key} must define exactly {sorted(required)}")
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
                    "spec": "control-plane" if name == "controlPlane" else "envd",
                    "markdown": f"openapi/markdown/{metadata['id']}.md",
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
