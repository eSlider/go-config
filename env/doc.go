// Package env loads dotenv files and process environment variables into nested
// maps and Go structs. Multiple sources merge with later scalars overriding
// earlier ones; maps merge recursively.
//
// Example:
//
//	c := env.New(
//	    env.WithFile(".default.env"),
//	    env.WithFile(".env"),
//	    env.WithCurrentEnvironment(),
//	)
//	var cfg MyConfig
//	if err := c.Unmarshal(&cfg); err != nil { ... }
package env
