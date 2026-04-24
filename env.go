package env

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// Unmarshal decodes all process environment variables into dst. The variable
// name is split by '_' and the resulting path is matched against dst's fields
// case-insensitively.
//
// dst must be a pointer to a struct (or a map that mapstructure can populate).
// See the package documentation for details and examples.
func Unmarshal(dst any, opts ...Option) error {
	return UnmarshalPrefix(dst, "", opts...)
}

// UnmarshalPrefix is like Unmarshal but only considers environment variables
// that start with prefix. The prefix is stripped from each variable name
// before the path is built.
//
// An empty prefix is equivalent to Unmarshal.
func UnmarshalPrefix(dst any, prefix string, opts ...Option) error {
	o := defaultOptions()
	for _, f := range opts {
		f(&o)
	}

	data := asMapFromEnviron(os.Environ(), prefix)

	hooks := []mapstructure.DecodeHookFunc{}
	if o.trim {
		hooks = append(hooks, trimStringHook)
	}
	hooks = append(hooks, o.extraHooks...)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:           dst,
		WeaklyTypedInput: o.weaklyTyped,
		TagName:          o.tagName,
		DecodeHook:       mapstructure.ComposeDecodeHookFunc(hooks...),
	})
	if err != nil {
		return fmt.Errorf("env: configure decoder: %w", err)
	}
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("env: decode: %w", err)
	}
	return nil
}

// AsMap returns a nested map built from all current process environment
// variables, using '_' as path separator. Keys are lower-cased.
//
// It's primarily useful for debugging or for callers that want to plug
// their own decoder; normal code should prefer Unmarshal.
func AsMap() map[string]any {
	return asMapFromEnviron(os.Environ(), "")
}

// AsMapPrefix is the prefix-scoped variant of AsMap. Variables whose names
// don't start with prefix are skipped; the prefix itself is stripped before
// the map is built.
func AsMapPrefix(prefix string) map[string]any {
	return asMapFromEnviron(os.Environ(), prefix)
}

// asMapFromEnviron is the core building block. Exposed as an unexported
// function so tests can feed it a deterministic environ slice instead of
// mutating the real process env.
func asMapFromEnviron(environ []string, prefix string) map[string]any {
	out := make(map[string]any)

	for _, entry := range environ {
		eq := strings.IndexByte(entry, '=')
		if eq < 0 {
			// Malformed entry — shouldn't happen on Unix but guard anyway.
			continue
		}
		name, value := entry[:eq], entry[eq+1:]

		if prefix != "" {
			if !strings.HasPrefix(name, prefix) {
				continue
			}
			name = name[len(prefix):]
			if name == "" {
				continue
			}
		}

		insertPath(out, strings.Split(name, "_"), value)
	}

	return out
}

// insertPath walks path inside root, creating intermediate nested maps as
// needed, and stores value at the leaf.
//
// Collisions are resolved by preferring the first write: once a leaf value
// is set for a path, a later sibling with the same path prefix will NOT
// overwrite it. This matches the behaviour of the ai-fabric original so
// that tools migrating to go-env don't observe new surprises.
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

func trimStringHook(from, to reflect.Type, data any) (any, error) {
	if from.Kind() == reflect.String && to.Kind() == reflect.String {
		if s, ok := data.(string); ok {
			return strings.TrimSpace(s), nil
		}
	}
	return data, nil
}
