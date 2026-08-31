#!/usr/bin/env python3
"""Remove backend-only compatibility shapes from a bundled SDK OpenAPI file."""

from __future__ import annotations

import argparse
from pathlib import Path
import yaml

from public_openapi import filter_public_openapi


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("path", type=Path)
    args = parser.parse_args()

    document = yaml.safe_load(args.path.read_text())
    filter_public_openapi(document)
    args.path.write_text(yaml.safe_dump(document, sort_keys=False))


if __name__ == "__main__":
    main()
