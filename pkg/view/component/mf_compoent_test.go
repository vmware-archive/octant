package component_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/util/json"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_MFComponent_Marshal(t *testing.T) {
	mfc := component.MfComponentConfig{
		Name:          "template-mf-angular",
		RemoteEntry:   "http://localhost:3000/remoteEntry.js",
		RemoteName:    "templatemfangular",
		ExposedModule: "./web-components",
		ElementName:   "template-mf-angular",
	}
	cases := []struct {
		name         string
		input        *component.MfComponent
		expectedPath string
	}{
		{
			name:         "in general",
			input:        component.NewMFComponent(mfc),
			expectedPath: "mfcomponent.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := json.Marshal(tc.input)
			require.NoError(t, err)

			expected, err := ioutil.ReadFile(filepath.Join("testdata", tc.expectedPath))
			require.NoError(t, err)

			assert.JSONEq(t, string(expected), string(got))
		})
	}
}
