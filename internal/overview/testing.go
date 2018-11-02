package overview

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func assertViewInvalidObject(t *testing.T, v View) {
	ctx := context.Background()
	_, err := v.Content(ctx, nil, nil)
	require.Error(t, err)
}
