package structconv

import (
	"fmt"
	"testing"
)

func TestDecodeEncodeRoundTrip(t *testing.T) {
	type nested struct {
		V int `mapstructure:"v"`
	}
	type root struct {
		A nested `mapstructure:"a"`
	}
	in := map[string]any{"a": map[string]any{"v": 7}}
	var dst root
	if err := Decode(in, &dst, Options{WeaklyTyped: true, Trim: true}); err != nil {
		t.Fatal(err)
	}
	out, err := Encode(&dst)
	if err != nil {
		t.Fatal(err)
	}
	am := out["a"].(map[string]any)
	if fmt.Sprint(am["v"]) != "7" {
		t.Fatalf("%#v", out)
	}
}
