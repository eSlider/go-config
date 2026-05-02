package yaml

import (
	"context"
	"reflect"
	"testing"

	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/testfixtures"
	yaml3 "gopkg.in/yaml.v3"
)

func TestCodec_IdentityFixture(t *testing.T) {
	c := New(WithBytes(testfixtures.Load(t, "identity", "service.yaml")))
	m, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	exp, err := loadExpectedMap(t)
	if err != nil {
		t.Fatal(err)
	}
	keymap.Walk(exp, keymap.LowerAlnum)
	if !mapsEqual(m, exp) {
		t.Fatalf("map mismatch\n got: %#v\n exp: %#v", m, exp)
	}
}

func TestCodec_MergeReplaceSlices(t *testing.T) {
	c := New(
		WithBytes(testfixtures.Load(t, "merge", "defaults.yaml")),
		WithBytes(testfixtures.Load(t, "merge", "overlay.yaml")),
	)
	got, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	exp := map[string]any{}
	if err := yaml3.Unmarshal(testfixtures.Load(t, "merge", "expected-replace.yaml"), &exp); err != nil {
		t.Fatal(err)
	}
	keymap.Walk(exp, keymap.LowerAlnum)
	if !mapsEqual(got, exp) {
		t.Fatalf("merge mismatch\n got: %#v\n exp: %#v", got, exp)
	}
}

func TestCodec_MergeConcatSlices(t *testing.T) {
	c := New(
		WithBytes(testfixtures.Load(t, "merge", "defaults.yaml")),
		WithBytes(testfixtures.Load(t, "merge", "overlay.yaml")),
		WithSliceMerge(merge.Concat),
	)
	got, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	exp := map[string]any{}
	if err := yaml3.Unmarshal(testfixtures.Load(t, "merge", "expected-concat.yaml"), &exp); err != nil {
		t.Fatal(err)
	}
	keymap.Walk(exp, keymap.LowerAlnum)
	if !mapsEqual(got, exp) {
		t.Fatalf("concat mismatch\n got: %#v\n exp: %#v", got, exp)
	}
}

func loadExpectedMap(t *testing.T) (map[string]any, error) {
	t.Helper()
	var m map[string]any
	err := yaml3.Unmarshal(testfixtures.Load(t, "identity", "service.yaml"), &m)
	return m, err
}

func mapsEqual(a, b map[string]any) bool {
	return reflect.DeepEqual(a, b)
}
