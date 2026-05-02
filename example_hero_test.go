package config_test

import (
	"fmt"

	"github.com/eslider/go-config/env"
	"github.com/eslider/go-config/yaml"
)

func Example_hero_offline() {
	// README hero pattern using offline sources (no network, no .env files).
	y := yaml.New(yaml.WithBytes([]byte(`service: { name: demo }`)))
	e := env.New(env.WithBytes([]byte("SERVICE_PORT=8080")))
	var m map[string]any
	_ = y.Unmarshal(&m)
	var flat map[string]any
	_ = e.Unmarshal(&flat)
	fmt.Println(m["service"].(map[string]any)["name"], flat["service"].(map[string]any)["port"])
	// Output: demo 8080
}
