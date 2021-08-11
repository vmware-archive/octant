package printer

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/testutil"
)

func Test_GetImageManifest(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		manifestPath string
		configPath   string
		hostOS       string
		error        string
	}{
		{
			name:         "alpine",
			input:        "docker://alpine:3.14.0",
			manifestPath: "alpine_manifest.json",
			configPath:   "alpine_config.json",
			hostOS:       "linux",
		},
		{
			name:         "nginx",
			input:        "docker://nginx:1.16.1",
			manifestPath: "nginx_manifest.json",
			configPath:   "nginx_config.json",
			hostOS:       "linux",
		},
		{
			name:         "nginx without docker protocol",
			input:        "nginx:1.16.1",
			manifestPath: "nginx_manifest.json",
			configPath:   "nginx_config.json",
			hostOS:       "linux",
		},
		{
			name:         "alpine windows",
			input:        "docker://alpine:3.14.0",
			manifestPath: "alpine_manifest.json",
			configPath:   "alpine_config.json",
			hostOS:       "windows",
			error:        "error parsing manifest for for image docker://alpine:3.14.0: choosing image instance: no image found in manifest list for architecture amd64, variant \"\", OS windows",
		},
	}
	mc := NewManifestConfiguration()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := testutil.LoadTestData(t, tt.manifestPath)
			expected := string(input)[:len(input)-1] // remove last `\n` added by editor

			inputConfig := testutil.LoadTestData(t, tt.configPath)
			expectedConfig := string(inputConfig)[:len(inputConfig)-1]

			manifest, config, err := mc.GetImageManifest(context.Background(), tt.hostOS, tt.input)
			fmt.Println("err", err)
			if len(tt.error) > 0 {
				assert.EqualError(t, err, tt.error)
				return
			}

			require.NoError(t, err)
			testutil.AssertJSONEqual(t, expected, manifest)
			testutil.AssertJSONEqual(t, expectedConfig, config)
		})
	}
}
