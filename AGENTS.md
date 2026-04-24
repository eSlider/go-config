# AGENTS.md — `go-env`

This module is part of the eSlider `go-*` library standard (inventar ASR-0008).

## Purpose

Decode `os.Environ()` into Go structs. Extracted from
`produktor.io/ai-fabric/pkg/env` on 2026-04-24.

## Public API surface

- `env.Unmarshal(dst any, opts ...Option) error`
- `env.UnmarshalPrefix(dst any, prefix string, opts ...Option) error`
- `env.AsMap() map[string]any`
- `env.AsMapPrefix(prefix string) map[string]any`
- Options: `WithTrim`, `WithWeaklyTyped`, `WithTagName`, `WithDecodeHook`

Breaking changes require a new major version tag (SemVer). Internal helpers
(`asMapFromEnviron`, `trimStringHook`, `insertPath`) are unexported and
may change without notice.

## Testing policy

Follows the eSlider "no synthetic mocks" policy:

- **Unit tests** (`env_test.go`): pure inputs only. We pass synthetic
  `environ` slices to `asMapFromEnviron` — that is not a mock of anything
  external; it's the normal way to test a pure function. Use `t.Setenv`
  for `Unmarshal*` happy-path tests since Go's `os.Environ` lookup is
  well-defined local behaviour, not a vendor protocol.
- **No `httptest`** — this library has no HTTP surface.
- **No integration-test build tag** — everything is in-process.

## Checklist before release

```sh
cd go-env
go mod tidy
go vet ./...
go test -race -count=1 ./...
golangci-lint run --timeout 5m   # same preset as go-onlyoffice
```

Bump `CHANGELOG.md`, tag `vX.Y.Z`, push.

## Related

- `inventar/docs/asr/ASR-0008.md` — Go library module conventions
- `inventar/docs/asr/ASR-0008-ai-fabric-audit.md` — why this module exists
- `go-onlyoffice/AGENTS.md` — reference for the eSlider library template
