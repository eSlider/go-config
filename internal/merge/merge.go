// Package merge provides deep map merging with configurable slice behaviour.
package merge

// SliceStrategy controls how two []any values are combined at the same key.
type SliceStrategy int

const (
	// Replace overwrites the destination slice with the source slice.
	Replace SliceStrategy = iota
	// Concat appends source elements to the destination slice.
	Concat
)

// Option configures DeepMerge.
type Option func(*config)

type config struct {
	slice SliceStrategy
}

// WithSliceStrategy sets slice merge behaviour (default Replace).
func WithSliceStrategy(s SliceStrategy) Option {
	return func(c *config) { c.slice = s }
}

// DeepMerge folds src into dst in place. Maps recurse; scalar leaves use
// last-write-wins (src overwrites dst). For slice values, behaviour follows cfg.slice.
func DeepMerge(dst, src map[string]any, opts ...Option) {
	cfg := config{slice: Replace}
	for _, o := range opts {
		o(&cfg)
	}
	if dst == nil || src == nil {
		return
	}
	for k, sv := range src {
		dv, ok := dst[k]
		if !ok {
			dst[k] = cloneValue(sv)
			continue
		}
		dm, dIsMap := dv.(map[string]any)
		sm, sIsMap := sv.(map[string]any)
		if dIsMap && sIsMap {
			DeepMerge(dm, sm, opts...)
			dst[k] = dm
			continue
		}
		if dsa, dIsSlice := dv.([]any); dIsSlice {
			if ssa, sIsSlice := sv.([]any); sIsSlice {
				switch cfg.slice {
				case Concat:
					dst[k] = append(append([]any{}, dsa...), ssa...)
				default:
					dst[k] = append([]any{}, ssa...)
				}
				continue
			}
		}
		dst[k] = cloneValue(sv)
	}
}

func cloneValue(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		DeepMerge(out, t)
		return out
	case []any:
		cp := make([]any, len(t))
		copy(cp, t)
		return cp
	default:
		return v
	}
}
