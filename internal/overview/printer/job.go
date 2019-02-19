package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/conversion"
	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
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

		row["Name"] = link.ForObject(&job, job.Name)
		row["Labels"] = component.NewLabels(job.Labels)
		row["Completions"] = component.NewText(conversion.PtrInt32ToString(job.Spec.Completions))
		succeeded := fmt.Sprintf("%d", job.Status.Succeeded)
		row["Successful"] = component.NewText(succeeded)
		row["Age"] = component.NewTimestamp(job.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

// JobHandler printers a job.
func JobHandler(job *batchv1.Job, opts Options) (component.ViewComponent, error) {
	o := NewObject(job)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return createJobConfiguration(*job)
	}, 12)

	o.RegisterSummary(func() (component.ViewComponent, error) {
		return createJobStatus(*job)
	}, 12)

	o.EnablePodTemplate(job.Spec.Template)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return createPodListView(job, opts)
		},
		Width: 24,
	})

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return createJobConditions(job.Status.Conditions)
		},
		Width: 24,
	})

	o.EnableEvents()

	return o.ToComponent(opts)
}

func createJobConfiguration(job batchv1.Job) (*component.Summary, error) {
	var sections component.SummarySections

	sections.Add("Back Off Limit", component.NewText(conversion.PtrInt32ToString(job.Spec.BackoffLimit)))
	sections.Add("Completions", component.NewText(conversion.PtrInt32ToString(job.Spec.Completions)))
	sections.Add("Parallelism", component.NewText(conversion.PtrInt32ToString(job.Spec.Parallelism)))

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

func createJobStatus(job batchv1.Job) (*component.Summary, error) {
	var sections component.SummarySections

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

func createJobListView(object runtime.Object, options Options) (component.ViewComponent, error) {
	options.DisableLabels = true

	jobList := &batchv1.JobList{}

	if options.Cache == nil {
		return nil, errors.New("cache is nil")
	}

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

	key := cache.Key{
		Namespace:  namespace,
		APIVersion: "batch/v1beta1",
		Kind:       "Job",
	}

	list, err := options.Cache.List(key)
	if err != nil {
		return nil, errors.Wrapf(err, "list all objects for key %+v", key)
	}

	for _, u := range list {
		job := &batchv1.Job{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, job)
		if err != nil {
			return nil, err
		}

		if err := copyObjectMeta(job, u); err != nil {
			return nil, errors.Wrap(err, "copy object metadata")
		}

		for _, ownerReference := range job.OwnerReferences {
			if ownerReference.APIVersion == apiVersion &&
				ownerReference.Kind == kind &&
				ownerReference.Name == name {
				jobList.Items = append(jobList.Items, *job)
			}
		}
	}

	return JobListHandler(jobList, options)
}
