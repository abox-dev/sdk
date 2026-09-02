from typing import get_args

import pytest

from agentbox.connection_config import ConnectionConfig

from agentbox_code_interpreter import AsyncSandbox, RunCodeLanguage, Sandbox
from agentbox_code_interpreter.constants import DEFAULT_TEMPLATE


def make_sandbox(cls, **config_kwargs):
    # Constructing a sandbox instance makes no network requests, so URL
    # resolution can be tested without a live sandbox.
    return cls(
        sandbox_id="test-sandbox-id",
        sandbox_domain=None,
        envd_version="0.2.0",
        envd_access_token=None,
        connection_config=ConnectionConfig(**config_kwargs),
    )


def test_default_template_matches_agentbox_runtime():
    assert DEFAULT_TEMPLATE == "code-interpreter"
    assert Sandbox.default_template == DEFAULT_TEMPLATE
    assert AsyncSandbox.default_template == DEFAULT_TEMPLATE


def test_public_language_contract_matches_agentbox_runtime():
    assert get_args(RunCodeLanguage) == (
        "python",
        "javascript",
        "typescript",
        "js",
        "ts",
    )


@pytest.mark.parametrize("cls", [Sandbox, AsyncSandbox])
async def test_jupyter_url_points_to_sandbox_host_by_default(cls, monkeypatch):
    monkeypatch.delenv("AGENTBOX_SANDBOX_URL", raising=False)
    monkeypatch.delenv("AGENTBOX_DEBUG", raising=False)
    sandbox = make_sandbox(cls, domain="example.dev")
    assert sandbox._jupyter_url == "https://49999-test-sandbox-id.example.dev"


@pytest.mark.parametrize("cls", [Sandbox, AsyncSandbox])
async def test_jupyter_url_routes_through_sandbox_url_proxy(cls):
    sandbox = make_sandbox(cls, sandbox_url="https://proxy.example.com")
    assert sandbox._jupyter_url == "https://49999-test-sandbox-id.proxy.example.com"


@pytest.mark.parametrize("cls", [Sandbox, AsyncSandbox])
async def test_jupyter_url_routes_through_sandbox_url_env_var_proxy(cls, monkeypatch):
    monkeypatch.setenv("AGENTBOX_SANDBOX_URL", "https://env.example.com")
    sandbox = make_sandbox(cls)
    assert sandbox._jupyter_url == "https://49999-test-sandbox-id.env.example.com"
