package merge

import (
	"reflect"
	"testing"
)

func TestDeepMerge_MapRecursion(t *testing.T) {
	dst := map[string]any{
		"a": map[string]any{"x": "1"},
	}
	src := map[string]any{
		"a": map[string]any{"y": "2"},
	}
	DeepMerge(dst, src)
	want := map[string]any{
		"a": map[string]any{"x": "1", "y": "2"},
	}
	if !reflect.DeepEqual(dst, want) {
		t.Fatalf("got %#v want %#v", dst, want)
	}
}

func TestDeepMerge_ScalarLastWins(t *testing.T) {
	dst := map[string]any{"k": "first"}
	src := map[string]any{"k": "second"}
	DeepMerge(dst, src)
	if dst["k"] != "second" {
		t.Fatalf("got %v", dst["k"])
	}
}

func TestDeepMerge_SliceReplace(t *testing.T) {
	dst := map[string]any{"s": []any{"a", "b"}}
	src := map[string]any{"s": []any{"c"}}
	DeepMerge(dst, src, WithSliceStrategy(Replace))
	got := dst["s"].([]any)
	if len(got) != 1 || got[0] != "c" {
		t.Fatalf("got %#v", got)
	}
}

func TestDeepMerge_SliceConcat(t *testing.T) {
	dst := map[string]any{"s": []any{"a", "b"}}
	src := map[string]any{"s": []any{"c"}}
	DeepMerge(dst, src, WithSliceStrategy(Concat))
	got := dst["s"].([]any)
	want := []any{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}
