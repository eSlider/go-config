package env

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/eslider/go-config/internal/bytesutil"
	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
	"github.com/eslider/go-config/internal/structconv"
	"github.com/joho/godotenv"
)

type loadedLayer struct {
	load   func(context.Context) (map[string]string, error)
	jqPath string // empty: merge full nested tree; e.g. ".service" merges only that subtree
}

// Codec loads environment-style key/value data from multiple sources.
type Codec struct {
	layers     []loadedLayer
	prefix     string
	normalizer keymap.Normalizer
	sliceStrat merge.SliceStrategy
	structOpts structconv.Options
	mergeOpts  []merge.Option
	httpClient *http.Client
	urlHeader  http.Header
}

// New builds a Codec from options. Sources are merged left-to-right; later
// sources override earlier scalar leaves (see package merge).
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

// Map returns the merged nested map after all sources are loaded and keys normalized.
func (c *Codec) Map(ctx context.Context) (map[string]any, error) {
	if len(c.layers) == 0 {
		return nil, fmt.Errorf("env: no sources configured")
	}
	acc := make(map[string]any)
	for _, layer := range c.layers {
		flat, err := layer.load(ctx)
		if err != nil {
			return nil, err
		}
		nested := nestedFromFlat(flat, c.prefix)
		if segs := splitJQPath(layer.jqPath); len(segs) > 0 {
			sub := nestedUnderJQPath(nested, layer.jqPath)
			nested = buildNestedTreeAtPath(segs, sub)
		}
		merge.DeepMerge(acc, nested, c.mergeOpts...)
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

// UnmarshalContext decodes the merged map into dst using mapstructure.
func (c *Codec) UnmarshalContext(ctx context.Context, dst any) error {
	m, err := c.Map(ctx)
	if err != nil {
		return err
	}
	return structconv.Decode(m, dst, c.structOpts)
}

// Marshal encodes src into dotenv-style bytes (KEY=value lines, sorted keys).
// src may be a struct or map[string]any.
func (c *Codec) Marshal(src any) ([]byte, error) {
	m, err := encodeToMap(src)
	if err != nil {
		return nil, fmt.Errorf("env: marshal encode: %w", err)
	}
	flat := flattenMap(nil, m)
	var buf bytes.Buffer
	keys := make([]string, 0, len(flat))
	for k := range flat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(&buf, "%s=%s\n", k, flat[k])
	}
	return buf.Bytes(), nil
}

// WriteTo writes dotenv-style output to w.
func (c *Codec) WriteTo(w io.Writer, src any) (int64, error) {
	b, err := c.Marshal(src)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}

func flattenMap(prefix []string, m map[string]any) map[string]string {
	out := make(map[string]string)
	for k, v := range m {
		path := append(prefix, k)
		switch t := v.(type) {
		case map[string]any:
			for kk, vv := range flattenMap(path, t) {
				out[kk] = vv
			}
		case []any:
			parts := make([]string, 0, len(t))
			for _, e := range t {
				parts = append(parts, fmt.Sprint(e))
			}
			key := strings.Join(path, "_")
			out[strings.ToUpper(key)] = strings.Join(parts, ",")
		default:
			key := strings.Join(path, "_")
			out[strings.ToUpper(key)] = fmt.Sprint(t)
		}
	}
	return out
}

func encodeToMap(src any) (map[string]any, error) {
	if m, ok := src.(map[string]any); ok {
		return m, nil
	}
	return structconv.Encode(src)
}

func flatFromEnviron(environ []string) map[string]string {
	out := make(map[string]string)
	for _, entry := range environ {
		eq := strings.IndexByte(entry, '=')
		if eq < 0 {
			continue
		}
		out[entry[:eq]] = entry[eq+1:]
	}
	return out
}

func withSource(s source.Source, label string) func(*Codec) {
	return withSourceAtJQ(s, label, "")
}

func withSourceAtJQ(s source.Source, label, jqPath string) func(*Codec) {
	return func(c *Codec) {
		c.layers = append(c.layers, loadedLayer{
			jqPath: jqPath,
			load: func(ctx context.Context) (map[string]string, error) {
				b, err := bytesutil.ReadAll(ctx, s)
				if err != nil {
					return nil, fmt.Errorf("env: read %s: %w", label, err)
				}
				m, err := godotenv.UnmarshalBytes(b)
				if err != nil {
					return nil, fmt.Errorf("env: parse %s: %w", label, err)
				}
				return m, nil
			},
		})
	}
}
