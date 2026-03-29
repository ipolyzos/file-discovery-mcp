"""Tests for MCP application wiring (no live stdio transport)."""

from __future__ import annotations

import asyncio
from pathlib import Path

from fd_mcp.config import Settings
from fd_mcp.server import create_app


def test_create_app_registers_tools() -> None:
    app = create_app()
    assert app.name == "fd-mcp"


def test_search_files_tool_json_payload(tmp_path: Path) -> None:
    root = tmp_path / "proj"
    root.mkdir()
    (root / "a.py").write_text("x")
    cfg = Settings(
        roots=(root.resolve(),),
        default_max_results=50,
        absolute_max_results=500,
        default_max_depth=None,
    )
    app = create_app(cfg)

    async def _run() -> dict:
        _blocks, data = await app.call_tool(
            "search_files",
            {"pattern": "*.py", "root": str(root)},
        )
        return data

    data = asyncio.run(_run())
    assert data["ok"] is True
    assert any(p.endswith("a.py") for p in data["matches"])
