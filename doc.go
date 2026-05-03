// Package config is the module root for go-config (import path
// github.com/eslider/go-config). Use the subpackages [env], [yaml], [json], [toml],
// and [ini] for format-specific codecs and the [cmd/envc] command for CLI
// conversion.
//
// Hero workflow (YAML + layered dotenv + process env):
//
//	yamlCfg := yaml.New(yaml.WithURL("https://example.com/config.yaml"))
//	envCfg := env.New(
//	    env.WithFile(".default.env"),
//	    env.WithFile(".env"),
//	    env.WithCurrentEnvironment(),
//	)
//	var svc MyService
//	_ = yamlCfg.Unmarshal(&svc)
//	_ = envCfg.Unmarshal(&svc)
//
// [cmd/envc]: https://pkg.go.dev/github.com/eslider/go-config/cmd/envc
// [env]: https://pkg.go.dev/github.com/eslider/go-config/env
// [yaml]: https://pkg.go.dev/github.com/eslider/go-config/yaml
// [json]: https://pkg.go.dev/github.com/eslider/go-config/json
// [toml]: https://pkg.go.dev/github.com/eslider/go-config/toml
// [ini]: https://pkg.go.dev/github.com/eslider/go-config/ini
package config
