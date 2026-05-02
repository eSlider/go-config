// Package structconv wraps mapstructure for map<->struct conversion.
package structconv

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
)

// Options configures Decode and Encode.
type Options struct {
	TagName     string
	WeaklyTyped bool
	Trim        bool
	ExtraHooks  []mapstructure.DecodeHookFunc
}

func trimStringHook(from, to reflect.Type, data any) (any, error) {
	if from.Kind() == reflect.String && to.Kind() == reflect.String {
		if s, ok := data.(string); ok {
			return strings.TrimSpace(s), nil
		}
	}
	return data, nil
}

func composeHooks(o Options) mapstructure.DecodeHookFunc {
	var hooks []mapstructure.DecodeHookFunc
	if o.Trim {
		hooks = append(hooks, trimStringHook)
	}
	hooks = append(hooks, o.ExtraHooks...)
	if len(hooks) == 0 {
		return nil
	}
	return mapstructure.ComposeDecodeHookFunc(hooks...)
}

// Decode maps in into struct or map dst (dst must be pointer).
func Decode(in map[string]any, dst any, o Options) error {
	tag := o.TagName
	if tag == "" {
		tag = "mapstructure"
	}
	cfg := &mapstructure.DecoderConfig{
		Result:           dst,
		WeaklyTypedInput: o.WeaklyTyped,
		TagName:          tag,
		DecodeHook:       composeHooks(o),
	}
	dec, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return fmt.Errorf("structconv: decoder: %w", err)
	}
	if err := dec.Decode(in); err != nil {
		return fmt.Errorf("structconv: decode: %w", err)
	}
	return nil
}

// Encode flattens src (struct or map) into a new map[string]any.
func Encode(src any) (map[string]any, error) {
	out := make(map[string]any)
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{Result: &out})
	if err != nil {
		return nil, fmt.Errorf("structconv: encoder: %w", err)
	}
	if err := dec.Decode(src); err != nil {
		return nil, fmt.Errorf("structconv: encode: %w", err)
	}
	return out, nil
}
