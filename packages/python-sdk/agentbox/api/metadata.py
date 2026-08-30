import platform

from importlib import metadata

package_version = metadata.version("abox-sdk")

default_headers = {
    "lang": "python",
    "lang_version": platform.python_version(),
    "package_version": package_version,
    "publisher": "agentbox",
    "sdk_runtime": "python",
    "system": platform.system(),
}
