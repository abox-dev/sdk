"""Numeric AGENTBOX_* env vars are parsed at import time; an empty-string value
(e.g. `AGENTBOX_KEEPALIVE_EXPIRY=` in a dotenv file) must fall back to the default
instead of raising ValueError when the module is imported."""

import importlib

import agentbox.api

_ENV_VARS = (
    "AGENTBOX_KEEPALIVE_EXPIRY",
    "AGENTBOX_MAX_KEEPALIVE_CONNECTIONS",
    "AGENTBOX_CONNECTION_RETRIES",
)


def _reload():
    return importlib.reload(agentbox.api)


def test_empty_env_vars_fall_back_to_defaults(monkeypatch):
    for var in _ENV_VARS:
        monkeypatch.setenv(var, "")
    try:
        api = _reload()
        assert api.pool_idle_timeout == 300
        assert api.pool_max_idle_per_host == 20
        assert api.connection_retries == 3
    finally:
        monkeypatch.undo()
        _reload()


def test_set_env_vars_are_honored(monkeypatch):
    monkeypatch.setenv("AGENTBOX_KEEPALIVE_EXPIRY", "42")
    monkeypatch.setenv("AGENTBOX_CONNECTION_RETRIES", "5")
    try:
        api = _reload()
        assert api.pool_idle_timeout == 42
        assert api.connection_retries == 5
    finally:
        monkeypatch.undo()
        _reload()
