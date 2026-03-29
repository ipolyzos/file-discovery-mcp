"""FastMCP application: read-only search tool over scoped roots."""

from __future__ import annotations

from typing import Any

from mcp.server.fastmcp import FastMCP

from fd_mcp.config import Settings, load_settings
from fd_mcp.paths import resolve_under_roots
from fd_mcp.search import search_with_backend


def create_app(settings: Settings | None = None) -> FastMCP:
    cfg = settings or load_settings()
    instructions = (
        "Read-only filesystem discovery. All searches are confined to FD_MCP_ROOTS "
        "(default: process working directory). Requested paths must stay within those roots. "
        "Results are capped by FD_MCP_MAX_RESULTS and FD_MCP_ABSOLUTE_MAX."
    )
    mcp = FastMCP("fd-mcp", instructions=instructions)

    @mcp.tool()
    async def search_files(
        pattern: str,
        root: str,
        max_results: int | None = None,
        max_depth: int | None = None,
    ) -> dict[str, Any]:
        """Search for files matching a glob pattern under an allowed directory *root*.

        Uses ``fd`` when installed; otherwise scans with a pure Python walker.
        """
        if not pattern or not pattern.strip():
            return {"ok": False, "code": "INVALID_PATTERN", "error": "pattern must be non-empty"}
        try:
            root_path = resolve_under_roots(root, cfg.roots)
        except PermissionError as e:
            return {"ok": False, "code": "PATH_NOT_ALLOWED", "error": str(e)}
        limit = max_results if max_results is not None else cfg.default_max_results
        limit = max(1, min(limit, cfg.absolute_max_results))
        if max_depth is not None and max_depth < 0:
            return {"ok": False, "code": "INVALID_DEPTH", "error": "max_depth must be >= 0"}
        try:
            matches, backend, truncated = search_with_backend(
                pattern.strip(),
                root_path,
                limit,
                max_depth,
                cfg,
            )
        except OSError as e:
            return {"ok": False, "code": "SEARCH_FAILED", "error": str(e)}
        return {
            "ok": True,
            "matches": matches,
            "truncated": truncated,
            "backend": backend,
            "root": str(root_path),
            "pattern": pattern.strip(),
            "max_results": limit,
        }

    return mcp
