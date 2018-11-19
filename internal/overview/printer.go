package overview

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/heptio/developer-dash/internal/content"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/apimachinery/pkg/util/sets"
	resourcehelper "k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/helper"
	"k8s.io/kubernetes/pkg/apis/core/helper/qos"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/apis/rbac"
	"k8s.io/kubernetes/pkg/printers/internalversion"
)

func printCronJobSummary(cronJob *batch.CronJob, jobs []*batch.Job) (content.Section, error) {
	section := content.NewSection()

	section.AddText("Name", cronJob.GetName())
	section.AddText("Namespace", cronJob.GetNamespace())
	section.AddLabels("Labels", cronJob.GetLabels())
	section.AddLabels("Annotations", cronJob.GetAnnotations())
	section.AddTimestamp("Create Time", formatTime(&cronJob.CreationTimestamp))

	active := fmt.Sprintf("%d", len(cronJob.Status.Active))
	section.AddText("Active", active)

	section.AddText("Schedule", cronJob.Spec.Schedule)

	var suspend string
	if cronJob.Spec.Suspend == nil {
		suspend = "false"
	} else {
		suspend = fmt.Sprintf("%t", *cronJob.Spec.Suspend)
	}
	section.AddText("Suspend", suspend)

	section.AddTimestamp("Last Schedule", formatTime(cronJob.Status.LastScheduleTime))

	section.AddText("Concurrency Policy", string(cronJob.Spec.ConcurrencyPolicy))

	var startingDeadLine string
	if cronJob.Spec.StartingDeadlineSeconds != nil {
		startingDeadLine = fmt.Sprintf("%ds", *cronJob.Spec.StartingDeadlineSeconds)
	} else {
		startingDeadLine = "<unset>"
	}
	section.AddText("Starting Deadline Seconds", startingDeadLine)

	return section, nil
}

func printDeploymentSummary(deployment *extensions.Deployment) (content.Section, error) {
	section := content.NewSection()

	section.AddText("Name", deployment.GetName())
	section.AddText("Namespace", deployment.GetNamespace())
	section.AddLabels("Labels", deployment.GetLabels())
	section.AddLabels("Annotations", deployment.GetAnnotations())
	section.AddTimestamp("Creation Time", deployment.CreationTimestamp.Time.UTC().Format(time.RFC3339))

	selector, err := metav1.LabelSelectorAsSelector(deployment.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}
	section.AddText("Selector", selector.String())

	section.AddText("Strategy", string(deployment.Spec.Strategy.Type))

	minReadySeconds := fmt.Sprintf("%d", deployment.Spec.MinReadySeconds)
	section.AddText("Min Ready Seconds", minReadySeconds)

	var revisionHistoryLimit string
	if rhl := deployment.Spec.RevisionHistoryLimit; rhl != nil {
		revisionHistoryLimit = fmt.Sprintf("%d", *rhl)
	}
	section.AddText("Revision History Limit", revisionHistoryLimit)

	var rollingUpdateStrategy string
	if rus := deployment.Spec.Strategy.RollingUpdate; rus != nil {
		rollingUpdateStrategy = fmt.Sprintf("Max Surge: %s, Max unavailable: %s",
			rus.MaxSurge.String(), rus.MaxUnavailable.String())
	}
	section.AddText("Rolling Update Strategy", rollingUpdateStrategy)

	status := fmt.Sprintf("%d updated, %d total, %d available, %d unavailable",
		deployment.Status.UpdatedReplicas,
		deployment.Status.Replicas,
		deployment.Status.AvailableReplicas,
		deployment.Status.UnavailableReplicas,
	)
	section.AddText("Status", status)

	return section, nil
}

