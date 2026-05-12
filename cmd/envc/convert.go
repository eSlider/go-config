package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/eslider/go-config/env"
	"github.com/eslider/go-config/ini"
	"github.com/eslider/go-config/internal/bytesutil"
	libjson "github.com/eslider/go-config/json"
	"github.com/eslider/go-config/toml"
	"github.com/eslider/go-config/yaml"
	yaml3 "gopkg.in/yaml.v3"
)

// RunConvert implements: envc convert --from F --to T [--input] [--output].
func RunConvert(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("convert", flag.ContinueOnError)
	fs.SetOutput(stderr)
	from := fs.String("from", "", "source format: yaml|json|toml|ini|env")
	to := fs.String("to", "", "target format: yaml|json|toml|ini|env")
	input := fs.String("input", "-", "input path, URL, - for stdin, or \"environ\" for current process env (requires --from env)")
	output := fs.String("output", "-", "output path or - for stdout")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *from == "" || *to == "" {
		_, _ = fmt.Fprintln(stderr, "convert: --from and --to are required")
		return 2
	}
	ctx := context.Background()
	var m map[string]any
	var err error
	if *input == "environ" {
		if *from != "env" {
			_, _ = fmt.Fprintln(stderr, "convert: --input environ only works with --from env")
			return 2
		}
		m, err = env.New(env.WithCurrentEnvironment()).Map(ctx)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "convert: environ: %v\n", err)
			return 1
		}
	} else {
		src := openSource(*input, stdin)
		var b []byte
		b, err = bytesutil.ReadAll(ctx, src)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "convert: read: %v\n", err)
			return 1
		}
		m, err = loadMapFromBytes(ctx, *from, b)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "convert: parse: %v\n", err)
			return 1
		}
	}
	outb, err := marshalMap(*to, m)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "convert: marshal: %v\n", err)
		return 1
	}
	out, err := openOutput(*output, stdout)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "convert: output: %v\n", err)
		return 1
	}
	defer func() { _ = out.Close() }()
	if _, err := out.Write(outb); err != nil {
		_, _ = fmt.Fprintf(stderr, "convert: write: %v\n", err)
		return 1
	}
	if *output == "-" || *output == "" {
		if len(outb) > 0 && outb[len(outb)-1] != '\n' {
			_, _ = out.Write([]byte("\n"))
		}
	}
	return 0
}

func loadMapFromBytes(ctx context.Context, format string, b []byte) (map[string]any, error) {
	switch format {
	case "yaml":
		return yaml.New(yaml.WithBytes(b)).Map(ctx)
	case "json":
		return libjson.New(libjson.WithBytes(b)).Map(ctx)
	case "ini":
		return ini.New(ini.WithBytes(b)).Map(ctx)
	case "env":
		return env.New(env.WithBytes(b)).Map(ctx)
	case "toml":
		return toml.New(toml.WithBytes(b)).Map(ctx)
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

func marshalMap(format string, m map[string]any) ([]byte, error) {
	switch format {
	case "yaml":
		return yaml3.Marshal(m)
	case "json":
		return json.MarshalIndent(m, "", "  ")
	case "ini":
		return ini.New().Marshal(m)
	case "env":
		return env.New().Marshal(m)
	case "toml":
		return toml.New().Marshal(m)
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}
