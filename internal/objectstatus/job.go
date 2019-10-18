/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package objectstatus

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware-tanzu/octant/pkg/store"
)

func runJobStatus(_ context.Context, object runtime.Object, o store.Store) (ObjectStatus, error) {
	if object == nil {
		return ObjectStatus{}, errors.Errorf("job is nil")
	}

	job := &batchv1.Job{}

	if err := scheme.Scheme.Convert(object, job, 0); err != nil {
		return ObjectStatus{}, errors.Wrap(err, "convert object to job")
	}

	os := ObjectStatus{}

	if succeeded := job.Status.Succeeded; succeeded > 0 {
		os.AddDetailf("Job has succeeded %s", pluralizeTime(succeeded))
	}

	if failed := job.Status.Failed; failed > 0 {
		os.AddDetailf("Job has failed %s", pluralizeTime(failed))
	}

	switch {
	case hasJobCondition(*job, batchv1.JobComplete):
		status := job.Status
		if complete, start := status.CompletionTime, status.StartTime; complete != nil && start != nil {
			elapsed := complete.Sub(start.Time)
			os.AddDetailf("Job completed in %s", elapsed)
		}

	case hasJobCondition(*job, batchv1.JobFailed):
		os.SetError()
		for _, condition := range job.Status.Conditions {
			os.AddDetail(condition.Message)
		}

	default:
		os.SetWarning()
		os.AddDetail("Job is in progress")
	}

	return os, nil
}

func pluralizeTime(count int32) string {
	modifier := "time"
	if count > 1 {
		modifier = "times"
	}

	return fmt.Sprintf("%d %s", count, modifier)
}

func hasJobCondition(job batchv1.Job, conditionType batchv1.JobConditionType) bool {
	for _, condition := range job.Status.Conditions {
		if condition.Type == conditionType {
			return true
		}
	}
	return false
}
