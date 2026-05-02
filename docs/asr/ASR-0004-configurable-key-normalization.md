# ASR-0004: Configurable key normalization

## Context

Real-world configuration mixes `kebab-case`, `snake_case`, `SCREAMING_SNAKE`, and struct field names. Requiring `mapstructure` tags on every field is noisy.

## Decision

1. After parsing and merging, codecs run **`internal/keymap.Walk`** with a **`Normalizer`** (default **`LowerAlnum`**: lowercase + strip non `[a-z0-9]`).
2. Callers may set **`WithKeyNormalizer(nil)`** to disable walking, or supply a custom function.

## Consequences

- POSIX env keys without `-` still align with YAML keys that contain hyphens once normalized.
- CLI `get` applies the same normalizer to each path segment for consistency.

## Status

Accepted — 2026-05-02.

## References

- ASR-0003
