// Package source provides pluggable configuration inputs (bytes, files, URLs).
package source

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Source opens a byte stream for reading. Callers must close the ReadCloser.
type Source interface {
	Open(ctx context.Context) (io.ReadCloser, error)
	String() string
}

// Bytes holds raw bytes in memory.
type Bytes struct {
	Data []byte
	Name string
}

// Open returns a ReadCloser over a copy of Data.
func (b Bytes) Open(_ context.Context) (io.ReadCloser, error) {
	return io.NopCloser(bytes.NewReader(b.Data)), nil
}

// String returns a label for errors (defaults to "bytes").
func (b Bytes) String() string {
	if b.Name != "" {
		return b.Name
	}
	return "bytes"
}

// Reader wraps an io.Reader. The reader is consumed once per Open call.
type Reader struct {
	R    io.Reader
	Name string
}

// Open returns the reader as a ReadCloser when possible, otherwise wraps r.R.
func (r Reader) Open(_ context.Context) (io.ReadCloser, error) {
	rc, ok := r.R.(io.ReadCloser)
	if ok {
		return rc, nil
	}
	return io.NopCloser(r.R), nil
}

// String returns a label for errors (defaults to "reader").
func (r Reader) String() string {
	if r.Name != "" {
		return r.Name
	}
	return "reader"
}

// File opens a path on the local filesystem.
type File struct {
	Path string
}

// Open opens the file path.
func (f File) Open(_ context.Context) (io.ReadCloser, error) {
	return os.Open(f.Path)
}

// String returns the file path.
func (f File) String() string {
	return f.Path
}

// URL fetches remote content over HTTP or HTTPS.
type URL struct {
	Raw        string
	Header     http.Header
	HTTPClient *http.Client
}

func defaultClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}

// Open performs a GET request and returns the response body.
func (u URL) Open(ctx context.Context) (io.ReadCloser, error) {
	client := u.HTTPClient
	if client == nil {
		client = defaultClient()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.Raw, nil)
	if err != nil {
		return nil, fmt.Errorf("source url: new request: %w", err)
	}
	if u.Header != nil {
		req.Header = u.Header.Clone()
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("source url %q: %w", u.Raw, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("source url %q: status %s", u.Raw, resp.Status)
	}
	return resp.Body, nil
}

func (u URL) String() string {
	return u.Raw
}

// JoinNames joins source labels for error messages.
func JoinNames(srcs []Source) string {
	if len(srcs) == 0 {
		return ""
	}
	parts := make([]string, 0, len(srcs))
	for _, s := range srcs {
		if s != nil {
			parts = append(parts, s.String())
		}
	}
	return strings.Join(parts, ", ")
}
