package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructors(t *testing.T) {
	err := fmt.Errorf("encountered an error while streaming")
	table := []struct {
		name  string
		sErr  *StreamError
		fatal bool
	}{
		{"NewStreamError", NewStreamError(err), false},
		{"FatalStreamError", FatalStreamError(err), true},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			assert.NotEmpty(t, test.sErr.Timestamp())
			assert.Equal(t, test.sErr.Name(), StreamingConnectionError)
			assert.NotZero(t, test.sErr.ID())
			assert.Equal(t, test.fatal, test.sErr.Fatal)
			assert.Equal(t, err.Error(), test.sErr.Error())
		})
	}
}

func TestIsFatalStreamError(t *testing.T) {
	table := []struct {
		name     string
		sErr     error
		expected bool
	}{
		{
			"Expect IsFatalStreamError to be false for a standard stream error",
			NewStreamError(fmt.Errorf("")),
			false,
		},
		{
			"Expect IsFatalStreamError to be true for fatal stream error",
			FatalStreamError(fmt.Errorf("")),
			true,
		},
		{
			"A non pointer fatal stream error is provided",
			*FatalStreamError(fmt.Errorf("")),
			true,
		},
		{
			"A generic error is provided",
			fmt.Errorf(""),
			false,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsFatalStreamError(test.sErr))
		})
	}
}
