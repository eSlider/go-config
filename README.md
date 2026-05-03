# go-config

[![Go Reference](https://pkg.go.dev/badge/github.com/eslider/go-config.svg)](https://pkg.go.dev/github.com/eslider/go-config)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/eSlider/go-config)](https://github.com/eSlider/go-config/releases/latest)
[![Tests](https://github.com/eSlider/go-config/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/eSlider/go-config/actions/workflows/test.yml)
[![Lint](https://github.com/eSlider/go-config/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/eSlider/go-config/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/eslider/go-config)](https://goreportcard.com/report/github.com/eslider/go-config)
[![GitHub Stars](https://img.shields.io/github/stars/eSlider/go-config?style=social)](https://github.com/eSlider/go-config/stargazers)

Convert **env**, **YAML**, **JSON**, and **INI** to and from Go `map[string]any` and structs. Multi-source inputs merge with **deep map merge**: nested maps combine, **scalar leaves are last-write-wins**, and **slices** default to **replace** (opt-in **concat** via `WithSliceMerge`). Keys are normalized with a configurable **lower+alnum** rule so `sub-service`, `SUB_SERVICE`, and `SubService` line up across formats. Built on [go-viper/mapstructure/v2](https://github.com/go-viper/mapstructure).

## Architecture

```mermaid
flowchart TB
  Sources["Sources\nbytes reader file URL process env"]
  Parser["Parser\ngodotenv yaml ini json"]
  Norm["keymap.Walk\nNormalizer"]
  MergeOp["merge.DeepMerge"]
  Map["map string any"]
  MS["structconv\nmapstructure v2"]
  Struct["Go struct"]
  Sources --> Parser --> Norm --> MergeOp --> Map
  Map -->|Unmarshal| MS --> Struct
  Struct -->|Marshal| MS --> Map
  Map -->|WriteTo Marshal| Parser
```

## Hero example

```go
yamlCfg := yaml.New(yaml.WithURL("https://raw.githubusercontent.com/eSlider/mail-archive/refs/heads/master/docker-compose.yml"))
envCfg := env.New(
	env.WithFile(".default.env"), // lowest priority
	env.WithFile(".env"),
	env.WithCurrentEnvironment(), // highest priority — process env wins
)

var svc MyService
_ = yamlCfg.Unmarshal(&svc)
_ = envCfg.Unmarshal(&svc) // later sources override earlier scalar leaves; maps recurse
```

- **Maps** merge recursively (sub-trees are combined, not replaced wholesale).
- **Scalar leaves**: last-write-wins when you list sources lowest → highest priority.
- **Slices**: `merge.Replace` by default; use `WithSliceMerge(merge.Concat)` to append.

Runnable offline variant: see `Example_hero_offline` in [example_hero_test.go](example_hero_test.go).

## Install

```sh
go get github.com/eslider/go-config
go install github.com/eslider/go-config/cmd/envc@latest
```

## Quick start

### 1. Single YAML file

```go
c := yaml.New(yaml.WithFile("config.yaml"))
var cfg AppConfig
if err := c.Unmarshal(&cfg); err != nil { /* ... */ }
```

### 2. JSON over HTTPS with a header

```go
c := json.New(
	json.WithURL("https://api.example.com/v1/config.json"),
	json.WithHTTPHeader("Authorization", "Bearer "+token),
)
var cfg AppConfig
_ = c.Unmarshal(&cfg)
```

### 3. Cross-format conversion (YAML → JSON)

```go
ctx := context.Background()
m, _ := yaml.New(yaml.WithFile("in.yaml")).Map(ctx)
b, _ := json.New().Marshal(m)
os.WriteFile("out.json", b, 0o644)
```

## CLI: `envc`

Formats for `--from` / `--to`: `yaml`, `json`, `ini`, `env`. Run `envc help` or
`envc <command> -h` for flags.

| Command                                                            | Purpose                                                         |
| ------------------------------------------------------------------ | --------------------------------------------------------------- |
| `envc convert --from yaml --to env --input - --output -`           | Convert stdin YAML to dotenv on stdout                          |
| `envc get --from yaml --path service.name config.yaml`             | Print one path (dot segments map to lower+alnum keys)         |
| `envc merge --from yaml --to json --output out.json a.yaml b.yaml` | Deep-merge multiple YAML files, emit JSON                       |

### Help and version

```sh
envc help
envc version
```

### `convert` examples

```sh
# Pipe YAML in, JSON on stdout (default --input - and --output -)
printf 'app:\n  port: 8080\n' | envc convert --from yaml --to json

# File → file
envc convert --from json --to yaml --input settings.json --output settings.yaml

# Remote YAML → local dotenv
envc convert --from yaml --to env \
  --input https://example.com/config.yaml \
  --output .env.generated
```

### Apply YAML or INI to the **current** bash session

`envc convert … --to env` prints **dotenv-style** lines (`KEY=value`). They are normal
shell assignments, not `export` lines, so use **`set -a`** (allexport) if child processes
must see the variables. **Process substitution** `<(…)` needs **bash** (not plain `sh`).

```bash
# YAML → current shell (and export to children while sourcing)
set -a
source <(envc convert --from yaml --to env --input config.yaml)
set +a

# INI → current shell
set -a
source <(envc convert --from ini --to env --input app.ini)
set +a
```

Same idea from stdin:

```bash
set -a
source <(cat deploy.yaml | envc convert --from yaml --to env)
set +a
```

Only do this with **trusted** config files (same caution as `source` on any generated
script): values are expanded by the shell when you `source` them.

### `get` examples

Path segments use the same **lower+alnum** rules as the library (e.g. `sub-service`
and `SubService` both address `subservice`).

```sh
# Value from a file
envc get --from yaml --path app.port config.yaml

# Nested key; stdin when the last argument is `-` or omitted with a pipe
cat config.yaml | envc get --from yaml --path service.subservice.name -
```

### `merge` examples

Sources are merged **in order** (later files override scalar leaves; maps recurse).
Use `-` once to read **one** merged stdin blob as a source (same format as the others).

```sh
envc merge --from yaml --to json defaults.yaml overrides.yaml

# Write merged JSON to stdout, then save
envc merge --from yaml --to json base.yaml local.yaml | tee merged.json

# Stdin plus files: first load stdin as YAML, then merge each file
cat patch.yaml | envc merge --from yaml --to yaml - base.yaml
```

### Stdin and `-`

With `--input -` (convert) or `-` as the get/merge input, `envc` reads **until EOF**.
From an interactive terminal with no pipe, that waits until you press **Ctrl-D**
(end of input). Prefer `--input path` or a URL when scripting.

## API (all codecs)

| Method                                          | Description                       |
| ----------------------------------------------- | --------------------------------- |
| `New(opts...)`                                  | Construct codec                   |
| `Map(ctx)`                                      | Merged `map[string]any`           |
| `Unmarshal(dst)` / `UnmarshalContext(ctx, dst)` | Decode into struct (or map)       |
| `Marshal(src)` / `WriteTo(w, src)`              | Encode struct or `map[string]any` |

Shared options (each subpackage): `WithBytes`, `WithReader`, `WithFile`, `WithURL`, `WithHTTPHeader`, `WithHTTPClient`, `WithKeyNormalizer`, `WithSliceMerge`, `WithTrim`, `WithWeaklyTyped`, `WithTagName`, `WithDecodeHook`.

`env` adds: `WithCurrentEnvironment`, `WithPrefix`.

## Cross-format mapping

| Go                        | YAML                       | ENV                       |
| ------------------------- | -------------------------- | ------------------------- |
| `Service.SubService.Name` | `service.sub-service.name` | `SERVICE_SUBSERVICE_NAME` |

INI uses dotted sections, e.g. `[service.subservice]` with `name=...`.

## Related libraries

| Module                                                    | Role           |
| --------------------------------------------------------- | -------------- |
| [go-matrix-bot](https://github.com/eSlider/go-matrix-bot) | Matrix bots    |
| [go-onlyoffice](https://github.com/eSlider/go-onlyoffice) | OnlyOffice API |
| [go-ollama](https://github.com/eSlider/go-ollama)         | Ollama client  |

## Contributing

Testing expectations, local commands, commit message conventions, and how **release-please**
and **GoReleaser** publish tags and `envc` binaries are documented in
**[CONTRIBUTING.md](CONTRIBUTING.md)**.

## Decisions

Repo-local ASRs: [docs/asr/README.md](docs/asr/README.md).

## License

MIT © Andriy Oblivantsev
