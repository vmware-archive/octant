package overview

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_stringText(t *testing.T) {
	st := newStringText("foo")

	data, err := json.Marshal(st)
	require.NoError(t, err)

	expected := `{"text":"foo","type":"string"}`

	assert.Equal(t, expected, string(data))
}

func Test_linkText(t *testing.T) {
	lt := newLinkText("foo", "/bar")

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"ref":"/bar","text":"foo","type":"link"}`

	assert.Equal(t, expected, string(data))
}

func Test_labelsText(t *testing.T) {
	m := map[string]string{
		"foo": "bar",
	}
	lt := newLabelsText(m)

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"labels":{"foo":"bar"},"type":"labels"}`

	assert.Equal(t, expected, string(data))
}
