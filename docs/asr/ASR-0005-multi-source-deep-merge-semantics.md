# ASR-0005: Multi-source deep-merge semantics

## Context

Layered configuration (defaults file < local overrides < process environment) requires predictable composition rules.

## Decision

1. Sources are applied in **option order** (left → right). Each source parses to a `map[string]any`, then **`merge.DeepMerge`** folds them.
2. **Maps recurse** — nested keys are unioned; existing nested maps are not replaced wholesale by a sibling scalar.
3. **Scalar leaves** use **last-write-wins** (later source wins).
4. **Slices** default to **`merge.Replace`**; **`WithSliceMerge(merge.Concat)`** appends prior and new slice elements.

## Consequences

- Env files and `WithCurrentEnvironment()` compose when ordered **lowest → highest** priority.
- Tests lock behaviour via `fixtures/merge/*`.

## Status

Accepted — 2026-05-02.

## References

- ASR-0003
- `internal/merge`
