package yaml

import (
	"context"
	"testing"

	"github.com/eslider/go-config/testfixtures"
)

func TestCodec_InvalidYAML(t *testing.T) {
	c := New(WithBytes(testfixtures.Load(t, "invalid", "malformed.yaml")))
	_, err := c.Map(context.Background())
	if err == nil {
		t.Fatal("expected parse error")
	}
}
