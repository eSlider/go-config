# go-env

[![Go Reference](https://pkg.go.dev/badge/github.com/eslider/go-env.svg)](https://pkg.go.dev/github.com/eslider/go-env)

Tiny, zero-ceremony library for decoding process environment variables into
Go structs.

- Uses `_` as a path separator: `SERVICE_HTTP_PORT` → `Service.HTTP.Port`.
- Weakly-typed by default (`"1"` → `int`, `"true"` → `bool`).
- Optional prefix filter with automatic stripping.
- Pluggable decode hooks via [`mitchellh/mapstructure`][mapstructure].

## Install

```sh
go get github.com/eslider/go-env
```

## Quick start

```go
package main

import (
	"fmt"

	"github.com/eslider/go-env"
)

type Config struct {
	Service struct {
		HTTP struct {
			Port int
		}
		Key string
	}
}

func main() {
	var cfg Config
	if err := env.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", cfg)
}
```

With a prefix:

```go
// Only looks at APP_* variables; strips the APP_ prefix before decoding.
_ = env.UnmarshalPrefix(&cfg, "APP_")
```

## API

| Function | Purpose |
|---|---|
| `Unmarshal(dst, opts...)` | Decode all env vars into `dst`. |
| `UnmarshalPrefix(dst, prefix, opts...)` | Same, but only vars starting with `prefix`. |
| `AsMap()` / `AsMapPrefix(prefix)` | Return the nested `map[string]any` used by the decoder (debugging, custom decoders). |

### Options

| Option | Default | Purpose |
|---|---|---|
| `WithTrim(bool)` | `true` | TrimSpace every string value. |
| `WithWeaklyTyped(bool)` | `true` | `mapstructure`'s weakly-typed coercion. |
| `WithTagName(string)` | `"mapstructure"` | Struct tag name for field overrides. |
| `WithDecodeHook(h)` | — | Append a `mapstructure.DecodeHookFunc` to the chain. |

## Semantics

- **Path collisions**: first write wins. If both `FOO=1` and `FOO=2` exist
  in the environ, only `FOO=1` is kept. This matches the original
  `ai-fabric/pkg/env` behaviour.
- **Case-insensitive keys**: all path components are lower-cased; struct
  fields are matched via `mapstructure` which is also case-insensitive.
- **`_`-only delimiter**: there's no escape — if a variable legitimately
  contains `_` inside a "leaf" name, you must restructure your struct to
  match the nested layout.

## Status

Extracted from `produktor.io/ai-fabric` as part of the eSlider `go-*`
library standard (ASR-0008). Merges the best parts of three previously
divergent copies:

- `produktor.io/ai-fabric/pkg/env`
- `markets-platform/TP-general-code/pkg/system/env.go`
- the various `pkg/system/env.go` snapshots inside `var/agents/issue-*/`

## License

MIT © Andriy Oblivantsev

[mapstructure]: https://github.com/mitchellh/mapstructure
