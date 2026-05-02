# Go config: many good pieces, no whole picture

*Why I keep three different config libraries open in three different services — and what a "unified" tool would actually have to do.*

---

I have been writing Go services for years and I still do not have one obvious answer to: *"how does this `config.yaml`, that `.env`, the process environment, and my nested `Config` struct stay one consistent story across dev, staging, and prod?"*

The Go ecosystem has plenty of config libraries. Most of them are good at their slice. Almost none of them advertise honestly which slice that is, and the gap shows the moment your `Config` is more than three flat fields and your sources are more than one file.

This is not a benchmark shoot-out. It is a map of where the seams are, what the popular libraries actually own, and a pattern — *"map in the middle"* — that keeps reappearing in every working setup I have shipped, regardless of which library is on top.

---

## The problem in one scene

You clone a repo. There is a `config.yaml` for local dev, a `.env.example` nobody copies correctly, and production sets secrets through Kubernetes or systemd. Somewhere in code:

```go
type Config struct {
    HTTP struct {
        Listen string
        TLS    bool
    }
    Database struct {
        URL string
    }
}
```

Compose says `http.listen`. Kubernetes sets `HTTP_LISTEN`. Your teammate adds `database.url` in YAML, but the managed-Postgres vendor's docs say `DATABASE_URL`, so `.env` uses that.

None of these choices is wrong. They are how humans write configuration. The pain starts when one binary has to accept all of them, decide *who wins*, and decode into one typed struct without 40 lines of `os.Getenv` glued together with manual `strconv` calls.

That is the moment you notice "configuration management" is three jobs sold as one:

1. Parsing a format (YAML, JSON, INI, dotenv).
2. Composing several sources with predictable precedence.
3. Bridging a `map[string]any` into a typed struct, with weak-typing and tag rules.

A handful of services need only one of those. Most production services need at least two. Nested structs are where the seams between them tear.

---

## What "unified" should mean

If "unified" is going to mean anything testable, here is the checklist I use:

