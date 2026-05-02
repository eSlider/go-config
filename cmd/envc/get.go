package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/eslider/go-config/internal/bytesutil"
	"github.com/eslider/go-config/internal/keymap"
)

// RunGet implements: envc get --from F --path p [input].
func RunGet(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("get", flag.ContinueOnError)
	fs.SetOutput(stderr)
	from := fs.String("from", "", "format: yaml|json|ini|env")
	path := fs.String("path", "", "dot-separated path (e.g. service.subservice.name)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *from == "" || *path == "" {
		_, _ = fmt.Fprintln(stderr, "get: --from and --path are required")
		return 2
	}
	pos := fs.Args()
	input := "-"
	if len(pos) > 0 {
		input = pos[0]
	}
	ctx := context.Background()
	src := openSource(input, stdin)
	b, err := bytesutil.ReadAll(ctx, src)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "get: read: %v\n", err)
		return 1
	}
	m, err := loadMapFromBytes(ctx, *from, b)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "get: parse: %v\n", err)
		return 1
	}
	v, err := walkPath(m, strings.Split(*path, "."))
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "get: %v\n", err)
		return 1
	}
	switch val := v.(type) {
	case string, bool, float64, int, int64, nil:
		_, _ = fmt.Fprintln(stdout, val)
	default:
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(val); err != nil {
			_, _ = fmt.Fprintf(stderr, "get: encode: %v\n", err)
			return 1
		}
	}
	return 0
}

func walkPath(m map[string]any, parts []string) (any, error) {
	var cur any = m
	for i, p := range parts {
		if p == "" {
			continue
		}
		pk := keymap.LowerAlnum(p)
		switch t := cur.(type) {
		case map[string]any:
			nx, ok := t[pk]
			if !ok {
				return nil, fmt.Errorf("missing key %q (normalized %q) at segment %d", p, pk, i)
			}
			cur = nx
		default:
			return nil, fmt.Errorf("not a map at segment %d", i)
		}
	}
	return cur, nil
}
