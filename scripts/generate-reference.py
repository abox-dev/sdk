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


def merge_schema(target: dict, source: dict) -> None:
    """Merge one allOf branch into a documentation view of the schema."""
    for key, value in source.items():
        if key == "properties" and isinstance(value, dict):
            target.setdefault("properties", {}).update(value)
        elif key == "required" and isinstance(value, list):
            target["required"] = list(
                dict.fromkeys([*target.get("required", []), *value])
            )
        elif key == "nullable":
            target[key] = bool(target.get(key)) or bool(value)
        elif key not in target:
            target[key] = value


def normalize_schema(
    document: dict,
    schema: dict,
    ancestors: frozenset[tuple[str, object]] = frozenset(),
) -> tuple[dict, str | None]:
    """Resolve refs and flatten allOf while preserving wrapper metadata."""
    resolved, reference_name = resolve_ref(document, schema)
    if not isinstance(resolved, dict):
        return {}, reference_name
    normalized = dict(resolved)
    if schema is not resolved:
        normalized.update(
            {key: value for key, value in schema.items() if key != "$ref"}
        )

    all_of = normalized.pop("allOf", None)
    if not isinstance(all_of, list):
        return normalized, reference_name

    marker = (
        ("reference", reference_name) if reference_name else ("schema", id(resolved))
    )
    if marker in ancestors:
        return normalized, reference_name
    descendants = ancestors | {marker}
    branch_names = []
    for branch in all_of:
        if not isinstance(branch, dict):
            continue
        branch_schema, branch_name = normalize_schema(document, branch, descendants)
        merge_schema(normalized, branch_schema)
        if branch_name:
            branch_names.append(branch_name)
    if reference_name is None and len(branch_names) == 1:
        reference_name = branch_names[0]
    return normalized, reference_name


def schema_type(document: dict, schema: dict) -> str:
    resolved, reference_name = normalize_schema(document, schema)
    if reference_name:
        type_name = reference_name
    elif "oneOf" in resolved:
        type_name = " | ".join(
            schema_type(document, item) for item in resolved["oneOf"]
        )
    elif "anyOf" in resolved:
        type_name = " | ".join(
            schema_type(document, item) for item in resolved["anyOf"]
        )
    elif resolved.get("type") == "array":
        type_name = f"array<{schema_type(document, resolved.get('items', {}))}>"
    elif resolved.get("enum"):
        type_name = " | ".join(str(item) for item in resolved["enum"])
    elif resolved.get("properties"):
        type_name = "object"
    else:
        type_name = resolved.get("type", "object")
    if resolved.get("nullable") and "null" not in type_name.split(" | "):
        return f"{type_name} | null"
    return type_name


def render_named_item(
    name: str, type_name: str, qualifiers: list[str], description: object
) -> list[str]:
    details = " · ".join([f"`{type_name}`", *qualifiers])
    lines = [f"- **`{name}`** · {details}", ""]
    normalized_description = markdown_text(description)
    if normalized_description:
        lines.extend([f"  {normalized_description}", ""])
    return lines


def schema_value(value: object) -> str:
    if isinstance(value, (dict, list, bool)) or value is None:
        return json.dumps(value, sort_keys=True)
    return str(value)


def render_schema_metadata(document: dict, schema: dict, indent: str = "") -> list[str]:
    """Render validation and representation details without losing referenced enums."""
    resolved, reference_name = normalize_schema(document, schema)
    if not isinstance(resolved, dict):
        return []

    lines = []
    if resolved.get("format"):
        lines.extend([f"{indent}Format: `{resolved['format']}`", ""])
    if resolved.get("enum"):
        label = f" for `{reference_name}`" if reference_name else ""
        values = " | ".join(f"`{schema_value(value)}`" for value in resolved["enum"])
        lines.extend([f"{indent}Allowed values{label}: {values}", ""])
    if "default" in resolved:
        lines.extend([f"{indent}Default: `{schema_value(resolved['default'])}`", ""])
    for key, label in (
        ("minimum", "Minimum"),
        ("maximum", "Maximum"),
        ("exclusiveMinimum", "Exclusive minimum"),
        ("exclusiveMaximum", "Exclusive maximum"),
        ("minLength", "Minimum length"),
        ("maxLength", "Maximum length"),
        ("minItems", "Minimum items"),
        ("maxItems", "Maximum items"),
        ("pattern", "Pattern"),
    ):
        if key in resolved:
            lines.extend([f"{indent}{label}: `{schema_value(resolved[key])}`", ""])
    if "uniqueItems" in resolved:
        unique = "yes" if resolved["uniqueItems"] else "no"
        lines.extend([f"{indent}Unique items: `{unique}`", ""])
    if resolved.get("type") == "array" and isinstance(resolved.get("items"), dict):
        lines.extend(render_schema_metadata(document, resolved["items"], indent))
    return lines


