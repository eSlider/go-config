package config_test

import (
	"testing"

	"github.com/eslider/go-config/env"
	"github.com/eslider/go-config/yaml"
)

type motivationConfig struct {
	HTTP struct {
		Listen string `mapstructure:"listen"`
		TLS    bool   `mapstructure:"tls"`
	} `mapstructure:"http"`
	Database struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"database"`
}

func TestMotivationConfig_YAMLThenEnvThenProcessEnv(t *testing.T) {
	t.Setenv("HTTP_LISTEN", ":10000")
	t.Setenv("HTTP_TLS", "true")
	t.Setenv("DATABASE_URL", "postgres://process-env")

	yamlCfg := yaml.New(yaml.WithBytes([]byte(`
http:
  listen: ":8080"
  tls: false
database:
  url: "postgres://yaml"
`)))

	envCfg := env.New(
		env.WithBytes([]byte(`
HTTP_LISTEN=:9090
DATABASE_URL=postgres://dotenv
`)),
		env.WithCurrentEnvironment(),
	)

	var cfg motivationConfig
	if err := yamlCfg.Unmarshal(&cfg); err != nil {
		t.Fatalf("yaml unmarshal: %v", err)
	}
	if err := envCfg.Unmarshal(&cfg); err != nil {
		t.Fatalf("env unmarshal: %v", err)
	}

	if cfg.HTTP.Listen != ":10000" {
		t.Fatalf("HTTP.Listen=%q", cfg.HTTP.Listen)
	}
	if !cfg.HTTP.TLS {
		t.Fatalf("HTTP.TLS=%v", cfg.HTTP.TLS)
	}
	if cfg.Database.URL != "postgres://process-env" {
		t.Fatalf("Database.URL=%q", cfg.Database.URL)
	}
}
