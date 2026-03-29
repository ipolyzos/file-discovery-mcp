"""Smoke tests for package metadata and entry behavior."""

from __future__ import annotations

import subprocess
import sys

import fd_mcp


def test_version_matches_pyproject_pattern() -> None:
    assert fd_mcp.__version__
    parts = fd_mcp.__version__.split(".")
    assert len(parts) >= 2
    for p in parts:
        assert p.isdigit()


def test_main_exits_with_code_2() -> None:
    proc = subprocess.run(
        [sys.executable, "-m", "fd_mcp"],
        check=False,
        capture_output=True,
        text=True,
    )
    assert proc.returncode == 2
    assert "not yet implemented" in proc.stderr
