package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
)

func TestRouter_Match(t *testing.T) {
	genContentResponse := func(name string) component.ContentResponse {
		return component.ContentResponse{
			Components: []component.Component{
				component.NewText(name),
			},
		}

	}

	router := NewRouter()
	router.HandleFunc("/nested1/nested2", func(request *Request) (component.ContentResponse, error) {
		return genContentResponse("nested2"), nil
	})
	router.HandleFunc("/nested1", func(request *Request) (component.ContentResponse, error) {
		return genContentResponse("nested1"), nil
	})
	router.HandleFunc("/glob*", func(request *Request) (component.ContentResponse, error) {
		return genContentResponse("glob"), nil
	})
	router.HandleFunc("/", func(request *Request) (component.ContentResponse, error) {
		return genContentResponse("root"), nil
	})

	tests := []struct {
		name     string
		path     string
		isFound  bool
		expected component.ContentResponse
	}{
		{
			name:     "path matches",
			path:     "/",
			isFound:  true,
			expected: genContentResponse("root"),
		},
		{
			name:     "glob",
			path:     "/glob/foo",
			isFound:  true,
			expected: genContentResponse("glob"),
		},
		{
			name:     "nested content",
			path:     "/nested1",
			isFound:  true,
			expected: genContentResponse("nested1"),
		},
		{
			name:     "deeply nested content",
			path:     "/nested1/nested2",
			isFound:  true,
			expected: genContentResponse("nested2"),
		},
		{
			name:    "not found",
			path:    "/invalid",
			isFound: false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			handleFunc, ok := router.Match(test.path)
			if !test.isFound {
				require.False(t, ok)
				return
			}
			require.True(t, ok)

			request := &Request{
				baseRequest:     newBaseRequest(context.Background(), "plugin-name"),
				dashboardClient: nil,
				Path:            test.path,
			}
			got, err := handleFunc(request)
			require.NoError(t, err)

			testutil.AssertJSONEqual(t, test.expected, got)
		})
	}
}
