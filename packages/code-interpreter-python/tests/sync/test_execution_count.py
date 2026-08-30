import pytest

from agentbox_code_interpreter.code_interpreter_sync import Sandbox


@pytest.mark.skip_debug()
def test_execution_count(sandbox: Sandbox):
    sandbox.run_code("echo 'AgentBox is awesome!'")
    result = sandbox.run_code("!pwd")
    assert result.execution_count == 2
