package printer

import (
	"fmt"
	"strconv"

	"github.com/heptio/developer-dash/internal/overview/link"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/pkg/errors"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

// CronJobListHandler is a printFunc that lists cronjobs
func CronJobListHandler(list *batchv1beta1.CronJobList, opts Options) (component.ViewComponent, error) {
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
func CronJobHandler(c *batchv1beta1.CronJob, opts Options) (component.ViewComponent, error) {
	fl := flexlayout.New()

	configSection := fl.AddSection()
	cronjobConfigGen := NewCronJobConfiguration(c)
	configView, err := cronjobConfigGen.Create()
	if err != nil {
		return nil, err
	}

	if err := configSection.Add(configView, 14); err != nil {
		return nil, errors.Wrap(err, "add cronjob config to layout")
	}

	jobListSection := fl.AddSection()
	jobListTable, err := createJobListView(c, opts)
	if err != nil {
		return nil, errors.Wrap(err, "create job list for cronjob")
	}
	if err := jobListSection.Add(jobListTable, 24); err != nil {
		return nil, errors.Wrap(err, "add job list to layout")
	}

	podTemplate := NewJobTemplate(c, c.Spec.JobTemplate)
	if err = podTemplate.AddToFlexLayout(fl); err != nil {
		return nil, errors.Wrap(err, "add job template to layout")
	}

	if err := createEventsForObject(fl, c, opts); err != nil {
		return nil, errors.Wrap(err, "add events to layout")
	}

	view := fl.ToComponent("Summary")
	return view, nil
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

	creationTimestamp := cc.cronjob.CreationTimestamp.Time
	sections = append(sections, component.SummarySection{
		Header:  "Age",
		Content: component.NewTimestamp(creationTimestamp),
	})

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
