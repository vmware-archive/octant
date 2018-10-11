package overview

import (
	"fmt"
	"net/url"
	"path"

	"github.com/pkg/errors"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

// CronJobsDescriber creates content for a list of cron jobs.
type CronJobsDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

// NewCronJobsDescriber creates an instance of CronJobsDescriber.
func NewCronJobsDescriber() *CronJobsDescriber {
	return &CronJobsDescriber{
		baseDescriber: newBaseDescriber(),
		cacheKeys: []CacheKey{
			{
				APIVersion: "batch/v1beta1",
				Kind:       "CronJob",
			},
		},
	}
}

// Describe creates content.
func (d *CronJobsDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	var contents []Content

	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	if len(objects) < 1 {
		return contents, nil
	}

	t := newCronJobTable("Cron Jobs")
	for _, object := range objects {
		cj := &batchv1beta1.CronJob{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, cj)
		if err != nil {
			return nil, err
		}

		t.Rows = append(t.Rows, printCronJob(cj, prefix, namespace, d.clock()))
	}

	contents = append(contents, t)

	return contents, nil
}

// CronJobDescriber creates content for a single cron job by name.
type CronJobDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

// NewCronJobDescriber creates an instance of CronJobDescriber.
func NewCronJobDescriber() *CronJobDescriber {
	return &CronJobDescriber{
		baseDescriber: newBaseDescriber(),
		cacheKeys: []CacheKey{
			{
				APIVersion: "batch/v1beta1",
				Kind:       "CronJob",
			},
		},
	}
}

// Describe creates content.
func (d *CronJobDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]Content, error) {
	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	var contents []Content

	t := newCronJobTable("Cron Job")

	if len(objects) != 1 {
		return nil, errors.Errorf("expected 1 cron job")
	}

	object := objects[0]

	cj := &batchv1beta1.CronJob{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, cj)
	if err != nil {
		return nil, err
	}

	t.Rows = append(t.Rows, printCronJob(cj, prefix, namespace, d.clock()))

	contents = append(contents, t)

	eventsTable, err := eventsForObject(object, cache, prefix, namespace, d.clock())
	if err != nil {
		return nil, err
	}

	contents = append(contents, eventsTable)

	return contents, nil
}

func newCronJobTable(name string) table {
	t := newTable(name)

	t.Columns = []tableColumn{
		{Name: "Name", Accessor: "name"},
		{Name: "Labels", Accessor: "labels"},
		{Name: "Schedule", Accessor: "schedule"},
		{Name: "Suspend", Accessor: "suspend"},
		{Name: "Active", Accessor: "active"},
		{Name: "Last Schedule", Accessor: "last_schedule"},
		{Name: "Age", Accessor: "age"},
	}

	return t
}

func printCronJob(cj *batchv1beta1.CronJob, prefix, namespace string, c clock.Clock) tableRow {
	lastScheduleTime := "<none>"
	if cj.Status.LastScheduleTime != nil {
		lastScheduleTime = translateTimestamp(*cj.Status.LastScheduleTime, c)
	}

	values := url.Values{}
	values.Set("namespace", namespace)

	cjPath := fmt.Sprintf("%s?%s",
		path.Join(prefix, "/workloads/cron-jobs", cj.GetName()),
		values.Encode(),
	)

	return tableRow{
		"name":          newLinkText(cj.GetName(), cjPath),
		"labels":        newLabelsText(cj.GetLabels()),
		"schedule":      newStringText(cj.Spec.Schedule),
		"active":        newStringText(fmt.Sprintf("%d", int64(len(cj.Status.Active)))),
		"last_schedule": newStringText(lastScheduleTime),
		"age":           newStringText(translateTimestamp(cj.CreationTimestamp, c)),
	}
}
