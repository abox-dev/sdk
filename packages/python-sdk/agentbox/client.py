from typing import Dict, Type, TypeVar, cast

from typing_extensions import Unpack

from agentbox.connection_config import ApiParams
from agentbox.sandbox_async.main import AsyncSandbox
from agentbox.sandbox_sync.main import Sandbox
from agentbox.template_async.main import AsyncTemplate
from agentbox.template_sync.main import Template

T = TypeVar("T")


class AgentBoxClientParams(ApiParams, total=False):
    """Params bound to an :class:`AgentBox` client, used as the defaults for every
    call made through its resource classes. Same shape as :class:`ApiParams`."""


def _bind(cls: Type[T], api_params: ApiParams) -> Type[T]:
    """Generate a subclass of ``cls`` carrying ``api_params`` as its bound params."""
    return cast(
        Type[T],
        type(cls.__name__, (cls,), {"_bound_api_params": api_params}),
    )


class AgentBox:
    """
    AgentBox client with an explicitly bound connection configuration.

    The resource classes exposed by the client (`Sandbox`, `AsyncSandbox`,
    `Template`, and `AsyncTemplate`) behave exactly like the top-level exports of the same name,
    except the params passed to the client are used as the defaults instead of
    the environment variables.
    Per-call params still take precedence over the client's params.

    Multiple clients are fully isolated from each other and from the top-level
    env-configured exports.

    Example:
    ```python
    from agentbox import AgentBox

    client = AgentBox(api_key="ab_...")

    sandbox = client.Sandbox.create()
    exists = client.Template.exists("my-template")
    ```
    """

    def __init__(self, **opts: Unpack[AgentBoxClientParams]):
        """
        Create a new client with the API params bound to it.

        :param opts: API params used as the defaults for every call made
            through this client's resource classes.
        """
        # Params are copied so later mutations of the caller's dicts cannot
        # change the bound configuration.
        api_params = cast(ApiParams, dict(cast(Dict[str, object], opts)))

        self.Sandbox = _bind(Sandbox, api_params)
        """`Sandbox` class bound to this client's connection configuration."""

        self.AsyncSandbox = _bind(AsyncSandbox, api_params)
        """`AsyncSandbox` class bound to this client's connection configuration."""

        self.Template = _bind(Template, api_params)
        """`Template` class bound to this client's connection configuration."""

        self.AsyncTemplate = _bind(AsyncTemplate, api_params)
        """`AsyncTemplate` class bound to this client's connection configuration."""
