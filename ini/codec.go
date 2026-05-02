package ini

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/eslider/go-config/internal/bytesutil"
	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
	"github.com/eslider/go-config/internal/structconv"
	iniv1 "gopkg.in/ini.v1"
)

// Codec loads INI from multiple sources.
type Codec struct {
	sources    []source.Source
	normalizer keymap.Normalizer
	sliceStrat merge.SliceStrategy
	structOpts structconv.Options
	mergeOpts  []merge.Option
	httpClient *http.Client
	urlHeader  http.Header
}

// New creates an INI codec.
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

// Map merges all INI sources into one map.
func (c *Codec) Map(ctx context.Context) (map[string]any, error) {
	if len(c.sources) == 0 {
		return nil, fmt.Errorf("ini: no sources")
	}
	acc := make(map[string]any)
	for _, s := range c.sources {
		b, err := bytesutil.ReadAll(ctx, s)
		if err != nil {
			return nil, fmt.Errorf("ini: read %s: %w", s.String(), err)
		}
		f, err := iniv1.LoadSources(iniv1.LoadOptions{Loose: true, Insensitive: false}, b)
		if err != nil {
			return nil, fmt.Errorf("ini: parse %s: %w", s.String(), err)
		}
		parsed := mapFromINIFile(f)
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

// UnmarshalContext decodes merged INI into dst.
func (c *Codec) UnmarshalContext(ctx context.Context, dst any) error {
	m, err := c.Map(ctx)
	if err != nil {
		return err
	}
	return structconv.Decode(m, dst, c.structOpts)
}

// Marshal encodes src to INI bytes (struct or map[string]any).
func (c *Codec) Marshal(src any) ([]byte, error) {
	m, err := encodeToMap(src)
	if err != nil {
		return nil, fmt.Errorf("ini: encode: %w", err)
	}
	f := iniv1.Empty()
	if err := emitMap(nil, m, f); err != nil {
		return nil, err
	}
	var buf strings.Builder
	if _, err := f.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("ini: write: %w", err)
	}
	return []byte(buf.String()), nil
}

// WriteTo writes INI to w.
func (c *Codec) WriteTo(w io.Writer, src any) (int64, error) {
	b, err := c.Marshal(src)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

func emitMap(prefix []string, m map[string]any, f *iniv1.File) error {
	for k, v := range m {
		switch t := v.(type) {
		case map[string]any:
			if err := emitMap(append(prefix, k), t, f); err != nil {
				return err
			}
		case []any:
			parts := make([]string, 0, len(t))
			for _, e := range t {
				parts = append(parts, fmt.Sprint(e))
			}
			secName := strings.Join(prefix, ".")
			if secName == "" {
				secName = iniv1.DefaultSection
			}
			sec, err := f.GetSection(secName)
			if err != nil {
				sec, err = f.NewSection(secName)
				if err != nil {
					return err
				}
			}
			if _, err := sec.NewKey(k, strings.Join(parts, ",")); err != nil {
				return err
			}
		default:
			secName := strings.Join(prefix, ".")
			if secName == "" {
				secName = iniv1.DefaultSection
			}
			sec, err := f.GetSection(secName)
			if err != nil {
				sec, err = f.NewSection(secName)
				if err != nil {
					return err
				}
			}
			if _, err := sec.NewKey(k, fmt.Sprint(t)); err != nil {
				return err
			}
		}
	}
	return nil
}

func encodeToMap(src any) (map[string]any, error) {
	if m, ok := src.(map[string]any); ok {
		return m, nil
	}
	return structconv.Encode(src)
}
