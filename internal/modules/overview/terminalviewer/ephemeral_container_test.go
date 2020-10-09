package terminalviewer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"

	clusterFake "github.com/vmware-tanzu/octant/internal/cluster/fake"
	configFake "github.com/vmware-tanzu/octant/internal/config/fake"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/testutil"
)

func Test_NewEphemeralContainerGenerator(t *testing.T) {
	ctx := context.Background()
	_, err := NewEphemeralContainerGenerator(ctx, nil, log.NopLogger(), nil)
	require.Error(t, err)
}

//func Test_UpdateObject(t *testing.T) {
//	// TODO: Have mocks that allow interface conversion to EphemeralContainers
//	// See https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/component-base/featuregate/testing/feature_gate.go
//	//defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, "", true)()
//
//	controller := gomock.NewController(t)
//	defer controller.Finish()
//
//	pod := testutil.CreatePod("pod")
//	object := runtime.Object(pod)
//
//	dashConfig := configFake.NewMockDash(controller)
//	ctx := context.Background()
//
//	clusterClient := clusterFake.NewMockClientInterface(controller)
//	dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()
//	kubernetesClient := clusterFake.NewMockKubernetesInterface(controller)
//	clusterClient.EXPECT().KubernetesClient().Return(kubernetesClient, nil).AnyTimes()
//
//	fakeClientSet := testClient.NewSimpleClientset(&corev1.PodList{Items: []corev1.Pod{*pod}})
//	kubernetesClient.EXPECT().CoreV1().AnyTimes().Return(fakeClientSet.CoreV1())
//
//	ecg, err := NewEphemeralContainerGenerator(ctx, dashConfig, log.NopLogger(), object)
//	require.NoError(t, err)
//
//	err = ecg.UpdateObject(ctx, object)
//	require.NoError(t, err)
//}

func Test_FeatureEnabled_false(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	dashConfig := configFake.NewMockDash(controller)
	ctx := context.Background()

	clusterClient := clusterFake.NewMockClientInterface(controller)
	dashConfig.EXPECT().ClusterClient().Return(clusterClient).AnyTimes()

	discoveryInterface := clusterFake.NewMockDiscoveryInterface(controller)
	clusterClient.EXPECT().DiscoveryClient().Return(discoveryInterface, nil).AnyTimes()

	discoveryInterface.EXPECT().ServerGroupsAndResources().Return(nil, nil, nil).AnyTimes()

	pod := testutil.CreatePod("pod")
	object := runtime.Object(pod)
	ecg, err := NewEphemeralContainerGenerator(ctx, dashConfig, log.NopLogger(), object)
	require.NoError(t, err)

	enabled := ecg.FeatureEnabled()
	require.False(t, enabled)
}
