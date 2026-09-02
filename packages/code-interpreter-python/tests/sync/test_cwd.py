import pytest

from agentbox_code_interpreter.code_interpreter_sync import Sandbox


@pytest.mark.skip_debug()
def test_cwd_python(sandbox: Sandbox):
    result = sandbox.run_code("from pathlib import Path; print(Path.cwd())")
    assert "".join(result.logs.stdout).strip() == "/home/user"


@pytest.mark.skip_debug()
def test_cwd_javascript(sandbox: Sandbox):
    result = sandbox.run_code("process.cwd()", language="js")
    assert result.text == "/home/user"


@pytest.mark.skip_debug()
def test_cwd_typescript(sandbox: Sandbox):
    result = sandbox.run_code("process.cwd()", language="ts")
    assert result.text == "/home/user"
