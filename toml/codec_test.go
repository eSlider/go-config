package toml

import (
	"context"
	"testing"

	"github.com/eslider/go-config/internal/testfixtures"
)

func TestCodec_IdentityTOMLFixture(t *testing.T) {
	c := New(WithBytes(testfixtures.Load(t, "identity", "service.toml")))
	m, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if m["service"].(map[string]any)["name"] != "my-service" {
		t.Fatalf("%#v", m)
	}
}

func TestCodec_MarshalRoundTripMap(t *testing.T) {
	in := map[string]any{
		"app": map[string]any{
			"name": "demo",
			"port": int64(8080),
		},
	}
	c := New()
	b, err := c.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	c2 := New(WithBytes(b))
	m, err := c2.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if m["app"].(map[string]any)["name"] != "demo" {
		t.Fatalf("%#v", m)
	}
}
