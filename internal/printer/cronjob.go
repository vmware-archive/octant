/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vmware-tanzu/octant/internal/octant"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/pkg/action"
	"github.com/vmware-tanzu/octant/pkg/view/component"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// CronJobListHandler is a printFunc that lists cronjobs
func CronJobListHandler(ctx context.Context, list *batchv1beta1.CronJobList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Schedule", "Age")
	ot := NewObjectTable("CronJobs", "We couldn't find any cron jobs!", cols, opts.DashConfig.ObjectStore())

	for _, c := range list.Items {
		row := component.TableRow{}

		nameLink, err := opts.Link.ForObject(&c, c.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink

		row["Labels"] = component.NewLabels(c.Labels)

		row["Schedule"] = component.NewText(c.Spec.Schedule)

		ts := c.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		if err := addCronJobActions(c, row); err != nil {
			return nil, err
		}

		if err := ot.AddRowForObject(ctx, &c, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

func addCronJobActions(c batchv1beta1.CronJob, row component.TableRow) error {
	payload := action.Payload{
		"namespace":  c.Namespace,
		"apiVersion": c.APIVersion,
		"kind":       c.Kind,
		"name":       c.Name,
	}

	actions := []component.GridAction{
		{
			Name:         "Trigger",
			ActionPath:   octant.ActionOverviewCronjob,
			Payload:      payload,
			Confirmation: nil,
			Type:         component.GridActionDanger,
		},
	}

	suspendAction := component.GridAction{
		Name:         "Suspend",
		ActionPath:   octant.ActionOverviewSuspendCronjob,
		Payload:      payload,
		Confirmation: nil,
		Type:         component.GridActionDanger,
	}
	resumeAction := component.GridAction{
		Name:         "Resume",
		ActionPath:   octant.ActionOverviewResumeCronjob,
		Payload:      payload,
		Confirmation: nil,
		Type:         component.GridActionDanger,
	}
	if c.Spec.Suspend != nil && *c.Spec.Suspend {
		actions = append(actions, resumeAction)
	} else {
		actions = append(actions, suspendAction)
	}

	for _, action := range actions {
		row.AddAction(action)
	}

	return nil
}

// CronJobHandler is a printFunc that prints a CronJob
func CronJobHandler(ctx context.Context, cronJob *batchv1beta1.CronJob, options Options) (component.Component, error) {
	o := NewObject(cronJob)
	o.EnableEvents()

	ch, err := newCronJobHandler(cronJob, o)
	if err != nil {
		return nil, err
	}

	if err := ch.Config(options); err != nil {
		return nil, errors.Wrap(err, "print cronjob configuration")
	}

	if err := ch.Jobs(ctx, cronJob, options); err != nil {
		return nil, errors.Wrap(err, "print cronjob job list")
	}

	return o.ToComponent(ctx, options)
}

// CronJobConfiguration generates cronjob configuration
type CronJobConfiguration struct {
	cronjob *batchv1beta1.CronJob
}

// NewCronJobConfiguration creates an instance of CronJobConfiguration
func NewCronJobConfiguration(c *batchv1beta1.CronJob) *CronJobConfiguration {
	return &CronJobConfiguration{
		cronjob: c,
	}
}

// Create creates a cronjob configuration summary
func (cc *CronJobConfiguration) Create() (*component.Summary, error) {
	if cc == nil || cc.cronjob == nil {
		return nil, errors.New("cronjob is nil")
	}

	sections := component.SummarySections{}

	sections.AddText("Schedule", cc.cronjob.Spec.Schedule)

	if suspend := cc.cronjob.Spec.Suspend; suspend != nil {
		sections.AddText("Suspend", strconv.FormatBool(*suspend))
	}

	sections.AddText("Concurrency Policy", string(cc.cronjob.Spec.ConcurrencyPolicy))

	if lastScheduleTime := cc.cronjob.Status.LastScheduleTime; lastScheduleTime != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Last Schedule Time",
			Content: component.NewTimestamp(lastScheduleTime.Time),
		})
	}

	if sdls := cc.cronjob.Spec.StartingDeadlineSeconds; sdls != nil {
		seconds := fmt.Sprintf("%ds", *sdls)
		sections = append(sections, component.SummarySection{
			Header:  "Starting Deadline Seconds",
			Content: component.NewText(seconds),
		})
	}

	sjhl := cc.cronjob.Spec.SuccessfulJobsHistoryLimit
	fjhl := cc.cronjob.Spec.FailedJobsHistoryLimit

	if sjhl != nil {
		sections.AddText("Successful Job History Limit", fmt.Sprintf("%d", *sjhl))
	}

	if fjhl != nil {
		sections.AddText("Failed Job History Limit", fmt.Sprintf("%d", *fjhl))
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

type cronJobObject interface {
	Config(options Options) error
	Jobs(ctx context.Context, object runtime.Object, options Options) error
}

type cronJobHandler struct {
	cronJob    *batchv1beta1.CronJob
	configFunc func(*batchv1beta1.CronJob, Options) (*component.Summary, error)
	jobFunc    func(context.Context, runtime.Object, Options) (component.Component, error)
	object     *Object
}

var _ cronJobObject = (*cronJobHandler)(nil)

func newCronJobHandler(cronJob *batchv1beta1.CronJob, object *Object) (*cronJobHandler, error) {
	if cronJob == nil {
		return nil, errors.New("can't print a nil cronjob")
	}

	if object == nil {
		return nil, errors.New("can't print cronjob using a nil object printer")
	}

	ch := &cronJobHandler{
		cronJob:    cronJob,
		configFunc: defaultCronJobConfig,
		jobFunc:    defaultCronJobJobs,
		object:     object,
	}
	return ch, nil
}

func (c *cronJobHandler) Config(options Options) error {
	out, err := c.configFunc(c.cronJob, options)
	if err != nil {
		return err
	}
	c.object.RegisterConfig(out)
	return nil
}

func defaultCronJobConfig(cronJob *batchv1beta1.CronJob, options Options) (*component.Summary, error) {
	return NewCronJobConfiguration(cronJob).Create()
}

func (c *cronJobHandler) Jobs(ctx context.Context, object runtime.Object, options Options) error {
	c.object.EnableJobTemplate(c.cronJob.Spec.JobTemplate)

	c.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return c.jobFunc(ctx, object, options)
		},
	})
	return nil
}

func defaultCronJobJobs(ctx context.Context, object runtime.Object, options Options) (component.Component, error) {
	return createJobListView(ctx, object, options)
}