- **Parse common formats** — Ops live in YAML, secrets live in env, legacy config still lives in INI.
- **Merge layers** — Defaults < repo config < local overrides < process env (the [12-factor](https://12factor.net/config) ordering). That top layer is how Docker, Kubernetes (ConfigMaps, Secrets, `env` / `envFrom`), systemd, CI, and similar inject values into the running process: the variables are there even when nothing on disk matches prod.
- **Normalize keys** — `sub-service`, `sub_service`, and `SUB_SERVICE` should not each need their own struct tag.
- **Decode to structs** — The end state is a typed `Config`, not a `map`.
- **Encode back** — Dump effective config for a support ticket; generate a starter `.env` from defaults.
- **Operational CLI** — SREs convert and inspect config without `go run ./cmd/debug-config`.
- **Testability** — No package-level globals; deterministic fixtures; merge rules covered by unit tests.

A library can win two items and stay out of scope for the rest. That is fine. The frustration is when the README says "all your config needs" and the **stress test** — nested structs, multiple files, and process env in the same binary — is left to the reader.

---

## The landscape

The Go ecosystem is not short of options. [Awesome Go — Configuration](https://awesome-go.com/configuration/) lists dozens. The interesting question is which items on the checklist each tool optimises for.

### Viper

[Viper](https://github.com/spf13/viper) is what most Go developers reach for first. It reads files, env, flags, and remote stores; `Unmarshal` decodes into a struct via [mapstructure](https://github.com/go-viper/mapstructure). The mindshare is enormous, which means error messages are searchable and tutorials exist.

The honest cost is the env story. Binding nested keys from environment variables — `AutomaticEnv` plus `SetEnvKeyReplacer` plus per-leaf `BindEnv` — is the source of recurring bug reports for years; see for instance [#641](https://github.com/spf13/viper/issues/641) and [#2001](https://github.com/spf13/viper/issues/2001). Those issues are not "Viper is broken"; they are evidence that flattening trees into env tokens and decoding back into trees is a hard problem that Viper solves under conventions you have to learn.

I still pick Viper when a team has already standardised on it and the config shape is shallow. I avoid it for greenfield code where the env layout is non-trivial.

### Koanf

[Koanf](https://github.com/knadh/koanf) is the lighter, more composable alternative: separate **providers** (where bytes come from) and **parsers** (how to read them), explicit merge order, no package-level singleton. The mental model is closer to how I think about config anyway, which makes it easier to reason about *which layer won*.

The trade-off is honest: key normalisation, dotenv quoting edge cases, and struct-decoding policy are choices Koanf will not pre-decide for you. That is the price of the pipeline being explicit, and for a greenfield service I usually want it.

### Env-first struct loaders

[caarlos0/env](https://github.com/caarlos0/env), [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig), and [cleanenv](https://github.com/ilyakaznacheev/cleanenv) shine when the source of truth is environment variables and the goal is a typed struct with a few validators and defaults. Tiny surface area, sane parsing of durations and lists, no global state.

What they are not trying to be is a YAML merge engine. If half your config is repo YAML and half is platform env, you will still hand-roll the parser, the merge, and the precedence yourself.

### Dotenv

[joho/godotenv](https://github.com/joho/godotenv) parses `.env` files: quoting, exports, expansion. That is the entire job, and it does it well. It is a building block, not a config story — pair it with a struct loader and a merge policy and you have something usable.

### Format-only parsers

The standard library's [`encoding/json`](https://pkg.go.dev/encoding/json), [`gopkg.in/yaml.v3`](https://github.com/go-yaml/yaml), and an INI reader of your choice each turn bytes into a value. They do not say how three files plus process env become one tree. "We use `encoding/json` for config" is fine for one file; it leaves merge order, env overrides, and key aliasing as application policy duplicated across every repo.

### mapstructure

[go-viper/mapstructure](https://github.com/go-viper/mapstructure) — the v2 line, forked from the (now archived) `mitchellh/mapstructure` and maintained inside the Viper org — is the de facto bridge from `map[string]any` into structs. Tags, `WeaklyTypedInput`, decode hooks, `squash` for embedded fields, time and net.IP parsing.

It is *not* a format parser. It does not define merge precedence. It does not decide whether `SUB_SERVICE_NAME` lines up with `SubService.Name` until you give it a map where those keys already match. Almost every working setup ends up using mapstructure as a stage, even if the outer library is Viper, Koanf, or hand-written glue.

---

## Why nested and embedded structs are the real stress test

Flat structs are easy. The pain shows up here:

**Embedding.** Go pushes composition via embedded types and `Config` structs follow. The moment you embed, you have to decide: are the embedded type's keys *flattened* into the outer namespace, or *namespaced* under a sub-key? Mapstructure's `squash` exists for exactly this question, but it only helps once the merged map already matches whichever choice you made. Two teammates can disagree by accident and the bug looks like "this field is silently empty".

**Multiple tag vocabularies.** `json:"..."` for the wire, `yaml:"..."` for files, `mapstructure:"..."` for the generic decoder, sometimes `env:"..."` on top. This is not a Go failure; it is several audiences (API clients, operators, decoder) sharing one struct. Without a convention, you end up with fields that decode from one format and silently miss in another.

**Slices and maps in env.** Environment variables are strings. Lists become comma-separated strings, JSON-in-a-string, or repeated keys with an index suffix — every library picks a different convention. Pick one team-wide and document it; otherwise you spend Mondays explaining why `FEATURES=a,b,c` produced a `[]string{"a,b,c"}`.

**Weak typing.** `encoding/json` decodes every JSON number into `float64` when the destination is `any`; YAML 1.1 happily turns `no` into a boolean `false`; envs are always strings until somebody coerces them. mapstructure's `WeaklyTypedInput` and decode hooks paper over a lot of this — but only on the typed-struct side, not when you read a value out of a generic map.

**Pointer vs value.** Optional sub-trees are often `*Section`. Merge policy and decode policy must agree on whether a missing key sets the pointer to `nil`, leaves it at the previous value, or constructs an empty `&Section{}`. The case you forget is the one staging hits at 2 a.m.

**JSON-in-an-env-var.** It works until shell quoting eats a `"`, or your log aggregator truncates the line, or somebody rotates a secret and reformats it. It is a perfectly fine escape hatch and a terrible default.

If a library makes flat env parsing delightful but treats nested trees as second-class, it is not failing — it is optimising for a different shape. My services keep landing in the nested zone, which is why I keep coming back to the same pattern.

---

## The missing middle: "map in the middle"

Once you separate the three jobs above, the same shape keeps showing up:

```text
bytes → nested map[string]any → mapstructure → struct
```

Parsers give you the left arrow. Mapstructure gives you the right arrow. The middle arrow — canonical nested maps, deep merge semantics, key normalisation — is where every team I have worked with reinvents the same private `deepMerge` and the same `strings.NewReplacer("-", "_")` helper.

Most libraries do *some* version of the middle arrow internally; few document it as the primary mental model for *all* formats. Once you do, three things get easier at once:

- **Cross-format equivalence becomes one function.** If `sub-service` and `SUB_SERVICE` collapse to the same path component before decoding, you stop maintaining parallel tag systems for the same field.
- **Tests improve.** You can snapshot the merged map before struct decode and bisect "parser bug" vs "tag bug" without reading files from disk in every case.
- **CLI tooling becomes cheap.** A `convert` subcommand is just *parse → map → re-encode*. A `merge` subcommand is *parse N → fold → re-encode*. No second framework needed.

This is not a call to rewrite everything in `map[string]any`. It is a call to *name the pipeline* so teams can argue about merge rules and key normalisation in the same vocabulary they already use for JSON and YAML.

---

## Closing: own the layer you actually use

There is good tooling in Go. What I have not found is a **default** that a new teammate can assume without reading the wiki, because the right answer depends on which items on the checklist your service is centred on.

Two practical recommendations from years of doing this wrong first:

1. **Pick a center of gravity, then accept the edges.** If your service is twelve-factor and env-first, an env struct loader is the whole story; do not bolt on YAML "just in case". If your service is YAML-first with a few env overrides, Viper or Koanf is fine — but write down the env-key flattening rule on day one, not after the first incident.
2. **Make the intermediate map explicit.** If you keep hitting merge-and-normalise bugs across formats, stop hiding the `map[string]any` step inside ad-hoc helpers. A single deep-merge function and a single key normaliser, both unit-tested, are worth more than any new dependency.

Other ecosystems have made the same observation under their own names — Rust's [Figment](https://docs.rs/figment/latest/figment/), used by Rocket, treats providers and profiles as first-class. The tension between "one source of bytes" and "one typed config" is not language-specific. Go just distributes the answer across more packages.

### Where `go-config` fits (this repository)

[go-config](https://github.com/eslider/go-config) is my attempt to optimise for the middle arrow with symmetric codecs across `env`, `yaml`, `json`, and `ini`. Same `Codec` shape per format; deep merge for maps with last-write-wins on scalars; `LowerAlnum` key normaliser so cross-format paths line up before struct tags do the fine work; a small `envc` CLI for `convert / get / merge`.

The `env` package layers dotenv files and then **`WithCurrentEnvironment()`**, which reads the live process environment (`os.Environ()` at `Map` time) and merges it as the usual highest-priority source—so the same code path covers a laptop `.env` and variables injected into a Docker or Kubernetes container.

It is *an* answer to the pain described here, not a claim that other libraries are obsolete. For API details see the [README](https://github.com/eSlider/go-config/blob/main/README.md); for the architectural decisions behind it (what is in scope, what is not, what was deliberately broken from `v0.1.0`) see [docs/asr/README.md](https://github.com/eSlider/go-config/blob/main/docs/asr/README.md).

---

## References

- 12-factor — Config: <https://12factor.net/config>
- Awesome Go — Configuration: <https://awesome-go.com/configuration/>
- Viper: <https://github.com/spf13/viper>
- Viper issue #641 — nested env binding: <https://github.com/spf13/viper/issues/641>
- Viper issue #2001 — nested struct decode from env: <https://github.com/spf13/viper/issues/2001>
- Koanf: <https://github.com/knadh/koanf>
- caarlos0/env: <https://github.com/caarlos0/env>
- kelseyhightower/envconfig: <https://github.com/kelseyhightower/envconfig>
- cleanenv: <https://github.com/ilyakaznacheev/cleanenv>
- joho/godotenv: <https://github.com/joho/godotenv>
- go-viper/mapstructure (v2): <https://github.com/go-viper/mapstructure>
- encoding/json (stdlib): <https://pkg.go.dev/encoding/json>
- gopkg.in/yaml.v3: <https://github.com/go-yaml/yaml>
- Figment (Rust): <https://docs.rs/figment/latest/figment/>
- go-config — README: <https://github.com/eSlider/go-config/blob/main/README.md>
- go-config — ASR index: <https://github.com/eSlider/go-config/blob/main/docs/asr/README.md>