func printJobSummary(job *batch.Job, pods []*core.Pod) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", job.GetName())
	section.AddText("Namespace", job.GetNamespace())

	selector, err := metav1.LabelSelectorAsSelector(job.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}

	section.AddText("Selector", selector.String())

	section.AddLabels("Labels", job.GetLabels())
	section.AddLabels("Annotations", job.GetAnnotations())

	if controllerRef := metav1.GetControllerOf(job); controllerRef != nil {
		section.AddLink("Controlled By", controllerRef.Name, controlledByPath(controllerRef))
	}

	if p := job.Spec.Parallelism; p != nil {
		section.AddText("Parallelism", fmt.Sprintf("%d", *p))
	}

	if c := job.Spec.Completions; c != nil {
		section.AddText("Completions", fmt.Sprintf("%d", *c))
	} else {
		section.AddText("Completions", "<unset>")
	}

	if st := job.Status.StartTime; st != nil {
		section.AddTimestamp("Start Time", formatTime(st))
	}

	if ads := job.Spec.ActiveDeadlineSeconds; ads != nil {
		section.AddText("Active Deadline Seconds", fmt.Sprintf("%d", *ads))
	}

	ps := createPodStatus(pods)
	section.AddText("Pod Statuses", fmt.Sprintf("%d Running / %d Succeeded / %d Failed",
		ps.Running, ps.Succeeded, ps.Failed))

	return section, nil
}

func printPodSummary(pod *core.Pod, c clock.Clock) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", pod.GetName())
	section.AddText("Namespace", pod.GetNamespace())

	if pod.Spec.Priority != nil {
		section.AddText("Priority", fmt.Sprintf("%d", pod.Spec.Priority))
		section.AddText("PriorityClassName", fmt.Sprintf("%s", stringOrNone(pod.Spec.PriorityClassName)))
	}

	section.AddText("Node", stringOrNone(pod.Spec.NodeName))
	section.AddTimestamp("Start Time", formatTime(pod.Status.StartTime))
	section.AddLabels("Labels", pod.GetLabels())
	section.AddLabels("Annotations", pod.GetAnnotations())

	if pod.DeletionTimestamp != nil {
		ts := translateTimestamp(*pod.DeletionTimestamp, c)
		section.AddText("Status", fmt.Sprintf("Terminating (lasts %s)\n", ts))
		section.AddText("Termination Grace Period", fmt.Sprintf("%d", pod.DeletionGracePeriodSeconds))
	} else {
		section.AddText("Status", string(pod.Status.Phase))
	}

	if pod.Status.Reason != "" {
		section.AddText("Reason", pod.Status.Reason)
	}
	if pod.Status.Message != "" {
		section.AddText("Message", pod.Status.Message)
	}

	section.AddText("IP", pod.Status.PodIP)

	if controllerRef := metav1.GetControllerOf(pod); controllerRef != nil {
		item := content.LinkItem("Controlled By", controllerRef.Name, controlledByPath(controllerRef))
		section.Items = append(section.Items, item)
	}

	if pod.Status.NominatedNodeName != "" {
		section.AddText("NominatedNodeName", pod.Status.NominatedNodeName)
	}

	if pod.Status.QOSClass != "" {
		section.AddText("QoS Class", string(pod.Status.QOSClass))
	} else {
		section.AddText("QoS Class", string(qos.GetPodQOS(pod)))
	}

	section.AddLabels("Node-Selectors", pod.Spec.NodeSelector)
	section.AddLink("Service Account", pod.Spec.ServiceAccountName,
		gvkPath("v1", "ServiceAccount", pod.Spec.ServiceAccountName))

	// TODO add tolerations printer

	return section, nil
}

func printDaemonSetSummary(replicaSet *extensions.DaemonSet, pods []*core.Pod) (content.Section, error) {
	section := content.NewSection()

	section.AddText("Name", replicaSet.GetName())
	section.AddText("Namespace", replicaSet.GetNamespace())

	selector, err := metav1.LabelSelectorAsSelector(replicaSet.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}
	section.AddText("Selector", selector.String())

	section.AddLabels("Labels", replicaSet.GetLabels())
	section.AddLabels("Annotations", replicaSet.GetAnnotations())

	ps := createPodStatus(pods)

	podStatus := fmt.Sprintf("%d Running / %d Waiting / %d Succeeded / %d Failed",
		ps.Running, ps.Waiting, ps.Succeeded, ps.Failed)
	section.AddText("Pod Status", podStatus)

	return section, nil
}

