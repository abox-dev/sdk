#!/usr/bin/env python3
"""Remove language-specific type overrides from the derived Go envd spec."""

from pathlib import Path
import re
import sys


path = Path(sys.argv[1])
contents = path.read_text()
filtered, replacements = re.subn(r"^\s*x-go-type:\s*SecureToken\s*$\n?", "", contents, flags=re.MULTILINE)
if replacements == 0:
    raise SystemExit("SecureToken x-go-type override was not found in the envd spec")
path.write_text(filtered)
