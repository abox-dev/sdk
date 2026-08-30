import pytest

from agentbox import ConnectionConfig
from agentbox.api import ApiClient
from agentbox.exceptions import AuthenticationException


def test_api_client_requires_api_key(monkeypatch):
    monkeypatch.delenv("AGENTBOX_API_KEY", raising=False)
    config = ConnectionConfig()
    with pytest.raises(AuthenticationException, match=r"API key is required"):
        ApiClient(config)