func printReplicaSetSummary(replicaSet *extensions.ReplicaSet, pods []*core.Pod) (content.Section, error) {
	section := content.NewSection()

	section.AddText("Name", replicaSet.GetName())
	section.AddText("Namespace", replicaSet.GetNamespace())

	selector, err := metav1.LabelSelectorAsSelector(replicaSet.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}
	section.AddText("Selector", selector.String())

	section.AddLabels("Labels", replicaSet.GetLabels())
	section.AddLabels("Annotations", replicaSet.GetAnnotations())

	ps := createPodStatus(pods)

	if controllerRef := metav1.GetControllerOf(replicaSet); controllerRef != nil {
		section.AddLink("Controlled By", controllerRef.Name, controlledByPath(controllerRef))
	}

	replicas := fmt.Sprintf("%d current / %d desired",
		replicaSet.Status.Replicas, replicaSet.Spec.Replicas)

	section.AddText("Replicas", replicas)

	podStatus := fmt.Sprintf("%d Running / %d Waiting / %d Succeeded / %d Failed",
		ps.Running, ps.Waiting, ps.Succeeded, ps.Failed)
	section.AddText("Pod Status", podStatus)

	return section, nil
}

func printReplicationControllerSummary(rc *core.ReplicationController, pods []*core.Pod) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", rc.GetName())
	section.AddText("Namespace", rc.GetNamespace())

	ls := &metav1.LabelSelector{MatchLabels: rc.Spec.Selector}
	selector, err := metav1.LabelSelectorAsSelector(ls)
	if err != nil {
		return content.Section{}, err
	}

	section.AddText("Selector", selector.String())

	section.AddLabels("Labels", rc.GetLabels())
	section.AddLabels("Annotations", rc.GetAnnotations())

	section.AddText("Replicas", fmt.Sprintf("%d current / %d desired",
		rc.Status.Replicas, rc.Spec.Replicas))

	ps := createPodStatus(pods)
	podStatus := fmt.Sprintf("%d Running / %d Waiting / %d Succeeded / %d Failed",
		ps.Running, ps.Waiting, ps.Succeeded, ps.Failed)
	section.AddText("Pod Status", podStatus)

	// TODO: add pod template

	return section, nil
}

func printStatefulSetSummary(ss *apps.StatefulSet, pods []*core.Pod) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", ss.GetName())
	section.AddText("Namespace", ss.GetNamespace())
	section.AddText("CreationTimestamp", formatTime(&ss.CreationTimestamp))

	selector, err := metav1.LabelSelectorAsSelector(ss.Spec.Selector)
	if err != nil {
		return content.Section{}, err
	}

	section.AddText("Selector", selector.String())

	section.AddLabels("Labels", ss.GetLabels())
	section.AddLabels("Annotations", ss.GetAnnotations())

	section.AddText("Replicas", fmt.Sprintf("%d current / %d desired",
		ss.Status.Replicas, ss.Spec.Replicas))

	section.AddText("Update Strategy", string(ss.Spec.UpdateStrategy.Type))

	if ss.Spec.UpdateStrategy.RollingUpdate != nil {
		ru := ss.Spec.UpdateStrategy.RollingUpdate
		if ru.Partition != 0 {
			section.AddText("Partition", fmt.Sprintf("%d", ru.Partition))
		}
	}

	ps := createPodStatus(pods)
	podStatus := fmt.Sprintf("%d Running / %d Waiting / %d Succeeded / %d Failed",
		ps.Running, ps.Waiting, ps.Succeeded, ps.Failed)
	section.AddText("Pod Status", podStatus)

	// TODO: add pod template

	return section, nil
}

func printServiceSummary(s *core.Service) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", s.GetName())
	section.AddText("Namespace", s.GetNamespace())

	section.AddLabels("Labels", s.GetLabels())
	section.AddLabels("Annotations", s.GetAnnotations())
	section.AddText("Type", string(s.Spec.Type))
	section.AddText("IP", s.Spec.ClusterIP)

	if len(s.Spec.ExternalIPs) > 0 {
		section.AddText("External IPs", strings.Join(s.Spec.ExternalIPs, ","))
	}

	if s.Spec.LoadBalancerIP != "" {
		section.AddText("LoadBalancer IP", s.Spec.LoadBalancerIP)
	}

	if s.Spec.ExternalName != "" {
		section.AddText("External Name", s.Spec.ExternalName)
	}

	if len(s.Status.LoadBalancer.Ingress) > 0 {
		list := buildIngressString(s.Status.LoadBalancer.Ingress)
		section.AddText("LoadBalancer Ingress", list)
	}

	section.AddText("Session Affinity", string(s.Spec.SessionAffinity))

	if s.Spec.ExternalTrafficPolicy != "" {
		section.AddText("External Traffic Policy", string(s.Spec.ExternalTrafficPolicy))
	}

	if s.Spec.HealthCheckNodePort != 0 {
		section.AddText("HealthCheck NodePort", fmt.Sprintf("%d", s.Spec.HealthCheckNodePort))
	}

	if len(s.Spec.LoadBalancerSourceRanges) > 0 {
		section.AddText("LoadBalancer Source Ranges", strings.Join(s.Spec.LoadBalancerSourceRanges, ","))
	}

	return section, nil
}

