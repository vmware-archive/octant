package printer

import (
	"fmt"

	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/heptio/developer-dash/internal/view/component"
)

// PodListHandler is a printFunc that prints pods
func PodListHandler(list *corev1.PodList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Ready", "Status", "Restarts", "Age")
	tbl := component.NewTable("Pods", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		row["Name"] = component.NewText("", d.Name)
		row["Labels"] = component.NewLabels(d.Labels)

		readyCounter := 0
		for _, c := range d.Status.ContainerStatuses {
			if c.Ready {
				readyCounter++
			}
		}
		ready := fmt.Sprintf("%d/%d", readyCounter, len(d.Spec.Containers))
		row["Ready"] = component.NewText("", ready)

		status := fmt.Sprintf("%s", d.Status.Phase)
		row["Status"] = component.NewText("", status)

		restartCounter := 0
		for _, c := range d.Status.ContainerStatuses {
			restartCounter += int(c.RestartCount)
		}
		restarts := fmt.Sprintf("%d", restartCounter)
		row["Restarts"] = component.NewText("", restarts)

		ts := d.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		tbl.Add(row)
	}

	return tbl, nil
}

// PodHandler is a printFunc that prints Pods
// TODO: This handler is incomplete
func PodHandler(p *corev1.Pod, opts Options) (component.ViewComponent, error) {
	grid := component.NewGrid("Summary")
	summary := component.NewSummary("Details")
	panel := component.NewPanel("", summary)

	grid.Add(*panel)

	list := component.NewList("", []component.ViewComponent{grid})

	return list, nil
}

type podStatus struct {
	Running   int
	Waiting   int
	Succeeded int
	Failed    int
}

func CreatePodStatus(c corev1client.PodInterface, selector labels.Selector) podStatus {
	var ps podStatus

	options := metav1.ListOptions{LabelSelector: selector.String()}
	pods, err := c.List(options)
	if err != nil {
		return ps
	}

	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			ps.Running++
		case corev1.PodPending:
			ps.Waiting++
		case corev1.PodSucceeded:
			ps.Succeeded++
		case corev1.PodFailed:
			ps.Failed++
		}
	}

	return ps
}

// PodConfiguration generates pod configuration.
type PodConfiguration struct {
	pod *corev1.Pod
}

// NewPodConfiguration creates an instance of PodConfiguration.
func NewPodConfiguration(p *corev1.Pod) *PodConfiguration {
	return &PodConfiguration{
		pod: p,
	}
}

// Create creates a pod configuration summary.
func (p *PodConfiguration) Create() (*component.Summary, error) {
	if p.pod == nil {
		return nil, errors.New("pod is nil")
	}
	pod := p.pod

	sections := component.SummarySections{}

	if pod.Spec.Priority != nil {
		sections.AddText("Priority", fmt.Sprintf("%d", *pod.Spec.Priority))
	}
	if pod.Spec.PriorityClassName != "" {
		sections.AddText("PriorityClassName", pod.Spec.PriorityClassName)
	}

	if pod.Status.StartTime != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Start Time",
			Content: component.NewTimestamp(pod.Status.StartTime.Time),
		})
	}

	if pod.DeletionTimestamp != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Status: Terminating",
			Content: component.NewTimestamp(pod.DeletionTimestamp.Time),
		})
		if pod.DeletionGracePeriodSeconds != nil {
			sections.AddText("Termination Grace Period", fmt.Sprintf("%ds", *pod.DeletionGracePeriodSeconds))
		}
	} else {
		sections.AddText("Status", string(pod.Status.Phase))
	}

	if pod.Status.Reason != "" {
		sections.AddText("Reason", pod.Status.Reason)
	}
	if pod.Status.Message != "" {
		sections.AddText("Message", pod.Status.Message)
	}

	if controllerRef := metav1.GetControllerOf(pod); controllerRef != nil {
		sections = append(sections, component.SummarySection{
			Header:  "Controlled By",
			Content: linkForOwner(controllerRef),
		})
	}

	if pod.Status.NominatedNodeName != "" {
		sections.AddText("NominatedNodeName", pod.Status.NominatedNodeName)
	}

	if pod.Status.QOSClass != "" {
		sections.AddText("QoS Class", string(pod.Status.QOSClass))
	}

	sections = append(sections, component.SummarySection{
		Header: "Service Account",
		Content: linkForObject("v1", "ServiceAccount", pod.Spec.ServiceAccountName,
			pod.Spec.ServiceAccountName),
	})

	summary := component.NewSummary("Configuration", sections...)
	return summary, nil
}
