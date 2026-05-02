package keymap

import (
	"reflect"
	"testing"
)

func TestLowerAlnum(t *testing.T) {
	if got, want := LowerAlnum("sub-Service"), "subservice"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestWalk(t *testing.T) {
	m := map[string]any{
		"Sub-Service": map[string]any{"Pool-Size": "10"},
	}
	Walk(m, LowerAlnum)
	want := map[string]any{
		"subservice": map[string]any{"poolsize": "10"},
	}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("got %#v want %#v", m, want)
	}
}