func printServiceAccountSummary(serviceAccount *core.ServiceAccount, tokens []*core.Secret, missingSecrets sets.String) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", serviceAccount.GetName())
	section.AddText("Namespace", serviceAccount.GetNamespace())

	section.AddLabels("Labels", serviceAccount.GetLabels())
	section.AddLabels("Annotations", serviceAccount.GetAnnotations())

	var (
		emptyHeader = ""
		pullHeader  = "Image pull secrets"
		mountHeader = "Mountable secrets"
		tokenHeader = "Tokens"

		pullSecretNames  = []string{}
		mountSecretNames = []string{}
		tokenSecretNames = []string{}
	)

	for _, s := range serviceAccount.ImagePullSecrets {
		pullSecretNames = append(pullSecretNames, s.Name)
	}
	for _, s := range serviceAccount.Secrets {
		mountSecretNames = append(mountSecretNames, s.Name)
	}
	for _, s := range tokens {
		tokenSecretNames = append(tokenSecretNames, s.Name)
	}

	types := map[string][]string{
		pullHeader:  pullSecretNames,
		mountHeader: mountSecretNames,
		tokenHeader: tokenSecretNames,
	}
	for _, header := range sets.StringKeySet(types).List() {
		names := types[header]
		if len(names) == 0 {
			section.AddText(header, "<none>")
		} else {
			prefix := header
			for _, name := range names {
				if missingSecrets.Has(name) {
					section.AddText(prefix, fmt.Sprintf("%s (not found)", name))
				} else {
					section.AddLink(prefix, name, gvkPath("v1", "Secret", name))
				}
				prefix = emptyHeader
			}
		}
	}

	return section, nil
}

func printIngressSummary(ingress *v1beta1.Ingress) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", ingress.GetName())
	section.AddText("Namespace", ingress.GetNamespace())

	section.AddLabels("Labels", ingress.GetLabels())
	section.AddLabels("Annotations", ingress.GetAnnotations())

	def := ingress.Spec.Backend
	// ns := ingress.Namespace
	if def == nil {
		def = &v1beta1.IngressBackend{
			ServiceName: "default-http-backend",
			ServicePort: intstr.IntOrString{Type: intstr.Int, IntVal: 80},
		}
		// TODO: re-enable this once there is a backend describer
		// ns = metav1.NamespaceSystem
	}

	// TODO: will need to know the service this points to.
	section.AddText("Default Backend", fmt.Sprintf("%s (%s)",
		backendStringer(def), "<none>"))

	return section, nil
}

func printConfigMapSummary(configMap *core.ConfigMap) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", configMap.GetName())
	section.AddText("Namespace", configMap.GetNamespace())

	section.AddLabels("Labels", configMap.GetLabels())
	section.AddLabels("Annotations", configMap.GetAnnotations())

	return section, nil
}

func printSecretSummary(secret *core.Secret) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", secret.GetName())
	section.AddText("Namespace", secret.GetNamespace())

	section.AddLabels("Labels", secret.GetLabels())
	section.AddLabels("Annotations", secret.GetAnnotations())
	section.AddText("Type", string(secret.Type))

	return section, nil
}

