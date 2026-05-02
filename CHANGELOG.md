# Changelog

All notable changes to `go-config` are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.1](https://github.com/eSlider/go-config/compare/v0.2.0...v0.2.1) (2026-05-02)


### Bug Fixes

* add motivation config precedence regression test ([d44b568](https://github.com/eSlider/go-config/commit/d44b56848bade38ff49fed0bac7cdf3e00dbc982))
* **ci:** make lint and windows tests reliable ([a6a51ec](https://github.com/eSlider/go-config/commit/a6a51ec3ef328fe403eccd205b38ed7f646819f1))
* **ci:** pin lint toolchain to Go 1.26 ([26aa513](https://github.com/eSlider/go-config/commit/26aa51360b6e1cd73049b65d8bbcea5bba73498c))
* **ci:** unblock release and cross-platform test pipeline ([76b4f49](https://github.com/eSlider/go-config/commit/76b4f499c0842788d7b1db9140064f84266507b8))
* **ci:** use golangci-lint v2 config format ([cb0ee8a](https://github.com/eSlider/go-config/commit/cb0ee8ac223cbec376291656d6e51be0b0f6f769))


### Code Refactoring

* move testfixtures to internal package and update imports ([eadcf9f](https://github.com/eSlider/go-config/commit/eadcf9fa14d8e3ff3c028e9fb1ee22d602f596e3))


### Documentation

* enhance release flow documentation and add version command to CLI ([0ae52ef](https://github.com/eSlider/go-config/commit/0ae52ef72d9b161020516a528ccc44bdc18247f5))

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
