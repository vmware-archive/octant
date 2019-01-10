package content

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStringText(t *testing.T) {
	st := NewStringText("foo")

	data, err := json.Marshal(st)
	require.NoError(t, err)

	expected := `{"config":{"value":"foo"},"metadata":{"type":"text"},"text":"foo","type":"string"}`

	assert.Equal(t, expected, string(data))
}

func TestLinkText(t *testing.T) {
	lt := NewLinkText("foo", "/bar")

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"config":{"ref":"/bar","value":"foo"},"metadata":{"type":"link"},"ref":"/bar","text":"foo","type":"link"}`

	assert.Equal(t, expected, string(data))
}

func TestLabelsText(t *testing.T) {
	m := map[string]string{
		"foo": "bar",
	}
	lt := NewLabelsText(m)

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"config":{"labels":{"foo":"bar"}},"labels":{"foo":"bar"},"metadata":{"type":"labels"},"type":"labels"}`

	assert.Equal(t, expected, string(data))
}

func TestListText(t *testing.T) {
	list := []string{"foo", "bar"}

	lt := NewListText(list)

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"list":["foo","bar"],"type":"list"}`

	assert.Equal(t, expected, string(data))

}
