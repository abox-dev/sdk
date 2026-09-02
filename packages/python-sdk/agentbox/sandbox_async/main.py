import datetime
import logging
from typing import Dict, List, Optional, Union, overload

import httpx
from packaging.version import Version
from typing_extensions import Self, Unpack

from agentbox.api.client.types import Unset
from agentbox.api.client_async import get_envd_api
from agentbox.connection_config import ApiParams, ConnectionConfig
from agentbox.envd.api import ENVD_API_HEALTH_ROUTE, ahandle_envd_api_exception
from agentbox.envd.versions import ENVD_DEBUG_FALLBACK
from agentbox.exceptions import TemplateException, format_request_timeout_error
from agentbox.sandbox.main import SandboxOpts
from agentbox.sandbox.sandbox_api import (
    SandboxIamOpts,
    SandboxLifecycle,
    SandboxMetrics,
    SandboxNetworkOpts,
    SandboxNetworkUpdate,
    SnapshotInfo,
)
from agentbox.sandbox.utils import class_method_variant
from agentbox.sandbox_async.commands.command import Commands
from agentbox.sandbox_async.commands.pty import Pty
from agentbox.sandbox_async.filesystem.filesystem import Filesystem
from agentbox.sandbox_async.sandbox_api import SandboxApi, SandboxInfo
from agentbox.sandbox_async.paginator import AsyncSnapshotPaginator

logger = logging.getLogger(__name__)


