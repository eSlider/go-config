# ASR-0008: GitHub presentation and release flow

## Context

The repository was renamed from `eSlider/go-env` to `eSlider/go-config`. The v0.1.0 tag existed locally before the rename. Presentation should match other eSlider `go-*` libraries (e.g. `go-matrix-bot`).

## Decision

1. **v0.1.0 preservation:** push tag `v0.1.0` and create a GitHub **Release** with notes extracted from `CHANGELOG.md` before landing breaking work.
2. **Development branch:** `release/v2` carries the `go-config` implementation until merged to `main`.
3. **Rename:** `gh repo rename go-config` from `go-env`; local `origin` URL updated to `https://github.com/eSlider/go-config.git`.
4. **Metadata:** `gh repo edit` sets description, **`homepage`** `https://pkg.go.dev/github.com/eslider/go-config`, and **topics** (`go`, `golang`, `config`, `yaml`, `json`, `ini`, `env`, `dotenv`, `mapstructure`, `codec`, `encoder`, `decoder`, `cli`, `library`).
5. **README:** badges row, mermaid architecture diagram, hero snippet, quick starts, CLI table, API tables, related libraries, link to `docs/asr/`.

## Consequences

- Old `github.com/eslider/go-env` URLs redirect for a limited GitHub window; module path for v0.1.0 code remains the historical `go-env` import.

## Status

Accepted — 2026-05-02.

## References

- [go-matrix-bot](https://github.com/eSlider/go-matrix-bot) — README pattern reference