def schema_fields(document: dict, schema: dict) -> tuple[dict, str | None]:
    """Return the object schema whose fields describe a value."""
    resolved, reference_name = normalize_schema(document, schema)
    visited = set()
    while isinstance(resolved, dict) and resolved.get("type") == "array":
        marker = id(resolved)
        if marker in visited:
            return {}, reference_name
        visited.add(marker)
        items = resolved.get("items")
        if not isinstance(items, dict):
            return {}, reference_name
        resolved, item_reference = normalize_schema(document, items)
        reference_name = item_reference or reference_name
    return (resolved if isinstance(resolved, dict) else {}), reference_name


def render_schema_fields(
    document: dict,
    schema: dict,
    prefix: str = "",
    ancestors: frozenset[tuple[str, object]] = frozenset(),
) -> list[str]:
    """Render public fields recursively with cycle-safe dotted paths."""
    fields, reference_name = schema_fields(document, schema)
    marker = ("reference", reference_name) if reference_name else ("schema", id(fields))
    if marker in ancestors:
        return []
    descendants = ancestors | {marker}
    lines = []

    variants = fields.get("oneOf") or fields.get("anyOf")
    if isinstance(variants, list):
        discriminator = fields.get("discriminator", {})
        discriminator_name = discriminator.get("propertyName")
        mapping = discriminator.get("mapping", {})
        for variant in variants:
            if not isinstance(variant, dict):
                continue
            _, variant_name = resolve_ref(document, variant)
            variant_name = variant_name or schema_type(document, variant)
            reference = variant.get("$ref")
            discriminator_value = next(
                (value for value, target in mapping.items() if target == reference),
                None,
            )
            qualifier = ""
            if discriminator_name and discriminator_value is not None:
                qualifier = (
                    f" · discriminator `{discriminator_name}` = `{discriminator_value}`"
                )
            lines.extend([f"- **Variant `{variant_name}`**{qualifier}", ""])
            variant_prefix = (
                f"{prefix}<{variant_name}>" if prefix else f"<{variant_name}>"
            )
            lines.extend(
                render_schema_fields(document, variant, variant_prefix, descendants)
            )
        return lines

    required = set(fields.get("required", []))
    for name, property_schema in fields.get("properties", {}).items():
        field_name = f"{prefix}.{name}" if prefix else name
        property_value, _ = normalize_schema(document, property_schema)
        lines.extend(
            render_named_item(
                field_name,
                schema_type(document, property_schema),
                ["required" if name in required else "optional"],
                property_value.get("description", ""),
            )
        )
        lines.extend(render_schema_metadata(document, property_schema, "  "))
        lines.extend(
            render_schema_fields(document, property_schema, field_name, descendants)
        )

    additional = fields.get("additionalProperties")
    if isinstance(additional, dict):
        field_name = f"{prefix}.*" if prefix else "*"
        additional_value, _ = normalize_schema(document, additional)
        lines.extend(
            render_named_item(
                field_name,
                schema_type(document, additional),
                ["additional property"],
                additional_value.get("description", ""),
            )
        )
        lines.extend(render_schema_metadata(document, additional, "  "))
        lines.extend(
            render_schema_fields(document, additional, field_name, descendants)
        )
    return lines


def render_schema(document: dict, schema: dict) -> list[str]:
    resolved, reference_name = normalize_schema(document, schema)
    lines = [f"Schema: `{reference_name or schema_type(document, resolved)}`", ""]
    if isinstance(resolved, dict) and resolved.get("description"):
        lines.extend([markdown_text(resolved["description"]), ""])
    lines.extend(render_schema_metadata(document, schema))
    lines.extend(render_schema_fields(document, resolved))
    return lines


