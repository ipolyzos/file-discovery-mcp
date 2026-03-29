"""Tests for path scoping and containment checks."""

from __future__ import annotations

from pathlib import Path

import pytest

from fd_mcp.paths import parse_roots_from_env, resolve_under_roots


def test_parse_roots_default_cwd(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    monkeypatch.chdir(tmp_path)
    monkeypatch.delenv("FD_MCP_ROOTS", raising=False)
    roots = parse_roots_from_env()
    assert len(roots) == 1
    assert roots[0] == Path(".").resolve()


def test_parse_roots_pathsep_split(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    a = tmp_path / "a"
    b = tmp_path / "b"
    a.mkdir()
    b.mkdir()
    sep = __import__("os").pathsep
    monkeypatch.setenv("FD_MCP_ROOTS", f"{a}{sep}{b}")
    roots = parse_roots_from_env()
    assert len(roots) == 2
    assert roots[0] == a.resolve()
    assert roots[1] == b.resolve()


def test_resolve_under_roots_accepts_inside(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    root = (tmp_path / "root").resolve()
    sub = root / "deep" / "here"
    sub.mkdir(parents=True)
    monkeypatch.chdir(root)
    resolved = resolve_under_roots("deep/here", [root])
    assert resolved == sub.resolve()


def test_resolve_under_roots_rejects_escape(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    root = (tmp_path / "root").resolve()
    root.mkdir()
    other = (tmp_path / "other").resolve()
    other.mkdir()
    monkeypatch.chdir(root)
    with pytest.raises(PermissionError, match="not contained"):
        resolve_under_roots("../other", [root])


def test_resolve_under_roots_rejects_nonexistent_outside(
    tmp_path: Path, monkeypatch: pytest.MonkeyPatch
) -> None:
    root = (tmp_path / "root").resolve()
    root.mkdir()
    monkeypatch.chdir(root)
    with pytest.raises(PermissionError, match="not contained"):
        resolve_under_roots("../nope", [root])
