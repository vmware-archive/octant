/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware-tanzu/octant/internal/conversion"
	"github.com/vmware-tanzu/octant/internal/util/kubernetes"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var (
	JobCols = component.NewTableCols("Name", "Labels", "Completions", "Successful", "Age")
)

// JobListHandler prints a job list.
func JobListHandler(ctx context.Context, list *batchv1.JobList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("job list is nil")
	}

	ot := NewObjectTable("Jobs", "We couldn't find any jobs!", JobCols, opts.DashConfig.ObjectStore(), opts.DashConfig.TerminateThreshold())
	ot.EnablePluginStatus(opts.DashConfig.PluginManager())
	for _, job := range list.Items {
		row := component.TableRow{}
		nameLink, err := opts.Link.ForObject(&job, job.Name)

		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(job.Labels)
		row["Completions"] = component.NewText(conversion.PtrInt32ToString(job.Spec.Completions))
		succeeded := fmt.Sprintf("%d", job.Status.Succeeded)
		row["Successful"] = component.NewText(succeeded)
		row["Age"] = component.NewTimestamp(job.CreationTimestamp.Time)

		if err := ot.AddRowForObject(ctx, &job, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// JobHandler printers a job.
func JobHandler(ctx context.Context, job *batchv1.Job, options Options) (component.Component, error) {
	o := NewObject(job)
	o.EnableEvents()

	jh, err := newJobHandler(job, o)
	if err != nil {
		return nil, err
	}

	if err := jh.Config(options); err != nil {
		return nil, errors.Wrap(err, "print job configuration")
	}

	if err := jh.Status(options); err != nil {
		return nil, errors.Wrap(err, "print job status")
	}

	if err := jh.Pods(ctx, job, options); err != nil {
		return nil, errors.Wrap(err, "print job pods")
	}

	if err := jh.Conditions(options); err != nil {
		return nil, errors.Wrap(err, "print job conditions")
	}

	return o.ToComponent(ctx, options)
}

// JobConfiguration generates a job configuration
type JobConfiguration struct {
	job *batchv1.Job
}

// NewJobConfiguration creates an instance of JobConfiguration
func NewJobConfiguration(job *batchv1.Job) *JobConfiguration {
	return &JobConfiguration{
		job: job,
	}
}

// Create creates a job configuration summary
func (j *JobConfiguration) Create(option Options) (*component.Summary, error) {
	if j == nil || j.job == nil {
		return nil, errors.New("job is nil")
	}

	job := j.job

	sections := component.SummarySections{}

	sections.Add("Back Off Limit", component.NewText(conversion.PtrInt32ToString(job.Spec.BackoffLimit)))
	sections.Add("Completions", component.NewText(conversion.PtrInt32ToString(job.Spec.Completions)))
	sections.Add("Parallelism", component.NewText(conversion.PtrInt32ToString(job.Spec.Parallelism)))

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func createJobStatus(job batchv1.Job) (*component.Summary, error) {
	sections := component.SummarySections{}

	if startTime := job.Status.StartTime; startTime != nil {
		sections.Add("Started", component.NewTimestamp(startTime.Time))
	}

	if completionTime := job.Status.CompletionTime; completionTime != nil {
		sections.Add("Completed", component.NewTimestamp(completionTime.Time))
	}

	sections.Add("Succeeded", component.NewText(fmt.Sprintf("%d", job.Status.Succeeded)))

	summary := component.NewSummary("Status", sections...)
	return summary, nil
}

func createJobConditions(conditions []batchv1.JobCondition) (*component.Table, error) {
	cols := component.NewTableCols("Type", "Last Probe", "Last Transition",
		"Status", "Message", "Reason")
	table := component.NewTable("Conditions", "There are no job conditions!", cols)

	for _, condition := range conditions {
		row := component.TableRow{}

		row["Type"] = component.NewText(string(condition.Type))
		row["Last Probe"] = component.NewTimestamp(condition.LastProbeTime.Time)
		row["Last Transition"] = component.NewTimestamp(condition.LastTransitionTime.Time)
		row["Status"] = component.NewText(string(condition.Status))
		row["Message"] = component.NewText(condition.Message)
		row["Reason"] = component.NewText(condition.Reason)

		table.Add(row)
	}

	return table, nil
}

func createJobListView(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	options.DisableLabels = true

	jobList := &batchv1.JobList{}

	objectStore := options.DashConfig.ObjectStore()
	accessor := meta.NewAccessor()

	namespace, err := accessor.Namespace(object)
	if err != nil {
		return nil, errors.Wrap(err, "get namespace for object")
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return nil, errors.Wrap(err, "Get apiVersion for object")
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return nil, errors.Wrap(err, "get kind for object")
	}

	name, err := accessor.Name(object)
	if err != nil {
		return nil, errors.Wrap(err, "get name for object")
	}

	key := store.Key{
		Namespace:  namespace,
		APIVersion: "batch/v1beta1",
		Kind:       "Job",
	}

	list, _, err := objectStore.List(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "list all objects for key %+v", key)
	}

	for i := range list.Items {
		job := &batchv1.Job{}
		err := kubernetes.FromUnstructured(&list.Items[i], job)
		if err != nil {
			return nil, err
		}

		for _, ownerReference := range job.OwnerReferences {
			if ownerReference.APIVersion == apiVersion &&
				ownerReference.Kind == kind &&
				ownerReference.Name == name {
				jobList.Items = append(jobList.Items, *job)
			}
		}
	}

	return JobListHandler(ctx, jobList, options)
}

type jobObject interface {
	Config(options Options) error
	Status(options Options) error
	Pods(ctx context.Context, object runtime.Object, options Options) error
	Conditions(options Options) error
}

type jobHandler struct {
	job            *batchv1.Job
	configFunc     func(*batchv1.Job, Options) (*component.Summary, error)
	statusFunc     func(*batchv1.Job, Options) (*component.Summary, error)
	podFunc        func(context.Context, runtime.Object, Options) (component.Component, error)
	conditionsFunc func(*batchv1.Job, Options) (*component.Table, error)
	object         *Object
}

var _ jobObject = (*jobHandler)(nil)

func newJobHandler(job *batchv1.Job, object *Object) (*jobHandler, error) {
	if job == nil {
		return nil, errors.New("can't print a nil job")
	}

	if object == nil {
		return nil, errors.New("can't print a job using a nil object printer")
	}

	jh := &jobHandler{
		job:            job,
		configFunc:     defaultJobConfig,
		statusFunc:     defaultJobStatus,
		podFunc:        defaultJobPods,
		conditionsFunc: defaultJobConditions,
		object:         object,
	}

	return jh, nil
}

func (j *jobHandler) Config(options Options) error {
	out, err := j.configFunc(j.job, options)
	if err != nil {
		return err
	}
	j.object.RegisterConfig(out)
	return nil
}

func defaultJobConfig(job *batchv1.Job, options Options) (*component.Summary, error) {
	return NewJobConfiguration(job).Create(options)
}

func (j *jobHandler) Status(options Options) error {
	out, err := j.statusFunc(j.job, options)
	if err != nil {
		return err
	}
	j.object.RegisterSummary(out)
	return nil
}

func defaultJobStatus(job *batchv1.Job, options Options) (*component.Summary, error) {
	return createJobStatus(*job)
}

func (j *jobHandler) Pods(ctx context.Context, object runtime.Object, options Options) error {
	j.object.EnablePodTemplate(j.job.Spec.Template)

	j.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return j.podFunc(ctx, object, options)
		},
	})
	return nil
}

func defaultJobPods(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createPodListView(ctx, object, options)
}

func (j *jobHandler) Conditions(options Options) error {
	if j.job == nil {
		return errors.New("can;t display conditions for nil job")
	}

	j.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return j.conditionsFunc(j.job, options)
		},
	})

	return nil
}

func defaultJobConditions(job *batchv1.Job, options Options) (*component.Table, error) {
	return createJobConditions(job.Status.Conditions)
}
