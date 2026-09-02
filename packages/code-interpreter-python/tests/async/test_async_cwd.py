import pytest

from agentbox_code_interpreter.code_interpreter_async import AsyncSandbox


@pytest.mark.skip_debug()
async def test_cwd_python(async_sandbox: AsyncSandbox):
    result = await async_sandbox.run_code("from pathlib import Path; print(Path.cwd())")
    assert "".join(result.logs.stdout).strip() == "/home/user"


@pytest.mark.skip_debug()
async def test_cwd_javascript(async_sandbox: AsyncSandbox):
    result = await async_sandbox.run_code("process.cwd()", language="js")
    assert result.text == "/home/user"


@pytest.mark.skip_debug()
async def test_cwd_typescript(async_sandbox: AsyncSandbox):
    result = await async_sandbox.run_code("process.cwd()", language="ts")
    assert result.text == "/home/user"
