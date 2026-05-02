# ASR-0006: Third-party library choices

## Context

We need battle-tested parsers and struct bridging without maintaining our own grammars.

## Decision

1. **Struct / map bridging:** `github.com/go-viper/mapstructure/v2` (maintained fork of `mitchellh/mapstructure`).
2. **YAML:** `gopkg.in/yaml.v3`.
3. **JSON:** `encoding/json` (stdlib).
4. **INI:** `gopkg.in/ini.v1` with `Loose` parsing for tolerant files.
5. **`.env` files:** `github.com/joho/godotenv` (`UnmarshalBytes`).
6. **CLI:** stdlib `flag` only (no Cobra).

## Consequences

- `go.mod` carries the above direct dependencies (plus transitive modules as resolved by MVS).

## Status

Accepted — 2026-05-02.

## References

- [go-viper/mapstructure](https://github.com/go-viper/mapstructure)
