# Using file-discovery-mcp

**Purpose:** How to **invoke** the server and what **hosts** must provide for stdio MCP integration. For a gentler introduction, see the [user guide](user-guide.md). For copy-paste samples, see [examples](examples.md).

---

## Invoking the MCP server

- **Command:** Run the `file-discovery-mcp` binary **with no arguments**. The process is an MCP **stdio** server only; there is no HTTP listen mode or subcommands in the stock binary.
- **Working directory:** Set explicitly in the host when possible. It matters in two ways:
  - If **`FILE_DISCOVERY_MCP_ROOTS`** is **unset**, the server’s initial working directory is loaded as the **sole** allowed root at startup.
  - Relative `find_files` **`root`** values are resolved against the process **current working directory** at the time of each call.
- **Environment:** Optional vars (**`FILE_DISCOVERY_MCP_ROOTS`**, **`FILE_DISCOVERY_MCP_MAX_RESULTS`**, **`FILE_DISCOVERY_MCP_ABSOLUTE_MAX`**) are read when the process starts (roots and defaults are fixed for the lifetime of that process).

**Manual smoke test (optional):** From a terminal, `./file-discovery-mcp` will wait on stdin; without a proper JSON-RPC client attached, it will not “do” anything visible. Real use is always **host + subprocess**.

---

## Integration expectations (stdio)

| Expectation | Detail |
|-------------|--------|
| **Transport** | **stdin** and **stdout** carry newline-delimited JSON-RPC (MCP). The host **must** connect both pipes to the child; they are not for interactive human typing under normal use. |
| **No alternate channel** | Do not assume a socket or TTY replaces stdio for the protocol in this project’s stock server. |
| **Binary-safe I/O** | Treat stdout as protocol-only; avoid mixing logging or extra prints into the server binary (the released server does not). |
| **Lifecycle** | The host **spawns** `file-discovery-mcp`, performs MCP **initialize** / capability exchange, then **`tools/list`** / **`tools/call`** as needed. On shutdown, closing stdin or sending **SIGINT** ends the session; the server exits **0** on SIGINT after cleanup (see repository `cmd/file-discovery-mcp`). |
| **Concurrency** | Assume tools may be called in sequence per MCP client rules; follow your SDK’s threading model if you embed a client. |

Hosts (Cursor, Claude Desktop, custom agents, etc.) differ in configuration file layout but agree on this model: **executable path**, optional **cwd**, optional **env map**, and **stdio** transport.

---

## Environment variables

| Variable | Role |
|----------|------|
| **`FILE_DISCOVERY_MCP_ROOTS`** | OS path-list (same separator as `PATH`) of directory roots under which **`find_files`** may run. Empty/unset → single root = cwd at server start. |
| **`FILE_DISCOVERY_MCP_MAX_RESULTS`** | Default per-call cap when **`max_results`** is omitted or non-positive (server default **500**). |
| **`FILE_DISCOVERY_MCP_ABSOLUTE_MAX`** | Upper bound the server applies after clamping (default **50000**). |

---

## Tools overview

| Tool | Role |
|------|------|
| **`get_server_info`** | Returns version and flags; **no** filesystem access. Arguments: `{}`. |
| **`find_files`** | Lists files whose **basename** matches **`pattern`** (glob), under directory **`root`** (must resolve under allowed roots). Optional **`max_results`**, **`max_depth`**. |

Structured success/failure bodies use **`ok`**, **`matches`**, **`truncated`**, **`code`**, etc., as in the [client contract](client-contract.md).

---

## Related documents

- [Examples](examples.md) — configs and sample tool I/O.
- [User guide](user-guide.md) — end-user narrative.
- [Client contract](client-contract.md) — normative schemas and error codes.
- [Security](security.md) — constraints and reporting.
