"""Path scoping: allowed roots from environment and containment checks."""

from __future__ import annotations

import os
from collections.abc import Sequence
from pathlib import Path


def parse_roots_from_env() -> list[Path]:
    """Return allowed directory roots from FD_MCP_ROOTS or the current working directory.

    Multiple entries use the OS path separator (``;`` on Windows, ``:`` on POSIX).
    """
    raw = os.environ.get("FD_MCP_ROOTS", "").strip()
    if not raw:
        return [Path.cwd().resolve()]
    parts = [p.strip() for p in raw.split(os.pathsep) if p.strip()]
    if not parts:
        return [Path.cwd().resolve()]
    return [Path(p).expanduser().resolve(strict=False) for p in parts]


def resolve_under_roots(root_arg: str, roots: Sequence[Path]) -> Path:
    """Resolve *root_arg* to an absolute path and ensure it stays within *roots*."""
    p = Path(root_arg).expanduser()
    if not p.is_absolute():
        p = (Path.cwd() / p).resolve()
    else:
        p = p.resolve()
    for rp in roots:
        rpv = rp.resolve()
        try:
            p.relative_to(rpv)
            return p
        except ValueError:
            continue
    msg = f"path {p} is not contained in allowed roots"
    raise PermissionError(msg)
