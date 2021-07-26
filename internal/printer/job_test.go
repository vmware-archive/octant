/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/octant/internal/conversion"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func Test_JobListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	validJobLabels := map[string]string{
		"app": "testing",
	}

	validJobCreationTime := testutil.Time()

	validJob := testutil.CreateJob("job")
	validJob.CreationTimestamp = *testutil.CreateTimestamp()
	validJob.Labels = validJobLabels
	validJob.Spec = batchv1.JobSpec{
		Completions: conversion.PtrInt32(1),
	}
	validJob.Status = batchv1.JobStatus{
		Succeeded: 1,
		Conditions: []batchv1.JobCondition{
			{
				Reason: "reason",
			},
		},
	}

	validJobList := &batchv1.JobList{
		Items: []batchv1.Job{
			*validJob,
		},
	}

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	tpo.PathForObject(validJob, validJob.Name, "/job")

	ctx := context.Background()
	got, err := JobListHandler(ctx, validJobList, printOptions)
	require.NoError(t, err)

	expected := component.NewTable("Jobs", "We couldn't find any jobs!", JobCols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "job", "/job",
			genObjectStatus(component.TextStatusWarning, []string{
				"Job has succeeded 1 time",
				"Job is in progress",
			})),
		"Labels":      component.NewLabels(validJobLabels),
		"Completions": component.NewText("1"),
		"Successful":  component.NewText("1"),
		"Age":         component.NewTimestamp(validJobCreationTime),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, validJob),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_JobConfiguration(t *testing.T) {
	var backofflimit int32 = 4
	var completions int32 = 1
	var parallelism int32 = 1

	job := testutil.CreateJob("job")
	job.Spec.BackoffLimit = &backofflimit
	job.Spec.Completions = &completions
	job.Spec.Parallelism = &parallelism

	cases := []struct {
		name     string
		job      *batchv1.Job
		isErr    bool
		expected *component.Summary
	}{
		{
			name: "general",
			job:  job,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Back Off Limit",
					Content: component.NewText("4"),
				},
				{
					Header:  "Completions",
					Content: component.NewText("1"),
				},
				{
					Header:  "Parallelism",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name:  "job is nil",
			job:   nil,
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			tpo := newTestPrinterOptions(controller)
			printOptions := tpo.ToOptions()

			jh := NewJobConfiguration(tc.job)

			summary, err := jh.Create(printOptions)
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createJobStatus(t *testing.T) {
	job := testutil.CreateJob("job")
	job.Status.Succeeded = int32(1)
	job.Status.StartTime = &metav1.Time{Time: testutil.Time()}
	job.Status.CompletionTime = &metav1.Time{Time: time.Now()}

	got, err := createJobStatus(*job)
	require.NoError(t, err)

	sections := component.SummarySections{
		{Header: "Started", Content: component.NewTimestamp(testutil.Time())},
		{Header: "Completed", Content: component.NewTimestamp(time.Now())},
		{Header: "Succeeded", Content: component.NewText("1")},
	}
	expected := component.NewSummary("Status", sections...)

	assert.Equal(t, expected, got)
}

func Test_createJobConditions(t *testing.T) {
	now := metav1.Time{Time: time.Now()}

	job := testutil.CreateJob("job")
	job.Status.Conditions = []batchv1.JobCondition{
		{
			Type:               batchv1.JobComplete,
			LastProbeTime:      now,
			LastTransitionTime: now,
			Status:             corev1.ConditionTrue,
			Message:            "message",
			Reason:             "reason",
		},
	}

	got, err := createJobConditions(job)
	require.NoError(t, err)

	cols := component.NewTableCols("Type", "Last Probe", "Last Transition",
		"Status", "Message", "Reason")
	expected := component.NewTable("Conditions", "There are no conditions!", cols)
	expected.Add([]component.TableRow{
		{
			"Type":            component.NewText("Complete"),
			"Last Probe":      component.NewTimestamp(now.Time),
			"Last Transition": component.NewTimestamp(now.Time),
			"Status":          component.NewText("True"),
			"Message":         component.NewText("message"),
			"Reason":          component.NewText("reason"),
		},
	}...)

	component.AssertEqual(t, expected, got)
}
