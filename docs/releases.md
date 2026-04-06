# Releases

**Purpose:** Maintainer checklist for promoting **`develop`** to **`main`**, tagging, and publishing **GitHub Releases**, consistent with [contributing](contributing.md) and [development](development.md).

Routine feature work **does not** use this flow; it merges to **`develop`** via topic PRs only.

---

## Version scheme

- **Format:** **`vMAJOR.MINOR.PATCH`** for Git tags (e.g. **`v0.2.0`**). The Programmatic string in **`internal/meta.Version()`** matches **`MAJOR.MINOR.PATCH`** **without** the `v` prefix (e.g. **`0.2.0`**).
- **Early `0.x`:** Breaking or additive MCP/tool changes should be documented in [CHANGELOG.md](../CHANGELOG.md); strict semver for every JSON field is not guaranteed beyond the [client contract](client-contract.md) and changelog notes.
- **Single source for runtime:** `get_server_info` and MCP **`Implementation.version`** read **`meta.Version()`**. That value **must** match the tagged release users build or download.

---

## Branch flow

| Branch | Role |
|--------|------|
| **`develop`** | Default integration target for **`feature/*`**, **`bugfix/*`**, **`docs/*`**, **`chore/*`**, etc. |
| **`main`** | **Release line:** commit history that tagged versions point at. |

**Promotion path:** **`develop` → `main`** only through a **release pull request** (or equivalent reviewed merge) that brings `main` up to the commit (or set of commits) you intend to tag. **Do not** push directly to **`main`** outside the agreed release process ([contributing](contributing.md)).

After `main` contains the release commits, create the **Git tag** on the release commit on **`main`**, then publish the **GitHub Release** (below).

---

## PR requirements before a release

A release PR (**`develop` → `main`**) should not merge until:

1. **CI green** on that PR — same checks as [development](development.md): `gofmt`, **`go build -v ./...`**, **`go test -v ./...`**, **`go vet ./...`** on **Go oldstable** and **stable** (`.github/workflows/ci.yml`).
2. **Review** — maintainer approval per [contributing](contributing.md); no open blockers on behavioral or security impact.
3. **Version metadata updated** on the branch being released (typically already on `develop` before promotion):
   - **`internal/meta.Version()`** equals the new **`x.y.z`**.
   - **[CHANGELOG.md](../CHANGELOG.md)** has a **`## [x.y.z] - YYYY-MM-DD`** section describing user-visible changes.
   - **[README.md](../README.md)** **Current release** line updated to **`vx.y.z`** if maintainers keep it in sync.
4. **No drive-by scope** — the PR should be release-oriented (promotion + version/changelog commits as needed), not a grab bag of unrelated work.

Topic PRs merged **to `develop`** already required green CI before merge; the release PR is a final gate on the exact snapshot going to **`main`**.

---

## Tagging strategy

1. Merge the release PR so **`main`** equals the intended tree.
2. On **`main`**, at the release commit:  
   `git tag vX.Y.Z`  
   Use **annotated** tags if your signing policy requires it (`git tag -a vX.Y.Z -m "vX.Y.Z"`).
3. **`git push origin vX.Y.Z`** (or push all tags per team habit).

Tags are **immutable**: if a tag is wrong, prefer a new patch (**`v0.2.1`**) rather than re-pointing an existing public tag.

---

## Release notes expectations

- **CHANGELOG.md** is the canonical long-lived record. Each release adds a dated section with subsections (e.g. **Added**, **Fixed**, **Known limitations**) as appropriate.
- Entries should be **factual** and **user-facing** where possible (tools, env vars, security-relevant behavior). Link to issues or PRs when helpful.
- **Known limitations** and breaking contract changes **must** be called out explicitly.

---

## GitHub releases

1. On GitHub: **Releases → Draft a new release**.
2. **Choose tag:** `vX.Y.Z` (create from target `main` if not pushed yet—prefer tagging locally then pushing as above).
3. **Title:** typically **`vX.Y.Z`**.
4. **Description:** paste or summarize the **`[x.y.z]`** section from **CHANGELOG.md**; add install pointers ([deploying](deploying.md)) or security notes ([security](security.md)) if the release warrants it.
5. **Assets** — optional: attach cross-built **`file-discovery-mcp`** binaries. The repository CI **does not** build or upload them automatically today; add artifacts manually or via a separate pipeline if needed.

---

## Related documents

- [Contributing](contributing.md) — PR and branch rules.
- [Development](development.md) — CI matrix and local checks.
- [Deploying](deploying.md) — what operators build or pin after a tag exists.
