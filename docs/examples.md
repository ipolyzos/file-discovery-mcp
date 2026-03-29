# Examples

**Purpose:** Minimal, realistic **`find_files`** inputs and **structured outputs**, plus a representative MCP **stdio** server snippet. Exact host JSON varies by product; adapt keys to your client’s docs.

---

## MCP client configuration (stdio)

Illustrative shape: **command** (absolute path recommended), optional **environment**, optional **working directory** if the host supports it.

```json
{
  "mcpServers": {
    "fd-mcp": {
      "command": "/home/alice/bin/fd-mcp",
      "env": {
        "FD_MCP_ROOTS": "/home/alice/work/myproject:/home/alice/oss/contrib"
      }
    }
  }
}
```

- If **`FD_MCP_ROOTS`** is omitted, only the server’s **startup cwd** is allowed (set cwd in the host when you can).
- Replace paths and top-level key (`mcpServers` vs `servers`) with whatever your editor or runtime expects.

---

## Simple search

Find all Go files under a repository root (basename glob only).

**`tools/call`** — `name`: `find_files`, `arguments`:

```json
{
  "pattern": "*.go",
  "root": "/home/alice/work/myproject"
}
```

**Structured result** (success; `matches` order not guaranteed):

```json
{
  "ok": true,
  "matches": [
    "/home/alice/work/myproject/cmd/fd-mcp/main.go",
    "/home/alice/work/myproject/internal/findfiles/findfiles.go"
  ],
  "truncated": false,
  "backend": "fd",
  "root": "/home/alice/work/myproject",
  "pattern": "*.go",
  "max_results": 500
}
```

`backend` may be **`walk`** if `fd` is not used. `max_results` reflects the **effective** cap after defaults and clamping (`500` here = default when the client omitted `max_results`).

---

## Filtered search

Same tree, but only **two directory levels** below `root`, and at most **20** paths.

**Arguments:**

```json
{
  "pattern": "*.md",
  "root": "/home/alice/work/myproject",
  "max_depth": 2,
  "max_results": 20
}
```

**Structured result** (success with truncation):

```json
{
  "ok": true,
  "matches": [
    "/home/alice/work/myproject/README.md",
    "/home/alice/work/myproject/docs/README.md"
  ],
  "truncated": false,
  "backend": "walk",
  "root": "/home/alice/work/myproject",
  "pattern": "*.md",
  "max_results": 20
}
```

If more than 20 files matched, you would see **`truncated": true`** and at most 20 entries in **`matches`**.

---

## Edge cases (structured failure)

### Root outside configured roots

**Arguments** (repo allowed only under `/home/alice/work/myproject`, but `root` points elsewhere):

```json
{
  "pattern": "*.go",
  "root": "/etc"
}
```

**Structured result:**

```json
{
  "ok": false,
  "code": "PATH_NOT_ALLOWED",
  "error": "allowlist: path \"/etc\" is not under allowed roots"
}
```

Exact **`error`** text may differ slightly; clients should branch on **`code`**, not parse **`error`**.

### Empty pattern

```json
{
  "pattern": "   ",
  "root": "/home/alice/work/myproject"
}
```

**Structured result:**

```json
{
  "ok": false,
  "code": "INVALID_PATTERN",
  "error": "pattern is required"
}
```

---

## `get_server_info`

**Arguments:** `{}`

**Structured result** (representative):

```json
{
  "version": "0.2.0",
  "description": "fd-mcp: read-only filesystem discovery MCP server",
  "read_only": true,
  "filesystem_search_wip": false
}
```

---

## Related documents

- [Using](using.md) — invocation and stdio expectations.
- [User guide](user-guide.md) — usage narrative.
- [Troubleshooting](troubleshooting.md) — when searches fail unexpectedly.
- [Client contract](client-contract.md) — full field definitions.
