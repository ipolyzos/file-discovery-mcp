# Troubleshooting

**Purpose:** Fast checks for typical **file-discovery-mcp** issues. For tool JSON and error codes, see the [client contract](client-contract.md) and [examples](examples.md).

---

## `fd` not installed

**Symptom:** Searches still work; successful responses show `"backend": "walk"` instead of `"fd"`.

**What to know:** The server **does not require** `fd`. It looks up `fd` on **`PATH`** at search time; if missing, it uses an in-process directory walk. That is expected, not a failure.

**If you want `fd`:** Install it and ensure the MCP host runs `file-discovery-mcp` with a **`PATH`** that includes the `fd` binary (some GUIs strip or override env—set `PATH` in the server’s `env` block if needed).

**If you see `SEARCH_FAILED` with `fd` installed:** The subprocess exited non-zero. The `error` string may include **`fd` stderr**. Fix `fd` flags compatibility, permissions on the `fd` binary, or a broken `PATH` entry pointing at a bad shim.

---

## Permissions issues

**Symptom:** `SEARCH_FAILED`, partial results, or opaque walk errors when traversing the tree.

**Checks:**

1. **`root` readable:** The server user must be able to traverse `root` and list directories along the path to matches.
2. **Symlinks / mounts:** Unreadable mount points or loops can confuse walks; try a `root` inside a fully readable subtree.
3. **`PATH_NOT_ALLOWED`:** Resolved `root` is outside **`FILE_DISCOVERY_MCP_ROOTS`** (or outside the startup cwd when that env is unset). Widen allowlist or fix `root`—see [using](using.md).

**Note:** The server never escalates privileges; fix OS-level permissions or choose a different `root`.

---

## Empty `matches`

**Symptom:** `"ok": true` but `"matches": []` (or missing).

**Checks (in order):**

| Check | Action |
|--------|--------|
| Wrong tree | Confirm `root` points at the directory that actually contains files (try an absolute path). |
| Pattern too narrow | Remember: **basename glob only** (e.g. `*.go` works; path-style patterns do not). |
| **`max_depth`** | A small depth excludes deeper files; omit `max_depth` or increase it. |
| **Filesystem state** | No files match yet, or everything matching is a directory (directories are not listed). |
| **Truncation confusion** | Empty list with `"truncated": false` means zero matches, not “hit limit.” |

Compare with the same `pattern` / `root` using your shell’s `fd` or `find` **only** to validate expectations—do not paste untrusted patterns into a shell on production systems.

---

## Performance concerns

**Symptom:** Slow or heavy CPU on huge directories; large `matches` payloads.

**Mitigations:**

- Set a lower **`max_results`** so the backend stops earlier (watch **`truncated`: true**).
- Prefer **`fd`** on **`PATH`** for large trees; it is usually faster than the pure Go walk.
- Narrow **`root`** to the smallest subtree you need.
- Use **`max_depth`** when you only care about shallow files.

**Caveat:** The protocol returns a single list per call—there is no paging in the current tools.

---

## Debugging tips

1. Call **`get_server_info`** first to confirm the session and version.
2. Run **`find_files`** on a **tiny** `root` (e.g. a temp dir with one known file) using the **same** env the host uses (`FILE_DISCOVERY_MCP_ROOTS`, `PATH`).
3. Inspect the response **`backend`** field (`fd` vs `walk`) to see which path is active.
4. On failure, read **`code`** (`PATH_NOT_ALLOWED`, `SEARCH_FAILED`, etc.) and the human **`error`** line; for `fd` failures, stderr is often appended there.
5. If the process exits immediately: check **stderr** from the host—startup errors (e.g. allowlist load / cwd) print `file-discovery-mcp: …` and exit non-zero (**`cmd/file-discovery-mcp`**).

---

## Logging guidance

**Today:** The stock server has **no debug or verbose logging flag**. It writes to **stderr** only for:

- **Startup failure** (`BuildServer` / config), and  
- **Non-zero exit** from the MCP run loop,

both prefixed with `file-discovery-mcp:`.

**MCP stdout/stderr:** **stdout** is reserved for the MCP stream. Avoid wrapping the binary in scripts that print to stdout. **`find_files`** does not log per-request detail; the only structured diagnostic is the tool result (`ok`, `code`, `error`).

If you need traces, use your **host’s** MCP or process logging, or wrap stderr capture at the supervisor level—**do not** expect line-delimited request logs from file-discovery-mcp itself.

---

## Related documents

- [Using](using.md) — invocation and env.
- [Examples](examples.md) — sample inputs/outputs.
- [Security](security.md) — reporting and constraints.
