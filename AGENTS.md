# AGENTS.md — `go-config`

This module is part of the eSlider `go-*` library standard (inventar ASR-0008).

## Purpose

Convert **env**, **YAML**, **JSON**, and **INI** into nested `map[string]any` and Go structs (and back), with multi-source merging and the `envc` CLI.

## Public API surface

- Subpackages: `env`, `yaml`, `json`, `ini` — each exports `New`, `(*Codec).Map`, `Unmarshal`, `UnmarshalContext`, `Marshal`, `WriteTo`, and format-specific options.
- `cmd/envc` — binary `envc`: `convert`, `get`, `merge`.
- Internals under `internal/` are not stable API.

Breaking changes require a new major SemVer tag (or `/v2` module path if the policy changes).

## Testing policy

Follows the eSlider "no synthetic mocks" policy (see [.cursor/rules/no-synthetic-mocks.mdc](.cursor/rules/no-synthetic-mocks.mdc)):

- **Unit tests** — pure inputs for merge, keymap, structconv, env flat-map helpers.
- **`httptest`** — allowed only to test **our** HTTP client behaviour in `internal/source` (not third-party API emulation).
- **Fixtures** — real files under `fixtures/`; `testfixtures` resolves paths from module root.

## Decisions

Architecture Significant Requirements: [docs/asr/README.md](docs/asr/README.md).

## Checklist before release

Local sanity check (CI runs the same on every PR):

```sh
go mod tidy
go vet ./...
go test -race -shuffle=on -count=1 ./...
golangci-lint run --timeout 5m
```

Versioning is **automated via [release-please](https://github.com/googleapis/release-please-action)**
and [GoReleaser](https://goreleaser.com) — do **not** hand-edit version strings or tag manually:

1. Commit using [Conventional Commits](https://www.conventionalcommits.org/)
   (`feat:`, `fix:`, `feat!:` for breaking, etc.).
2. `.github/workflows/release-please.yml` opens a release PR on each push to `main`
   with the computed next SemVer and an updated `CHANGELOG.md`.
3. Merging that PR creates the `vX.Y.Z` tag. `.github/workflows/release.yml` then runs
   GoReleaser to publish cross-platform `envc` binaries and a GitHub Release.

## Related

- `inventar/docs/asr/ASR-0008.md` — Go library module conventions
- `go-onlyoffice/AGENTS.md` — reference for the eSlider library template
