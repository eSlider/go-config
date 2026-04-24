# Changelog

All notable changes to `go-env` are documented here.
The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-04-24

Initial release. Extracted from `produktor.io/ai-fabric/pkg/env` per
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
