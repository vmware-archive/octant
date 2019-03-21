package printer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/heptio/developer-dash/internal/overview/link"

	"github.com/heptio/developer-dash/pkg/view/component"
	"github.com/pkg/errors"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

// CronJobListHandler is a printFunc that lists cronjobs
func CronJobListHandler(ctx context.Context, list *batchv1beta1.CronJobList, opts Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Schedule", "Age")
	tbl := component.NewTable("CronJobs", cols)

	for _, c := range list.Items {
		row := component.TableRow{}
		row["Name"] = link.ForObject(&c, c.Name)
		row["Labels"] = component.NewLabels(c.Labels)

		row["Schedule"] = component.NewText(c.Spec.Schedule)

		ts := c.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// CronJobHandler is a printFunc that prints a CronJob
func CronJobHandler(ctx context.Context, c *batchv1beta1.CronJob, opts Options) (component.Component, error) {
	o := NewObject(c)

	cronjobConfigGen := NewCronJobConfiguration(c)
	summary, err := cronjobConfigGen.Create()
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(summary)
	o.EnableJobTemplate(c.Spec.JobTemplate)
	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return createJobListView(ctx, c, opts)
		},
		Width: component.WidthFull,
	})
	o.EnableEvents()

	return o.ToComponent(ctx, opts)
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
