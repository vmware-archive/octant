/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// ServiceListHandler is a printFunc that lists services
func ServiceListHandler(_ context.Context, list *corev1.ServiceList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("nil list")
	}

	cols := component.NewTableCols("Name", "Labels", "Type", "Cluster IP", "External IP", "Target Ports", "Age", "Selector")
	tbl := component.NewTable("Services", "We couldn't find any services!", cols)

	for _, s := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&s, s.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(s.Labels)
		row["Type"] = component.NewText(string(s.Spec.Type))
		row["Cluster IP"] = component.NewText(s.Spec.ClusterIP)
		row["External IP"] = component.NewText(describeExternalIPs(s))
		row["Target Ports"] = printServicePorts(s.Spec.Ports)

		ts := s.CreationTimestamp.Time
		row["Age"] = component.NewTimestamp(ts)

		row["Selector"] = printSelectorMap(s.Spec.Selector)

		tbl.Add(row)
	}
	return tbl, nil
}

// ServiceHandler is a printFunc that prints a Services.
func ServiceHandler(ctx context.Context, service *corev1.Service, options Options) (component.Component, error) {
	o := NewObject(service)

	configSummary, err := serviceConfiguration(ctx, service, options)
	if err != nil {
		return nil, err
	}

	serviceSummary, err := serviceSummary(service)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)
	o.RegisterSummary(serviceSummary)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.Component, error) {
			return serviceEndpoints(ctx, options, service)
		},
		Width: component.WidthFull,
	})

	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

func printServicePorts(ports []corev1.ServicePort) component.Component {
	out := make([]string, len(ports))
	for i, port := range ports {
		out[i] = describeTargetPort(port)
	}

	return component.NewText(strings.Join(out, ", "))
}

func serviceConfiguration(ctx context.Context, service *corev1.Service, options Options) (*component.Summary, error) {
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
		Content: component.NewText(string(service.Spec.Type)),
	})

	var ports []string
	for _, port := range service.Spec.Ports {
		ports = append(ports, describePort(port))
	}
	sections = append(sections, component.SummarySection{
		Header:  "Ports",
		Content: component.NewText(strings.Join(ports, ", ")),
	})

	sections = append(sections, component.SummarySection{
		Header:  "Session Affinity",
		Content: component.NewText(string(service.Spec.SessionAffinity)),
	})

	if service.Spec.ExternalTrafficPolicy != "" {
		sections = append(sections, component.SummarySection{
			Header:  "External Traffic Policy",
			Content: component.NewText(string(service.Spec.ExternalTrafficPolicy)),
		})
	}

	if service.Spec.HealthCheckNodePort != 0 {
		sections = append(sections, component.SummarySection{
			Header:  "Health Check Node Port",
			Content: component.NewText(fmt.Sprintf("%d", service.Spec.HealthCheckNodePort)),
		})
	}

	if len(service.Spec.LoadBalancerSourceRanges) > 0 {
		sections = append(sections, component.SummarySection{
			Header:  "Load Balancer Source Ranges",
			Content: component.NewText(strings.Join(service.Spec.LoadBalancerSourceRanges, ", ")),
		})

	}

	summary := component.NewSummary("Configuration", sections...)

	configEditor, err := editServiceAction(ctx, service, options)
	if err != nil {
		return nil, err
	}
	summary.AddAction(configEditor)

	return summary, nil
}

var (
	selectorKeyPrefixSkipList = []string{
		"pod-template-hash",
	}
)

