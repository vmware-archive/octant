package printer

import (
	"fmt"

	"github.com/heptio/developer-dash/internal/view/flexlayout"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
)

// StatefulSetListHandler is a printFunc that lists deployments
func StatefulSetListHandler(list *appsv1.StatefulSetList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Desired", "Current", "Age", "Selector")
	tbl := component.NewTable("StatefulSets", cols)

	for _, sts := range list.Items {
		row := component.TableRow{}
		statefulsetPath := gvkPath(sts.Namespace, sts.TypeMeta.APIVersion, sts.TypeMeta.Kind, sts.Name)
		row["Name"] = component.NewLink("", sts.Name, statefulsetPath)
		row["Labels"] = component.NewLabels(sts.Labels)

		desired := fmt.Sprintf("%d", *sts.Spec.Replicas)
		row["Desired"] = component.NewText(desired)

		current := fmt.Sprintf("%d", sts.Status.Replicas)
		row["Current"] = component.NewText(current)

		ts := sts.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		row["Selector"] = printSelector(sts.Spec.Selector)

		tbl.Add(row)
	}

	return tbl, nil
}

// StatefulSetHandler is a printFunc that prints a StatefulSet
func StatefulSetHandler(sts *appsv1.StatefulSet, options Options) (component.ViewComponent, error) {
	fl := flexlayout.New()

	configSection := fl.AddSection()

	stsConfigGen := NewStatefulSetConfiguration(sts)
	configView, err := stsConfigGen.Create()
	if err != nil {
		return nil, err
	}

	if err := configSection.Add(configView, 16); err != nil {
		return nil, errors.Wrap(err, "add replicaset config to layout")
	}

	stsSummaryGen := NewStatefulSetStatus(sts)
	statusView, err := stsSummaryGen.Create(options.Cache)
	if err != nil {
		return nil, err
	}

	if err := configSection.Add(statusView, 8); err != nil {
		return nil, errors.Wrap(err, "add statefulset summary to layout")
	}

	PodTemplate := NewPodTemplate(sts, sts.Spec.Template)
	if err = PodTemplate.AddToFlexLayout(fl); err != nil {
		return nil, errors.Wrap(err, "add pod template to layout")
	}

	view := fl.ToComponent("Summary")

	return view, nil
}

// StatefulSetConfiguration generates a statefulset configuration
type StatefulSetConfiguration struct {
	statefulset *appsv1.StatefulSet
}

// NewStatefulSetConfiguration creates an insteance of StatefulSetconfiguration
func NewStatefulSetConfiguration(sts *appsv1.StatefulSet) *StatefulSetConfiguration {
	return &StatefulSetConfiguration{
		statefulset: sts,
	}
}

// Create generates a statefulset configuration summary
func (sc *StatefulSetConfiguration) Create() (*component.Summary, error) {
	if sc == nil || sc.statefulset == nil {
		return nil, errors.New("statefulset is nil")
	}

	sts := sc.statefulset

	sections := component.SummarySections{}

	sections.AddText("Update Strategy", string(sts.Spec.UpdateStrategy.Type))

	if selector := sts.Spec.Selector; selector != nil {
		var selectors []component.Selector

		for _, lsr := range selector.MatchExpressions {
			o, err := component.MatchOperator(string(lsr.Operator))
			if err != nil {
				return nil, err
			}

			es := component.NewExpressionSelector(lsr.Key, o, lsr.Values)
			selectors = append(selectors, es)
		}

		for k, v := range selector.MatchLabels {
			ls := component.NewLabelSelector(k, v)
			selectors = append(selectors, ls)
		}

		sections = append(sections, component.SummarySection{
			Header:  "Selectors",
			Content: component.NewSelectors(selectors),
		})
	}

	total := fmt.Sprintf("%d", sts.Status.Replicas)

	if desired := sts.Spec.Replicas; desired != nil {
		desired := fmt.Sprintf("%d", *desired)
		status := fmt.Sprintf("%s Desired / %s Total", desired, total)
		sections.AddText("Replicas", status)
	}

	sections.AddText("Pod Management Policy", string(sts.Spec.PodManagementPolicy))

	createTimestamp := sts.CreationTimestamp.Time
	sections = append(sections, component.SummarySection{
		Header:  "Age",
		Content: component.NewTimestamp(createTimestamp),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}

// StatefulSetStatus generates a statefulset status
type StatefulSetStatus struct {
	statefulset *appsv1.StatefulSet
}

// NewStatefulSetStatus creates an instance of StatefulSetStatus
func NewStatefulSetStatus(sts *appsv1.StatefulSet) *StatefulSetStatus {
	return &StatefulSetStatus{
		statefulset: sts,
	}
}

// Create generates a statefulset status quadrant
func (sts *StatefulSetStatus) Create(c cache.Cache) (*component.Quadrant, error) {
	if sts.statefulset == nil {
		return nil, errors.New("statefulset is nil")
	}

	pods, err := listPods(sts.statefulset.ObjectMeta.Namespace, sts.statefulset.Spec.Selector, sts.statefulset.GetUID(), c)
	if err != nil {
		return nil, err
	}

	ps := createPodStatus(pods)

	quadrant := component.NewQuadrant()
	if err := quadrant.Set(component.QuadNW, "Running", fmt.Sprintf("%d", ps.Running)); err != nil {
		return nil, errors.New("unable to set quadrant nw")
	}
	if err := quadrant.Set(component.QuadNE, "Waiting", fmt.Sprintf("%d", ps.Waiting)); err != nil {
		return nil, errors.New("unable to set quadrant ne")
	}
	if err := quadrant.Set(component.QuadSW, "Succeeded", fmt.Sprintf("%d", ps.Succeeded)); err != nil {
		return nil, errors.New("unable to set quadrant sw")
	}
	if err := quadrant.Set(component.QuadSE, "Failed", fmt.Sprintf("%d", ps.Failed)); err != nil {
		return nil, errors.New("unable to set quadrant se")
	}

	return quadrant, nil
}
