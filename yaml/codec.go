package yaml

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/eslider/go-config/internal/bytesutil"
	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
	"github.com/eslider/go-config/internal/structconv"
	yaml3 "gopkg.in/yaml.v3"
)

// Codec loads YAML from multiple sources.
type Codec struct {
	sources    []source.Source
	normalizer keymap.Normalizer
	sliceStrat merge.SliceStrategy
	structOpts structconv.Options
	mergeOpts  []merge.Option
	httpClient *http.Client
	urlHeader  http.Header
}

// New creates a YAML codec.
func New(opts ...Option) *Codec {
	c := &Codec{
		normalizer: keymap.LowerAlnum,
		sliceStrat: merge.Replace,
		structOpts: structconv.Options{
			TagName:     "mapstructure",
			WeaklyTyped: true,
			Trim:        true,
		},
	}
	for _, o := range opts {
		o(c)
	}
	c.mergeOpts = []merge.Option{merge.WithSliceStrategy(c.sliceStrat)}
	return c
}

// Map merges all YAML sources into one map.
func (c *Codec) Map(ctx context.Context) (map[string]any, error) {
	if len(c.sources) == 0 {
		return nil, fmt.Errorf("yaml: no sources")
	}
	acc := make(map[string]any)
	for _, s := range c.sources {
		b, err := bytesutil.ReadAll(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("yaml: read %s: %w", s.String(), err)
		}
		var parsed map[string]any
		if err := yaml3.Unmarshal(b, &parsed); err != nil {
			return nil, fmt.Errorf("yaml: parse %s: %w", s.String(), err)
		}
		if parsed == nil {
			parsed = map[string]any{}
		}
		merge.DeepMerge(acc, parsed, c.mergeOpts...)
	}
	if c.normalizer != nil {
		keymap.Walk(acc, c.normalizer)
	}
	return acc, nil
}

// Unmarshal decodes using context.Background.
func (c *Codec) Unmarshal(dst any) error {
	return c.UnmarshalContext(context.Background(), dst)
}

// UnmarshalContext decodes merged YAML into dst.
func (c *Codec) UnmarshalContext(ctx context.Context, dst any) error {
	m, err := c.Map(ctx)
	if err != nil {
		return err
	}
	return structconv.Decode(m, dst, c.structOpts)
}

// Marshal encodes src to YAML bytes (struct or map[string]any).
func (c *Codec) Marshal(src any) ([]byte, error) {
	m, err := encodeToMap(src)
	if err != nil {
		return nil, fmt.Errorf("yaml: encode: %w", err)
	}
	b, err := yaml3.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("yaml: marshal: %w", err)
	}
	return b, nil
}

// WriteTo writes YAML to w.
func (c *Codec) WriteTo(w io.Writer, src any) (int64, error) {
	b, err := c.Marshal(src)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

func encodeToMap(src any) (map[string]any, error) {
	if m, ok := src.(map[string]any); ok {
		return m, nil
	}
	return structconv.Encode(src)
}
