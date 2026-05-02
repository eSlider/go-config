// Package testfixtures resolves paths to files under the module's fixtures/ directory.
package testfixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	dir := filepath.Dir(file)
	for {
		data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil && strings.Contains(string(data), "module github.com/eslider/go-config") {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("module root not found")
		}
		dir = parent
	}
}

// Path joins repo fixtures/ with parts.
func Path(t *testing.T, parts ...string) string {
	t.Helper()
	return filepath.Join(append([]string{moduleRoot(t), "fixtures"}, parts...)...)
}

// Load reads a fixture file relative to fixtures/.
func Load(t *testing.T, parts ...string) []byte {
	t.Helper()
	b, err := os.ReadFile(Path(t, parts...))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
