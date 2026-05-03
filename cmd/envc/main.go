package main

import (
	"fmt"
	"io"
	"os"
)

// Build metadata injected at release time via -ldflags.
// Defaults keep local `go build` / `go install` builds identifiable as unversioned.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// printUsage writes the root help text to w.
func printUsage(w io.Writer) {
	_, _ = io.WriteString(w, `envc — convert, query, and merge configuration across env, YAML, JSON, and INI.

Usage:
  envc <command> [arguments]

Commands:
  convert   Read one source, normalize keys, write another format (--from, --to)
  get       Read one source and print the value at a dot path (--from, --path)
  merge     Deep-merge multiple sources in order, then write one output (--from, --to)
  version   Print version and build metadata (also -v, --version)

Each command accepts its own flags; run:
  envc <command> -h
`)
}

func main() {
	stdin, stdout, stderr := os.Stdin, os.Stdout, os.Stderr
	args := os.Args[1:]
	if len(args) < 1 {
		printUsage(stderr)
		os.Exit(2)
	}
	var code int
	switch args[0] {
	case "help", "-h", "--help":
		printUsage(stdout)
		code = 0
	case "convert":
		code = RunConvert(args[1:], stdin, stdout, stderr)
	case "get":
		code = RunGet(args[1:], stdin, stdout, stderr)
	case "merge":
		code = RunMerge(args[1:], stdin, stdout, stderr)
	case "version", "--version", "-v":
		_, _ = fmt.Fprintf(stdout, "envc %s (commit %s, built %s)\n", version, commit, date)
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n\n", args[0])
		printUsage(stderr)
		code = 2
	}
	os.Exit(code)
}
