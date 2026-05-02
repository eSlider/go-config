package env

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/eslider/go-config/internal/keymap"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
	"github.com/go-viper/mapstructure/v2"
)

// Option configures a Codec.
type Option func(*Codec)

// WithCurrentEnvironment appends the process environment as a source (read at Map time).
func WithCurrentEnvironment() Option {
	return func(c *Codec) {
		c.layers = append(c.layers, func(_ context.Context) (map[string]string, error) {
			return flatFromEnviron(os.Environ()), nil
		})
	}
}

// WithFile appends a dotenv file path as a source.
func WithFile(path string) Option {
	return withSource(source.File{Path: path}, path)
}

// WithBytes appends raw dotenv bytes as a source.
func WithBytes(b []byte) Option {
	return withSource(source.Bytes{Data: b, Name: "bytes"}, "bytes")
}

// WithReader appends dotenv content from r (read fully on each Map call).
func WithReader(r io.Reader) Option {
	return withSource(source.Reader{R: r, Name: "reader"}, "reader")
}

// WithURL appends an HTTP(S) URL returning dotenv content. Use WithHTTPHeader /
// WithHTTPClient before WithURL so they apply to this request.
func WithURL(raw string) Option {
	return func(c *Codec) {
		var hdr http.Header
		if c.urlHeader != nil {
			hdr = c.urlHeader.Clone()
		}
		s := source.URL{Raw: raw, Header: hdr, HTTPClient: c.httpClient}
		withSource(s, raw)(c)
	}
}

// WithHTTPClient sets the HTTP client used by subsequent WithURL sources.
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

// WithPrefix strips prefix from variable names before path splitting.
func WithPrefix(prefix string) Option {
	return func(c *Codec) { c.prefix = prefix }
}

// WithKeyNormalizer sets key normalizer after merge (nil disables).
func WithKeyNormalizer(n keymap.Normalizer) Option {
	return func(c *Codec) { c.normalizer = n }
}

// WithSliceMerge sets slice merge strategy when merging sources.
func WithSliceMerge(s merge.SliceStrategy) Option {
	return func(c *Codec) { c.sliceStrat = s }
}

// WithTrim toggles struct string trim via mapstructure hook.
func WithTrim(enable bool) Option {
	return func(c *Codec) { c.structOpts.Trim = enable }
}

// WithWeaklyTyped toggles weak typing for struct decode.
func WithWeaklyTyped(enable bool) Option {
	return func(c *Codec) { c.structOpts.WeaklyTyped = enable }
}

// WithTagName sets the struct tag name for mapstructure.
func WithTagName(name string) Option {
	return func(c *Codec) { c.structOpts.TagName = name }
}

// WithDecodeHook appends a decode hook.
func WithDecodeHook(h mapstructure.DecodeHookFunc) Option {
	return func(c *Codec) {
		c.structOpts.ExtraHooks = append(c.structOpts.ExtraHooks, h)
	}
}
