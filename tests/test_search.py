"""Tests for filesystem search backends."""

from __future__ import annotations

from pathlib import Path
from unittest.mock import patch

from fd_mcp.search import search_files_walk


def test_search_files_walk_glob_and_limit(tmp_path: Path) -> None:
    (tmp_path / "a.py").write_text("x")
    (tmp_path / "sub").mkdir()
    (tmp_path / "sub" / "b.py").write_text("y")
    (tmp_path / "sub" / "c.txt").write_text("z")
    matches, truncated = search_files_walk("*.py", tmp_path, max_results=10, max_depth=None)
    assert not truncated
    assert sorted(matches) == sorted(
        [str((tmp_path / "a.py").resolve()), str((tmp_path / "sub" / "b.py").resolve())]
    )


def test_search_files_walk_respects_max_depth(tmp_path: Path) -> None:
    (tmp_path / "root.py").write_text("x")
    d1 = tmp_path / "one"
    d1.mkdir()
    (d1 / "two.py").write_text("y")
    d2 = d1 / "two"
    d2.mkdir()
    (d2 / "deep.py").write_text("z")
    matches, _trunc = search_files_walk("*.py", tmp_path, max_results=50, max_depth=1)
    rel = {Path(p).name for p in matches}
    assert "root.py" in rel
    assert "two.py" in rel
    assert "deep.py" not in rel


def test_search_files_walk_truncates(tmp_path: Path) -> None:
    for i in range(5):
        (tmp_path / f"f{i}.txt").write_text("x")
    matches, truncated = search_files_walk("*.txt", tmp_path, max_results=3, max_depth=None)
    assert len(matches) == 3
    assert truncated is True


@patch("fd_mcp.search.shutil.which", return_value="/bin/fd")
@patch("fd_mcp.search.subprocess.run")
def test_search_with_backend_prefers_fd(mock_run: object, tmp_path: Path) -> None:
    from fd_mcp.config import Settings
    from fd_mcp.search import search_with_backend

    root = tmp_path.resolve()
    cfg = Settings(
        roots=[root],
        default_max_results=100,
        absolute_max_results=10_000,
        default_max_depth=None,
    )
    mock_run.return_value.returncode = 0
    mock_run.return_value.stdout = f"{root / 'a'}\n".encode()
    matches, backend, truncated = search_with_backend(
        pattern_glob="*.py",
        root=root,
        max_results=5,
        max_depth=None,
        settings=cfg,
    )
    assert backend == "fd"
    assert truncated is False
    assert matches == [str(root / "a")]
    invoked = mock_run.call_args[0][0]
    assert invoked[0] == "/bin/fd"
    assert "--glob" in invoked
