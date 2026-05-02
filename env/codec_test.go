package env

import (
	"context"
	"testing"

	"github.com/eslider/go-config/internal/testfixtures"
)

type identityRoot struct {
	Service svc `mapstructure:"service"`
}

type svc struct {
	Name       string
	Subservice sub `mapstructure:"subservice"`
	Database   db  `mapstructure:"database"`
}

type sub struct {
	Name    string
	Key     string
	Timeout int
	Enabled bool
}

type db struct {
	URL      string
	PoolSize int `mapstructure:"poolsize"`
}

func TestCodec_IdentityEnvFixture(t *testing.T) {
	c := New(WithBytes(testfixtures.Load(t, "identity", "service.env")))
	var got identityRoot
	if err := c.Unmarshal(&got); err != nil {
		t.Fatal(err)
	}
	assertIdentity(t, got)
}

func assertIdentity(t *testing.T, got identityRoot) {
	t.Helper()
	if got.Service.Name != "my-service" {
		t.Fatalf("Name=%q", got.Service.Name)
	}
	if got.Service.Subservice.Name != "abc" || got.Service.Subservice.Key != "abs" {
		t.Fatalf("sub %+v", got.Service.Subservice)
	}
	if got.Service.Subservice.Timeout != 30 || !got.Service.Subservice.Enabled {
		t.Fatalf("sub scalars %+v", got.Service.Subservice)
	}
	if got.Service.Database.URL != "postgres://localhost:5432/db" || got.Service.Database.PoolSize != 10 {
		t.Fatalf("db %+v", got.Service.Database)
	}
}

func TestCodec_MultiSourceLastWins(t *testing.T) {
	c := New(
		WithBytes(testfixtures.Load(t, "merge", "defaults.env")),
		WithBytes(testfixtures.Load(t, "merge", "overlay.env")),
	)
	m, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := m["service"].(map[string]any)["subservice"].(map[string]any)["name"]; got != "overlay-name" {
		t.Fatalf("subservice.name = %v", got)
	}
}

func TestCodec_WithCurrentEnvironment(t *testing.T) {
	t.Setenv("SERVICE_KEY", "abc")
	t.Setenv("SERVICE_PORT", "8080")

	var cfg struct {
		Service struct {
			Key  string `mapstructure:"key"`
			Port int    `mapstructure:"port"`
		} `mapstructure:"service"`
	}
	c := New(WithCurrentEnvironment())
	if err := c.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Service.Key != "abc" || cfg.Service.Port != 8080 {
		t.Fatalf("%+v", cfg)
	}
}

func TestCodec_WithTrimWeaklyTyped(t *testing.T) {
	t.Setenv("NAME", "  x  ")
	var cfg struct{ Name string }
	c := New(WithCurrentEnvironment(), WithTrim(true))
	if err := c.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "x" {
		t.Fatalf("%q", cfg.Name)
	}
}

func TestCodec_WithWeaklyTypedOff(t *testing.T) {
	t.Setenv("PORT", "8080")
	var cfg struct{ Port int }
	c := New(WithCurrentEnvironment(), WithWeaklyTyped(false))
	if err := c.Unmarshal(&cfg); err == nil {
		t.Fatal("expected error")
	}
}

func TestCodec_WithTagName(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")
	var cfg struct {
		DB struct {
			URL string `mapstructure:"url"`
		} `mapstructure:"database" env:"database"`
	}
	c := New(WithCurrentEnvironment(), WithTagName("env"))
	if err := c.Unmarshal(&cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.DB.URL != "postgres://x" {
		t.Fatalf("%q", cfg.DB.URL)
	}
}
