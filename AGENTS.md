# AGENTS.md — `go-config`

This module is part of the eSlider `go-*` library standard (inventar ASR-0008).

## Purpose

Convert **env**, **YAML**, **JSON**, **TOML**, and **INI** into nested `map[string]any` and Go structs (and back), with multi-source merging and the `envc` CLI.

## Public API surface

- Subpackages: `env`, `yaml`, `json`, `toml`, `ini` — each exports `New`, `(*Codec).Map`, `Unmarshal`, `UnmarshalContext`, `Marshal`, `WriteTo`, and format-specific options.
- `cmd/envc` — binary `envc`: `convert`, `get`, `merge`.
- Internals under `internal/` are not stable API.

Breaking changes require a new major SemVer tag (or `/v2` module path if the policy changes).

## Testing policy, PR checklist, and releases

Human-oriented detail lives in **[CONTRIBUTING.md](CONTRIBUTING.md)** (testing rules,
`go test` / lint commands, Conventional Commits, the release-please / GoReleaser flow, and
**architecture decisions** / repo ASRs). Follow that document for any change that will ship
in a versioned release.

## Related

- `inventar/docs/asr/ASR-0008.md` — Go library module conventions
- `go-onlyoffice/AGENTS.md` — reference for the eSlider library template
