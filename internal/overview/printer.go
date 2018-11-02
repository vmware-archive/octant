package overview

import (
	"fmt"
	"path"
	"time"

	"github.com/heptio/developer-dash/internal/content"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func printDeploymentSummary(deployment *extensions.Deployment) (content.Section, error) {
	minReadySeconds := fmt.Sprintf("%d", deployment.Spec.MinReadySeconds)

	var revisionHistoryLimit string
	if rhl := deployment.Spec.RevisionHistoryLimit; rhl != nil {
		revisionHistoryLimit = fmt.Sprintf("%d", *rhl)
	}

	var rollingUpdateStrategy string
	if rus := deployment.Spec.Strategy.RollingUpdate; rus != nil {
		rollingUpdateStrategy = fmt.Sprintf("Max Surge: %s, Max unavailable: %s",
			rus.MaxSurge.String(), rus.MaxUnavailable.String())
	}

	status := fmt.Sprintf("%d updated, %d total, %d available, %d unavailable",
		deployment.Status.UpdatedReplicas,
		deployment.Status.Replicas,
		deployment.Status.AvailableReplicas,
		deployment.Status.UnavailableReplicas,
	)

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}

	section := content.Section{
		Items: []content.Item{
			content.TextItem("Name", deployment.GetName()),
			content.TextItem("Namespace", deployment.GetNamespace()),
			content.LabelsItem("Labels", deployment.GetLabels()),
			content.ListItem("Annotations", deployment.GetAnnotations()),
			content.TextItem("Creation Time", deployment.CreationTimestamp.Time.UTC().Format(time.RFC1123Z)),
			content.TextItem("Selector", selector.String()),
			content.TextItem("Strategy", string(deployment.Spec.Strategy.Type)),
			content.TextItem("Min Ready Seconds", minReadySeconds),
			content.TextItem("Revision History Limit", revisionHistoryLimit),
			content.TextItem("Rolling Update Strategy", rollingUpdateStrategy),
			content.TextItem("Status", status),
		},
	}

	return section, nil
}

func printReplicaSetSummary(replicaSet *extensions.ReplicaSet, pods []*core.Pod) (content.Section, error) {
	selector, err := metav1.LabelSelectorAsSelector(replicaSet.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}

	replicas := fmt.Sprintf("%d current / %d desired",
		replicaSet.Status.Replicas, replicaSet.Spec.Replicas)

	ps := createPodStatus(pods)

	podStatus := fmt.Sprintf("%d Running / %d Waiting / %d Succeeded / %d Failed",
		ps.Running, ps.Waiting, ps.Succeeded, ps.Failed)

	section := content.Section{
		Items: []content.Item{
			content.TextItem("Name", replicaSet.GetName()),
			content.TextItem("Namespace", replicaSet.GetNamespace()),
			content.TextItem("Selector", selector.String()),
			content.LabelsItem("Labels", replicaSet.GetLabels()),
			content.LabelsItem("Annotations", replicaSet.GetAnnotations()),
		},
	}

	if controllerRef := metav1.GetControllerOf(replicaSet); controllerRef != nil {
		linkPath := path.Join("/content/overview/workloads/deployments", controllerRef.Name)
		item := content.LinkItem("Controlled By", controllerRef.Name, linkPath)
		section.Items = append(section.Items, item)
	}

	section.Items = append(section.Items, []content.Item{
		content.TextItem("Replicas", replicas),
		content.TextItem("Pod Status", podStatus),
	}...)

	return section, nil
}