func printPersistentVolumeClaimSummary(pvc *core.PersistentVolumeClaim) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", pvc.GetName())
	section.AddText("Namespace", pvc.GetNamespace())
	section.AddText("StorageClass", getPersistentVolumeClaimClass(pvc))

	if pvc.ObjectMeta.DeletionTimestamp != nil {
		section.AddText("Status", fmt.Sprintf("Terminating (lasts %s)",
			translateTimestamp(*pvc.ObjectMeta.DeletionTimestamp, &clock.RealClock{})),
		)
	} else {
		section.AddText("Status", string(pvc.Status.Phase))
	}

	section.AddText("Volume", pvc.Spec.VolumeName)

	section.AddLabels("Labels", pvc.GetLabels())
	section.AddLabels("Annotations", pvc.GetAnnotations())

	section.AddText("Finalizers", strings.Join(pvc.ObjectMeta.Finalizers, ", "))

	storage := pvc.Spec.Resources.Requests[core.ResourceStorage]
	capacity := ""
	accessModes := ""
	if pvc.Spec.VolumeName != "" {
		accessModes = helper.GetAccessModesAsString(pvc.Status.AccessModes)
		storage = pvc.Status.Capacity[core.ResourceStorage]
		capacity = storage.String()
	}

	section.AddText("Capacity", capacity)
	section.AddText("Access Modes", accessModes)

	if pvc.Spec.VolumeMode != nil {
		section.AddText("VolumeMode", string(*pvc.Spec.VolumeMode))
	}

	return section, nil
}

func printRoleSummary(role *rbac.Role) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", role.GetName())
	section.AddText("Namespace", role.GetNamespace())

	section.AddLabels("Labels", role.GetLabels())
	section.AddLabels("Annotations", role.GetAnnotations())

	return section, nil
}

func printRoleRule(role *rbac.Role) (content.Table, error) {
	table := content.NewTable("Rules")

	columnNames := []string{
		"Resources",
		"Non-Resource URLs",
		"Resource Names",
		"Verbs",
	}

	for _, name := range columnNames {
		table.Columns = append(table.Columns, tableCol(name))
	}

	for _, rule := range role.Rules {
		resources := strings.Join(rule.Resources, ", ")
		nonResourceURLs := "[]"
		if len(rule.NonResourceURLs) > 0 {
			nonResourceURLs = strings.Join(rule.NonResourceURLs, ", ")
		}
		resourceNames := "[]"
		if len(rule.ResourceNames) > 0 {
			nonResourceURLs = strings.Join(rule.ResourceNames, ", ")
		}
		verbs := "[]"
		if len(rule.Verbs) > 0 {
			verbs = strings.Join(rule.Verbs, ", ")
		}

		table.AddRow(content.TableRow{
			columnNames[0]: content.NewStringText(resources),
			columnNames[1]: content.NewStringText(nonResourceURLs),
			columnNames[2]: content.NewStringText(resourceNames),
			columnNames[3]: content.NewStringText(verbs),
		})
	}

	return table, nil
}

func printRoleBindingSummary(roleBinding *rbac.RoleBinding, role *rbac.Role) (content.Section, error) {
	section := content.NewSection()
	section.AddText("Name", roleBinding.GetName())
	section.AddText("Namespace", roleBinding.GetNamespace())

	section.AddLabels("Labels", roleBinding.GetLabels())
	section.AddLabels("Annotations", roleBinding.GetAnnotations())

	section.AddLink("Role", role.GetName(), gvkPath(role.APIVersion, role.Kind, role.Name))

	return section, nil
}

func printRoleBindingSubjects(roleBinding *rbac.RoleBinding) (content.Table, error) {
	table := content.NewTable("Rules")

	columnNames := []string{
		"Kind",
		"Name",
		"Namespace",
	}

	for _, name := range columnNames {
		table.Columns = append(table.Columns, tableCol(name))
	}

	for _, subject := range roleBinding.Subjects {
		kind := subject.Kind
		name := subject.Name
		namespace := subject.Namespace

		table.AddRow(content.TableRow{
			columnNames[0]: content.NewStringText(kind),
			columnNames[1]: content.NewStringText(name),
			columnNames[2]: content.NewStringText(namespace),
		})
	}

	return table, nil
}

func printPodTemplate(template *core.PodTemplateSpec, containerStatuses []core.ContainerStatus) ([]content.Content, error) {

	templateSection := content.NewSection()
	templateSection.AddLabels("Labels", template.Labels)
	if len(template.Annotations) > 0 {
		templateSection.AddList("Annotations", template.Annotations)
	}
	if template.Spec.ServiceAccountName != "" {
		templateSection.AddLink("Service Account", template.Spec.ServiceAccountName,
			gvkPath("v1", "ServiceAccount", template.Spec.ServiceAccountName))
	}

	podTemplateSections := []content.Section{templateSection}
	summary := content.NewSummary("Pod Template", podTemplateSections)

	contents := []content.Content{
		&summary,
	}

	if len(template.Spec.InitContainers) > 0 {
		initContainerSections := describeContainers(template.Spec.InitContainers, nil, nil)
		initContainerSummary := content.NewSummary("Init Container Template", initContainerSections)
		contents = append(contents, &initContainerSummary)
	}

	containerSections := describeContainers(template.Spec.Containers, containerStatuses, nil)
	containerSummary := content.NewSummary("Container Template", containerSections)
	contents = append(contents, &containerSummary)

	containersFromEnv := describeContainersEnvFrom(template.Spec.Containers)
	contents = append(contents, &containersFromEnv)

	return contents, nil
}

