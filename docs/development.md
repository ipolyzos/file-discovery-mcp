# Development

**Purpose:** Clone fd-mcp, build and test locally, run the stdio server safely, and match what GitHub Actions enforces—so your PR is review-ready on the first push.

---

## Prerequisites

- **Go** — Version in **`go.mod`** (currently **1.25.x**). Use [go.dev/dl](https://go.dev/dl/) or your distro packages; `go version` should satisfy the toolchain line in the repo root.
- **Git** — For branching and PRs.
- **Optional: `fd`** — Not required to build or test; improves the `find_files` backend when on `PATH`. Without it, tests force the **walk** fallback where needed.

No Docker or Node is required for core development.

---

## Install (clone and toolchain)

From your workspace:

```bash
git clone https://github.com/ipolyzos/fd-mcp.git
cd fd-mcp
go env GOPATH   # optional sanity check
go mod download
```

There is no `make` target yet; everything goes through the **`go`** tool.

---

## Build

From the repository root:

```bash
go build -o fd-mcp ./cmd/fd-mcp
```

Install into `$GOBIN` or `$GOPATH/bin` if you prefer:

```bash
go install ./cmd/fd-mcp
```

---

## Test

Run the same quality checks you should run before every push:

```bash
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
```

- **`gofmt`** — CI fails if any `.go` file would change; format before committing.
- **`go vet`** — Static checks for common mistakes.
- **`go test ./...`** — Unit and integration tests under `cmd/...` and `internal/...`.

For a quick loop during a single package:

```bash
go test -count=1 ./internal/mcpapp/...
```

---

## Run the stdio server locally

The server speaks MCP on **stdin/stdout**. Do **not** rely on an interactive TTY attached to the same stream the client expects to use.

**Smoke run (will idle until stdin closes or SIGINT):**

```bash
go run ./cmd/fd-mcp
```

**With a built binary and scoped roots (example, Unix):**

```bash
FD_MCP_ROOTS="$PWD" ./fd-mcp
```

- **SIGINT** — Clean shutdown exits **0** (see `cmd/fd-mcp` exit mapping).
- To exercise tools end-to-end, use an **MCP client** or inspector that launches the binary as a subprocess with pipes—not `fd-mcp` in a raw shell session for protocol debugging unless you know the client owns stdio.

---

## Project structure (overview)

| Path | Role |
|------|------|
| `cmd/fd-mcp` | `main`, stdio `Run`, SIGINT-aware exit codes. |
| `internal/mcpapp` | MCP server, tools (`get_server_info`, `find_files`), limits. |
| `internal/allowlist` | `FD_MCP_ROOTS` and path containment. |
| `internal/findfiles` | `fd` subprocess (fixed argv) or `WalkDir` search. |
| `internal/meta` | Version string. |
| `docs/` | Long-form documentation. |
| `.github/workflows/` | CI workflow definitions. |

See [architecture.md](architecture.md) for layers and data flow.

---

## Branching strategy

| Branch / pattern | Use |
|------------------|-----|
| **`develop`** | Integration branch; open PRs here for day-to-day work. |
| **`main`** | Release line; promote from `develop` via release PRs when cutting versions. |
| **`feature/*`** | New capability (e.g. `feature/foo`). |
| **`bugfix/*`** | Fixes (e.g. `bugfix/exit-code`). |
| **`chore/*`** | Tooling, deps, CI, docs scaffolding. |
| **`docs/*`** | Documentation-focused branches. |

**Flow:** Branch off **`develop`**, push your topic branch, open a **PR into `develop`**, wait for **green CI**, then merge (merge strategy as maintainers prefer for the repo).

---

## Commit expectations

- **Small and focused** — One logical change per commit when possible; easy to bisect and revert.
- **Revertible** — Avoid mixing unrelated refactors with behavior changes.
- **Formatted** — `gofmt` on Go files; no churn-only formatting mixed with large features unless necessary.
- **Message style** — Short, imperative subject (per common hooks and reviewer habits), e.g. “Add X”, “Fix Y”, “Document Z”.

Push after each commit on long-running branches if that matches team practice; at minimum push before opening/updating a PR.

---

## CI expectations (GitHub Actions)

Workflow: **`.github/workflows/ci.yml`**.

| Trigger | Branches |
|---------|----------|
| **push** | `main`, `develop`, `bugfix/**`, `chore/**`, `feature/**` |
| **pull_request** | Target `main` or `develop` |

**Job `go` (Ubuntu):** matrix **Go oldstable** and **Go stable** — steps:

1. **Verify formatting** — `test -z "$(gofmt -l .)"`
2. **Build** — `go build -v ./...`
3. **Test** — `go test -v ./...`
4. **Vet** — `go vet ./...`

Your PR should pass both matrix legs locally (use a Go version close to **stable** if you only install one SDK).

---

## Related documents

- [Contributing](contributing.md) — review norms and broader project expectations.
- [Architecture](architecture.md) — package roles and safety boundaries.
- [Documentation index](index.md)
