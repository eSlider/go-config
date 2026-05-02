package main

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/eslider/go-config/internal/bytesutil"
	"github.com/eslider/go-config/internal/merge"
	"github.com/eslider/go-config/internal/source"
)

// RunMerge implements: envc merge --from F --to T [--output] inputs...
func RunMerge(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	fs.SetOutput(stderr)
	from := fs.String("from", "", "input format: yaml|json|ini|env")
	to := fs.String("to", "", "output format: yaml|json|ini|env")
	output := fs.String("output", "-", "output path or - for stdout")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *from == "" || *to == "" {
		_, _ = fmt.Fprintln(stderr, "merge: --from and --to are required")
		return 2
	}
	inputs := fs.Args()
	if len(inputs) == 0 {
		_, _ = fmt.Fprintln(stderr, "merge: at least one input file or URL required")
		return 2
	}
	ctx := context.Background()
	acc := make(map[string]any)
	for _, in := range inputs {
		var src source.Source
		if in == "-" {
			b, err := io.ReadAll(stdin)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "merge: stdin: %v\n", err)
				return 1
			}
			src = source.Bytes{Data: b, Name: "stdin"}
		} else {
			src = openSource(in, stdin)
		}
		b, err := bytesutil.ReadAll(ctx, src)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "merge: read %s: %v\n", in, err)
			return 1
		}
		m, err := loadMapFromBytes(ctx, *from, b)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "merge: parse %s: %v\n", in, err)
			return 1
		}
		merge.DeepMerge(acc, m)
	}
	outb, err := marshalMap(*to, acc)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "merge: marshal: %v\n", err)
		return 1
	}
	out, err := openOutput(*output, stdout)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "merge: output: %v\n", err)
		return 1
	}
	defer func() { _ = out.Close() }()
	if _, err := out.Write(outb); err != nil {
		_, _ = fmt.Fprintf(stderr, "merge: write: %v\n", err)
		return 1
	}
	if *output == "-" || *output == "" {
		if len(outb) > 0 && outb[len(outb)-1] != '\n' {
			_, _ = out.Write([]byte("\n"))
		}
	}
	return 0
}
