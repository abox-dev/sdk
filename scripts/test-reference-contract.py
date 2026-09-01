#!/usr/bin/env python3
"""Fast assertions for the generated public documentation contract."""

import hashlib
import json
import re
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
    "AdminApiKeyAuth",
    "AdminTeamAuth",
    "AuthProviderBearerAuth",
    "AuthProviderTeamAuth",
}
FORBIDDEN_COMPONENTS = {
    "AdminSandboxKillResult",
    "AdminBuildCancelResult",
    "TeamAPIKey",
    "CreatedTeamAPIKey",
    "NewTeamAPIKey",
    "UpdateTeamAPIKey",
    "apiKeyID",
    "accessTokenID",
    "volumeID",
    "VolumeMount",
}


def markdown_section(document: str, heading: str) -> str:
    marker = f"### {heading}\n"
    start = document.index(marker) + len(marker)
    boundaries = [
        position
        for prefix in ("\n### ", "\n## ")
        if (position := document.find(prefix, start)) >= 0
    ]
    end = min(boundaries) if boundaries else len(document)
    return document[start:end]


def markdown_table_names(section: str) -> list[str]:
    return re.findall(r"^\| `([^`]+)` \|", section, re.M)


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


def component_references(value) -> set[tuple[str, str]]:
    references = set()
    if isinstance(value, dict):
        reference = value.get("$ref")
        if isinstance(reference, str) and reference.startswith("#/components/"):
            section, name, *_ = reference.removeprefix("#/components/").split("/")
            references.add((section, name))
        for child in value.values():
            references.update(component_references(child))
    elif isinstance(value, list):
        for child in value:
            references.update(component_references(child))
    return references


