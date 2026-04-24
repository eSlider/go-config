// Package env reads process environment variables and decodes them into
// Go structs using underscore-delimited paths.
//
// A variable named "SERVICE_HTTP_PORT" becomes the path Service.HTTP.Port
// (case-insensitive, '_' is a path separator). Values are decoded via
// github.com/mitchellh/mapstructure, so numeric, boolean and slice
// conversions happen automatically.
//
//	var cfg struct {
//		Service struct {
//			HTTP struct {
//				Port int
//			}
//			Key string
//		}
//	}
//	if err := env.Unmarshal(&cfg); err != nil { ... }
//
// UnmarshalPrefix ignores variables that don't start with the given
// prefix and strips it before building the path:
//
//	_ = os.Setenv("APP_DB_HOST", "localhost")
//	_ = env.UnmarshalPrefix(&cfg, "APP_")
//	// cfg.Db.Host == "localhost"
//
// String values are TrimSpace-trimmed by default; pass WithTrim(false) to
// opt out.
package env
