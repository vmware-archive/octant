/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"

	"github.com/vmware/octant/internal/conversion"
	"github.com/vmware/octant/internal/testutil"
	"github.com/vmware/octant/pkg/view/component"
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
		"Name":        component.NewLink("", "job", "/job"),
		"Labels":      component.NewLabels(validJobLabels),
		"Completions": component.NewText("1"),
		"Successful":  component.NewText("1"),
		"Age":         component.NewTimestamp(validJobCreationTime),
	})

	component.AssertEqual(t, expected, got)
}