func describeContainersEnvFrom(containers []core.Container) content.Table {
	table := content.NewTable("Environment From")

	table.Columns = tableCols("Name", "From", "Prefix", "Optional")

	for _, container := range containers {
		for _, e := range container.EnvFrom {
			from := ""
			name := ""
			optional := false

			if e.ConfigMapRef != nil {
				from = "ConfigMap"
				name = e.ConfigMapRef.Name
				optional = e.ConfigMapRef.Optional != nil && *e.ConfigMapRef.Optional
			} else if e.SecretRef != nil {
				from = "Secret"
				name = e.SecretRef.Name
				optional = e.SecretRef.Optional != nil && *e.SecretRef.Optional
			}

			table.AddRow(content.TableRow{
				"Name":     content.NewStringText(name),
				"From":     content.NewStringText(from),
				"Prefix":   content.NewStringText(e.Prefix),
				"Optional": content.NewStringText(fmt.Sprintf("%t", optional)),
			})
		}
	}

	return table
}

func describeContainers(containers []core.Container, containerStatuses []core.ContainerStatus,
	resolverFn internalversion.EnvVarResolverFunc) []content.Section {
	statuses := make(map[string]core.ContainerStatus)
	for _, status := range containerStatuses {
		statuses[status.Name] = status
	}

	var sections []content.Section
	for _, container := range containers {
		status, ok := statuses[container.Name]
		section := describeContainer(container, status, ok, resolverFn)
		sections = append(sections, section)
	}

	return sections
}

