# Architecture

**Purpose:** Describe how fd-mcp is structured: layers, data flow, principles, and how new capability should plug in without breaking safety.

---

## Overview

fd-mcp is a **stdio MCP server**: a long-lived process that speaks the [Model Context Protocol](https://modelcontextprotocol.io/) over **stdin/stdout** (newline-delimited JSON/RPC). The binary entrypoint is `cmd/fd-mcp`; it builds an MCP `Server` from `internal/mcpapp`, runs it on `mcp.StdioTransport`, and maps process exit status so normal **SIGINT** shutdown is **exit 0** (`cmd/fd-mcp` exit policy).

The product surface today is **tools only** (no resources/prompts): metadata (`get_server_info`) and **read-only** path search (`find_files`). All filesystem work is **scoped**, **bounded**, and implemented **without a shell**.

---

## High-level components

### MCP interface layer

- **Role:** Protocol boundary—tool names, JSON Schemas (via go-sdk), typed inputs/outputs, session lifecycle.
- **Location:** `internal/mcpapp` registers tools on `mcp.Server` (`BuildServer`); `github.com/modelcontextprotocol/go-sdk/mcp` handles wire format and validation for declared tool args.
- **Behavior:** Each tool handler is thin: parse/validate (often delegated to SDK), call domain packages, return **structured JSON** (`FindFilesOutput`, `ServerInfo`) or logical errors encoded in the output shape (`ok: false`, `code`, `error`) where appropriate.

### Command execution layer (fd / walk abstraction)

- **Role:** Execute **read-only discovery** under an already-validated absolute root.
- **Location:** `internal/findfiles`—`Run` selects **fd** via `exec.LookPath` + `exec.CommandContext` with a **fixed argv template** (`--glob`, `--max-results`, optional `--max-depth`, then `root`), or falls back to **`filepath.WalkDir`** + **`filepath.Match`** on basenames.
- **Important:** User input becomes **argument values**, never a command line string interpreted by `sh`. There is **no** `find` subprocess in the current codepath; the in-process walk provides a portable fallback where `fd` is absent.

### Safety layer (scoping and limits)

- **Role:** Prevent searches outside operator-chosen roots and cap blast radius (result size).
- **Location:** `internal/allowlist` loads **`FD_MCP_ROOTS`** (path-list like `PATH`), resolves the caller’s `root`, and enforces **containment** relative to configured roots. `internal/mcpapp/limits.go` clamps **`max_results`** using **`FD_MCP_MAX_RESULTS`** and **`FD_MCP_ABSOLUTE_MAX`**.
- **Behavior:** Denials are explicit (`PATH_NOT_ALLOWED`, invalid pattern/root/depth) and returned as structured tool output rather than crashing the server.

### Response formatting layer

- **Role:** Stable, machine-readable tool results and metadata for clients.
- **Location:** Structs in `internal/mcpapp` (`FindFilesOutput`, `ServerInfo`); go-sdk maps them to MCP **`structuredContent`** (and text content when applicable).
- **Fields:** Success paths expose `matches`, `truncated`, `backend` (`fd` vs `walk`), effective caps, and echo of `root` / `pattern` for traceability.

---

## Data flow (find_files)

End-to-end path for a single tool call:

1. **Request** — Client sends `tools/call` with JSON args (`pattern`, `root`, optional `max_results`, `max_depth`).
2. **Validation** — SDK validates against tool schema; handler trims strings, rejects empty pattern/root, rejects negative `max_depth`.
3. **Scoping** — Resolve `root` to an absolute path; **`MustBeUnderRoots`** ensures it lies under **`FD_MCP_ROOTS`** (or the process cwd default).
4. **Limits** — `clampMaxResults` applies environment defaults and absolute ceiling.
5. **Execution** — `findfiles.Run` invokes **fd** or **walk** under the validated root; context cancellation propagates to the `fd` child process.
6. **Transformation** — Paths collected; `truncated` set when the cap may have hidden further matches (`fd` line count vs cap; walk stops with `fs.SkipAll`).
7. **Response** — Structured JSON returned to the client (`ok: true` and fields above, or `ok: false` with `code` / `error`).

`get_server_info` skips steps 3–6 and returns static capability metadata.

---

## Design principles

| Principle | Meaning for fd-mcp |
|-----------|---------------------|
| **Read-only** | No tool creates, mutates, or deletes files; no write side-effects on the tree beyond OS read access for directory traversal. |
| **Deterministic policy** | Given the same roots, env caps, and inputs, behavior is predictable: containment and clamps are code-defined, not host-shell aliases. |
| **Safe execution** | No arbitrary subprocesses; only **`fd`** with frozen flags, or pure Go walking. Patterns are data, not script. |
| **Composability** | New tools should follow the same pattern: **allowlist** → optional **findfiles** (or future engines) → **typed output**—keep MCP handlers thin. |

---

## Non-goals

- **Arbitrary shell** — No `sh -c`, pipeline, or user-defined executable from tool parameters.
- **Write / admin operations** — No file editing, chmod, package install, or network egress as part of core tools.
- **Full query language** — Basename globs only today; not a general SQL/Ripgrep-style grammar in this slice.
- **Multi-session HTTP/SSE server** — Stdio-first; other transports would be a separate concern.

---

## Future extensibility

- **Alternate backends** — Additional discoverers can sit behind the same **validation + limits** contract as `findfiles.Run` (or a small interface) without changing MCP tool names if outputs stay compatible.
- **Richer containment** — Optional symlink canonicalization, platform-specific root semantics, or per-tool roots could extend `allowlist` with explicit compatibility notes.
- **Streaming / pagination** — Large result sets might adopt chunked responses while preserving caps and cancellation.
- **Resources / prompts** — Could be added alongside tools; keep **read-only** and **scoped** invariants explicit in new handlers.

---

## Source layout (reference)

| Path | Role |
|------|------|
| `cmd/fd-mcp` | Process entry: stdio transport, `BuildServer`, exit status mapping (`runExitCode`). |
| `internal/mcpapp` | MCP server assembly, tool handlers, `FindFilesOutput` / `ServerInfo`, env-based caps. |
| `internal/allowlist` | `FD_MCP_ROOTS` parsing, path resolution, containment checks. |
| `internal/findfiles` | fd subprocess (fixed argv) or `WalkDir` discovery; truncation signaling. |
| `internal/meta` | Version string for implementation and tools. |

---

## Related documents

- [Client contract](client-contract.md) — observable tool and error semantics.
- [Security](security.md) — trust boundaries and reporting.
- [Development](development.md) — build, test, and layout.
- [Documentation index](index.md)