func editServiceAction(ctx context.Context, service *corev1.Service, options Options) (component.Action, error) {
	if service == nil {
		return component.Action{}, errors.New("service is nil")
	}

	var choices []component.InputChoice
	seenSelectors := make(map[string]bool)

	// add current selectors to list
	for k, v := range service.Spec.Selector {
		value := fmt.Sprintf("%s:%s", k, v)
		choice := component.InputChoice{
			Label:   value,
			Value:   value,
			Checked: true,
		}
		choices = append(choices, choice)
		seenSelectors[value] = true
	}

	// find other possible selectors
	objectStore := options.DashConfig.ObjectStore()
	key := store.Key{
		Namespace:  service.Namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	podList, _, err := objectStore.List(ctx, key)
	if err != nil {
		return component.Action{}, err
	}

	for _, item := range podList.Items {
		for k, v := range item.GetLabels() {
			value := fmt.Sprintf("%s:%s", k, v)
			if _, ok := seenSelectors[value]; ok {
				continue
			}

			skipped := false
			for i := range selectorKeyPrefixSkipList {
				if strings.HasPrefix(value, selectorKeyPrefixSkipList[i]) {
					skipped = true
				}
			}

			if skipped {
				continue
			}

			choice := component.InputChoice{
				Label:   value,
				Value:   value,
				Checked: false,
			}
			choices = append(choices, choice)
			seenSelectors[value] = true
		}
	}

	form, err := component.CreateFormForObject("overview/serviceEditor", service,
		component.NewFormFieldSelect("Selectors", "selectors", choices, true))
	if err != nil {
		return component.Action{}, err
	}

	action := component.Action{
		Name:  "Edit",
		Title: "Service Editor",
		Form:  form,
	}

	return action, nil
}

func serviceSummary(service *corev1.Service) (*component.Summary, error) {
	if service == nil {
		return nil, errors.New("service is nil")
	}

	var sections component.SummarySections

	sections = append(sections, component.SummarySection{
		Header:  "Cluster IP",
		Content: component.NewText(service.Spec.ClusterIP),
	})

	if externalIPs := describeExternalIPs(*service); len(externalIPs) > 0 {
		sections = append(sections, component.SummarySection{
			Header:  "External IPs",
			Content: component.NewText(externalIPs),
		})
	}

	if service.Spec.LoadBalancerIP != "" {
		sections = append(sections, component.SummarySection{
			Header:  "Load Balancer IP",
			Content: component.NewText(service.Spec.LoadBalancerIP),
		})
	}

	if service.Spec.ExternalName != "" {
		sections = append(sections, component.SummarySection{
			Header:  "External Name",
			Content: component.NewText(service.Spec.ExternalName),
		})
	}

	summary := component.NewSummary("Status", sections...)

	return summary, nil
}

func serviceEndpoints(ctx context.Context, options Options, service *corev1.Service) (*component.Table, error) {
	o := options.DashConfig.ObjectStore()

	if o == nil {
		return nil, errors.New("object store is nil")
	}

	if service == nil {
		return nil, errors.New("service is nil")
	}

	key := store.Key{
		Namespace:  service.Namespace,
		APIVersion: "v1",
		Kind:       "Endpoints",
		Name:       service.Name,
	}

	object, found, err := o.Get(ctx, key)
	if err != nil {
		return nil, errors.Wrapf(err, "get endpoints for service %s", service.Name)
	}
	if !found {
		return nil, errors.Errorf("no endpoints for service %s", service.Name)
	}

	cols := component.NewTableCols("Target", "IP", "Node Name")
	table := component.NewTable("Endpoints", "There are no endpoints!", cols)

	endpoints := &corev1.Endpoints{}
	if err := scheme.Scheme.Convert(object, endpoints, 0); err != nil {
		return nil, errors.Wrap(err, "convert unstructured object to endpoints")
	}

	for _, subset := range endpoints.Subsets {
		for _, address := range subset.Addresses {
			row := component.TableRow{}

			var target component.Component = component.NewText("No target")
			if targetRef := address.TargetRef; targetRef != nil {
				// Only references to v1/Pod are possible here
				target, err = options.Link.ForGVK(service.Namespace, "v1", targetRef.Kind,
					targetRef.Name, targetRef.Name)
				if err != nil {
					return nil, err
				}
			}

			row["Target"] = target
			row["IP"] = component.NewText(address.IP)

			nodeName := ""
			if address.NodeName != nil {
				nodeName = *address.NodeName
			}
			row["Node Name"] = component.NewText(nodeName)

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

func describeExternalIPs(service corev1.Service) string {
	externalIPs := make([]string, 0, len(service.Status.LoadBalancer.Ingress))

	if len(service.Spec.ExternalIPs) > 0 {
		return strings.Join(service.Spec.ExternalIPs, ", ")
	}

	for _, ingress := range service.Status.LoadBalancer.Ingress {
		if ingress.Hostname != "" {
			externalIPs = append(externalIPs, ingress.Hostname)
		}
		if ingress.IP != "" {
			externalIPs = append(externalIPs, ingress.IP)
		}
	}

	// TODO: Display if pending
	if len(externalIPs) == 0 {
		return "<none>"
	}
	return strings.Join(externalIPs, ", ")
}