func describeContainer(container core.Container, status core.ContainerStatus, ok bool,
	resolverFn internalversion.EnvVarResolverFunc) content.Section {
	section := content.NewSection()
	section.Title = container.Name

	if ok {
		section.AddText("Container ID", status.ContainerID)
	}
	section.AddText("Image", container.Image)
	if ok {
		section.AddText("Image ID", status.ImageID)
	}

	portString := describeContainerPorts(container.Ports)
	if strings.Contains(portString, ",") {
		section.AddText("Ports", portString)
	} else {
		section.AddText("Port", stringOrNone(portString))
	}

	hostPortString := describeContainerHostPorts(container.Ports)
	if strings.Contains(hostPortString, ",") {
		section.AddText("Host Ports", hostPortString)
	} else {
		section.AddText("Host Port", stringOrNone(hostPortString))
	}

	if len(container.Command) > 0 {
		section.AddText("Command", fmt.Sprintf("[%s]", strings.Join(container.Command, ", ")))
	}

	if len(container.Args) > 0 {
		section.AddText("Args", fmt.Sprintf("[%s]", strings.Join(container.Args, ", ")))
	}

	resources := container.Resources

	if len(resources.Limits) > 0 {
		limits := make(map[string]string)
		for _, name := range internalversion.SortedResourceNames(resources.Limits) {
			quantity := resources.Limits[name]
			limits[string(name)] = quantity.String()
		}

		section.AddList("Limits", limits)
	}

	if len(resources.Requests) > 0 {
		requests := make(map[string]string)
		for _, name := range internalversion.SortedResourceNames(resources.Requests) {
			quantity := resources.Limits[name]
			requests[string(name)] = quantity.String()
		}

		section.AddList("Requests", requests)
	}

	if container.LivenessProbe != nil {
		probe := internalversion.DescribeProbe(container.LivenessProbe)
		section.AddText("Liveness", probe)
	}

	if container.ReadinessProbe != nil {
		probe := internalversion.DescribeProbe(container.ReadinessProbe)
		section.AddText("Readiness", probe)
	}

	if len(container.Env) == 0 {
		section.AddText("Environment", "<none>")
	} else {
		containerEnv := make(map[string]string)

		for _, e := range container.Env {
			if e.ValueFrom == nil {
				containerEnv[e.Name] = e.Value
				continue
			}

			switch {
			case e.ValueFrom.FieldRef != nil:
				var valueFrom string
				if resolverFn != nil {
					valueFrom = resolverFn(e)
				}
				containerEnv[e.Name] = fmt.Sprintf("%s (%s:%s)",
					valueFrom, e.ValueFrom.FieldRef.APIVersion, e.ValueFrom.FieldRef.FieldPath)
			case e.ValueFrom.ResourceFieldRef != nil:
				valueFrom, err := resourcehelper.ExtractContainerResourceValue(e.ValueFrom.ResourceFieldRef, &container)
				if err != nil {
					valueFrom = ""
				}
				resource := e.ValueFrom.ResourceFieldRef.Resource
				if valueFrom == "0" && (resource == "limits.cpu" || resource == "limits.memory") {
					valueFrom = "node allocatable"
				}
				containerEnv[e.Name] = fmt.Sprintf("%s (%s)", valueFrom, resource)
			case e.ValueFrom.SecretKeyRef != nil:
				optional := e.ValueFrom.SecretKeyRef.Optional != nil && *e.ValueFrom.SecretKeyRef.Optional
				containerEnv[e.Name] = fmt.Sprintf("<set to the key '%s' in secret '%s'> Optional: %t",
					e.ValueFrom.SecretKeyRef.Key, e.ValueFrom.SecretKeyRef.Name, optional)
			case e.ValueFrom.ConfigMapKeyRef != nil:
				optional := e.ValueFrom.ConfigMapKeyRef.Optional != nil && *e.ValueFrom.ConfigMapKeyRef.Optional
				containerEnv[e.Name] = fmt.Sprintf("<set to the key '%s' in config map '%s'> Optional: %t",
					e.ValueFrom.ConfigMapKeyRef.Key, e.ValueFrom.ConfigMapKeyRef.Name, optional)
			}
		}

		section.AddList("Environment", containerEnv)
	}

	if len(container.VolumeMounts) == 0 {
		section.AddText("Mounts", "<none>")
	} else {
		mounts := make(map[string]string)
		sort.Sort(internalversion.SortableVolumeMounts(container.VolumeMounts))
		for _, mount := range container.VolumeMounts {
			flags := []string{}
			switch {
			case mount.ReadOnly:
				flags = append(flags, "ro")
			case !mount.ReadOnly:
				flags = append(flags, "rw")
			case mount.SubPath != "":
				flags = append(flags, fmt.Sprintf("path=%q", mount.SubPath))
			}
			mounts[mount.MountPath] = fmt.Sprintf("%s (%s)", mount.Name, strings.Join(flags, ","))
		}
		section.AddList("Mounts", mounts)
	}

	if len(container.VolumeDevices) > 0 {
		sort.Sort(internalversion.SortableVolumeDevices(container.VolumeDevices))
		devices := make(map[string]string)
		for _, device := range container.VolumeDevices {
			devices[device.DevicePath] = device.Name
		}
		section.AddList("Devices", devices)
	}

	return section
}

func describeContainerPorts(cPorts []core.ContainerPort) string {
	ports := make([]string, 0, len(cPorts))
	for _, cPort := range cPorts {
		ports = append(ports, fmt.Sprintf("%d/%s", cPort.ContainerPort, cPort.Protocol))
	}
	return strings.Join(ports, ", ")
}

func describeContainerHostPorts(cPorts []core.ContainerPort) string {
	ports := make([]string, 0, len(cPorts))
	for _, cPort := range cPorts {
		ports = append(ports, fmt.Sprintf("%d/%s", cPort.HostPort, cPort.Protocol))
	}
	return strings.Join(ports, ", ")
}

func stringOrNone(s string) string {
	return stringOrDefaultValue(s, "<none>")
}

func stringOrDefaultValue(s, defaultValue string) string {
	if len(s) > 0 {
		return s
	}
	return defaultValue
}

