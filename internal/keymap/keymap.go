// Package keymap normalizes map keys recursively for cross-format matching.
package keymap

import "strings"

// Normalizer transforms a single path segment or key string.
type Normalizer func(string) string

// LowerAlnum lowercases and strips any rune that is not [a-z0-9].
func LowerAlnum(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Identity returns s unchanged.
func Identity(s string) string { return s }

// Walk recursively rewrites every map key in place using n. Slices of maps are walked.
func Walk(m map[string]any, n Normalizer) {
	if m == nil || n == nil {
		return
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for _, k := range keys {
		v := m[k]
		delete(m, k)
		nk := n(k)
		if nk == "" {
			nk = k
		}
		switch vv := v.(type) {
		case map[string]any:
			Walk(vv, n)
			m[nk] = vv
		case []any:
			for i := range vv {
				if sm, ok := vv[i].(map[string]any); ok {
					Walk(sm, n)
					vv[i] = sm
				}
			}
			m[nk] = vv
		default:
			m[nk] = v
		}
	}
}
