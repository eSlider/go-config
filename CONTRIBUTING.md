# Contributing to go-config

Thank you for helping improve this module. This document covers **testing expectations**,
**local checks**, and **how releases are cut** (so nobody tags or edits versions by hand).

For repository purpose, public API notes, and agent-oriented context, see [AGENTS.md](AGENTS.md).

## Testing policy

This project follows the eSlider **no synthetic mocks** policy (see [.cursor/rules/no-synthetic-mocks.mdc](.cursor/rules/no-synthetic-mocks.mdc)):

- **Unit tests** — pure inputs for merge, keymap, structconv, env flat-map helpers.
- **`httptest`** — allowed only to test **our** HTTP client behaviour in `internal/source` (not third-party API emulation).
- **Fixtures** — real files under `fixtures/`; `testfixtures` resolves paths from module root.

## Before you open a pull request

Run the same checks CI will run (adjust paths if your checkout layout differs):

```sh
go mod tidy
go vet ./...
go test -race -shuffle=on -count=1 ./...
golangci-lint run --timeout 5m
```

## Commit messages

Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `perf:`,
`feat!:` or `BREAKING CHANGE:` for breaking changes, etc.). Release automation reads these
messages to choose the next SemVer.

## Release and versioning

Versions are **computed from git history**, not edited in Go source. Do **not** create
release tags manually or bump version constants in code.

### CI on every change

1. Every push and pull request runs [`test.yml`](.github/workflows/test.yml) (matrix: Go
   1.22, 1.23, stable × linux / macOS / windows, with `-race -shuffle=on`) and
   [`lint.yml`](.github/workflows/lint.yml) (`go vet` + `golangci-lint`).

### Release Please and GoReleaser

2. On pushes to `main`, [`release-please.yml`](.github/workflows/release-please.yml) parses
   commits since the last tag and opens (or updates) a **release pull request** that bumps
   [`CHANGELOG.md`](CHANGELOG.md) and [`.release-please-manifest.json`](.release-please-manifest.json)
   to the next SemVer.

3. **Merging that release PR** creates the git tag `vX.Y.Z` and a GitHub Release.

4. The tag triggers [`release.yml`](.github/workflows/release.yml), which runs
   [GoReleaser](https://goreleaser.com) to cross-compile the `envc` binary (linux / darwin /
   windows × amd64 / arm64), injects `-X main.version={{.Version}}` (and related ldflags) at
   link time, and attaches archives plus `checksums.txt` to the release.

Local `go install` builds of `envc` report `dev` for `envc version` unless you pass your own
`-ldflags`; release binaries show the released version.

## Architecture decisions

Significant design choices for this repository are captured as **Architecture Significant
Requirements** (ASRs):

- **Repo-local ASRs:** [docs/asr/README.md](docs/asr/README.md)
- **eSlider Go library conventions (inventar):** `inventar/docs/asr/ASR-0008.md` — also linked
  from [AGENTS.md](AGENTS.md) under Related.
