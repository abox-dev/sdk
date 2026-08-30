#!/usr/bin/env python3
from __future__ import annotations

import hashlib
import json
import sys
import urllib.error
import urllib.request
from pathlib import Path


def main() -> None:
    directory = Path(sys.argv[1])
    project = sys.argv[2]
    version = sys.argv[3]

    try:
        with urllib.request.urlopen(
            f"https://pypi.org/pypi/{project}/{version}/json", timeout=30
        ) as response:
            published = {
                item["filename"]: item["digests"]["sha256"]
                for item in json.load(response)["urls"]
            }
    except urllib.error.HTTPError as error:
        if error.code != 404:
            raise
        return

    for artifact in directory.iterdir():
        if not artifact.is_file():
            continue
        digest = hashlib.sha256(artifact.read_bytes()).hexdigest()
        current = published.get(artifact.name)
        if current is not None and current != digest:
            raise RuntimeError(
                f"{project}=={version} file {artifact.name} exists with a different digest"
            )


if __name__ == "__main__":
    main()
