package main

import (
	"fmt"
	"os"
)

// Build metadata injected at release time via -ldflags.
// Defaults keep local `go build` / `go install` builds identifiable as unversioned.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	stdin, stdout, stderr := os.Stdin, os.Stdout, os.Stderr
	args := os.Args[1:]
	if len(args) < 1 {
		_, _ = fmt.Fprintln(stderr, "usage: envc <convert|get|merge|version> [flags]")
		os.Exit(2)
	}
	var code int
	switch args[0] {
	case "convert":
		code = RunConvert(args[1:], stdin, stdout, stderr)
	case "get":
		code = RunGet(args[1:], stdin, stdout, stderr)
	case "merge":
		code = RunMerge(args[1:], stdin, stdout, stderr)
	case "version", "--version", "-v":
		_, _ = fmt.Fprintf(stdout, "envc %s (commit %s, built %s)\n", version, commit, date)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		code = 2
	}
	os.Exit(code)
}
