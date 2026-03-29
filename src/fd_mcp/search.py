"""Filesystem search: optional `fd` subprocess or pure Python os.walk."""

from __future__ import annotations

import fnmatch
import os
import shutil
import subprocess
from pathlib import Path
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from fd_mcp.config import Settings


def search_files_walk(
    pattern_glob: str,
    root: Path,
    max_results: int,
    max_depth: int | None,
) -> tuple[list[str], bool]:
    """Glob-match filenames under *root* without invoking external tools."""
    base = root.resolve()
    matches: list[str] = []
    truncated = False

    for dirpath, dirnames, filenames in os.walk(base, topdown=True):
        dir_path = Path(dirpath).resolve()
        rel = dir_path.relative_to(base)
        depth = 0 if rel == Path(".") else len(rel.parts)
        if max_depth is not None:
            if depth > max_depth:
                dirnames.clear()
                continue
            if depth == max_depth:
                dirnames.clear()
        for fn in filenames:
            if not fnmatch.fnmatch(fn, pattern_glob):
                continue
            full = dir_path / fn
            matches.append(str(full))
            if len(matches) >= max_results:
                return matches, True
    return matches, truncated


def _search_fd(
    pattern_glob: str,
    root: Path,
    max_results: int,
    max_depth: int | None,
    fd_executable: str,
    timeout_s: float = 120.0,
) -> tuple[list[str], bool] | None:
    cmd: list[str] = [
        fd_executable,
        "--glob",
        pattern_glob,
        "--max-results",
        str(max_results),
    ]
    if max_depth is not None:
        cmd.extend(["--max-depth", str(max_depth)])
    cmd.append(str(root.resolve()))
    try:
        proc = subprocess.run(
            cmd,
            capture_output=True,
            check=False,
            timeout=timeout_s,
        )
    except subprocess.TimeoutExpired:
        return None
    if proc.returncode != 0:
        return None
    text = proc.stdout.decode(errors="replace")
    lines = [ln for ln in text.splitlines() if ln.strip()]
    truncated = len(lines) >= max_results
    return lines, truncated


def search_with_backend(
    pattern_glob: str,
    root: Path,
    max_results: int,
    max_depth: int | None,
    settings: Settings,
) -> tuple[list[str], str, bool]:
    """Prefer ``fd`` when available; otherwise use :func:`search_files_walk`."""
    depth = settings.default_max_depth if max_depth is None else max_depth
    fd_path = shutil.which("fd")
    if fd_path:
        out = _search_fd(pattern_glob, root, max_results, depth, fd_path)
        if out is not None:
            matches, truncated = out
            return matches, "fd", truncated
    matches, truncated = search_files_walk(pattern_glob, root, max_results, depth)
    return matches, "walk", truncated
