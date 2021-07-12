package api_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/testutil"

	"github.com/golang/mock/gomock"

	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	octantFake "github.com/vmware-tanzu/octant/internal/octant/fake"
	"github.com/vmware-tanzu/octant/pkg/api/fake"
)

func Test_TerminalStateManager(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)

	state := octantFake.NewMockState(controller)
	octantClient := fake.NewMockOctantClient(controller)

	tsm := api.NewTerminalStateManager(dashConfig)

	ctx := context.Background()
	tsm.Start(ctx, state, octantClient)
}

func Test_isWindowsContainer(t *testing.T) {
	windowsPod := testutil.CreatePod("pod")
	windowsPod.Spec.Tolerations = []corev1.Toleration{
		{
			Key:      "os",
			Operator: corev1.TolerationOpEqual,
			Value:    "windows",
			Effect:   corev1.TaintEffectNoSchedule,
		},
	}
	windowsPod.Spec.NodeSelector = map[string]string{
		"kubernetes.io/os": "windows",
	}

	pod := testutil.CreatePod("pod")

	cases := []struct {
		name      string
		pod       *corev1.Pod
		isWindows bool
	}{
		{
			name:      "windows container",
			pod:       windowsPod,
			isWindows: true,
		},
		{
			name:      "not windows container",
			pod:       pod,
			isWindows: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.isWindows, api.IsWindowsContainer(tc.pod))
		})
	}

}
