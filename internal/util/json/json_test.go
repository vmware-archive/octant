package json

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_JSON_ConversionSortsProperly(t *testing.T) {
	input := map[string]interface{}{
		"a": 5,
		"c": "test",
		"d": true,
		"b": 95.5,
	}
	output, err := Marshal(input)
	require.NoError(t, err)

	require.Equal(t, `{"a":5,"b":95.5,"c":"test","d":true}`, string(output))
}
