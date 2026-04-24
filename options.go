package env

import "github.com/mitchellh/mapstructure"

// Option customises Unmarshal/UnmarshalPrefix behaviour.
type Option func(*options)

type options struct {
	trim        bool
	weaklyTyped bool
	tagName     string
	extraHooks  []mapstructure.DecodeHookFunc
}

func defaultOptions() options {
	return options{
		trim:        true,
		weaklyTyped: true,
		tagName:     "mapstructure",
	}
}

// WithTrim toggles automatic TrimSpace on string values. Default: true.
// Use WithTrim(false) when whitespace is semantically meaningful.
func WithTrim(enable bool) Option {
	return func(o *options) { o.trim = enable }
}

// WithWeaklyTyped toggles mapstructure's WeaklyTypedInput. Default: true.
// When true, "1" decodes into int, "true" into bool, etc. Disable for
// strict string-only decoding.
func WithWeaklyTyped(enable bool) Option {
	return func(o *options) { o.weaklyTyped = enable }
}

// WithTagName selects the struct tag mapstructure uses for field names.
// Default: "mapstructure". Set "env" to use `env:"FIELD_NAME"` tags.
func WithTagName(name string) Option {
	return func(o *options) { o.tagName = name }
}

// WithDecodeHook appends a user hook to the decoder chain. Hooks run after
// the built-in trim hook (unless trimming is disabled) and in the order
// they're added.
func WithDecodeHook(h mapstructure.DecodeHookFunc) Option {
	return func(o *options) { o.extraHooks = append(o.extraHooks, h) }
}
