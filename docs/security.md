# Security

**Purpose:** Document what **file-discovery-mcp** does and does not protect against, how risks are reduced in code, and how to report defects that weaken those guarantees.

This is not legal advice or a certification. Deployers remain responsible for OS hardening, MCP host configuration, and trust in any assistant that can invoke tools.

---

## Scope and trust boundaries

**In scope for this document:** The `file-discovery-mcp` process: MCP stdio server, tool handlers, allowlist, search backends (`fd` subprocess or in-process walk).

**Out of scope (assumed adversarial or trusted per your deployment):**

- The **MCP host** (it can start `file-discovery-mcp` with arbitrary env, cwd, and attach stdio).
- The **`fd` binary** resolved from `PATH` when that backend is used (supply-chain and `PATH` integrity).
- **Users or models** who legitimately call tools; the design goal is to limit blast radius when inputs are malicious or mistaken, not to authenticate callers.

---

## Threat model

### Arbitrary command execution

**Risk:** A tool invocation could be crafted to run shell commands, invoke arbitrary binaries, or pass user text to an interpreter (`sh -c`, `eval`, etc.).

**Relevant attackers:** Untrusted content in tool arguments (e.g. model-generated `pattern` / `root`), a compromised client sending MCP `tools/call` payloads, or a malicious wrapper in `PATH` masquerading as `fd`.

**Design position:** The stock server **does not** spawn a shell. `find_files` passes **`pattern`** only into `fd --glob <pattern>` with a fixed argv template or into **`filepath.Match`** on basenames during a normal directory walk. There is no “run this string as a command” surface in application code.

**Residual risk:** Execution of the **`fd`** executable discovered via **`exec.LookPath("fd")`**. That is a **single fixed program name**, not user-controlled argv[0], but a hostile **`PATH`** could substitute a trojan `fd`. The process also runs with the **same privileges** as the MCP host (no privilege drop in `file-discovery-mcp`).

### Path traversal and unauthorized filesystem read

**Risk:** A caller supplies a `root` that escapes intended directories (e.g. symlinks, `..` sequences, or confusing absolute paths) and reads paths outside the operator’s intent.

**Relevant attackers:** Same as above; also misconfiguration (overly broad `FILE_DISCOVERY_MCP_ROOTS`).

**Design position:**

1. **`root`** is trimmed, cleaned, made absolute, and checked with **`MustBeUnderRoots`** using **`filepath.Clean`** and **`filepath.Rel`**-based containment against the allowlist loaded at startup (`internal/allowlist`).
2. **`pattern`** is used only for **basename** matching, not as a second path component. It does not select arbitrary absolute paths by itself.
3. **Reads** are implicit in directory traversal and metadata (`fd` or `WalkDir`); there is no file-content tool.

**Residual risk:** **`filepath` containment** is not the same as full **filesystem canonicalization** (e.g. symlink-aware policy). A `root` inside an allowed directory that resolves via symlinks to another volume may produce matches the operator did not mentally include. Heavy symlink layouts are called out as a known edge case in repository documentation. **Misconfiguration** (e.g. setting `FILE_DISCOVERY_MCP_ROOTS` to `/`) intentionally grants wide read scope; the software enforces **containment under configured roots**, not **least privilege** unless the operator configures tight roots.

### Denial of service and resource exhaustion

**Risk:** Large trees or huge result sets exhaust CPU, memory, or response size.

**Mitigation:** **`max_results`** is clamped to an absolute ceiling (**`FILE_DISCOVERY_MCP_ABSOLUTE_MAX`**, default **50000**) and defaults when unset (**`FILE_DISCOVERY_MCP_MAX_RESULTS`**, default **500**). Deep trees can still be expensive before the cap is hit; **`max_depth`** reduces depth. This limits **returned rows** and backend **`--max-results`** for `fd`; it does not guarantee **constant-time** or **bounded scan** of the entire disk.

### Information disclosure via responses

**Risk:** Returned paths reveal repository layout, usernames in paths, or existence of sensitive files under allowed roots.

**Position:** This is **intended behavior** for a discovery tool. Scope sensitivity via **`FILE_DISCOVERY_MCP_ROOTS`** and host policies, not by expecting `file-discovery-mcp` to redact paths.

---

## Mitigations (implemented)

| Control | Mechanism |
|---------|-----------|
| **No shell for tool arguments** | `exec.CommandContext(fdExe, args...)` with fixed flags; no `sh -c`. Walk path uses Go stdlib only. |
| **Structured validation** | Non-empty trimmed `pattern` / `root`; **`max_depth` ≥ 0** when set; JSON types enforced by the MCP SDK layer. |
| **Path scoping** | **`FILE_DISCOVERY_MCP_ROOTS`** (or startup cwd only if unset); each call’s resolved **`root`** must lie under an allowlisted root or **`PATH_NOT_ALLOWED`**. |
| **Result bounding** | **`clampMaxResults`**; success includes effective **`max_results`** and **`truncated`**. |
| **Read-only tools** | No APIs to write, delete, or chmod; **`get_server_info`** does not touch the filesystem. |

---

## Safe defaults (operational)

| Default | Effect |
|---------|--------|
| **`FILE_DISCOVERY_MCP_ROOTS` unset** | Only the **current working directory at server start** is an allowed root. Narrower than `/` but **only as safe as that cwd** (wrong cwd → wrong boundary). |
| **`FILE_DISCOVERY_MCP_MAX_RESULTS` unset** | Per-call default **500** matches (after handler logic). |
| **`FILE_DISCOVERY_MCP_ABSOLUTE_MAX` unset** | Hard cap **50000** on requested/ effective result count. |

**Recommendation:** Set **`FILE_DISCOVERY_MCP_ROOTS`** explicitly to the smallest directory set required. Set **`cwd`** in the MCP host to a known project root when you rely on the implicit single-root mode.

---

## Future considerations

Non-committing list of hardening directions the project may pursue:

- **Symlink / canonical path policy** — optional `EvalSymlinks` or platform-specific containment rules with clear contract updates.
- **`fd` pinning** — configurable absolute path to `fd` to avoid `PATH` substitution (trade-off: portability).
- **Audit / debug mode** — structured, opt-in logging without breaking stdio MCP on stdout.
- **Rate limiting or scan budgets** — wall-clock or inode budgets beyond row caps.
- **Supply-chain** — reproducible builds, signed releases (process outside this repo’s current scope).

---

## Reporting security issues

Do **not** open a **public** issue for an unfixed vulnerability.

Report privately so maintainers can assess impact and release a fix before wide disclosure. Use the repository host’s **private security advisory** mechanism (e.g. GitHub **Security → Report a vulnerability**) if enabled, or contact the maintainers through an agreed private channel.

Include: affected **version/commit**, **environment**, **reproduction** steps, and whether **`FILE_DISCOVERY_MCP_ROOTS`** / **`PATH`** were non-default.

---

## Related documents

- [Architecture](architecture.md) — component boundaries.
- [Client contract](client-contract.md) — normative tool behavior and error codes.
- [Troubleshooting](troubleshooting.md) — permissions and backend behavior.
