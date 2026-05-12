package config_test

import (
	"context"
	"encoding/json"
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

	ctx := context.Background()
	m, err := envCfg.Map(ctx)
	if err != nil {
		t.Fatalf("env marshal: %v", err)
	}
	if m == nil {
		t.Fatalf("env marshal returned nil map")
	}

	// convert m to json and back to map
	js, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	m2 := make(map[string]interface{})
	if err := json.Unmarshal(js, &m2); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}

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
