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

	Service struct {
		Database struct {
			Poolsize string `json:"poolsize"`
			Url      string `json:"url"`
		} `json:"database"`
		Name       string `json:"name"`
		Subservice struct {
			Enabled string `json:"enabled"`
			Key     string `json:"key"`
			Name    string `json:"name"`
			Timeout string `json:"timeout"`
		} `json:"subservice"`
	} `json:"service"`
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
		env.WithFile("fixtures/identity/service.env", ".service"),
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

	js, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	var roundTrip map[string]interface{}
	if err := json.Unmarshal(js, &roundTrip); err != nil {
		t.Fatalf("json unmarshal: %v", err)
	}
	svc, ok := roundTrip["service"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected service in merged map, got %s", string(js))
	}
	if svc["name"] != "my-service" {
		t.Fatalf("service.name=%v", svc["name"])
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
	if cfg.Service.Name != "my-service" {
		t.Fatalf("Service.Name=%q", cfg.Service.Name)
	}
	if cfg.Service.Database.Url != "postgres://localhost:5432/db" {
		t.Fatalf("Service.Database.Url=%q", cfg.Service.Database.Url)
	}
	if cfg.Service.Database.Poolsize != "10" {
		t.Fatalf("Service.Database.Poolsize=%q", cfg.Service.Database.Poolsize)
	}
	if cfg.Service.Subservice.Name != "abc" {
		t.Fatalf("Service.Subservice.Name=%q", cfg.Service.Subservice.Name)
	}
}
