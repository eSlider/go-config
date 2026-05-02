package env

import (
	"strings"
)

// insertPath walks path inside root, creating intermediate nested maps as
// needed, and stores value at the leaf. First write wins for a given leaf key
// within a single source map.
func insertPath(root map[string]any, path []string, value string) {
	current := root
	last := len(path) - 1
	for i, raw := range path {
		k := strings.ToLower(raw)

		if i == last {
			if _, exists := current[k]; exists {
				return
			}
			current[k] = value
			return
		}

		next, ok := current[k].(map[string]any)
		if !ok {
			next = make(map[string]any)
			current[k] = next
		}
		current = next
	}
}

func nestedFromFlat(flat map[string]string, prefix string) map[string]any {
	out := make(map[string]any)
	for k, v := range flat {
		name := k
		if prefix != "" {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			name = name[len(prefix):]
			if name == "" {
				continue
			}
		}
		insertPath(out, strings.Split(name, "_"), v)
	}
	return out
}
