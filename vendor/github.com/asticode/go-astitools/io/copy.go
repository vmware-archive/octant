package astiio

import (
	"context"
	"io"
)

// Copy represents a cancellable copy
func Copy(ctx context.Context, src io.Reader, dst io.Writer) (int64, error) {
	return io.Copy(dst, NewReader(ctx, src))
}
