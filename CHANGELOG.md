# Changelog

All notable changes to `go-config` are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-05-02

### Added

- Module rename to `github.com/eslider/go-config` with subpackages `env`, `yaml`, `json`, `ini`.
- Uniform `Codec` API: `New`, `Map`, `Unmarshal`, `UnmarshalContext`, `Marshal`, `WriteTo`.
- `internal/source` for bytes, readers, files, and HTTP(S) URLs with optional headers and custom `http.Client`.
- `internal/keymap` recursive key walk with default `LowerAlnum` normalizer.
- `internal/merge` deep merge: recursive maps, last-write-wins scalars, configurable slice replace vs concat.
- `internal/structconv` around `github.com/go-viper/mapstructure/v2`.
- CLI `envc` (`cmd/envc`): `convert`, `get`, `merge`.
- Fixtures under `fixtures/` for identity, merge, edge, and invalid parser cases.
- `testfixtures` helper for tests.
- Repo-local ASRs in `docs/asr/`.

### Changed

- **Breaking:** the v0.1.0 `github.com/eslider/go-env` API (`Unmarshal`, `UnmarshalPrefix`, `AsMap`, …) is not re-exported; use `env.New(...).Unmarshal(...)`.

### Dependencies

- `github.com/go-viper/mapstructure/v2`, `github.com/joho/godotenv`, `gopkg.in/yaml.v3`, `gopkg.in/ini.v1`.

## [0.1.0] - 2026-04-24

Initial release as `github.com/eslider/go-env`. Extracted from `produktor.io/ai-fabric/pkg/env` per
`inventar/docs/asr/ASR-0008-ai-fabric-audit.md`.

### Added

- `Unmarshal(dst, opts...)` — decode all process env vars into a Go struct.
- `UnmarshalPrefix(dst, prefix, opts...)` — prefix-scoped variant that
  strips the prefix before building the path.
- `AsMap()` / `AsMapPrefix(prefix)` — escape hatch returning the nested
  `map[string]any` built from `os.Environ()`.
- Options: `WithTrim`, `WithWeaklyTyped`, `WithTagName`, `WithDecodeHook`.
- TrimSpace on string values is enabled by default; disable with
  `WithTrim(false)`.
- Error wrapping with `%w` on both decoder configuration and decode
  failures — matches `markets-platform/TP-general-code/pkg/system`
  behaviour, improves on the original ai-fabric version.

### Tests

- 15 pure-Go unit tests. No mocks. `asMapFromEnviron` is tested with
  synthetic `environ` slices (not a mock — a pure input). Follows the
  eSlider `go-*` no-synthetic-mocks policy.
