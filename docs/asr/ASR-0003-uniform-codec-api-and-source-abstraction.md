# ASR-0003: Uniform codec API and Source abstraction

## Context

Users switch between YAML, JSON, INI, and env-file formats; a different function surface per format would increase cognitive load and test matrices.

## Decision

1. Every format package exports **`New(opts...) *Codec`** with **`Map`**, **`Unmarshal`**, **`UnmarshalContext`**, **`Marshal`**, **`WriteTo`**.
2. Bytes, readers, files, and URLs are modeled by **`internal/source.Source`** with **`Open(ctx) (io.ReadCloser, error)`**.
3. HTTP(S) sources support **`WithHTTPHeader`** and **`WithHTTPClient`**.

## Consequences

- Cross-format examples in the README stay structurally identical.
- URL behaviour is testable with `httptest` against our client only.

## Status

Accepted — 2026-05-02.

## References

- ASR-0005 — merge semantics across multiple sources
- ASR-0008 — httptest policy alignment
