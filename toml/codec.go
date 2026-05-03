package toml

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/eslider/go-config/internal/bytesutil"
	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
	"github.com/eslider/go-config/internal/structconv"
	toml2 "github.com/pelletier/go-toml/v2"
)

// Codec loads TOML from multiple sources.
type Codec struct {
	sources    []source.Source
	normalizer keymap.Normalizer
	sliceStrat merge.SliceStrategy
	structOpts structconv.Options
	mergeOpts  []merge.Option
	httpClient *http.Client
	urlHeader  http.Header
}

// New creates a TOML codec.
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

// Map merges all TOML sources into one map.
func (c *Codec) Map(ctx context.Context) (map[string]any, error) {
	if len(c.sources) == 0 {
		return nil, fmt.Errorf("toml: no sources")
	}
	acc := make(map[string]any)
	for _, s := range c.sources {
		b, err := bytesutil.ReadAll(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("toml: read %s: %w", s.String(), err)
		}
		var parsed map[string]any
		if err := toml2.Unmarshal(b, &parsed); err != nil {
			return nil, fmt.Errorf("toml: parse %s: %w", s.String(), err)
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

// UnmarshalContext decodes merged TOML into dst.
func (c *Codec) UnmarshalContext(ctx context.Context, dst any) error {
	m, err := c.Map(ctx)
	if err != nil {
		return err
	}
	return structconv.Decode(m, dst, c.structOpts)
}

// Marshal encodes src to TOML bytes (struct or map[string]any).
func (c *Codec) Marshal(src any) ([]byte, error) {
	m, err := encodeToMap(src)
	if err != nil {
		return nil, fmt.Errorf("toml: encode: %w", err)
	}
	var buf bytes.Buffer
	enc := toml2.NewEncoder(&buf)
	enc.SetIndentTables(true)
	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("toml: marshal: %w", err)
	}
	return buf.Bytes(), nil
}

// WriteTo writes TOML to w.
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
