# ASR-0001: Module rename and multi-package layout

## Context

The original `github.com/eslider/go-env` module only decoded process environment variables. The product scope expanded to YAML, JSON, INI, URLs, multi-source merging, and a CLI. The old name and single-package layout no longer matched the capability surface.

## Decision

1. The canonical module path is **`github.com/eslider/go-config`**.
2. Format-specific code lives in first-class subpackages: **`env/`**, **`yaml/`**, **`json/`**, **`ini/`**.
3. Shared building blocks live under **`internal/{source,keymap,merge,structconv,bytesutil}`**.
4. The CLI lives at **`cmd/envc/`** (one binary per inventar **ASR-0008** §2 — CLI inside the same module).
5. The GitHub repository is named **`eSlider/go-config`** (renamed from `go-env`).

## Consequences

- Consumers must change their `import` paths and `go get` target.
- CI, badges, and documentation reference `go-config` and `eslider/go-config`.

## Status

Accepted — 2026-05-02.

## References

- inventar ASR-0008 — Go library module conventions
- ASR-0002 — API clean break
