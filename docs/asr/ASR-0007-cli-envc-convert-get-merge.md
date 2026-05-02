# ASR-0007: CLI `envc` — convert, get, merge

## Context

Operators need a small tool to translate configuration between formats and to inspect merged trees without writing Go.

## Decision

1. Binary name **`envc`**; module path **`github.com/eslider/go-config/cmd/envc`** (`go install …/cmd/envc@latest`).
2. Subcommands: **`convert`**, **`get`**, **`merge`** — each implemented as a **`Run*(args, stdin, stdout, stderr) int`** function for in-process tests.
3. **No** third-party CLI framework in v1 of the CLI.

## Consequences

- `cmd/envc` may grow flags; keep business logic in subpackages, not in `main` beyond wiring.

## Status

Accepted — 2026-05-02.

## References

- inventar ASR-0008 §2 — CLI placement
- ASR-0003
