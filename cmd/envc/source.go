package main

import (
	"io"
	"os"
	"strings"

	"github.com/eslider/go-config/internal/source"
)

func openSource(pathOrURL string, stdin io.Reader) source.Source {
	if pathOrURL == "" || pathOrURL == "-" {
		return source.Reader{R: stdin, Name: "stdin"}
	}
	if strings.HasPrefix(pathOrURL, "http://") || strings.HasPrefix(pathOrURL, "https://") {
		return source.URL{Raw: pathOrURL}
	}
	return source.File{Path: pathOrURL}
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }

func openOutput(path string, stdout io.Writer) (io.WriteCloser, error) {
	if path == "" || path == "-" {
		if wc, ok := stdout.(io.WriteCloser); ok {
			return wc, nil
		}
		return nopCloser{stdout}, nil
	}
	return os.Create(path)
}