def render_parameter_serialization(document: dict, parameter: dict) -> list[str]:
    if "style" not in parameter and "explode" not in parameter:
        return []
    location = parameter.get("in")
    style = parameter.get(
        "style",
        {"query": "form", "cookie": "form", "path": "simple", "header": "simple"}.get(
            location
        ),
    )
    if not style:
        raise SystemExit(
            f"Cannot infer serialization style for parameter {parameter.get('name')}"
        )
    explode = parameter.get("explode", style == "form")
    details = [f"style `{style}`", f"explode `{'true' if explode else 'false'}`"]
    schema, _ = normalize_schema(document, parameter.get("schema", {}))
    if location == "query" and schema.get("type") == "array":
        name = parameter.get("name", "parameter")
        if style == "form":
            wire_format = (
                f"{name}=value1&{name}=value2" if explode else f"{name}=value1,value2"
            )
            details.append(f"wire format `{wire_format}`")
        elif style == "spaceDelimited":
            details.append(f"wire format `{name}=value1%20value2`")
        elif style == "pipeDelimited":
            details.append(f"wire format `{name}=value1|value2`")
    return [f"  Serialization: {'; '.join(details)}", ""]


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
            parameter_schema = parameter.get("schema", {})
            lines.extend(
                render_named_item(
                    parameter.get("name", ""),
                    schema_type(document, parameter_schema),
                    [
                        parameter.get("in", ""),
                        "required" if parameter.get("required") else "optional",
                    ],
                    parameter.get("description", ""),
                )
            )
            lines.extend(render_schema_metadata(document, parameter_schema, "  "))
            lines.extend(render_parameter_serialization(document, parameter))

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
        headers = resolved.get("headers", {})
        if headers:
            lines.extend(["#### Response headers", ""])
            for name, header in headers.items():
                resolved_header, _ = resolve_ref(document, header)
                header_schema = resolved_header.get("schema", {})
                lines.extend(
                    render_named_item(
                        name,
                        schema_type(document, header_schema),
                        ["response header"],
                        resolved_header.get("description", ""),
                    )
                )
                lines.extend(render_schema_metadata(document, header_schema, "  "))
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


def apply_envd_routing(public: dict, records: list[dict]) -> None:
    """Publish the stable sandbox proxy contract required by envd operations."""
    public["servers"] = [
        {
            "url": "https://sandbox.agentbox-runtime.ru",
            "description": "AgentBox sandbox proxy. Routing headers are required.",
        }
    ]
    parameters = {
        "AgentboxSandboxId": {
            "name": "Agentbox-Sandbox-Id",
            "in": "header",
            "required": True,
            "description": "Identifier of the sandbox that receives the request.",
            "schema": {"type": "string"},
        },
        "AgentboxSandboxPort": {
            "name": "Agentbox-Sandbox-Port",
            "in": "header",
            "required": True,
            "description": "Internal envd HTTP port exposed through the sandbox proxy.",
            "schema": {"type": "integer", "default": 49983},
        },
    }
    public.setdefault("components", {}).setdefault("parameters", {}).update(parameters)
    for record in records:
        operation = public["paths"][record["path"]][record["method"].lower()]
        operation["parameters"] = [
            {"$ref": "#/components/parameters/AgentboxSandboxId"},
            {"$ref": "#/components/parameters/AgentboxSandboxPort"},
            *operation.get("parameters", []),
        ]


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
    if name == "envd":
        apply_envd_routing(public, records)
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


PROTO_TOKEN = re.compile(
    r"""
    (?P<whitespace>\s+)
    |(?P<line_comment>//[^\n]*)
    |(?P<block_comment>/\*.*?\*/)
    |(?P<string>"(?:\\.|[^"\\])*")
    |(?P<identifier>[A-Za-z_][A-Za-z0-9_]*)
    |(?P<number>-?(?:0[xX][0-9A-Fa-f]+|\d+))
    |(?P<symbol>[{}()\[\]<>=;,.])
    |(?P<unknown>.)
    """,
    re.S | re.X,
)


