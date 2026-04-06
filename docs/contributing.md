# Contributing to file-discovery-mcp

Thank you for taking the time to improve file-discovery-mcp. This document describes how we work together: branching, commits, pull requests, code expectations, documentation, and reporting issues.

Participation is governed by the repository [**Code of Conduct**](../CODE_OF_CONDUCT.md).

For day-to-day setup and CI details, see [development.md](development.md). For system design and safety boundaries, see [architecture.md](architecture.md) and [security.md](security.md).

---

## How to contribute

1. **Open or follow an issue** when the change is non-trivial, so scope and approach are visible before a large PR lands.
2. **Branch from `develop`**, implement the change, keep commits focused, and open a **pull request into `develop`**.
3. **Ensure tests and formatting pass** locally (and in CI) before requesting review.
4. **Respond to review feedback** promptly; maintainers may merge when CI is green and concerns are addressed.

We welcome bug fixes, documentation improvements, tests, and carefully scoped features that respect the project’s **read-only** and **no arbitrary execution** constraints.

---

## Branch naming conventions

Use a **short, descriptive slug** after the prefix. Prefix indicates intent:

| Prefix | Use |
|--------|-----|
| `feature/` | New user-visible capability or substantial behavior change |
| `bugfix/` | Fixes incorrect behavior or regressions |
| `docs/` | Documentation-only (including guides under `docs/`) |
| `chore/` | Tooling, CI, dependency bumps, non-user-facing maintenance |
| `refactor/` | Internal restructuring without intended behavior change |

Examples: `feature/stream-results`, `bugfix/exit-code`, `docs/contributing`.

**Integration branches:** **`develop`** receives topic PRs; **`main`** reflects released history. Do not push directly to `main` except via agreed release process.

---

## Commit messages

- Use a **short, imperative subject line** (e.g. “Add X”, “Fix Y”, “Document Z”), roughly **50–72 characters** when practical.
- **One logical change per commit** when possible so history is easy to bisect and revert.
- **Body (optional):** explain *why* when the subject alone is insufficient; avoid repeating the diff.
- **No** proprietary or confidential content in public commits.

If your environment uses a commit-msg hook, follow its rules (for example, an imperative first line).

---

## Pull request process

### Before you open a PR

Run the same checks CI enforces (from [development.md](development.md)):

```bash
gofmt -w .
test -z "$(gofmt -l .)"
go vet ./...
go test ./...
```

Fix failures **before** opening or updating the PR so reviewers spend time on design, not on avoidable CI noise.

### PR scope

- **Prefer small PRs** that are straightforward to review: one concern, minimal unrelated diffs.
- Large changes should be **split** or accompanied by a design note in the PR description.
- **Avoid drive-by refactors** mixed with feature work unless necessary for correctness.

### Description and review

- Summarize **what** changed and **why**; link related issues.
- Call out **behavioral or security** implications explicitly when relevant.
- **Maintainers review** for correctness, safety, clarity, and fit with the architecture. Approval is required before merge per repository policy.
- Keep discussion **professional and constructive**.

Target branch for routine work is **`develop`**. Releases are promoted to **`main`** separately.

---

## Code standards

### Readability

- Prefer **clear names** and **small functions** over cleverness.
- Match **existing style** in surrounding files (formatting, error handling patterns, package layout).
- **Exported APIs** and public tool contracts deserve concise comments where behavior is not obvious.

### Safety constraints

file-discovery-mcp is a **read-only** server with **scoped** filesystem access. Contributions must preserve:

- **No arbitrary shell** — user input must not be interpolated into shell commands; use fixed argv, `exec` with explicit fields, or pure Go I/O as today’s code does.
- **Path scoping** — new tools that touch the filesystem should respect **`FILE_DISCOVERY_MCP_ROOTS`** (or extend `allowlist` in a documented, reviewable way).
- **Bounded work** — respect **timeouts**, **context cancellation**, and **result limits** where applicable.

If you are unsure whether a change weakens these guarantees, **ask in the PR** early.

---

## Documentation requirements

- **User-visible behavior**, new environment variables, or **tool fields** should be reflected in **README** or the appropriate **`docs/`** page in the same PR when practical.
- **Architecture or workflow** shifts should update [architecture.md](architecture.md) or [development.md](development.md) as appropriate.
- Keep docs **accurate and concise**; avoid duplicating the same long prose in multiple places—**link** instead.

---

## Reporting issues

### Bug reports

Include where possible:

- **Version** (git tag or commit) and **Go version**
- **Steps to reproduce** and **expected vs actual** behavior
- **Environment** (OS, relevant env vars such as `FILE_DISCOVERY_MCP_ROOTS`)
- **Minimal** logs or snippets (no secrets)

### Security issues

Do **not** open public issues for unfixed vulnerabilities. Follow [security.md](security.md) for responsible disclosure.

### Feature requests

Describe the **problem** and **constraints** (read-only, stdio MCP, etc.). Large features may need discussion before a PR.

---

## Related documents

- [Code of Conduct](../CODE_OF_CONDUCT.md) — community standards and enforcement
- [Development](development.md) — local build, test, CI matrix, running the server
- [Architecture](architecture.md) — layers and data flow
- [Security](security.md) — reporting and trust model
- [Documentation index](index.md)
