# ASR-0002: Clean break from go-env v0.1.0

## Context

`go-env` v0.1.0 exposed package-level helpers (`Unmarshal`, `UnmarshalPrefix`, `AsMap`, …). The new design is `Codec`-centric with multiple sources per format.

## Decision

1. **No compatibility shims** in `go-config` re-exporting the old API.
2. v0.1.0 remains tagged on git history as `github.com/eslider/go-env@v0.1.0` for archival consumers.
3. Migration path is documented in `CHANGELOG.md` and the README.

## Consequences

- Single-call-site migrations must switch to `env.New(env.WithCurrentEnvironment(), …).Unmarshal(&cfg)`.

## Status

Accepted — 2026-05-02.

## References

- CHANGELOG `[0.2.0]`
- ASR-0001
