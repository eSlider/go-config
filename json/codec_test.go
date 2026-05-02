package json

import (
	"context"
	"testing"

	"github.com/eslider/go-config/internal/testfixtures"
)

func TestCodec_IdentityJSONFixture(t *testing.T) {
	c := New(WithBytes(testfixtures.Load(t, "identity", "service.json")))
	m, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if m["service"].(map[string]any)["name"] != "my-service" {
		t.Fatalf("%#v", m)
	}
}
