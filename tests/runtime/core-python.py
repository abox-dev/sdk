import asyncio
import time

from agentbox import (
    AgentBox,
    AsyncSandbox,
    CommandExitException,
    Sandbox,
    SandboxQuery,
    TimeoutException,
)


MARKER = f"sdk-runtime-python-{time.time_ns()}"


def sync_smoke() -> None:
    client = AgentBox()
    sandbox = client.Sandbox.create(
        "base",
        metadata={"sdkRuntime": MARKER},
        timeout=180,
    )
    fork = None
    try:
        assert sandbox.get_info().sandbox_id == sandbox.sandbox_id
        paginator = Sandbox.list(query=SandboxQuery(metadata={"sdkRuntime": MARKER}))
        assert any(
            item.sandbox_id == sandbox.sandbox_id for item in paginator.next_items()
        )

        stdout: list[str] = []
        stderr: list[str] = []
        result = sandbox.commands.run(
            "printf 'python-out'; printf 'python-err' >&2",
            on_stdout=stdout.append,
            on_stderr=stderr.append,
        )
        assert result.stdout == "python-out" == "".join(stdout)
        assert result.stderr == "python-err" == "".join(stderr)

        try:
            sandbox.commands.run("sleep 2", timeout=0.1)
            raise AssertionError("command timeout was not enforced")
        except TimeoutException:
            pass

        try:
            sandbox.commands.run("printf 'expected-error' >&2; exit 7")
            raise AssertionError("non-zero command unexpectedly succeeded")
        except CommandExitException as error:
            assert error.exit_code == 7

        sandbox.files.write("/tmp/agentbox/text.txt", "hello-python")
        sandbox.files.write("/tmp/agentbox/data.bin", bytes([0, 1, 2, 127, 255]))
        assert sandbox.files.read("/tmp/agentbox/text.txt") == "hello-python"
        assert sandbox.files.read("/tmp/agentbox/data.bin", format="bytes") == bytes(
            [0, 1, 2, 127, 255]
        )
        assert any(
            entry.name == "data.bin" for entry in sandbox.files.list("/tmp/agentbox")
        )
        sandbox.files.remove("/tmp/agentbox/data.bin")
        assert not sandbox.files.exists("/tmp/agentbox/data.bin")

        sandbox.set_timeout(180)
        sandbox_id = sandbox.sandbox_id
        assert sandbox.pause()
        sandbox = Sandbox.connect(sandbox_id, timeout=180)
        assert sandbox.get_info().state.value == "running"

        fork_result = sandbox.fork(count=1, timeout=180)[0]
        if isinstance(fork_result, Exception):
            raise fork_result
        fork = fork_result
        assert fork.files.read("/tmp/agentbox/text.txt") == "hello-python"
    finally:
        if fork is not None:
            fork.kill()
        sandbox.kill()


async def async_smoke() -> None:
    sandbox = await AsyncSandbox.create(
        "base",
        metadata={"sdkRuntime": f"{MARKER}-async"},
        timeout=180,
    )
    try:
        info = await sandbox.get_info()
        assert info.sandbox_id == sandbox.sandbox_id
        connected = await AsyncSandbox.connect(sandbox.sandbox_id, timeout=180)
        result = await connected.commands.run("printf 'async-agentbox'")
        assert result.stdout == "async-agentbox"
        await connected.files.write("/tmp/async.txt", "async-file")
        assert await connected.files.read("/tmp/async.txt") == "async-file"
        paginator = AsyncSandbox.list(
            query=SandboxQuery(metadata={"sdkRuntime": f"{MARKER}-async"})
        )
        assert any(
            item.sandbox_id == sandbox.sandbox_id
            for item in await paginator.next_items()
        )
    finally:
        await sandbox.kill()


sync_smoke()
asyncio.run(async_smoke())
print("Python sync/async core runtime smoke passed")