class AsyncSandbox(SandboxApi):
    """
    AgentBox cloud sandbox is a secure and isolated cloud environment.

    The sandbox allows you to:
    - Access Linux OS
    - Create, list, and delete files and directories
    - Run commands
    - Run isolated code
    - Access the internet

    See the sandbox documentation at https://docs.agentbox.ru/en/sdk/sandboxes/.

    Use the `AsyncSandbox.create()` to create a new sandbox.

    Example:
    ```python
    from agentbox import AsyncSandbox

    sandbox = await AsyncSandbox.create()
    ```
    """

    @property
    def files(self) -> Filesystem:
        """
        Module for interacting with the sandbox filesystem.
        """
        return self._filesystem

    @property
    def commands(self) -> Commands:
        """
        Module for running commands in the sandbox.
        """
        return self._commands

    @property
    def pty(self) -> Pty:
        """
        Module for interacting with the sandbox pseudo-terminal.
        """
        return self._pty

    def __init__(
        self,
        **opts: Unpack[SandboxOpts],
    ):
        """
        Initialize a connected sandbox instance.

        Applications should use `AsyncSandbox.create()` or `AsyncSandbox.connect()`.
        """
        super().__init__(**opts)

        self._envd_api = get_envd_api(self.connection_config, self.envd_api_url)
        self._filesystem = Filesystem(
            self.envd_api_url,
            self._envd_version,
            self.connection_config,
            self._envd_api,
        )
        self._commands = Commands(
            self.envd_api_url,
            self.connection_config,
            self._envd_version,
            self._envd_api,
        )
        self._pty = Pty(
            self.envd_api_url,
            self.connection_config,
            self._envd_version,
            self._envd_api,
        )

    async def is_running(self, request_timeout: Optional[float] = None) -> bool:
        """
        Check if the sandbox is running.

        :param request_timeout: Timeout for the request in **seconds**

        :return: `True` if the sandbox is running, `False` otherwise

        Example
        ```python
        sandbox = await AsyncSandbox.create()
        await sandbox.is_running() # Returns True

        await sandbox.kill()
        await sandbox.is_running() # Returns False
        ```
        """
        try:
            r = await self._envd_api.get(
                ENVD_API_HEALTH_ROUTE,
                timeout=self.connection_config.get_request_timeout(request_timeout),
            )

            if r.status_code == 502:
                return False

            err = await ahandle_envd_api_exception(r)

            if err:
                raise err

        except httpx.TimeoutException:
            raise format_request_timeout_error()

        return True

    @classmethod
    async def create(
        cls,
        template: Optional[str] = None,
        timeout: Optional[int] = None,
        metadata: Optional[Dict[str, str]] = None,
        envs: Optional[Dict[str, str]] = None,
        secure: bool = True,
        allow_internet_access: bool = True,
        network: Optional[SandboxNetworkOpts] = None,
        iam: Optional[SandboxIamOpts] = None,
        lifecycle: Optional[SandboxLifecycle] = None,
        logger: Optional[logging.Logger] = None,
        **opts: Unpack[ApiParams],
    ) -> Self:
        """
        Create a new sandbox.

        By default, the sandbox is created from the default `base` sandbox template.

        :param template: Sandbox template name or ID
        :param timeout: Timeout for the sandbox in **seconds**, default to 300 seconds. The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.
        :param metadata: Custom metadata for the sandbox
        :param envs: Custom environment variables for the sandbox
        :param secure: Envd is secured with access token and cannot be used without it, defaults to `True`.
        :param allow_internet_access: Allow sandbox to access the internet, defaults to `True`. If set to `False`, it works the same as setting network `deny_out` to `[0.0.0.0/0]`.
        :param network: Sandbox network configuration. ``allow_out``/``deny_out`` may also be a callable receiving a :class:`SandboxNetworkSelectorContext` (``ctx.all_traffic``, ``ctx.rules``) and returning a list of strings. Per-host transform rules are nested under ``network.rules``; a rule's ``transform`` may be a callable receiving a :class:`SandboxNetworkTransformContext` of placeholder strings (``ctx.iam.tokens[name]``).
        :param iam: Sandbox workload identity configuration. Each token contains ``audience`` and ``token_type``. Registered tokens are exposed to ``network.rules`` ``transform`` callables as ``ctx.iam.tokens[name]`` placeholders, which the egress proxy resolves per request
        :param lifecycle: Sandbox lifecycle configuration — ``on_timeout``: ``"kill"`` or ``"pause"`` (omitted from the request when unset, leaving the API's default, currently ``"kill"``, in effect), or an object ``{"action": "pause"|"kill", "keep_memory": bool}`` where ``keep_memory`` set to ``False`` makes a timeout auto-pause filesystem-only (cold-boots on resume; cannot be combined with ``auto_resume``); an omitted ``keep_memory`` leaves the snapshot kind to the API; ``auto_resume``: leave unset to let the API pick the behavior, set ``False`` to opt out explicitly, or ``True`` (only when ``on_timeout`` action is ``"pause"``). Example: ``{"on_timeout": {"action": "pause", "keep_memory": False}}``
        :param logger: Logger used for request and response logging for this sandbox. Accepts any standard library `logging.Logger`. When omitted, no request/response logging is emitted.

        :return: A Sandbox instance for the new sandbox

        Use this method instead of using the constructor to create a new sandbox.
        """
        if not template:
            template = cls.default_template

        sandbox = await cls._create(
            template=template,
            timeout=timeout,
            metadata=metadata,
            envs=envs,
            secure=secure,
            allow_internet_access=allow_internet_access,
            network=network,
            iam=iam,
            lifecycle=lifecycle,
            logger=logger,
            **opts,
        )

        return sandbox

    @overload
    async def connect(
        self,
        timeout: Optional[int] = None,
        **opts: Unpack[ApiParams],
    ) -> Self:
        """
        Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
        Sandbox must be either running or be paused.

        With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

        :param timeout: Timeout for the sandbox in **seconds**
            For running sandboxes, the timeout will update only if the new timeout is longer than the existing one.
        :return: A running sandbox instance

        @example
        ```python
        sandbox = await AsyncSandbox.create()
        await sandbox.pause()

        # Another code block
        same_sandbox = await sandbox.connect()
        ```
        """
        ...

    @overload
    @staticmethod
    async def connect(
        sandbox_id: str,
        timeout: Optional[int] = None,
        logger: Optional[logging.Logger] = None,
        **opts: Unpack[ApiParams],
    ) -> "AsyncSandbox":
        """
        Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
        Sandbox must be either running or be paused.

        With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

        :param sandbox_id: Sandbox ID
        :param timeout: Timeout for the sandbox in **seconds**
            For running sandboxes, the timeout will update only if the new timeout is longer than the existing one.
        :param logger: Logger used for request and response logging for this sandbox. Accepts any standard library `logging.Logger`. When omitted, no request/response logging is emitted.
        :return: A running sandbox instance

        @example
        ```python
        sandbox = await AsyncSandbox.create()
        await AsyncSandbox.pause(sandbox.sandbox_id)

        # Another code block
        same_sandbox = await AsyncSandbox.connect(sandbox.sandbox_id)
        ```
        """
        ...

    @class_method_variant("_cls_connect_sandbox")
    async def connect(
        self,
        timeout: Optional[int] = None,
        **opts: Unpack[ApiParams],
    ) -> Self:
        """
        Connect to a sandbox. If the sandbox is paused, it will be automatically resumed.
        Sandbox must be either running or be paused.

        With sandbox ID you can connect to the same sandbox from different places or environments (serverless functions, etc).

        :param timeout: Timeout for the sandbox in **seconds**
            For running sandboxes, the timeout will update only if the new timeout is longer than the existing one.
        :return: A running sandbox instance

        @example
        ```python
        sandbox = await AsyncSandbox.create()
        await sandbox.pause()

        # Another code block
        same_sandbox = await sandbox.connect()
        ```
        """
        if self.connection_config.debug:
            # Skip connecting to the sandbox in debug mode
            return self

        await SandboxApi._cls_connect(
            sandbox_id=self.sandbox_id,
            timeout=timeout,
            **self.connection_config.get_api_params(**opts),
        )

        return self

    @overload
    async def fork(
        self,
        timeout: Optional[int] = None,
        count: Optional[int] = None,
        **opts: Unpack[ApiParams],
    ) -> List[Union[Self, Exception]]:
        """
        Fork the sandbox.

        The sandbox is checkpointed in place (briefly paused, snapshotted with
        its full memory state, and resumed — its ID and expiration stay
        untouched) and `count` new sandboxes are created from that snapshot.
        All forks boot from the same snapshot, so the snapshot is captured once
        regardless of count.

        Each fork succeeds or fails independently — the returned list contains
        one entry per requested fork, either a running `AsyncSandbox` instance
        or an exception describing why that fork failed to start. Per-fork
        error codes map to the same exception classes as other API errors
        (e.g. 429 to `RateLimitException`).

        :param timeout: Timeout for the forked sandboxes in **seconds**, defaults to 300 seconds
        :param count: Number of forked sandboxes to create, defaults to 1

        :return: List with one entry per requested fork — a sandbox instance or an exception

        @example
        ```python
        sandbox = await AsyncSandbox.create()

        fork1, fork2 = await sandbox.fork(count=2)
        ```
        """
        ...

    @overload
    @staticmethod
    async def fork(
        sandbox_id: str,
        timeout: Optional[int] = None,
        count: Optional[int] = None,
        logger: Optional[logging.Logger] = None,
        **opts: Unpack[ApiParams],
    ) -> List[Union["AsyncSandbox", Exception]]:
        """
        Fork a running sandbox specified by sandbox ID.

        The sandbox is checkpointed in place (briefly paused, snapshotted with
        its full memory state, and resumed — its ID and expiration stay
        untouched) and `count` new sandboxes are created from that snapshot.
        All forks boot from the same snapshot, so the snapshot is captured once
        regardless of count.

        Each fork succeeds or fails independently — the returned list contains
        one entry per requested fork, either a running `AsyncSandbox` instance
        or an exception describing why that fork failed to start. Per-fork
        error codes map to the same exception classes as other API errors
        (e.g. 429 to `RateLimitException`).

        :param sandbox_id: Sandbox ID
        :param timeout: Timeout for the forked sandboxes in **seconds**, defaults to 300 seconds
        :param count: Number of forked sandboxes to create, defaults to 1
        :param logger: Logger used for request and response logging for the forked sandboxes. Accepts any standard library `logging.Logger`. When omitted, no request/response logging is emitted.

        :return: List with one entry per requested fork — a sandbox instance or an exception

        @example
        ```python
        sandbox = await AsyncSandbox.create()

        fork1, fork2 = await AsyncSandbox.fork(sandbox.sandbox_id, count=2)
        ```
        """
        ...

    @class_method_variant("_cls_fork_sandbox")
    async def fork(
        self,
        timeout: Optional[int] = None,
        count: Optional[int] = None,
        **opts: Unpack[ApiParams],
    ) -> List[Union[Self, Exception]]:
        """
        Fork the sandbox.

        The sandbox is checkpointed in place (briefly paused, snapshotted with
        its full memory state, and resumed — its ID and expiration stay
        untouched) and `count` new sandboxes are created from that snapshot.
        All forks boot from the same snapshot, so the snapshot is captured once
        regardless of count.

        Each fork succeeds or fails independently — the returned list contains
        one entry per requested fork, either a running `AsyncSandbox` instance
        or an exception describing why that fork failed to start. Per-fork
        error codes map to the same exception classes as other API errors
        (e.g. 429 to `RateLimitException`).

        :param timeout: Timeout for the forked sandboxes in **seconds**, defaults to 300 seconds
        :param count: Number of forked sandboxes to create, defaults to 1

        :return: List with one entry per requested fork — a sandbox instance or an exception

        @example
        ```python
        sandbox = await AsyncSandbox.create()

        fork1, fork2 = await sandbox.fork(count=2)
        ```
        """
        return await type(self)._cls_fork_sandbox(
            self.sandbox_id,
            timeout=timeout,
            count=count,
            **self.connection_config.get_api_params(**opts),
        )

    async def __aenter__(self):
        return self

    async def __aexit__(self, exc_type, exc_value, traceback):
        await self.kill()

    @overload
    async def kill(
        self,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Kill the sandbox.

        :return: `True` if the sandbox was killed, `False` if the sandbox was not found
        """
        ...

    @overload
    @staticmethod
    async def kill(
        sandbox_id: str,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Kill the sandbox specified by sandbox ID.

        :param sandbox_id: Sandbox ID

        :return: `True` if the sandbox was killed, `False` if the sandbox was not found
        """
        ...

    @class_method_variant("_cls_kill")
    async def kill(
        self,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Kill the sandbox specified by sandbox ID.

        :return: `True` if the sandbox was killed, `False` if the sandbox was not found
        """
        if self.connection_config.debug:
            # Skip killing the sandbox in debug mode
            return True

        return await SandboxApi._cls_kill(
            sandbox_id=self.sandbox_id,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    async def set_timeout(
        self,
        timeout: int,
        **opts: Unpack[ApiParams],
    ) -> None:
        """
        Set the timeout of the sandbox.
        This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to `.set_timeout`.

        The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.

        :param timeout: Timeout for the sandbox in **seconds**
        """
        ...

    @overload
    @staticmethod
    async def set_timeout(
        sandbox_id: str,
        timeout: int,
        **opts: Unpack[ApiParams],
    ) -> None:
        """
        Set the timeout of the specified sandbox.
        This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to `.set_timeout`.

        The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.

        :param sandbox_id: Sandbox ID
        :param timeout: Timeout for the sandbox in **seconds**
        """
        ...

    @class_method_variant("_cls_set_timeout")
    async def set_timeout(
        self,
        timeout: int,
        **opts: Unpack[ApiParams],
    ) -> None:
        """
        Set the timeout of the specified sandbox.
        This method can extend or reduce the sandbox timeout set when creating the sandbox or from the last call to `.set_timeout`.

        The maximum time a sandbox can be kept alive is 24 hours (86_400 seconds) for Pro users and 1 hour (3_600 seconds) for Hobby users.

        :param timeout: Timeout for the sandbox in **seconds**
        """
        await SandboxApi._cls_set_timeout(
            sandbox_id=self.sandbox_id,
            timeout=timeout,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    async def update_network(
        self,
        network: SandboxNetworkUpdate,
        **opts: Unpack[ApiParams],
    ) -> None:
        """
        Update the network configuration of the sandbox.

        Replaces the current egress configuration atomically — fields that are
        omitted are cleared on the server.

        :param network: New network configuration.
        """
        ...

    @overload
    @staticmethod
    async def update_network(
        sandbox_id: str,
        network: SandboxNetworkUpdate,
        **opts: Unpack[ApiParams],
    ) -> None:
        """
        Update the network configuration of the sandbox specified by sandbox ID.

        Replaces the current egress configuration atomically — fields that are
        omitted are cleared on the server.

        :param sandbox_id: Sandbox ID.
        :param network: New network configuration.
        """
        ...

    @class_method_variant("_cls_update_network")
    async def update_network(
        self,
        network: SandboxNetworkUpdate,
        **opts: Unpack[ApiParams],
    ) -> None:
        """
        Update the network configuration of the sandbox.

        Replaces the current egress configuration atomically — fields that are
        omitted are cleared on the server.

        :param network: New network configuration.
        """
        await SandboxApi._cls_update_network(
            sandbox_id=self.sandbox_id,
            network=network,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    async def get_info(
        self,
        **opts: Unpack[ApiParams],
    ) -> SandboxInfo:
        """
        Get sandbox information like sandbox ID, template, metadata, started at/end at date.

        :return: Sandbox info
        """
        ...

    @overload
    @staticmethod
    async def get_info(
        sandbox_id: str,
        **opts: Unpack[ApiParams],
    ) -> SandboxInfo:
        """
        Get sandbox information like sandbox ID, template, metadata, started at/end at date.
        :param sandbox_id: Sandbox ID

        :return: Sandbox info
        """
        ...

    @class_method_variant("_cls_get_info")
    async def get_info(
        self,
        **opts: Unpack[ApiParams],
    ) -> SandboxInfo:
        """
        Get sandbox information like sandbox ID, template, metadata, started at/end at date.

        :return: Sandbox info
        """

        return await SandboxApi._cls_get_info(
            sandbox_id=self.sandbox_id,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    async def get_metrics(
        self,
        start: Optional[datetime.datetime] = None,
        end: Optional[datetime.datetime] = None,
        **opts: Unpack[ApiParams],
    ) -> List[SandboxMetrics]:
        """
        Get the metrics of the current sandbox.

        :param start: Start time for the metrics, defaults to the start of the sandbox
        :param end: End time for the metrics, defaults to the current time

        :return: List of sandbox metrics containing CPU, memory and disk usage information
        """
        ...

    @overload
    @staticmethod
    async def get_metrics(
        sandbox_id: str,
        start: Optional[datetime.datetime] = None,
        end: Optional[datetime.datetime] = None,
        **opts: Unpack[ApiParams],
    ) -> List[SandboxMetrics]:
        """
        Get the metrics of the sandbox specified by sandbox ID.

        :param sandbox_id: Sandbox ID
        :param start: Start time for the metrics, defaults to the start of the sandbox
        :param end: End time for the metrics, defaults to the current time

        :return: List of sandbox metrics containing CPU, memory and disk usage information
        """
        ...

    @class_method_variant("_cls_get_metrics")
    async def get_metrics(
        self,
        start: Optional[datetime.datetime] = None,
        end: Optional[datetime.datetime] = None,
        **opts: Unpack[ApiParams],
    ) -> List[SandboxMetrics]:
        """
        Get the metrics of the current sandbox.

        :param start: Start time for the metrics, defaults to the start of the sandbox
        :param end: End time for the metrics, defaults to the current time

        :return: List of sandbox metrics containing CPU, memory and disk usage information
        """
        if self.connection_config.debug:
            # Skip getting the metrics in debug mode
            return []

        if self._envd_version < Version("0.1.5"):
            raise TemplateException(
                "You need to update the template to use the new SDK."
            )

        if self._envd_version < Version("0.2.4"):
            logger.warning(
                "Disk metrics are not supported in this version of the sandbox, please rebuild the template to get disk metrics."
            )

        return await SandboxApi._cls_get_metrics(
            sandbox_id=self.sandbox_id,
            start=start,
            end=end,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    async def pause(
        self,
        keep_memory: bool = True,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Pause the sandbox.

        :param keep_memory: When `False`, the in-memory state is dropped and only the filesystem is persisted (no memory snapshot); resuming such a sandbox cold-boots (reboots) it from disk. Defaults to `True`.

        :return: `True` if the sandbox got paused, `False` if the sandbox was already paused
        """
        ...

    @overload
    @staticmethod
    async def pause(
        sandbox_id: str,
        keep_memory: bool = True,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Pause the sandbox specified by sandbox ID.

        :param sandbox_id: Sandbox ID
        :param keep_memory: When `False`, the in-memory state is dropped and only the filesystem is persisted (no memory snapshot); resuming such a sandbox cold-boots (reboots) it from disk. Defaults to `True`.

        :return: `True` if the sandbox got paused, `False` if the sandbox was already paused
        """
        ...

    @class_method_variant("_cls_pause")
    async def pause(
        self,
        keep_memory: bool = True,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Pause the sandbox.

        :param keep_memory: When `False`, the in-memory state is dropped and only the filesystem is persisted (no memory snapshot); resuming such a sandbox cold-boots (reboots) it from disk, losing running processes and open connections. Defaults to `True` (full memory snapshot).

        :return: `True` if the sandbox got paused, `False` if the sandbox was already paused
        """

        return await SandboxApi._cls_pause(
            sandbox_id=self.sandbox_id,
            keep_memory=keep_memory,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    async def create_snapshot(
        self,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> SnapshotInfo:
        """
        Create a snapshot of the sandbox's current state.

        The sandbox will be paused while the snapshot is being created.
        The snapshot can be used to create new sandboxes with the same filesystem and state.
        Snapshots are persistent and survive sandbox deletion.

        Use the returned `snapshot_id` with `AsyncSandbox.create(snapshot_id)` to create a new sandbox from the snapshot.

        :param name: Optional name for the snapshot template. If a snapshot template with this name already exists, a new build will be assigned to the existing template instead of creating a new one.

        :return: Snapshot information including the snapshot ID and names
        """
        ...

    @overload
    @staticmethod
    async def create_snapshot(
        sandbox_id: str,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> SnapshotInfo:
        """
        Create a snapshot from the sandbox specified by sandbox ID.

        The sandbox will be paused while the snapshot is being created.

        :param sandbox_id: Sandbox ID
        :param name: Optional name for the snapshot template. If a snapshot template with this name already exists, a new build will be assigned to the existing template instead of creating a new one.

        :return: Snapshot information including the snapshot ID and names
        """
        ...

    @class_method_variant("_cls_create_snapshot")
    async def create_snapshot(
        self,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> SnapshotInfo:
        """
        Create a snapshot of the sandbox's current state.

        The sandbox will be paused while the snapshot is being created.
        The snapshot can be used to create new sandboxes with the same filesystem and state.
        Snapshots are persistent and survive sandbox deletion.

        Use the returned `snapshot_id` with `AsyncSandbox.create(snapshot_id)` to create a new sandbox from the snapshot.

        :param name: Optional name for the snapshot template. If a snapshot template with this name already exists, a new build will be assigned to the existing template instead of creating a new one.

        :return: Snapshot information including the snapshot ID and names
        """
        return await SandboxApi._cls_create_snapshot(
            sandbox_id=self.sandbox_id,
            name=name,
            **self.connection_config.get_api_params(**opts),
        )

    @overload
    def list_snapshots(
        self,
        limit: Optional[int] = None,
        next_token: Optional[str] = None,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> AsyncSnapshotPaginator:
        """
        List snapshots for this sandbox.

        :param limit: Maximum number of snapshots to return per page
        :param next_token: Token for pagination
        :param name: Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-project/my-snapshot" or "my-snapshot:v1")

        :return: Paginator for listing snapshots
        """
        ...

    @overload
    @staticmethod
    def list_snapshots(
        sandbox_id: Optional[str] = None,
        limit: Optional[int] = None,
        next_token: Optional[str] = None,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> AsyncSnapshotPaginator:
        """
        List all snapshots.

        :param sandbox_id: Filter snapshots by source sandbox ID
        :param limit: Maximum number of snapshots to return per page
        :param next_token: Token for pagination
        :param name: Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-project/my-snapshot" or "my-snapshot:v1")

        :return: Paginator for listing snapshots
        """
        ...

    @class_method_variant("_cls_list_snapshots")
    def list_snapshots(
        self,
        limit: Optional[int] = None,
        next_token: Optional[str] = None,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> AsyncSnapshotPaginator:
        """
        List snapshots for this sandbox.

        :param limit: Maximum number of snapshots to return per page
        :param next_token: Token for pagination
        :param name: Filter snapshots by name or ID, optionally tag-qualified (e.g. "my-snapshot", "my-project/my-snapshot" or "my-snapshot:v1")

        :return: Paginator for listing snapshots
        """
        return AsyncSnapshotPaginator(
            sandbox_id=self.sandbox_id,
            name=name,
            limit=limit,
            next_token=next_token,
            **self.connection_config.get_api_params(**opts),
        )

    @classmethod
    def _cls_list_snapshots(
        cls,
        sandbox_id: Optional[str] = None,
        limit: Optional[int] = None,
        next_token: Optional[str] = None,
        name: Optional[str] = None,
        **opts: Unpack[ApiParams],
    ) -> AsyncSnapshotPaginator:
        return AsyncSnapshotPaginator(
            sandbox_id=sandbox_id,
            name=name,
            limit=limit,
            next_token=next_token,
            **cls._resolve_api_params(**opts),
        )

    @classmethod
    async def delete_snapshot(
        cls,
        snapshot_id: str,
        **opts: Unpack[ApiParams],
    ) -> bool:
        """
        Delete a snapshot.

        :param snapshot_id: Snapshot ID
        :return: `True` if the snapshot was deleted, `False` if it was not found
        """
        return await cls._cls_delete_snapshot(
            snapshot_id=snapshot_id,
            **opts,
        )

    @classmethod
    async def _cls_connect_sandbox(
        cls,
        sandbox_id: str,
        timeout: Optional[int] = None,
        logger: Optional[logging.Logger] = None,
        **opts: Unpack[ApiParams],
    ) -> Self:
        params = cls._resolve_api_params(**opts)
        debug = ConnectionConfig(**params).debug
        if debug:
            sandbox_domain = None
            envd_version = ENVD_DEBUG_FALLBACK
            envd_access_token = None
            traffic_access_token = None
        else:
            sandbox = await SandboxApi._cls_connect(
                sandbox_id=sandbox_id,
                timeout=timeout,
                logger=logger,
                **params,
            )

            sandbox_id = sandbox.sandbox_id
            sandbox_domain = sandbox.sandbox_domain
            envd_version = Version(sandbox.envd_version)
            envd_access_token = sandbox.envd_access_token
            traffic_access_token = sandbox.traffic_access_token

        sandbox_headers = {
            "Agentbox-Sandbox-Id": sandbox_id,
            "Agentbox-Sandbox-Port": str(ConnectionConfig.envd_port),
        }
        if envd_access_token is not None and not isinstance(envd_access_token, Unset):
            sandbox_headers["X-Access-Token"] = envd_access_token
        if traffic_access_token:
            sandbox_headers["Agentbox-Traffic-Access-Token"] = traffic_access_token

        connection_config = ConnectionConfig(
            extra_sandbox_headers=sandbox_headers,
            logger=logger,
            **params,
        )

        return cls(
            sandbox_id=sandbox_id,
            sandbox_domain=sandbox_domain,
            envd_version=envd_version,
            envd_access_token=envd_access_token,
            traffic_access_token=traffic_access_token,
            connection_config=connection_config,
        )

    @classmethod
    async def _cls_fork_sandbox(
        cls,
        sandbox_id: str,
        timeout: Optional[int] = None,
        count: Optional[int] = None,
        logger: Optional[logging.Logger] = None,
        **opts: Unpack[ApiParams],
    ) -> List[Union[Self, Exception]]:
        params = cls._resolve_api_params(**opts)
        responses = await SandboxApi._cls_fork(
            sandbox_id=sandbox_id,
            timeout=timeout,
            count=count,
            logger=logger,
            **params,
        )

        sandboxes: List[Union[Self, Exception]] = []
        for response in responses:
            if isinstance(response, Exception):
                sandboxes.append(response)
                continue

            sandbox_headers = {
                "Agentbox-Sandbox-Id": response.sandbox_id,
                "Agentbox-Sandbox-Port": str(ConnectionConfig.envd_port),
            }
            if response.envd_access_token is not None:
                sandbox_headers["X-Access-Token"] = response.envd_access_token
            if response.traffic_access_token:
                sandbox_headers["Agentbox-Traffic-Access-Token"] = (
                    response.traffic_access_token
                )

            connection_config = ConnectionConfig(
                extra_sandbox_headers=sandbox_headers,
                logger=logger,
                **params,
            )

            sandboxes.append(
                cls(
                    sandbox_id=response.sandbox_id,
                    sandbox_domain=response.sandbox_domain,
                    envd_version=Version(response.envd_version),
                    envd_access_token=response.envd_access_token,
                    traffic_access_token=response.traffic_access_token,
                    connection_config=connection_config,
                )
            )

        return sandboxes

    @classmethod
    async def _create(
        cls,
        template: Optional[str],
        timeout: Optional[int],
        metadata: Optional[Dict[str, str]],
        envs: Optional[Dict[str, str]],
        secure: bool,
        allow_internet_access: bool,
        network: Optional[SandboxNetworkOpts] = None,
        iam: Optional[SandboxIamOpts] = None,
        lifecycle: Optional[SandboxLifecycle] = None,
        logger: Optional[logging.Logger] = None,
        **opts: Unpack[ApiParams],
    ) -> Self:
        params = cls._resolve_api_params(**opts)
        extra_sandbox_headers = {}

        debug = ConnectionConfig(**params).debug
        if debug:
            sandbox_id = "debug_sandbox_id"
            sandbox_domain = None
            envd_version = ENVD_DEBUG_FALLBACK
            envd_access_token = None
            traffic_access_token = None
        else:
            response = await SandboxApi._create_sandbox(
                template=template or cls.default_template,
                timeout=timeout or cls.default_sandbox_timeout,
                metadata=metadata,
                env_vars=envs,
                secure=secure,
                allow_internet_access=allow_internet_access,
                network=network,
                iam=iam,
                lifecycle=lifecycle,
                logger=logger,
                **params,
            )

            sandbox_id = response.sandbox_id
            sandbox_domain = response.sandbox_domain
            envd_version = Version(response.envd_version)
            envd_access_token = response.envd_access_token
            traffic_access_token = response.traffic_access_token

            if envd_access_token is not None and not isinstance(
                envd_access_token, Unset
            ):
                extra_sandbox_headers["X-Access-Token"] = envd_access_token

        extra_sandbox_headers["Agentbox-Sandbox-Id"] = sandbox_id
        extra_sandbox_headers["Agentbox-Sandbox-Port"] = str(ConnectionConfig.envd_port)
        if traffic_access_token:
            extra_sandbox_headers["Agentbox-Traffic-Access-Token"] = (
                traffic_access_token
            )

        connection_config = ConnectionConfig(
            extra_sandbox_headers=extra_sandbox_headers,
            logger=logger,
            **params,
        )

        return cls(
            sandbox_id=sandbox_id,
            sandbox_domain=sandbox_domain,
            envd_version=envd_version,
            envd_access_token=envd_access_token,
            traffic_access_token=traffic_access_token,
            connection_config=connection_config,
        )
