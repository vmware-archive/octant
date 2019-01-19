package astios

import (
	"context"
	"os"
)

// Move is a cross partitions cancellable move even if files are on different partitions
func Move(ctx context.Context, src, dst string) (err error) {
	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Copy
	if err = Copy(ctx, src, dst); err != nil {
		return
	}

	// Check context
	if err = ctx.Err(); err != nil {
		return
	}

	// Delete
	err = os.Remove(src)
	return
}
