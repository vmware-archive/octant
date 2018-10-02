package cluster

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_FromKubeConfig(t *testing.T) {
	kubeconfig := filepath.Join("testdata", "kubeconfig.yaml")
	_, err := FromKubeconfig(kubeconfig)
	require.NoError(t, err)
}
