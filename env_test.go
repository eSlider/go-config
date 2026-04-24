package env

import (
	"reflect"
	"testing"
)

// asMapFromEnviron is tested with a synthetic environ slice — this is NOT a
// mock of anything external. It's a pure-input unit test that isolates the
// core mapping logic from process-global state.

func TestAsMapFromEnviron_FlatAndNested(t *testing.T) {
	got := asMapFromEnviron([]string{
		"FOO=bar",
		"SERVICE_HTTP_PORT=8080",
		"SERVICE_KEY=abc",
	}, "")

	want := map[string]any{
		"foo": "bar",
		"service": map[string]any{
			"http": map[string]any{"port": "8080"},
			"key":  "abc",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestAsMapFromEnviron_PreservesValuesWithEquals(t *testing.T) {
	got := asMapFromEnviron([]string{"TOKEN=a=b=c="}, "")
	if got["token"] != "a=b=c=" {
		t.Fatalf("expected full value preserved, got %#v", got["token"])
	}
}

func TestAsMapFromEnviron_PrefixFilterAndStrip(t *testing.T) {
	got := asMapFromEnviron([]string{
		"APP_DB_HOST=localhost",
		"APP_DB_PORT=5432",
		"OTHER_THING=skip",
	}, "APP_")

	want := map[string]any{
		"db": map[string]any{"host": "localhost", "port": "5432"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestAsMapFromEnviron_EmptyAndMalformedAreSkipped(t *testing.T) {
	got := asMapFromEnviron([]string{"NOEQUALS", ""}, "")
	if len(got) != 0 {
		t.Fatalf("expected no entries, got %#v", got)
	}
}

func TestAsMapFromEnviron_FirstWriteWinsOnCollision(t *testing.T) {
	// This matches the ai-fabric original: once a leaf is set for a path,
	// subsequent entries that would overwrite it are ignored.
	got := asMapFromEnviron([]string{
		"A=1",
		"A=2",
	}, "")
	if got["a"] != "1" {
		t.Fatalf("expected first-write-wins (a=1), got %#v", got["a"])
	}
}

func TestAsMapFromEnviron_PrefixOnlyVariableIsSkipped(t *testing.T) {
	// "APP_=xxx" with prefix "APP_" would strip to empty key — we skip it.
	got := asMapFromEnviron([]string{"APP_=oops"}, "APP_")
	if len(got) != 0 {
		t.Fatalf("expected empty map, got %#v", got)
	}
}

// ---------------------------------------------------------------------------
// Decoder / option tests
// ---------------------------------------------------------------------------

func TestUnmarshal_BasicStruct(t *testing.T) {
	t.Setenv("SERVICE_KEY", "abc")
	t.Setenv("SERVICE_PORT", "8080")

	var cfg struct {
		Service struct {
			Key  string
			Port int
		}
	}
	if err := Unmarshal(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Service.Key != "abc" {
		t.Errorf("Key = %q, want abc", cfg.Service.Key)
	}
	if cfg.Service.Port != 8080 {
		t.Errorf("Port = %d, want 8080 (weakly-typed string->int)", cfg.Service.Port)
	}
}

func TestUnmarshal_TrimsStringsByDefault(t *testing.T) {
	t.Setenv("NAME", "   hello   ")

	var cfg struct{ Name string }
	if err := Unmarshal(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "hello" {
		t.Fatalf("expected trimmed value, got %q", cfg.Name)
	}
}

func TestUnmarshal_WithTrimOff(t *testing.T) {
	t.Setenv("NAME", "  keep me  ")

	var cfg struct{ Name string }
	if err := Unmarshal(&cfg, WithTrim(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Name != "  keep me  " {
		t.Fatalf("expected untrimmed value, got %q", cfg.Name)
	}
}

func TestUnmarshalPrefix_StripsAndFilters(t *testing.T) {
	t.Setenv("APP_DB_HOST", "localhost")
	t.Setenv("APP_DB_PORT", "5432")
	t.Setenv("OTHER_DB_HOST", "should-be-ignored")

	var cfg struct {
		DB struct {
			Host string
			Port int
		} `mapstructure:"db"`
	}
	if err := UnmarshalPrefix(&cfg, "APP_"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DB.Host != "localhost" {
		t.Errorf("Host = %q, want localhost", cfg.DB.Host)
	}
	if cfg.DB.Port != 5432 {
		t.Errorf("Port = %d, want 5432", cfg.DB.Port)
	}
}

func TestUnmarshal_WithTagName(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://x")

	// The map is always nested by '_', so tags also reference the nested
	// path, one component per struct level. Here we override the outer
	// field name via a custom tag.
	var cfg struct {
		DB struct {
			URL string
		} `env:"database"`
	}
	if err := Unmarshal(&cfg, WithTagName("env")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DB.URL != "postgres://x" {
		t.Fatalf("URL = %q, want postgres://x", cfg.DB.URL)
	}
}

func TestUnmarshal_WeaklyTypedOff(t *testing.T) {
	t.Setenv("PORT", "8080")

	var cfg struct{ Port int }
	err := Unmarshal(&cfg, WithWeaklyTyped(false))
	// With strict typing, "8080" (string) cannot decode into an int.
	if err == nil {
		t.Fatal("expected strict typing to reject string->int coercion")
	}
}

func TestUnmarshal_SpecialCharactersSurvive(t *testing.T) {
	// Regression: keys with special chars in their values must pass through
	// unchanged (modulo TrimSpace).
	t.Setenv("SPECIAL", "@!#$%^&*()")

	var cfg struct{ Special string }
	if err := Unmarshal(&cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Special != "@!#$%^&*()" {
		t.Fatalf("Special = %q, want @!#$%%^&*()", cfg.Special)
	}
}

func TestUnmarshal_RejectsNonPointer(t *testing.T) {
	var cfg struct{ A string }
	err := Unmarshal(cfg) // no pointer
	if err == nil {
		t.Fatal("expected error when passing a non-pointer")
	}
}

func TestAsMap_ExposesCurrentEnv(t *testing.T) {
	t.Setenv("GO_ENV_ASMAP_PROBE", "hit")

	m := AsMap()
	goenv, ok := m["go"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map under 'go', got %T", m["go"])
	}
	env, ok := goenv["env"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map under 'go.env', got %T", goenv["env"])
	}
	asmap, ok := env["asmap"].(map[string]any)
	if !ok {
		t.Fatalf("expected nested map under 'go.env.asmap', got %T", env["asmap"])
	}
	if asmap["probe"] != "hit" {
		t.Fatalf("probe = %v, want hit", asmap["probe"])
	}
}
