# User guide

**Audience:** People who run **fd-mcp** alongside an MCP-capable application (IDE, assistant, or custom client) and want to search the filesystem **safely** and **read-only**. This guide is narrative and practical. For exact tool JSON and error codes, use the [client contract](client-contract.md).

---

## What fd-mcp is

**fd-mcp** is a small program that speaks the [Model Context Protocol](https://modelcontextprotocol.io/) over **standard input and output** (stdio). It exposes tools your client can call so an assistant can **list files** under directories you allow—using filename patterns similar to shell globs—without running a shell or modifying disk.

Under the hood it may use the **[`fd`](https://github.com/sharkdp/fd)** CLI when installed, or fall back to a built-in directory walk. You do not type `fd` yourself; the server invokes it with fixed arguments only.

---

## When to use it

Use fd-mcp when you want:

- **Scoped search:** Only paths under roots you configure (or under the server’s starting directory by default).
- **Structured answers:** File paths returned as data the model can reason about, not raw terminal text.
- **No shell from patterns:** Search patterns are **not** executed as shell commands.

Avoid relying on it for **full-path** matching, **content** search inside files, or **writing** files—it does not do those things.

---

## Basic usage patterns

### Running the server

The MCP client starts fd-mcp as a **child process** and connects stdin/stdout to the protocol. Typical pattern:

1. Install or build the `fd-mcp` binary (see the repository [README](../README.md) under **Build locally**).
2. In your client’s MCP settings, set the **command** to that binary (and optional **environment variables**).
3. Set the process **working directory** if your client allows it—it affects how **relative** `root` paths in `find_files` are resolved.

### Environment variables (what matters day to day)

| Variable | Effect |
|----------|--------|
| **`FD_MCP_ROOTS`** | List of directories (same separator style as your OS `PATH`) that searches may use as boundaries. If **unset**, only the directory the server was started in counts as an allowed root. |
| **`FD_MCP_MAX_RESULTS`** | Default cap on how many paths one `find_files` call returns when the client does not pass `max_results` (server default **500** if unset). |
| **`FD_MCP_ABSOLUTE_MAX`** | Hard ceiling the server applies even if the client asks for more (default **50000**). |

### Tools you will see

- **`get_server_info`** — Version and capability flags; does not touch the filesystem.
- **`find_files`** — Search for files whose **base name** matches a **glob**, under a directory **`root`** you choose (subject to the allowlist above).

---

## Example scenarios

Below are the **arguments** your client sends to **`find_files`** (not the full MCP wire format). Names and types match what assistants usually fill in.

### Find files by name

**Goal:** Locate files named exactly `README.md` somewhere under your project.

| Field | Example value |
|-------|----------------|
| `pattern` | `README.md` |
| `root` | `.` or an absolute path under an allowed root |

The server matches **file names**, not full paths. A pattern without `*` treats the name literally (subject to glob rules).

### Filter by extension

**Goal:** List Go source files in a tree.

| Field | Example value |
|-------|----------------|
| `pattern` | `*.go` |
| `root` | Path to the repository or subfolder you allow |

Use a single glob on the **basename**. Patterns like `**/*.go` are not “path-aware” in the contract: stick to `*.go` and point `root` at the subtree you care about.

### Limit how deep or how many results

Optional fields:

- **`max_depth`** — Non-negative integer; omit or leave unset for unlimited depth below `root`.
- **`max_results`** — Positive integer; omit to use the default from `FD_MCP_MAX_RESULTS` (still subject to the absolute cap).

---

## Expected outputs

### Successful search

The tool returns structured content including:

- **`matches`:** Absolute file paths (ordinary files; directories are not listed as matches).
- **`truncated`:** `true` if the server stopped because it hit the result limit—there may be more files than returned.
- **`backend`:** Whether **`fd`** or **`walk`** was used (informational).
- **`root`**, **`pattern`**, **`max_results`:** Echo of what was effectively applied after normalization and clamping.

Path order is **not** guaranteed; do not depend on sorting unless your client sorts.

### When something is wrong

Failures use the same structured shape with **`ok: false`**, a short **`code`** (for example `PATH_NOT_ALLOWED` if `root` resolves outside configured roots), and an **`error`** message intended for humans. Your host may also show transport or validation errors from the MCP layer; those are separate from these logical codes.

---

## Limitations

- **Basename glob only:** Patterns apply to the file **name**, not the whole path. You narrow location by choosing `root`, not by embedding `/` in the pattern.
- **Read-only:** No create, overwrite, or delete.
- **Roots are fixed at startup:** Allowlist comes from **`FD_MCP_ROOTS`** (or the initial working directory) when the process starts; changing env vars later requires restarting the server if your client cached the process.
- **Two backends:** With **`fd`** vs **`walk`**, edge cases (symlinks, exact glob edge cases) can differ slightly.
- **No file contents:** Search is by **path/name** only, not inside file bodies.
- **No MCP resources/prompts:** Only tools are implemented today.

---

## Related documents

- [Using](using.md) — Operator-oriented usage placeholders as docs grow.
- [Security](security.md) — Threat context and reporting.
- [Client contract](client-contract.md) — Machine-oriented API guarantees.
- [Examples](examples.md) — Sample configurations, if present.
