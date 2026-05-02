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

```sh
cd go-config
go mod tidy
go vet ./...
go test -race -count=1 ./...
golangci-lint run --timeout 5m
```

Bump `CHANGELOG.md`, tag `vX.Y.Z`, push.

## Related

- `inventar/docs/asr/ASR-0008.md` — Go library module conventions
- `go-onlyoffice/AGENTS.md` — reference for the eSlider library template
