/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
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

	"github.com/vmware/octant/internal/conversion"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	cols := component.NewTableCols("Name", "Labels", "Schedule", "Age")
	expected := component.NewTable("CronJobs", cols)
	expected.Add(component.TableRow{
		"Name":     component.NewLink("", "cron", "/cron"),
		"Labels":   component.NewLabels(labels),
		"Schedule": component.NewText("*/1 * * * *"),
		"Age":      component.NewTimestamp(now),
	})

	component.AssertEqual(t, expected, got)
}

func TestCronJobConfiguration(t *testing.T) {

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

			assert.Equal(t, tc.expected, summary)
		})
	}
}
