# Deploying

**Purpose:** How to **run**, **ship**, and **pin** fd-mcp for real use given the **current** repository—single Go binary, no first-party container or package recipes yet.

---

## Running locally

- **Build:** From repo root, `go build -o fd-mcp ./cmd/fd-mcp` or `go install ./cmd/fd-mcp` (needs Go **1.25+** per `go.mod`).
- **Invoke:** `./fd-mcp` with **no arguments**; protocol on **stdin/stdout**. Set **`FD_MCP_ROOTS`**, **`cwd`**, and **`PATH`** (for optional `fd`) in the MCP host config, not by hand in production shells unless you control stdio end-to-end.
- **Smoke check:** `go run ./cmd/fd-mcp` blocks until stdin closes or **SIGINT**; use an MCP client to exercise tools.

Details: [development.md](development.md) (build/test), [using.md](using.md) (stdio, env).

---

## Packaging options

**What exists today**

| Approach | Notes |
|----------|--------|
| **Local `go build` / `go install`** | Default; produces one static-ish binary (depends on Go build mode and cgo). |
| **Cross-compile** | e.g. `GOOS=linux GOARCH=amd64 go build -o fd-mcp ./cmd/fd-mcp` from any supported host. |
| **Git tag checkout** | Build from a **tagged** commit for reproducibility; see **Versioning**. |

**What does not exist in-repo yet**

- No **Dockerfile**, **Helm chart**, **deb/rpm**, or **GoReleaser** config.
- **Releases:** GitHub **Releases** may ship tag + artifact manually; operators can mirror the same `go build` steps in their own pipeline.

If you need a container, build a minimal image that **only** copies a prebuilt `fd-mcp` and sets `USER`/non-root as required by your policy; keep **stdin/stdout** as the MCP attach point the orchestrator expects.

---

## Integration into environments

- **MCP hosts (Cursor, Claude Desktop, custom):** Register **command** = absolute path to `fd-mcp`, optional **env** (`FD_MCP_ROOTS`, `FD_MCP_MAX_RESULTS`, `PATH`), optional **working directory**. Shape varies by product—see [examples.md](examples.md).
- **Process model:** The host **spawns** the binary; there is **no** long-lived listen port in stock `fd-mcp`. “Supervision” is the host’s job (restart on crash, stderr capture).
- **Hardening:** Tight **`FD_MCP_ROOTS`**, non-privileged user, stable **`PATH`** if you rely on `fd`. See [security.md](security.md).

---

## CI/CD expectations

**Upstream (this repository)**

- Workflow: [`.github/workflows/ci.yml`](../.github/workflows/ci.yml).
- **Triggers:** push to `main`, `develop`, `feature/**`, `bugfix/**`, `chore/**`; **pull_request** targeting `main` or `develop`.
- **Checks:** `gofmt` clean, `go build -v ./...`, `go test -v ./...`, `go vet ./...` on **Go oldstable** and **stable** (Ubuntu).

**Downstream (your deploy)**

- Reproduce the same four steps (or at least **build + test + vet**) on the **Go version** you use to produce binaries.
- **Pin** a **git tag** or **module pseudo-version** in your SBOM or build args; record it next to shipped artifacts.
- CI **does not** publish binaries; release automation is manual or external until added.

---

## Versioning approach (early)

- **User-visible string:** `get_server_info` / MCP `Implementation.version` come from **`internal/meta.Version()`** (currently **0.1.0**). It **must** be updated when cutting a release that changes behavior users rely on.
- **Tags:** Prefer **`v0.x.y`** Git tags aligned with that string and [CHANGELOG.md](../CHANGELOG.md) entries.
- **Branches:** Day-to-day merges target **`develop`**; **`main`** is the release line when maintainers promote—see [development.md](development.md). Early **0.x** releases may include compatible additions without strict semver API promises beyond the [client contract](client-contract.md).
- **No migration tooling:** Upgrades are “replace binary, restart host,” with contract docs as the compatibility reference.

---

## Related documents

- [Releases](releases.md) — placeholder for fuller release policy.
- [Development](development.md) — local workflow and CI mirror.
- [Security](security.md) — deployment-facing risks.
- [Using](using.md) — runtime env and stdio.
