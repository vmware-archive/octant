package printer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_createConditionsTableErrs(t *testing.T) {
	// No status found
	noStatus := unstructured.Unstructured{Object: map[string]interface{}{"noStatus": nil}}
	table, ok, err := createConditionsTable(noStatus, "", nil)
	assert.Nil(t, err)
	assert.False(t, ok)
	assert.Nil(t, table)

	// Bad status, not a map[string]interface{}
	badStatus := unstructured.Unstructured{Object: map[string]interface{}{"status": 1}}
	table, ok, err = createConditionsTable(badStatus, "", nil)
	assert.EqualError(t, err, ".status accessor error: 1 is of the type int, expected map[string]interface{}")
	assert.False(t, ok)
	assert.Nil(t, table)

	// No conditions found
	noConditions := unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{"noConditions": nil}}}
	table, ok, err = createConditionsTable(noConditions, "", nil)
	assert.Nil(t, err)
	assert.False(t, ok)
	assert.NotNil(t, table)
}

func Test_createPodConditionsView(t *testing.T) {
	now := metav1.Time{Time: time.Now()}

	pod := testutil.CreatePod("pod")
	pod.Status.Conditions = []corev1.PodCondition{
		{
			Type:               corev1.PodInitialized,
			Status:             corev1.ConditionTrue,
			LastTransitionTime: now,
			LastProbeTime:      now,
			Message:            "message",
			Reason:             "reason",
		},
	}

	u := toUnstructured(t, pod)
	got, ok, err := createConditionsTable(*u, "", podConditionColumns)
	require.NoError(t, err)
	assert.True(t, ok)

	cols := component.NewTableCols("Type", "Reason", "Status", "Message", "Last Probe", "Last Transition")
	expected := component.NewTable("Conditions", "There are no conditions!", cols)

	expected.Add([]component.TableRow{
		{
			"Type":            component.NewText("Initialized"),
			"Status":          component.NewText("True"),
			"Last Transition": component.NewTimestamp(now.Time),
			"Last Probe":      component.NewTimestamp(now.Time),
			"Message":         component.NewText("message"),
			"Reason":          component.NewText("reason"),
		},
	}...)

	component.AssertEqual(t, expected, got)
}
