package printer

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/heptio/developer-dash/internal/view/gridlayout"
	"github.com/pkg/errors"
)

// ServiceListHandler is a printFunc that lists deployments
func ServiceListHandler(list *corev1.ServiceList, opts Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Type", "Cluster IP", "External IP", "Target Ports", "Age", "Selector")
	tbl := component.NewTable("Services", cols)

	for _, d := range list.Items {
		row := component.TableRow{}
		row["Name"] = linkForObject("v1", "Service", d.Name, d.Name)
		row["Labels"] = component.NewLabels(d.Labels)
		row["Type"] = component.NewText("", string(d.Spec.Type))
		row["Cluster IP"] = component.NewText("", d.Spec.ClusterIP)
		row["External IP"] = component.NewText("", strings.Join(d.Spec.ExternalIPs, ","))
		row["Target Ports"] = printServicePorts(d.Spec.Ports)

		ts := d.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		row["Selector"] = printSelectorMap(d.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// ServiceHandler is a printFunc that prints a Services.
func ServiceHandler(service *corev1.Service, options Options) (component.ViewComponent, error) {
	gl := gridlayout.New()

	configSection := gl.CreateSection(9)

	configView, err := serviceConfiguration(service)
	if err != nil {
		return nil, err
	}

	configSection.Add(configView, 12)

	summaryView, err := serviceSummary(service)
	if err != nil {
		return nil, err
	}

	configSection.Add(summaryView, 12)

	endpointsSection := gl.CreateSection(8)
	endpointsView, err := serviceEndpoints(options.Cache, service)
	if err != nil {
		return nil, err
	}
	endpointsSection.Add(endpointsView, 24)

	grid := gl.ToGrid()

	return grid, nil
}

func printServicePorts(ports []corev1.ServicePort) component.ViewComponent {
	out := make([]string, len(ports))
	for i, port := range ports {
		out[i] = describeTargetPort(port)
	}

	return component.NewText("", strings.Join(out, ", "))
}

func serviceConfiguration(service *corev1.Service) (*component.Summary, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}

	var sections component.SummarySections

	var selectors []component.Selector
	for k, v := range service.Spec.Selector {
		ls := component.NewLabelSelector(k, v)
		selectors = append(selectors, ls)
	}

	sections = append(sections, component.SummarySection{
		Header:  "Selectors",
		Content: component.NewSelectors(selectors),
	})

	sections = append(sections, component.SummarySection{
		Header:  "Type",
		Content: component.NewText("", string(service.Spec.Type)),
	})

	var ports []string
	for _, port := range service.Spec.Ports {
		ports = append(ports, describePort(port))
	}
	sections = append(sections, component.SummarySection{
		Header:  "Ports",
		Content: component.NewText("", strings.Join(ports, ", ")),
	})

	sections = append(sections, component.SummarySection{
		Header:  "Session Affinity",
		Content: component.NewText("", string(service.Spec.SessionAffinity)),
	})

	if service.Spec.ExternalTrafficPolicy != "" {
		sections = append(sections, component.SummarySection{
			Header:  "External Traffic Policy",
			Content: component.NewText("", string(service.Spec.ExternalTrafficPolicy)),
		})
	}

	if service.Spec.HealthCheckNodePort != 0 {
		sections = append(sections, component.SummarySection{
			Header:  "Health Check Node Port",
			Content: component.NewText("", fmt.Sprintf("%d", service.Spec.HealthCheckNodePort)),
		})
	}

	if len(service.Spec.LoadBalancerSourceRanges) > 0 {
		sections = append(sections, component.SummarySection{
			Header:  "Load Balancer Source Ranges",
			Content: component.NewText("", strings.Join(service.Spec.LoadBalancerSourceRanges, ", ")),
		})

	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func serviceSummary(service *corev1.Service) (*component.Summary, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}

	var sections component.SummarySections

	sections = append(sections, component.SummarySection{
		Header:  "Cluster IP",
		Content: component.NewText("", service.Spec.ClusterIP),
	})

	if externalIPs := service.Spec.ExternalIPs; len(externalIPs) > 0 {
		sections = append(sections, component.SummarySection{
			Header:  "External IPs",
			Content: component.NewText("", strings.Join(externalIPs, ", ")),
		})
	}

	if service.Spec.LoadBalancerIP != "" {
		sections = append(sections, component.SummarySection{
			Header:  "Load Balancer IP",
			Content: component.NewText("", service.Spec.LoadBalancerIP),
		})
	}

	if service.Spec.ExternalName != "" {
		sections = append(sections, component.SummarySection{
			Header:  "External Name",
			Content: component.NewText("", service.Spec.ExternalName),
		})
	}

	summary := component.NewSummary("Status", sections...)

	return summary, nil
}

func serviceEndpoints(c cache.Cache, service *corev1.Service) (*component.Table, error) {
	if c == nil {
		return nil, errors.New("cache is nil")
	}

	if service == nil {
		return nil, errors.New("service is nil")
	}

	key := cache.Key{
		Namespace:  service.Namespace,
		APIVersion: "v1",
		Kind:       "Endpoints",
		Name:       service.Name,
	}

	object, err := c.Get(key)
	if err != nil {
		return nil, errors.Wrapf(err, "get endpoints for service %s", service.Name)
	}

	cols := component.NewTableCols("Target", "IP", "Node Name")
	table := component.NewTable("Endpoints", cols)

	endpoints := &corev1.Endpoints{}
	if err := scheme.Scheme.Convert(object, endpoints, 0); err != nil {
		return nil, errors.Wrap(err, "convert unstructured object to endpoints")
	}

	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			row := component.TableRow{}

			var target component.ViewComponent = component.NewText("", "No target")
			if targetRef := address.TargetRef; targetRef != nil {
				ref := gvkPath(targetRef.APIVersion, targetRef.Kind, targetRef.Name)
				target = component.NewLink("", targetRef.Name, ref)
			}

			row["Target"] = target
			row["IP"] = component.NewText("", address.IP)

			nodeName := ""
			if address.NodeName != nil {
				nodeName = *address.NodeName
			}
			row["Node Name"] = component.NewText("", nodeName)

			table.Add(row)
		}
	}

	return table, nil
}

func describeTargetPort(port corev1.ServicePort) string {
	if targetPort := port.TargetPort.String(); targetPort != "0" {
		return fmt.Sprintf("%s/%s", targetPort, port.Protocol)
	}

	return fmt.Sprintf("%d/%s", port.Port, port.Protocol)
}

func describePort(port corev1.ServicePort) string {
	var sb strings.Builder

	if port.Name != "" {
		sb.WriteString(fmt.Sprintf("%s ", port.Name))
	}

	sb.WriteString(fmt.Sprintf("%d", port.Port))

	if port.NodePort != 0 {
		sb.WriteString(fmt.Sprintf(":%d", port.NodePort))
	}

	protocol := port.Protocol
	if protocol == "" {
		protocol = "TCP"
	}
	sb.WriteString(fmt.Sprintf("/%s", protocol))

	if targetPort := port.TargetPort.String(); targetPort != "0" {
		sb.WriteString(fmt.Sprintf(" -> %s", targetPort))
	}

	return sb.String()
}