def assert_only_reachable_components(document: dict) -> None:
    components = document.get("components", {})
    roots = {key: value for key, value in document.items() if key != "components"}
    pending = list(component_references(roots))
    reachable = set()
    while pending:
        component = pending.pop()
        if component in reachable:
            continue
        reachable.add(component)
        section, name = component
        pending.extend(component_references(components[section][name]) - reachable)
    published = {
        (section, name)
        for section, definitions in components.items()
        if section != "securitySchemes"
        for name in definitions
    }
    assert published == reachable


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
        spec_name = name.removeprefix("openapi/").removesuffix(".yml")
        public_documents[spec_name] = document
        expected_scheme = (
            "ApiKeyAuth" if spec_name == "control-plane" else "AccessTokenAuth"
        )
        expected_header = (
            "X-API-Key" if spec_name == "control-plane" else "X-Access-Token"
        )
        assert set(document.get("components", {}).get("securitySchemes", {})) == {
            expected_scheme
        }
        assert document["components"]["securitySchemes"][expected_scheme] == {
            "type": "apiKey",
            "in": "header",
            "name": expected_header,
        }
        if spec_name == "envd":
            assert document["servers"] == [
                {
                    "url": "https://sandbox.agentbox-runtime.ru",
                    "description": (
                        "AgentBox sandbox proxy. Routing headers are required."
                    ),
                }
            ]
            for path_item in document["paths"].values():
                for method, operation in path_item.items():
                    if method not in {"get", "post", "put", "patch", "delete"}:
                        continue
                    headers = {
                        parameter["name"]: parameter
                        for parameter in (
                            resolve_local_ref(document, value)
                            for value in operation.get("parameters", [])
                        )
                        if parameter.get("in") == "header"
                    }
                    assert headers["Agentbox-Sandbox-Id"]["required"] is True
                    assert headers["Agentbox-Sandbox-Port"]["required"] is True
                    assert (
                        headers["Agentbox-Sandbox-Port"]["schema"]["default"] == 49983
                    )
        assert_public_schema(document)
        assert_local_refs_resolve(document)
        assert_only_reachable_components(document)
        assert not any(
            name in FORBIDDEN_COMPONENTS
            for definitions in document.get("components", {}).values()
            for name in definitions
        )

    for operation in operations:
        document = public_documents[operation["spec"]]
        rendered_operation = document["paths"][operation["path"]][
            operation["method"].lower()
        ]
        auth = operation["auth"]
        if auth is None:
            assert rendered_operation["security"] == []
            continue
        scheme = (
            "ApiKeyAuth" if operation["spec"] == "control-plane" else "AccessTokenAuth"
        )
        requirement = {scheme: []}
        expected = [requirement] if auth["required"] else [{}, requirement]
        assert rendered_operation["security"] == expected

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

    process_reference = (REFERENCE / "connect/process.md").read_text()
    assert markdown_table_names(markdown_section(process_reference, "PTY")) == ["size"]
    assert markdown_table_names(
        markdown_section(process_reference, "ProcessEvent")
    ) == ["start", "data", "end", "keepalive"]
    assert markdown_table_names(
        markdown_section(process_reference, "StreamInputRequest")
    ) == ["start", "data", "keepalive"]
    assert "| `envs` | `map<string, string>` | 3 |" in markdown_section(
        process_reference, "ProcessConfig"
    )
    assert "### PTY.Size" in process_reference
    assert "### ProcessEvent.StartEvent" in process_reference
    for value in (
        "SIGNAL_UNSPECIFIED",
        "SIGNAL_SIGTERM",
        "SIGNAL_SIGKILL",
    ):
        assert value in markdown_section(process_reference, "Signal")

    filesystem_reference = (REFERENCE / "connect/filesystem.md").read_text()
    assert "| `metadata` | `map<string, string>` | 11 |" in markdown_section(
        filesystem_reference, "EntryInfo"
    )
    for enum, values in {
        "FileType": (
            "FILE_TYPE_UNSPECIFIED",
            "FILE_TYPE_FILE",
            "FILE_TYPE_DIRECTORY",
            "FILE_TYPE_SYMLINK",
        ),
        "EventType": (
            "EVENT_TYPE_UNSPECIFIED",
            "EVENT_TYPE_CREATE",
            "EVENT_TYPE_WRITE",
            "EVENT_TYPE_REMOVE",
            "EVENT_TYPE_RENAME",
            "EVENT_TYPE_CHMOD",
        ),
    }.items():
        section = markdown_section(filesystem_reference, enum)
        for value in values:
            assert value in section

    javascript_files = {
        str(path.relative_to(REFERENCE))
        for path in (REFERENCE / "sdk/javascript").rglob("*.md")
    }
    assert not any("mapEntryInfo" in path for path in javascript_files)
    assert not any("formatSandboxTimeoutError" in path for path in javascript_files)
    assert not any("deserializeChart" in path for path in javascript_files)
    assert not any("parseOutput" in path for path in javascript_files)
    assert not any("extractError" in path for path in javascript_files)
    javascript_markdown = "\n".join(
        path.read_text() for path in (REFERENCE / "sdk/javascript").rglob("*.md")
    )
    assert "**`Internal`**" not in javascript_markdown
    assert "envdAccessToken" not in javascript_markdown
    assert "protected static createSandbox" not in javascript_markdown

    list_sandboxes = (REFERENCE / "openapi/markdown/listSandboxes.md").read_text()
    assert "Metadata query used to filter the sandboxes" in list_sandboxes
    assert "### 200" in list_sandboxes
    assert "Schema: `array<ListedSandbox>`" in list_sandboxes
    assert "- **`sandboxID`** · `string` · required" in list_sandboxes
    assert "- **`metadata`** · `string` · query · optional" in list_sandboxes
    for expected in (
        "Allowed values for `SandboxState`: `running` | `paused`",
        "#### Response headers",
        "- **`X-Next-Token`** · `string` · response header",
        "- **`X-Total-Running`** · `integer` · response header",
        "Format: `int32`",
        "Default: `100`",
        "Minimum: `1`",
        "Maximum: `100`",
    ):
        assert expected in list_sandboxes

    sandbox_logs = (REFERENCE / "openapi/markdown/getSandboxLogs.md").read_text()
    assert "Allowed values for `LogsDirection`: `forward` | `backward`" in sandbox_logs

    upload_file = (REFERENCE / "openapi/markdown/uploadFile.md").read_text()
    assert "### application/octet-stream" in upload_file
    assert "Raw file content. The 'path' query parameter is required" in upload_file
    assert "Format: `binary`" in upload_file

    start_build = (REFERENCE / "openapi/markdown/startTemplateBuild.md").read_text()
    for expected in (
        "Variant `AWSRegistry`",
        "discriminator `type` = `aws`",
        "fromImageRegistry<AWSRegistry>.awsAccessKeyId",
        "fromImageRegistry<AWSRegistry>.awsSecretAccessKey",
        "fromImageRegistry<AWSRegistry>.awsRegion",
        "Variant `GCPRegistry`",
        "discriminator `type` = `gcp`",
        "fromImageRegistry<GCPRegistry>.serviceAccountJson",
        "Variant `GeneralRegistry`",
        "discriminator `type` = `registry`",
        "fromImageRegistry<GeneralRegistry>.username",
        "fromImageRegistry<GeneralRegistry>.password",
    ):
        assert expected in start_build

    list_metrics = (REFERENCE / "openapi/markdown/listSandboxMetrics.md").read_text()
    assert "- **`sandbox_ids`** · `array<string>` · query · required" in list_metrics
    assert "Maximum items: `100`" in list_metrics
    assert "Unique items: `yes`" in list_metrics

    list_templates = (REFERENCE / "openapi/markdown/listTemplates.md").read_text()
    assert "- **`createdBy`** · `TeamUser | null` · required" in list_templates
    assert "- **`createdBy.id`** · `string` · required" in list_templates
    assert "Format: `uuid`" in list_templates
    assert "- **`lastSpawnedAt`** · `string | null` · required" in list_templates

    create_sandbox = (REFERENCE / "openapi/markdown/createSandbox.md").read_text()
    assert "Schema: `NewSandbox`" in create_sandbox
    for field in (
        "autoResume.enabled",
        "network.allowPublicTraffic",
        "network.maskRequestHost",
        "iam.tokens",
        "iam.tokens.*.audience",
    ):
        assert f"- **`{field}`** ·" in create_sandbox


if __name__ == "__main__":
    main()
