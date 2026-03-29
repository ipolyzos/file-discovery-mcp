"""Runtime configuration from environment variables."""

from __future__ import annotations

import os
from dataclasses import dataclass
from pathlib import Path

from fd_mcp.paths import parse_roots_from_env


def _env_int(name: str, default: int) -> int:
    raw = os.environ.get(name, "").strip()
    if not raw:
        return default
    try:
        return int(raw)
    except ValueError as e:
        msg = f"{name} must be an integer, got {raw!r}"
        raise ValueError(msg) from e


@dataclass(frozen=True)
class Settings:
    """Search limits and allowed filesystem roots."""

    roots: tuple[Path, ...]
    default_max_results: int
    absolute_max_results: int
    default_max_depth: int | None


def load_settings() -> Settings:
    roots_tuple = tuple(parse_roots_from_env())
    default_max = max(1, _env_int("FD_MCP_MAX_RESULTS", 1000))
    absolute_max = max(default_max, _env_int("FD_MCP_ABSOLUTE_MAX", 50_000))
    depth_raw = os.environ.get("FD_MCP_MAX_DEPTH", "").strip()
    max_depth: int | None
    if depth_raw:
        try:
            max_depth = int(depth_raw)
        except ValueError as e:
            msg = f"FD_MCP_MAX_DEPTH must be an integer, got {depth_raw!r}"
            raise ValueError(msg) from e
        if max_depth < 0:
            msg = "FD_MCP_MAX_DEPTH must be >= 0"
            raise ValueError(msg)
    else:
        max_depth = None
    return Settings(
        roots=roots_tuple,
        default_max_results=default_max,
        absolute_max_results=absolute_max,
        default_max_depth=max_depth,
    )
