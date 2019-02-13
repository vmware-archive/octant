package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/flexlayout"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
)

var (
	JobCols = component.NewTableCols("Name", "Labels", "Completions", "Successful", "Age")
)

// JobListHandler prints a job list.
func JobListHandler(list *batchv1.JobList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("job list is nil")
	}

	table := component.NewTable("Jobs", JobCols)

	for _, job := range list.Items {
		row := component.TableRow{}

		jobPath, err := gvkPathFromObject(&job)
		if err != nil {
			return nil, errors.Wrap(err, "get path for job")
		}

		row["Name"] = component.NewLink("", job.Name, jobPath)
		row["Labels"] = component.NewLabels(job.Labels)
		row["Completions"] = component.NewText(ptrInt32ToString(job.Spec.Completions))
		succeeded := fmt.Sprintf("%d", job.Status.Succeeded)
		row["Successful"] = component.NewText(succeeded)
		row["Age"] = component.NewTimestamp(job.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

// JobHandler printers a job.
func JobHandler(job *batchv1.Job, opts Options) (component.ViewComponent, error) {
	if job == nil {
		return nil, errors.New("job is nil")
	}

	fl := flexlayout.New()

	summarySection := fl.AddSection()
	jobConfigView, err := createJobConfiguration(*job)
	if err != nil {
		return nil, errors.Wrap(err, "create job configuration view")
	}
	if err := summarySection.Add(jobConfigView, 12); err != nil {
		return nil, errors.Wrap(err, "add job config to layout")
	}

	jobStatusView, err := createJobStatus(*job)
	if err != nil {
		return nil, errors.Wrap(err, "create job status view")
	}
	if err := summarySection.Add(jobStatusView, 12); err != nil {
		return nil, errors.Wrap(err, "add job status to layout")
	}

	podListSection := fl.AddSection()
	podListTable, err := createPodListView(job, opts)
	if err != nil {
		return nil, errors.Wrap(err, "create pod list for job")
	}
	if err := podListSection.Add(podListTable, 24); err != nil {
		return nil, errors.Wrap(err, "add pod list to layout")
	}

	conditionSection := fl.AddSection()
	conditionTable, err := createJobConditions(job.Status.Conditions)
	if err != nil {
		return nil, errors.Wrap(err, "create job conditions")
	}
	if err := conditionSection.Add(conditionTable, 24); err != nil {
		return nil, errors.Wrap(err, "add job status conditions to layout")
	}

	podTemplate := NewPodTemplate(job, job.Spec.Template)
	if err := podTemplate.AddToFlexLayout(fl); err != nil {
		return nil, errors.Wrap(err, "add pod template to layout")
	}

	if err := createEventsForObject(fl, job, opts); err != nil {
		return nil, errors.Wrap(err, "add events to layout")
	}

	return fl.ToComponent("Summary"), nil
}

func createJobConfiguration(job batchv1.Job) (*component.Summary, error) {
	var sections component.SummarySections

	sections.Add(component.SummarySection{
		Header:  "Back Off Limit",
		Content: component.NewText(ptrInt32ToString(job.Spec.BackoffLimit)),
	})

	sections.Add(component.SummarySection{
		Header:  "Completions",
		Content: component.NewText(ptrInt32ToString(job.Spec.Completions)),
	})

	sections.Add(component.SummarySection{
		Header:  "Parallelism",
		Content: component.NewText(ptrInt32ToString(job.Spec.Parallelism)),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func createJobStatus(job batchv1.Job) (*component.Summary, error) {
	var sections component.SummarySections

	if startTime := job.Status.StartTime; startTime != nil {
		sections.Add(component.SummarySection{
			Header:  "Started",
			Content: component.NewTimestamp(startTime.Time),
		})
	}

	if completionTime := job.Status.CompletionTime; completionTime != nil {
		sections.Add(component.SummarySection{
			Header:  "Completed",
			Content: component.NewTimestamp(completionTime.Time),
		})
	}

	sections.Add(component.SummarySection{
		Header:  "Succeeded",
		Content: component.NewText(fmt.Sprintf("%d", job.Status.Succeeded)),
	})

	summary := component.NewSummary("Status", sections...)
	return summary, nil
}

func createJobConditions(conditions []batchv1.JobCondition) (*component.Table, error) {
	cols := component.NewTableCols("Type", "Last Probe", "Last Transition",
		"Status", "Message", "Reason")
	table := component.NewTable("Conditions", cols)

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
