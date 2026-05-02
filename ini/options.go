package ini

import (
	"io"
	"net/http"

	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
	"github.com/go-viper/mapstructure/v2"
)

// Option configures a Codec.
type Option func(*Codec)

// WithBytes appends INI from bytes.
func WithBytes(b []byte) Option {
	return func(c *Codec) { c.sources = append(c.sources, source.Bytes{Data: b, Name: "bytes"}) }
}

// WithReader appends INI from r.
func WithReader(r io.Reader) Option {
	return func(c *Codec) { c.sources = append(c.sources, source.Reader{R: r, Name: "reader"}) }
}

// WithFile appends an INI file path.
func WithFile(path string) Option {
	return func(c *Codec) { c.sources = append(c.sources, source.File{Path: path}) }
}

// WithURL appends INI from an HTTP(S) URL.
func WithURL(raw string) Option {
	return func(c *Codec) {
		var hdr http.Header
		if c.urlHeader != nil {
			hdr = c.urlHeader.Clone()
		}
		c.sources = append(c.sources, source.URL{Raw: raw, Header: hdr, HTTPClient: c.httpClient})
	}
}

// WithHTTPClient sets the client for subsequent WithURL sources.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Codec) { c.httpClient = client }
}

// WithHTTPHeader adds a header for subsequent WithURL sources.
func WithHTTPHeader(k, v string) Option {
	return func(c *Codec) {
		if c.urlHeader == nil {
			c.urlHeader = make(http.Header)
		}
		c.urlHeader.Add(k, v)
	}
}

// WithKeyNormalizer sets key normalizer after merge (nil disables).
func WithKeyNormalizer(n keymap.Normalizer) Option {
	return func(c *Codec) { c.normalizer = n }
}

// WithSliceMerge sets slice merge strategy across sources.
func WithSliceMerge(s merge.SliceStrategy) Option {
	return func(c *Codec) { c.sliceStrat = s }
}

// WithTrim toggles string trim on struct decode.
func WithTrim(enable bool) Option {
	return func(c *Codec) { c.structOpts.Trim = enable }
}

// WithWeaklyTyped toggles weak typing for struct decode.
func WithWeaklyTyped(enable bool) Option {
	return func(c *Codec) { c.structOpts.WeaklyTyped = enable }
}

// WithTagName sets mapstructure tag name.
func WithTagName(name string) Option {
	return func(c *Codec) { c.structOpts.TagName = name }
}

// WithDecodeHook appends a decode hook.
func WithDecodeHook(h mapstructure.DecodeHookFunc) Option {
	return func(c *Codec) {
		c.structOpts.ExtraHooks = append(c.structOpts.ExtraHooks, h)
	}
}
