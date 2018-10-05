package overview

import (
	"fmt"
	"net/url"
	"path"
	"sort"
	"time"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/pkg/errors"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/duration"
)

type Describer interface {
	Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]content, error)
}

type baseDescriber struct{}

func newBaseDescriber() *baseDescriber {
	return &baseDescriber{}
}

func (d *baseDescriber) clock() clock.Clock {
	return &clock.RealClock{}
}

func loadObjects(cache Cache, namespace string, fields map[string]string, cacheKeys []CacheKey) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, cacheKey := range cacheKeys {
		cacheKey.Namespace = namespace

		if name, ok := fields["name"]; ok && name != "" {
			cacheKey.Name = name
		}

		objs, err := cache.Retrieve(cacheKey)
		if err != nil {
			return nil, err
		}

		objects = append(objects, objs...)
	}

	return objects, nil
}

type SectionDescriber struct {
	describers []Describer
}

func NewSectionDescriber(describers ...Describer) *SectionDescriber {
	return &SectionDescriber{
		describers: describers,
	}
}

func (d *SectionDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]content, error) {
	var contents []content

	for _, child := range d.describers {
		childContent, err := child.Describe(prefix, namespace, cache, fields)
		if err != nil {
			return nil, err
		}

		contents = append(contents, childContent...)
	}

	return contents, nil
}

type CronJobsDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

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

func (d *CronJobsDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]content, error) {
	var contents []content

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

type CronJobDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

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

func (d *CronJobDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]content, error) {
	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	var contents []content

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

	eventObjects, err := cache.Events(object)
	if err != nil {
		return nil, err
	}

	eventsTable := newEventTable("Events")
	for _, obj := range eventObjects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, event)
		if err != nil {
			return nil, err
		}

		eventsTable.AddRow(printEvent(event, prefix, namespace, d.clock()))
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

type EventsDescriber struct {
	*baseDescriber

	cacheKeys []CacheKey
}

func NewEventsDescriber() *EventsDescriber {
	return &EventsDescriber{
		baseDescriber: newBaseDescriber(),
		cacheKeys: []CacheKey{
			{
				APIVersion: "v1",
				Kind:       "Event",
			},
		},
	}
}

func (d *EventsDescriber) Describe(prefix, namespace string, cache Cache, fields map[string]string) ([]content, error) {
	objects, err := loadObjects(cache, namespace, fields, d.cacheKeys)
	if err != nil {
		return nil, err
	}

	var contents []content

	t := newEventTable("Events")

	sort.Slice(objects, func(i, j int) bool {
		tsI := objects[i].GetCreationTimestamp()
		tsJ := objects[j].GetCreationTimestamp()

		return tsI.Before(&tsJ)
	})

	for _, object := range objects {
		event := &corev1.Event{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(object.Object, event)
		if err != nil {
			return nil, err
		}

		t.Rows = append(t.Rows, printEvent(event, prefix, namespace, d.clock()))
	}

	contents = append(contents, t)

	return contents, nil
}

func newEventTable(name string) table {
	t := newTable(name)

	t.Columns = []tableColumn{
		{Name: "Message", Accessor: "message"},
		{Name: "Source", Accessor: "source"},
		{Name: "Sub-Object", Accessor: "sub_object"},
		{Name: "Count", Accessor: "count"},
		{Name: "First Seen", Accessor: "first_seen"},
		{Name: "Last Seen", Accessor: "last_seen"},
	}

	return t
}

func printEvent(event *corev1.Event, prefix, namespace string, c clock.Clock) tableRow {
	firstSeen := event.FirstTimestamp.Format(time.RFC3339)
	lastSeen := event.LastTimestamp.Format(time.RFC3339)

	return tableRow{
		"message":    newStringText(event.Message),
		"source":     newStringText(event.Source.Component),
		"sub_object": newStringText(""), // TODO: where does this come from?
		"count":      newStringText(fmt.Sprint(event.Count)),
		"first_seen": newStringText(firstSeen),
		"last_seen":  newStringText(lastSeen),
	}
}

// translateTimestamp returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestamp(timestamp metav1.Time, c clock.Clock) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.ShortHumanDuration(c.Since(timestamp.Time))
}
