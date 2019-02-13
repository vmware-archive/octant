package printer_test

import (
	"testing"
	"time"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/overview/printer"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_JobListHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	got, err := printer.JobListHandler(validJobList, printOptions)
	require.NoError(t, err)

	expected := component.NewTable("Jobs", printer.JobCols)
	expected.Add(component.TableRow{
		"Name":        component.NewLink("", "job", "/content/overview/namespace/default/workloads/jobs/job"),
		"Labels":      component.NewLabels(validJobLabels),
		"Completions": component.NewText("1"),
		"Successful":  component.NewText("1"),
		"Age":         component.NewTimestamp(validJobCreationTime),
	})

	assert.Equal(t, expected, got)
}

func Test_JobHandler(t *testing.T) {
	printOptions := printer.Options{
		Cache: cache.NewMemoryCache(),
	}

	got, err := printer.JobHandler(validJob, printOptions)
	require.NoError(t, err)

	layout, ok := got.(*component.FlexLayout)
	require.True(t, ok)

	assert.Len(t, layout.Config.Sections, 5)
}

var (
	validJobLabels = map[string]string{
		"app": "testing",
	}

	validJobCreationTime = time.Unix(1547211430, 0)

	validJob = &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1",
			Kind:       "Job",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "job",
			Namespace: "default",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: validJobLabels,
		},
		Spec: batchv1.JobSpec{
			Completions: ptrInt32(1),
		},
		Status: batchv1.JobStatus{
			Succeeded: 1,
			Conditions: []batchv1.JobCondition{
				{
					Reason: "reason",
				},
			},
		},
	}

	validJobList = &batchv1.JobList{
		Items: []batchv1.Job{
			*validJob,
		},
	}
)
