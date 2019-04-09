package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/heptio/developer-dash/internal/conversion"
	storefake "github.com/heptio/developer-dash/internal/objectstore/fake"
	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_JobListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	printOptions := Options{
		ObjectStore: storefake.NewMockObjectStore(controller),
	}

	ctx := context.Background()
	got, err := JobListHandler(ctx, validJobList, printOptions)
	require.NoError(t, err)

	expected := component.NewTable("Jobs", JobCols)
	expected.Add(component.TableRow{
		"Name":        component.NewLink("", "job", "/content/overview/namespace/default/workloads/jobs/job"),
		"Labels":      component.NewLabels(validJobLabels),
		"Completions": component.NewText("1"),
		"Successful":  component.NewText("1"),
		"Age":         component.NewTimestamp(validJobCreationTime),
	})

	assert.Equal(t, expected, got)
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
			Completions: conversion.PtrInt32(1),
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
