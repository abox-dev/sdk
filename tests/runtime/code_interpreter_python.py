import asyncio

from agentbox_code_interpreter import AsyncSandbox, Sandbox


def sync_smoke() -> None:
    sandbox = Sandbox.create("code-interpreter", timeout=180)
    try:
        stdout = []
        python = sandbox.run_code(
            "print('py-ci-stream')\n{'answer': 42}",
            language="python",
            on_stdout=lambda message: stdout.append(message.line),
        )
        assert "py-ci-stream" in "".join(stdout)
        assert python.error is None

        javascript = sandbox.run_code(
            'console.log("python-sdk-js-ok"); 6 * 7', language="javascript"
        )
        assert "python-sdk-js-ok" in "".join(javascript.logs.stdout)

        typescript = sandbox.run_code(
            "const value: number = 42; console.log(value)", language="typescript"
        )
        assert "42" in "".join(typescript.logs.stdout)

        failed = sandbox.run_code('raise RuntimeError("expected-python-ci-error")')
        assert failed.error is not None and failed.error.name == "RuntimeError"

        context = sandbox.create_code_context(language="python")
        sandbox.run_code("context_value = 41", context=context)
        contextual = sandbox.run_code("context_value + 1", context=context)
        assert contextual.error is None
        assert any(item.id == context.id for item in sandbox.list_code_contexts())
        sandbox.remove_code_context(context)

        data = sandbox.run_code(
            "import pandas as pd\ndisplay(pd.DataFrame({'x': [1, 2], 'y': [3, 4]}))"
        )
        assert any(
            result.data is not None or result.html is not None
            for result in data.results
        )

        chart = sandbox.run_code(
            "import matplotlib.pyplot as plt\nplt.plot([1, 2], [3, 4])\nplt.show()"
        )
        assert any(
            result.chart is not None or result.png is not None
            for result in chart.results
        )

        sandbox.files.write("/tmp/ci-python.txt", "code-interpreter-file")
        assert sandbox.files.read("/tmp/ci-python.txt") == "code-interpreter-file"
    finally:
        sandbox.kill()


async def async_smoke() -> None:
    sandbox = await AsyncSandbox.create("code-interpreter", timeout=180)
    try:
        execution = await sandbox.run_code("print('async-ci-ok')")
        assert "async-ci-ok" in "".join(execution.logs.stdout)
        await sandbox.files.write("/tmp/async-ci.txt", "async-ci-file")
        assert await sandbox.files.read("/tmp/async-ci.txt") == "async-ci-file"
    finally:
        await sandbox.kill()


sync_smoke()
asyncio.run(async_smoke())
print("Python sync/async Code Interpreter runtime smoke passed")
