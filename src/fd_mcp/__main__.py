"""Stdio MCP entrypoint."""

from __future__ import annotations

from fd_mcp.server import create_app


def main() -> None:
    create_app().run()


if __name__ == "__main__":
    main()
