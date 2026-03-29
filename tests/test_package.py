"""Smoke tests for package metadata and entry behavior."""

from __future__ import annotations

import fd_mcp


def test_version_matches_pyproject_pattern() -> None:
    assert fd_mcp.__version__
    parts = fd_mcp.__version__.split(".")
    assert len(parts) >= 2
    for p in parts:
        assert p.isdigit()


def test_main_and_create_app_are_wired() -> None:
    from fd_mcp.__main__ import main
    from fd_mcp.server import create_app

    assert callable(main)
    assert create_app().name == "fd-mcp"
