// Package bytesutil reads entire sources into memory.
package bytesutil

import (
	"context"
	"io"

	"github.com/eslider/go-config/internal/source"
)

// ReadAll reads and closes the stream opened by s.
func ReadAll(ctx context.Context, s source.Source) ([]byte, error) {
	rc, err := s.Open(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rc.Close() }()
	return io.ReadAll(rc)
}
