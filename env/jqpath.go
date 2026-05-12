package env

import (
	"strings"
)

// nestedUnderJQPath returns the value at jqPath inside nested as a sub-tree
// (just the inner fields — no enclosing key matching the path leaf). It
// implements a tiny subset of jq path selection: dotted segments matched on
// lowercased keys, consistent with nestedFromFlat / insertPath.
//
// Examples (jqPath = ".service"):
//
//	{"service": {"name": "x", "db": {...}}}              → {"name": "x", "db": {...}}
//	{"service": {"service": {"name": "inner"}}}          → {"name": "inner"}   (redundant
//	                                                       single-key wrapper repeating
//	                                                       the path leaf is peeled)
//	{"other": "v"}                                       → {}                   (missing)
//	{"service": "scalar"}                                → {}                   (scalar at
//	                                                       single-segment path is ignored)
//
// jqPath uses "." segments; a leading "." is optional. Empty / "." returns
// the input map untouched. The caller is responsible for re-wrapping the
// result at jqPath if it wants to merge the sub-tree at the same path in a
// larger document (see Codec.Map).
func nestedUnderJQPath(nested map[string]any, jqPath string) map[string]any {
	segs := splitJQPath(jqPath)
	if len(segs) == 0 {
		return nested
	}
	var cur any = nested
	for _, seg := range segs {
		m, ok := cur.(map[string]any)
		if !ok {
			return map[string]any{}
		}
		cur = m[seg]
		if cur == nil {
			return map[string]any{}
		}
	}
	leaf := segs[len(segs)-1]
	switch v := cur.(type) {
	case map[string]any:
		return peelRedundantPathLeaf(v, leaf)
	default:
		if len(segs) == 1 {
			return map[string]any{}
		}
		return buildNestedTreeAtPath(segs, cur)
	}
}

// splitJQPath parses a jq-style dotted path into lowercased segments. Empty
// path or "." yields a nil slice (interpreted as "no path").
func splitJQPath(jqPath string) []string {
	jqPath = strings.TrimSpace(jqPath)
	if jqPath == "" || jqPath == "." {
		return nil
	}
	jqPath = strings.TrimPrefix(jqPath, ".")
	if jqPath == "" {
		return nil
	}
	parts := strings.Split(jqPath, ".")
	segs := make([]string, 0, len(parts))
	for _, s := range parts {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		segs = append(segs, strings.ToLower(s))
	}
	if len(segs) == 0 {
		return nil
	}
	return segs
}

// peelRedundantPathLeaf strips outer {"<leaf>": {...}} shells while the only
// key in the current map matches pathLeaf. Stops as soon as the map has more
// than one key or the single key differs from pathLeaf.
func peelRedundantPathLeaf(m map[string]any, pathLeaf string) map[string]any {
	for len(m) == 1 {
		inner, ok := m[pathLeaf].(map[string]any)
		if !ok {
			break
		}
		m = inner
	}
	return m
}

// buildNestedTreeAtPath wraps leaf inside a chain of single-key maps following
// segments, e.g. (["a","b"], 1) → {"a": {"b": 1}}.
func buildNestedTreeAtPath(segments []string, leaf any) map[string]any {
	root := make(map[string]any)
	cur := root
	for i, seg := range segments {
		if i == len(segments)-1 {
			cur[seg] = leaf
			break
		}
		next := make(map[string]any)
		cur[seg] = next
		cur = next
	}
	return root
}
