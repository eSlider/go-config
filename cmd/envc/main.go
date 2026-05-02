package main

import (
	"fmt"
	"os"
)

func main() {
	stdin, stdout, stderr := os.Stdin, os.Stdout, os.Stderr
	args := os.Args[1:]
	if len(args) < 1 {
		_, _ = fmt.Fprintln(stderr, "usage: envc <convert|get|merge> [flags]")
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
	default:
		_, _ = fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		code = 2
	}
	os.Exit(code)
}
