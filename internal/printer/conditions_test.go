package printer

import (
	"fmt"
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
	cases := []struct {
		name          string
		status        unstructured.Unstructured
		err           error
		nilConditions bool
		conditions    []interface{}
	}{
		{
			name:          "no status",
			status:        unstructured.Unstructured{Object: map[string]interface{}{"noStatus": nil}},
			err:           fmt.Errorf("no status found for object"),
			nilConditions: true,
		},
		{
			name:          "bad status",
			status:        unstructured.Unstructured{Object: map[string]interface{}{"status": 1}},
			err:           fmt.Errorf(".status accessor error: 1 is of the type int, expected map[string]interface{}"),
			nilConditions: true,
		},
		{
			name:       "no conditions",
			status:     unstructured.Unstructured{Object: map[string]interface{}{"status": map[string]interface{}{"noConditions": nil}}},
			conditions: []interface{}{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conditions, err := parseConditions(tc.status)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.Errorf(t, tc.err, ".status accessor error: 1 is of the type int, expected map[string]interface{}", err)
			}
			if !tc.nilConditions {
				assert.NotNil(t, tc.conditions)
			} else {
				assert.Nil(t, conditions)
			}
		})
	}
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
	conditions, err := parseConditions(*u)
	require.NoError(t, err)
	got := createConditionsTable(conditions, "", podConditionColumns)
	require.NoError(t, err)

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
