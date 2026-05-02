package env

import (
	"reflect"
	"testing"
)

func TestNestedFromFlat(t *testing.T) {
	got := nestedFromFlat(map[string]string{
		"FOO":               "bar",
		"SERVICE_HTTP_PORT": "8080",
		"SERVICE_KEY":       "abc",
	}, "")
	want := map[string]any{
		"foo": "bar",
		"service": map[string]any{
			"http": map[string]any{"port": "8080"},
			"key":  "abc",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestNestedFromFlat_FirstWriteWins(t *testing.T) {
	flat := map[string]string{}
	flat["A"] = "1"
	got := nestedFromFlat(flat, "")
	// Second insert for same key is ignored by insertPath within one build.
	insertPath(got, []string{"a"}, "2")
	if got["a"] != "1" {
		t.Fatalf("got %v", got["a"])
	}
}
