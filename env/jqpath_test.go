package env

import (
	"reflect"
	"testing"
)

func TestNestedUnderJQPath(t *testing.T) {
	nested := map[string]any{
		"noise": map[string]any{"x": "ignore"},
		"service": map[string]any{
			"name": "my-service",
			"database": map[string]any{
				"url": "postgres://db",
			},
		},
	}
	got := nestedUnderJQPath(nested, ".service")
	want := map[string]any{
		"name": "my-service",
		"database": map[string]any{
			"url": "postgres://db",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestNestedUnderJQPath_EmptyPathNoop(t *testing.T) {
	n := map[string]any{"a": "1"}
	if got := nestedUnderJQPath(n, ""); !reflect.DeepEqual(got, n) {
		t.Fatalf("%#v", got)
	}
	if got := nestedUnderJQPath(n, "."); !reflect.DeepEqual(got, n) {
		t.Fatalf("%#v", got)
	}
}

func TestNestedUnderJQPath_MissingPath(t *testing.T) {
	got := nestedUnderJQPath(map[string]any{"other": "v"}, ".service")
	if len(got) != 0 {
		t.Fatalf("got %#v", got)
	}
}

func TestNestedUnderJQPath_PeelsRedundantServiceWrapper(t *testing.T) {
	nested := map[string]any{
		"service": map[string]any{
			"service": map[string]any{
				"name": "inner",
			},
		},
	}
	got := nestedUnderJQPath(nested, ".service")
	want := map[string]any{
		"name": "inner",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestNestedUnderJQPath_ScalarAtServiceIgnored(t *testing.T) {
	got := nestedUnderJQPath(map[string]any{"service": "x"}, ".service")
	if len(got) != 0 {
		t.Fatalf("got %#v", got)
	}
}
