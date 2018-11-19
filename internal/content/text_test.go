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

	expected := `{"text":"foo","type":"string"}`

	assert.Equal(t, expected, string(data))
}

func TestTimeText(t *testing.T) {
	cases := []struct {
		name     string
		timeText *TimeText
		expected string
		isErr    bool
	}{
		{
			name:     "RFC3339 string",
			timeText: NewTimeText("2018-11-08T17:55:45Z"),
			expected: `{"time":"2018-11-08T17:55:45Z","type":"time"}`,
		},
		{
			name:     "Zero time",
			timeText: NewTimeText("0001-01-01T00:00:00Z"),
			expected: `{"time":"","type":"time"}`,
		},
		{
			name:     "Invalid timestamp",
			timeText: NewTimeText("Tue Nov 10 23:00:00 2009"),
			isErr:    true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.timeText)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(data))
		})
	}
}

func TestLinkText(t *testing.T) {
	lt := NewLinkText("foo", "/bar")

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"ref":"/bar","text":"foo","type":"link"}`

	assert.Equal(t, expected, string(data))
}

func TestLabelsText(t *testing.T) {
	m := map[string]string{
		"foo": "bar",
	}
	lt := NewLabelsText(m)

	data, err := json.Marshal(lt)
	require.NoError(t, err)

	expected := `{"labels":{"foo":"bar"},"type":"labels"}`

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
