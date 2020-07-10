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
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/internal/conversion"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/internal/testutil"
	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_CronJobListHandler(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)
	printOptions := tpo.ToOptions()

	labels := map[string]string{
		"app": "myapp",
	}

	now := testutil.Time()

	cronJob := &batchv1beta1.CronJob{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "batch/v1beta1",
			Kind:       "CronJob",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cron",
			Namespace: "default",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
			Labels: labels,
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:    "*/1 * * * *",
			JobTemplate: batchv1beta1.JobTemplateSpec{},
		},
	}

	tpo.PathForObject(cronJob, cronJob.Name, "/cron")

	object := &batchv1beta1.CronJobList{
		Items: []batchv1beta1.CronJob{*cronJob},
	}

	ctx := context.Background()
	got, err := CronJobListHandler(ctx, object, printOptions)
	require.NoError(t, err)

	payload := action.Payload{
		"namespace":  cronJob.Namespace,
		"apiVersion": cronJob.APIVersion,
		"kind":       cronJob.Kind,
		"name":       cronJob.Name,
	}

	cols := component.NewTableCols("Name", "Labels", "Schedule", "Age")
	expected := component.NewTable("CronJobs", "We couldn't find any cron jobs!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "cron", "/cron", func(l *component.Link) {
			l.SetStatus(component.TextStatusOK,
				component.NewList(nil, []component.Component{
					component.NewText("batch/v1beta1 CronJob is OK"),
				},
				))
		}),
		"Labels":   component.NewLabels(labels),
		"Schedule": component.NewText("*/1 * * * *"),
		"Age":      component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			{
				Name:         "Trigger",
				ActionPath:   octant.ActionOverviewCronjob,
				Payload:      payload,
				Confirmation: nil,
				Type:         component.GridActionDanger,
			},
			{
				Name:         "Suspend",
				ActionPath:   octant.ActionOverviewSuspendCronjob,
				Payload:      payload,
				Confirmation: nil,
				Type:         component.GridActionDanger,
			},
			buildObjectDeleteAction(t, cronJob),
		}),
	})

	component.AssertEqual(t, expected, got)
}

func Test_CronJobConfiguration(t *testing.T) {

	now := time.Unix(1550075244, 0)
	lastScheduled := time.Unix(1550075184, 0)
	suspend := false

	cj := &batchv1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cron-test",
			CreationTimestamp: metav1.Time{
				Time: now,
			},
		},
		Spec: batchv1beta1.CronJobSpec{
			Schedule:                   "*/1 * * * *",
			StartingDeadlineSeconds:    conversion.PtrInt64(200),
			ConcurrencyPolicy:          batchv1beta1.ForbidConcurrent,
			Suspend:                    &suspend,
			SuccessfulJobsHistoryLimit: conversion.PtrInt32(3),
			FailedJobsHistoryLimit:     conversion.PtrInt32(1),
		},
		Status: batchv1beta1.CronJobStatus{
			LastScheduleTime: &metav1.Time{
				Time: lastScheduled,
			},
		},
	}

	cases := []struct {
		name     string
		cronjob  *batchv1beta1.CronJob
		isErr    bool
		expected *component.Summary
	}{
		{
			name:    "cronjob",
			cronjob: cj,
			expected: component.NewSummary("Configuration", []component.SummarySection{
				{
					Header:  "Schedule",
					Content: component.NewText("*/1 * * * *"),
				},
				{
					Header:  "Suspend",
					Content: component.NewText("false"),
				},
				{
					Header:  "Concurrency Policy",
					Content: component.NewText("Forbid"),
				},
				{
					Header:  "Last Schedule Time",
					Content: component.NewTimestamp(lastScheduled),
				},
				{
					Header:  "Starting Deadline Seconds",
					Content: component.NewText("200s"),
				},
				{
					Header:  "Successful Job History Limit",
					Content: component.NewText("3"),
				},
				{
					Header:  "Failed Job History Limit",
					Content: component.NewText("1"),
				},
			}...),
		},
		{
			name:    "cronjob is nil",
			cronjob: nil,
			isErr:   true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cc := NewCronJobConfiguration(tc.cronjob)

			summary, err := cc.Create()
			if tc.isErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			component.AssertEqual(t, tc.expected, summary)
		})
	}
}

func Test_createJobListView(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	tpo := newTestPrinterOptions(controller)

	ctx := context.Background()
	now := testutil.Time()
	labels := map[string]string{
		"foo": "bar",
	}

	cronJob := testutil.CreateCronJob("cronjob")
	job := testutil.CreateJob("job")

	job.SetOwnerReferences(testutil.ToOwnerReferences(t, cronJob))
	job.CreationTimestamp = metav1.Time{Time: now}
	job.Labels = labels
	job.Spec = batchv1.JobSpec{
		Completions: conversion.PtrInt32(1),
	}
	job.Status.Succeeded = 1

	jobs := &batchv1.JobList{
		Items: []batchv1.Job{*job},
	}

	tpo.PathForObject(job, job.Name, "/job")

	jobList := &unstructured.UnstructuredList{}
	for _, j := range jobs.Items {
		jobList.Items = append(jobList.Items, *testutil.ToUnstructured(t, &j))
	}
	key := store.Key{
		Namespace:  job.Namespace,
		APIVersion: "batch/v1beta1",
		Kind:       "Job",
	}

	tpo.objectStore.EXPECT().List(gomock.Any(), gomock.Eq(key)).Return(jobList, false, nil)

	printOptions := tpo.ToOptions()

	got, err := createJobListView(ctx, cronJob, printOptions)
	require.NoError(t, err)

	cols := component.NewTableCols("Name", "Labels", "Completions", "Successful", "Age")
	expected := component.NewTable("Jobs", "We couldn't find any jobs!", cols)
	expected.Add(component.TableRow{
		"Name": component.NewLink("", "job", "/job",
			genObjectStatus(component.TextStatusWarning, []string{
				"Job has succeeded 1 time",
				"Job is in progress",
			}),
		),
		"Labels":      component.NewLabels(labels),
		"Completions": component.NewText("1"),
		"Successful":  component.NewText("1"),
		"Age":         component.NewTimestamp(now),
		component.GridActionKey: gridActionsFactory([]component.GridAction{
			buildObjectDeleteAction(t, job),
		}),
	})

	component.AssertEqual(t, expected, got)
}