def proto_tokens(text: str) -> list[tuple[str, str]]:
    tokens = []
    for match in PROTO_TOKEN.finditer(text):
        kind = match.lastgroup
        value = match.group()
        if kind == "whitespace":
            continue
        if kind == "unknown":
            raise ValueError(f"Unsupported protobuf token: {value!r}")
        if kind == "line_comment":
            tokens.append(("comment", value.removeprefix("//").strip()))
        elif kind == "block_comment":
            body = value.removeprefix("/*").removesuffix("*/")
            tokens.append(
                (
                    "comment",
                    " ".join(
                        line.strip().removeprefix("*").strip()
                        for line in body.splitlines()
                    ),
                )
            )
        else:
            tokens.append((kind, value))
    return tokens


class ProtoParser:
    def __init__(self, text: str):
        self.tokens = proto_tokens(text)
        self.position = 0
        self.services = []
        self.messages = []
        self.enums = []

    def current(self, value: str | None = None) -> tuple[str, str] | bool | None:
        if self.position >= len(self.tokens):
            return None
        token = self.tokens[self.position]
        return token if value is None else token[1] == value

    def take(self, value: str | None = None) -> tuple[str, str]:
        token = self.current()
        if token is None or (value is not None and token[1] != value):
            found = token[1] if token else "end of file"
            raise ValueError(f"Expected {value or 'token'}, found {found}")
        self.position += 1
        return token

    def accept(self, value: str) -> bool:
        if self.current(value):
            self.position += 1
            return True
        return False

    def comments(self) -> str:
        comments = []
        while self.current() and self.current()[0] == "comment":
            comments.append(self.take()[1])
        return " ".join(filter(None, comments))

    def qualified_name(self) -> str:
        parts = []
        if self.accept("."):
            parts.append("")
        if not self.current() or self.current()[0] != "identifier":
            raise ValueError("Expected protobuf type name")
        parts.append(self.take()[1])
        while self.accept("."):
            if not self.current() or self.current()[0] != "identifier":
                raise ValueError("Expected protobuf type name after '.'")
            parts.append(self.take()[1])
        return ".".join(parts)

    def field_type(self) -> str:
        if self.accept("map"):
            self.take("<")
            key = self.qualified_name()
            self.take(",")
            value = self.qualified_name()
            self.take(">")
            return f"map<{key}, {value}>"
        return self.qualified_name()

    def skip_balanced(self, opening: str, closing: str) -> None:
        self.take(opening)
        depth = 1
        while depth:
            token = self.take()[1]
            if token == opening:
                depth += 1
            elif token == closing:
                depth -= 1

    def skip_statement(self) -> None:
        while self.current() is not None:
            value = self.current()[1]
            if value == "}":
                return
            if value == ";":
                self.take()
                return
            if value == "{":
                self.skip_balanced("{", "}")
                self.accept(";")
                return
            self.take()

    def field(self, comment: str, oneof: str | None = None) -> dict | None:
        start = self.position
        modifier = None
        if self.current() and self.current()[1] in {"optional", "repeated"}:
            modifier = self.take()[1]
        try:
            field_type = self.field_type()
            if not self.current() or self.current()[0] != "identifier":
                raise ValueError("missing field name")
            name = self.take()[1]
            self.take("=")
            if not self.current() or self.current()[0] != "number":
                raise ValueError("missing field number")
            number = self.take()[1]
        except ValueError:
            self.position = start
            return None
        if self.current("["):
            self.skip_balanced("[", "]")
        self.take(";")
        return {
            "name": name,
            "type": f"{modifier} {field_type}" if modifier else field_type,
            "number": number,
            "comment": comment,
            "oneof": oneof,
        }

    def oneof(self, message: dict) -> None:
        name = self.take()[1]
        self.take("{")
        while True:
            comment = self.comments()
            if self.accept("}"):
                break
            field = self.field(comment, name)
            if field:
                message["fields"].append(field)
            else:
                self.skip_statement()

    def message(self, comment: str, parents: tuple[str, ...] = ()) -> None:
        name = self.take()[1]
        qualified = ".".join((*parents, name))
        message = {"name": qualified, "comment": comment, "fields": []}
        self.messages.append(message)
        self.take("{")
        while True:
            child_comment = self.comments()
            if self.accept("}"):
                break
            if self.accept("message"):
                self.message(child_comment, (*parents, name))
            elif self.accept("enum"):
                self.enum(child_comment, (*parents, name))
            elif self.accept("oneof"):
                self.oneof(message)
            else:
                field = self.field(child_comment)
                if field:
                    message["fields"].append(field)
                else:
                    self.skip_statement()

    def enum(self, comment: str, parents: tuple[str, ...] = ()) -> None:
        name = self.take()[1]
        enum = {
            "name": ".".join((*parents, name)),
            "comment": comment,
            "values": [],
        }
        self.enums.append(enum)
        self.take("{")
        while True:
            value_comment = self.comments()
            if self.accept("}"):
                break
            if self.current() and self.current()[0] == "identifier":
                value = self.take()[1]
                self.take("=")
                number = self.take()[1]
                if self.current("["):
                    self.skip_balanced("[", "]")
                self.take(";")
                enum["values"].append(
                    {"name": value, "number": number, "comment": value_comment}
                )
            else:
                self.skip_statement()

    def rpc(self, comment: str) -> dict:
        name = self.take()[1]
        self.take("(")
        request_stream = self.accept("stream")
        request = self.qualified_name()
        self.take(")")
        self.take("returns")
        self.take("(")
        response_stream = self.accept("stream")
        response = self.qualified_name()
        self.take(")")
        self.skip_statement()
        return {
            "name": name,
            "comment": comment,
            "request": f"stream {request}" if request_stream else request,
            "response": f"stream {response}" if response_stream else response,
        }

    def service(self, comment: str) -> None:
        name = self.take()[1]
        service = {"name": name, "comment": comment, "rpcs": []}
        self.services.append(service)
        self.take("{")
        while True:
            rpc_comment = self.comments()
            if self.accept("}"):
                break
            if self.accept("rpc"):
                service["rpcs"].append(self.rpc(rpc_comment))
            else:
                self.skip_statement()

    def parse(self) -> dict:
        while self.current() is not None:
            comment = self.comments()
            if self.accept("service"):
                self.service(comment)
            elif self.accept("message"):
                self.message(comment)
            elif self.accept("enum"):
                self.enum(comment)
            else:
                self.skip_statement()
        return {
            "services": self.services,
            "messages": self.messages,
            "enums": self.enums,
        }


