# MCP client contract (fd-mcp)

**Purpose:** Normative expectations for MCP clients integrating with **fd-mcp**. This document specifies transport, tool names, JSON field semantics, error classification, safety guarantees, and versioning—not usage tutorials.

Conformance keywords are used in their usual sense: **must** (requirement), **should** (recommendation), **may** (optional).

---

## Transport and session

- **Transport:** **stdio** only (newline-delimited JSON-RPC over process **stdin** / **stdout**). Clients **must** own both streams; they **must not** assume a TTY for protocol I/O.
- **Process exit:** On **SIGINT** (or parent cancellation of the server context), the process **should** exit with status **0** after shutting down. Non-zero exit **may** indicate startup or runtime failure before or outside normal tool handling (see **System errors**).
- **Implementation identity:** The server advertises MCP `Implementation` with `name` **`fd-mcp`** and `version` equal to the string returned by tool output field **`version`** in **`get_server_info`** (same source as `internal/meta.Version` at build time).

---

## Tool naming

- Tool names are **stable**, **lowercase**, ASCII **`snake_case`**, and **must** match exactly:
  - `get_server_info`
  - `find_files`
- New tools **may** be added in minor releases; removals or renames **should** be treated as breaking and reflected in version policy (see **Stability**).

---

## Tools (current revision)

### `get_server_info`

**Invocation:** `tools/call` with `name`: **`get_server_info`**.

**Arguments:** Empty JSON object **`{}`** (no properties). Extra properties **should** be rejected by schema validation before the handler runs.

**Structured result** (object in `structuredContent`; field names JSON keys):

| Field | Type | Meaning |
|-------|------|---------|
| `version` | string | Server release / build version (non-empty). |
| `description` | string | Short human-readable description. |
| `read_only` | boolean | **Always `true`** for this implementation: tools do not mutate filesystem content. |
| `filesystem_search_wip` | boolean | **`false`** when `find_files` is implemented and advertised. Clients **may** treat `true` as “search not ready.” |

**Filesystem:** The handler **must not** read the filesystem.

---

### `find_files`

**Invocation:** `tools/call` with `name`: **`find_files`**.

#### Input (`arguments` object)

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `pattern` | string | yes | After trim of ASCII whitespace, **must not** be empty. Semantics: **basename glob** (e.g. `*.go`), not a full path pattern. Interpretation matches **`fd --glob`** when backend is `fd`, else **`filepath.Match`** on the file’s base name (minor differences between backends are possible). |
| `root` | string | yes | After trim, **must not** be empty. Interpreted as a directory path: if relative, resolved against the server process **current working directory** at call time, then cleaned and made absolute. **Must** fall under at least one entry of **`FD_MCP_ROOTS`** (see **Safety**). |
| `max_results` | integer | no | If **omitted or ≤ 0**, effective limit is **`FD_MCP_MAX_RESULTS`** env default (**500** if unset/invalid), then clamped. If **> 0**, used as requested limit then **clamped** so the effective value **does not exceed** **`FD_MCP_ABSOLUTE_MAX`** (default **50000** if unset/invalid env). Final effective limit **must** be ≥ **1**. |
| `max_depth` | integer \| null | no | If present, **must** be **≥ 0**. Omitted or null means **unlimited** depth under `root`. |

Invalid combinations **must** yield a **logical tool failure** (`ok: false`, see below)—not necessarily an MCP protocol error.

#### Output (`structuredContent` object)

**Success** (`ok` **`true`**):

| Field | Type | Meaning |
|-------|------|---------|
| `ok` | boolean | **`true`** |
| `matches` | array[string] | Absolute, cleaned file paths (directories excluded). Order **undefined**. |
| `truncated` | boolean | **`true`** iff results were cut off by **`max_results`** (or backend equivalent); **`false`** if the full set within the limit was returned. |
| `backend` | string | **`fd`** or **`walk`**. |
| `root` | string | Absolute cleaned directory that was searched. |
| `pattern` | string | Effective pattern after trim. |
| `max_results` | integer | Effective cap applied after clamping. |

**Logical failure** (`ok` **`false`**):

| Field | Type | Meaning |
|-------|------|---------|
| `ok` | boolean | **`false`** |
| `code` | string | Machine-readable code (see **Error model**). |
| `error` | string | Human-readable message; **must not** be parsed for control flow. |
| `matches` | — | Omitted or empty. |

---

## Error model

Errors fall into three classes:

### 1. Validation errors (tool arguments / schema)

- **MCP / JSON-RPC error:** Arguments **fail** JSON Schema validation enforced by the SDK (wrong types, unknown required shape). Client **must** treat as **invalid request**; no `find_files` structured body is defined by this contract for that case.
- **Logical validation** (handler): `pattern` / `root` empty after trim, unresolvable `root`, `max_depth < 0`. Response **must** be `ok: false` with:
  - `INVALID_PATTERN` — empty or missing effective pattern.
  - `INVALID_ROOT` — empty root or resolution failure.
  - `INVALID_DEPTH` — negative `max_depth`.

### 2. Execution errors (policy or I/O during search)

- **`PATH_NOT_ALLOWED`:** Resolved `root` is not contained in any configured allowlist root (`FD_MCP_ROOTS` semantics).
- **`SEARCH_FAILED`:** Search backend failed (e.g. `fd` non-zero exit, I/O error during walk). `error` **should** carry a short diagnostic string; clients **must not** assume a stable token format inside `error`.

These **must** use `ok: false` and the corresponding `code`.

### 3. System errors

- **Startup:** Server fails before advertising tools (e.g. cannot load roots env). Process **may** exit non-zero without MCP session.
- **Transport / protocol:** Connection loss, malformed RPC outside tool validation—standard MCP / JSON-RPC rules apply.
- **Cancellation:** Client or host cancels context; in-flight `find_files` **should** stop; response behavior **may** follow SDK defaults for cancelled calls.

---

## Safety guarantees (client-observable)

Clients **may** rely on the following for the current **fd-mcp** implementation:

1. **Read-only:** `find_files` does not create, modify, or delete files or directories.
2. **Path restriction:** Every successful `find_files` search uses a `root` verified against **`FD_MCP_ROOTS`** (path-list separator as for `PATH`); if unset at server start, only the process initial working directory is allowed (**implementation detail:** loaded once at `BuildServer`; document in ops guides).
3. **Bounded results:** The number of paths in `matches` **must not** exceed the effective **`max_results`** returned in the success object.
4. **No shell:** User-supplied `pattern` and `root` are **not** passed to `/bin/sh` or equivalent; `fd` is invoked with a fixed argv template.

**Not guaranteed without explicit future contract changes:** identical ordering of `matches`, identical glob semantics between `fd` and `walk`, symlink-normalized path containment.

---

## Stability and versioning

- **Server `version` string** (MCP + `get_server_info`): **should** change on each release that affects behavior or contract; format is **project-defined** (currently semantic-style **x.y.z**). Clients **may** log or gate features on this string but **must not** parse it as an API minor/major tuple unless documented.
- **Additive changes:** New optional fields on success objects **may** appear without a major bump; clients **should** ignore unknown keys.
- **Breaking changes:** Removing or renaming tools, removing fields, or changing `code` meanings **should** coincide with a **major** version bump and **changelog** entry.
- **Tool list:** Clients **should** call `tools/list` after connect rather than hard-coding only historical tool names.

---

## Related documents

- [Architecture](architecture.md) — internal layers and data flow.
- [Security](security.md) — reporting and threat context.
- [Using](using.md) — operator-oriented usage (non-normative).
- [Documentation index](index.md)
