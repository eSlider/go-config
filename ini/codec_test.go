package ini

import (
	"context"
	"testing"

	"github.com/eslider/go-config/testfixtures"
)

func TestCodec_IdentityINIFixture(t *testing.T) {
	c := New(WithBytes(testfixtures.Load(t, "identity", "service.ini")))
	m, err := c.Map(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	svc := m["service"].(map[string]any)
	if svc["name"] != "my-service" {
		t.Fatalf("%#v", svc)
	}
}