def markdown_cell(value: str) -> str:
    return " ".join(value.split()).replace("|", r"\|")


def proto_reference(source: Path, destination: Path, title: str) -> None:
    document = ProtoParser(source.read_text()).parse()
    if len(document["services"]) != 1:
        raise SystemExit(f"Expected exactly one service in {source}")
    service = document["services"][0]
    lines = [f"# {title}", "", f"Service: `{service['name']}`", ""]
    for rpc in service["rpcs"]:
        lines.extend(
            [
                f"## {rpc['name']}",
                "",
                rpc["comment"] or "Public RPC exposed by envd.",
                "",
                f"- Request: `{rpc['request']}`",
                f"- Response: `{rpc['response']}`",
                "",
            ]
        )
    lines.extend(["## Message types", ""])
    for message in document["messages"]:
        lines.extend([f"### {message['name']}", ""])
        if message["comment"]:
            lines.extend([message["comment"], ""])
        if message["fields"]:
            has_oneof = any(field["oneof"] for field in message["fields"])
            if has_oneof:
                lines.extend(
                    [
                        "| Field | Type | Number | Oneof | Description |",
                        "| --- | --- | ---: | --- | --- |",
                    ]
                )
            else:
                lines.extend(
                    [
                        "| Field | Type | Number | Description |",
                        "| --- | --- | ---: | --- |",
                    ]
                )
            for field in message["fields"]:
                cells = [
                    f"`{field['name']}`",
                    f"`{field['type']}`",
                    field["number"],
                ]
                if has_oneof:
                    cells.append(f"`{field['oneof']}`" if field["oneof"] else "")
                cells.append(markdown_cell(field["comment"]))
                lines.append(f"| {' | '.join(cells)} |")
            lines.append("")
    if document["enums"]:
        lines.extend(["## Enum types", ""])
    for enum in document["enums"]:
        lines.extend([f"### {enum['name']}", ""])
        if enum["comment"]:
            lines.extend([enum["comment"], ""])
        lines.extend(
            [
                "| Value | Number | Description |",
                "| --- | ---: | --- |",
            ]
        )
        for value in enum["values"]:
            lines.append(
                f"| `{value['name']}` | {value['number']} | {markdown_cell(value['comment'])} |"
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