func gvkPath(apiVersion, kind, name string) string {
	var p string

	switch {
	case apiVersion == "apps/v1" && kind == "DaemonSet":
		p = "/content/overview/workloads/daemon-sets"
	case apiVersion == "extensions/v1beta1" && kind == "ReplicaSet":
		p = "/content/overview/workloads/replica-sets"
	case apiVersion == "apps/v1" && kind == "StatefulSet":
		p = "/content/overview/workloads/stateful-sets"
	case apiVersion == "extensions/v1beta1" && kind == "Deployment":
		p = "/content/overview/workloads/deployments"
	case apiVersion == "batch/v1beta1" && kind == "CronJob":
		p = "/content/overview/workloads/cron-jobs"
	case (apiVersion == "batch/v1beta1" || apiVersion == "batch/v1") && kind == "Job":
		p = "/content/overview/workloads/jobs"
	case apiVersion == "v1" && kind == "ReplicationController":
		p = "/content/overview/workloads/replication-controllers"
	case apiVersion == "v1" && kind == "Secret":
		p = "/content/overview/config-and-storage/secrets"
	case apiVersion == "v1" && kind == "ServiceAccount":
		p = "/content/overview/config-and-storage/service-accounts"
	case apiVersion == "v1" && kind == "Service":
		p = "/content/overview/discovery-and-load-balancing/services"
	case apiVersion == "rbac.authorization.k8s.io/v1" && kind == "Role":
		p = "/content/overview/rbac/roles"
	default:
		panic(fmt.Sprintf("can't convert %s %s to path", apiVersion, kind))
	}

	return path.Join(p, name)

}

func controlledByPath(controllerRef *metav1.OwnerReference) string {
	return gvkPath(controllerRef.APIVersion, controllerRef.Kind, controllerRef.Name)
}

func buildIngressString(ingress []core.LoadBalancerIngress) string {
	var buffer bytes.Buffer

	for i := range ingress {
		if i != 0 {
			buffer.WriteString(", ")
		}
		if ingress[i].IP != "" {
			buffer.WriteString(ingress[i].IP)
		} else {
			buffer.WriteString(ingress[i].Hostname)
		}
	}
	return buffer.String()
}

// loadBalancerWidth is carried over from the kubectl describer
var loadBalancerWidth = 16

// loadBalancerStatusStringer behaves mostly like a string interface and converts the given status to a string.
// `wide` indicates whether the returned value is meant for --o=wide output. If not, it's clipped to 16 bytes.
func loadBalancerStatusStringer(s corev1.LoadBalancerStatus, wide bool) string {
	ingress := s.Ingress
	result := sets.NewString()
	for i := range ingress {
		if ingress[i].IP != "" {
			result.Insert(ingress[i].IP)
		} else if ingress[i].Hostname != "" {
			result.Insert(ingress[i].Hostname)
		}
	}

	r := strings.Join(result.List(), ",")
	if !wide && len(r) > loadBalancerWidth {
		r = r[0:(loadBalancerWidth-3)] + "..."
	}
	return r
}

// backendStringer behaves just like a string interface and converts the given backend to a string.
func backendStringer(backend *v1beta1.IngressBackend) string {
	if backend == nil {
		return ""
	}
	return fmt.Sprintf("%v:%v", backend.ServiceName, backend.ServicePort.String())
}

func getPersistentVolumeClaimClass(claim *core.PersistentVolumeClaim) string {
	// Use beta annotation first
	if class, found := claim.Annotations[core.BetaStorageClassAnnotation]; found {
		return class
	}

	if claim.Spec.StorageClassName != nil {
		return *claim.Spec.StorageClassName
	}

	return ""
}

func getAccessModesAsString(modes []core.PersistentVolumeAccessMode) string {
	modes = removeDuplicateAccessModes(modes)
	modesStr := []string{}
	if containsAccessMode(modes, core.ReadWriteOnce) {
		modesStr = append(modesStr, "RWO")
	}
	if containsAccessMode(modes, core.ReadOnlyMany) {
		modesStr = append(modesStr, "ROX")
	}
	if containsAccessMode(modes, core.ReadWriteMany) {
		modesStr = append(modesStr, "RWX")
	}
	return strings.Join(modesStr, ",")
}

func removeDuplicateAccessModes(modes []core.PersistentVolumeAccessMode) []core.PersistentVolumeAccessMode {
	accessModes := []core.PersistentVolumeAccessMode{}
	for _, m := range modes {
		if !containsAccessMode(accessModes, m) {
			accessModes = append(accessModes, m)
		}
	}
	return accessModes
}

func containsAccessMode(modes []core.PersistentVolumeAccessMode, mode core.PersistentVolumeAccessMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}
