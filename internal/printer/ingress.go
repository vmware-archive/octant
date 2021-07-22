/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// IngressListHandler is a printFunc that prints ingresses
func IngressListHandler(ctx context.Context, list *networkingv1.IngressList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("ingress list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Hosts", "Address", "Ports", "Age")
	ot := NewObjectTable("Ingresses", "We couldn't find any ingresses!", cols, options.DashConfig.ObjectStore())
	ot.EnablePluginStatus(options.DashConfig.PluginManager())
	for _, ingress := range list.Items {
		ports := "80"
		if len(ingress.Spec.TLS) > 0 {
			ports = "80, 443"
		}

		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&ingress, ingress.Name)
		if err != nil {
			return nil, err
		}

		isTLS := len(ingress.Spec.TLS) > 0
		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(ingress.Labels)
		row["Hosts"] = getHostComponent(formatIngressHosts(ingress.Spec.Rules), isTLS)
		row["Address"] = component.NewText(loadBalancerStatusStringer(ingress.Status.LoadBalancer))
		row["Ports"] = component.NewText(ports)
		row["Age"] = component.NewTimestamp(ingress.CreationTimestamp.Time)

		if err := ot.AddRowForObject(ctx, &ingress, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

func getHostComponent(link string, isTLS bool) component.Component {
	_, err := url.Parse(link)
	if strings.Contains(link, ".") && err == nil {
		prefix := "http://"
		if isTLS {
			prefix = "https://"
		}

		links := strings.Split(link, ",")
		if len(links) == 1 {
			return component.NewLink("", link, prefix+link)
		}

		// multiple urls, use Markdown instead
		urls := make([]string, len(links))
		for i := range urls {
			urls[i] = fmt.Sprintf("[%s](%s)", links[i], prefix+links[i])
		}
		return component.NewMarkdownText(strings.Join(urls, ", "))
	} else {
		return component.NewText(link)
	}
}

// IngressHandler is a printFunc that prints an Ingress
func IngressHandler(ctx context.Context, ingress *networkingv1.Ingress, options Options) (component.Component, error) {
	o := NewObject(ingress)
	o.EnableEvents()

	ih, err := newIngressHandler(ingress, o)
	if err != nil {
		return nil, err
	}

	if err := ih.Config(options); err != nil {
		return nil, errors.Wrap(err, "print ingress configuration")
	}

	if err := ih.Rules(options); err != nil {
		return nil, errors.Wrap(err, "print ingress rules")
	}

	return o.ToComponent(ctx, options)
}

// Create creates an ingress configuration summary
func (i *IngressConfiguration) Create(options Options) (*component.Summary, error) {
	if i.ingress == nil {
		return nil, errors.New("ingress is nil")
	}

	ingress := i.ingress

	sections := component.SummarySections{}

	if backend := ingress.Spec.DefaultBackend; backend != nil {
		backendPath, err := options.Link.ForGVK(ingress.Namespace, "v1", "Service",
			backend.Service.Name, backendStringer(backend))
		if err != nil {
			return nil, err
		}

		sections.Add("Default Backend", backendPath)
	} else {
		sections.AddText("Default Backend", "Default is not configured")
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.IngressRuleValue.HTTP == nil {
			continue
		}

		for _, path := range rule.IngressRuleValue.HTTP.Paths {
			if path.Backend.Service.Port.Name == "use-annotation" {
				if action, ok := ingress.Annotations["alb.ingress.kubernetes.io/actions."+path.Backend.Service.Name]; ok {
					sections.Add("Action: "+path.Backend.Service.Name, component.NewText(action))
				}
			}
		}
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func createIngressRulesView(ingress *networkingv1.Ingress, options Options) (*component.Table, error) {
	if ingress == nil {
		return nil, errors.New("ingress is nil")
	}

	cols := component.NewTableCols("Host", "Path", "Backends")
	table := component.NewTable("Rules", "There are no rules defined!", cols)

	ruleCount := 0
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}

		ruleCount++
		host := rule.Host
		if host == "" {
			host = "*"
		}

		for _, path := range rule.HTTP.Paths {
			var servicePath component.Component
			if path.Backend.Service.Port.Name == "use-annotation" {
				servicePath = component.NewMarkdownText("*defined via use-annotation*")
			} else {
				var err error
				servicePath, err = options.Link.ForGVK(ingress.Namespace, "v1", "Service",
					path.Backend.Service.Name, backendStringer(&path.Backend))
				if err != nil {
					return nil, err
				}
			}

			table.Add(component.TableRow{
				"Host":     component.NewText(host),
				"Path":     component.NewText(path.Path),
				"Backends": servicePath,
			})
		}
	}

	if backend := ingress.Spec.DefaultBackend; ruleCount == 0 && backend != nil {
		servicePath, err := options.Link.ForGVK(ingress.Namespace, "v1", "Service",
			backend.Service.Name, backendStringer(backend))
		if err != nil {
			return nil, err
		}

		table.Add(component.TableRow{
			"Host":     component.NewText("*"),
			"Path":     component.NewText("*"),
			"Backends": servicePath,
		})

	}

	return table, nil
}

// backendStringer behaves just like a string interface and converts the given backend to a string.
func backendStringer(backend *networkingv1.IngressBackend) string {
	if backend == nil {
		return ""
	}
	return fmt.Sprintf("%v:%v", backend.Service.Name, backend.Service.Port.Name)
}

func formatIngressHosts(rules []networkingv1.IngressRule) string {
	var list []string
	max := 3
	more := false
	for _, rule := range rules {
		if len(list) == max {
			more = true
		}
		if !more && len(rule.Host) != 0 {
			list = append(list, rule.Host)
		}
	}
	if len(list) == 0 {
		return "*"
	}
	ret := strings.Join(list, ",")
	if more {
		return fmt.Sprintf("%s + %d more...", ret, len(rules)-max)
	}
	return ret
}

// loadBalancerStatusStringer behaves mostly like a string interface and converts the given
// status to a string.
func loadBalancerStatusStringer(s corev1.LoadBalancerStatus) string {
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
	return r
}

// IngressConfiguration generates an ingress configuration
type IngressConfiguration struct {
	ingress *networkingv1.Ingress
}

// NewIngressConfiguration creates an instance of Ingressconfiguration
func NewIngressConfiguration(ingress *networkingv1.Ingress) *IngressConfiguration {
	return &IngressConfiguration{
		ingress: ingress,
	}
}

type ingressObject interface {
	Config(options Options) error
	Rules(options Options) error
}
type ingressHandler struct {
	ingress    *networkingv1.Ingress
	configFunc func(*networkingv1.Ingress, Options) (*component.Summary, error)
	rulesFunc  func(*networkingv1.Ingress, Options) (*component.Table, error)
	object     *Object
}

var _ ingressObject = (*ingressHandler)(nil)

func newIngressHandler(ingress *networkingv1.Ingress, object *Object) (*ingressHandler, error) {
	if ingress == nil {
		return nil, errors.New("can't print a nil ingress")
	}

	if object == nil {
		return nil, errors.New("can't print pod using a nil object")
	}

	ih := &ingressHandler{
		ingress:    ingress,
		configFunc: defaultIngressConfig,
		rulesFunc:  defaultIngressRules,
		object:     object,
	}

	return ih, nil
}

func (i *ingressHandler) Config(options Options) error {
	out, err := i.configFunc(i.ingress, options)
	if err != nil {
		return err
	}
	i.object.RegisterConfig(out)
	return nil
}

func defaultIngressConfig(ingress *networkingv1.Ingress, options Options) (*component.Summary, error) {
	return NewIngressConfiguration(ingress).Create(options)
}

func (i *ingressHandler) Rules(options Options) error {
	if i.ingress == nil {
		return errors.New("can't print rules for nil ingress")
	}

	i.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return i.rulesFunc(i.ingress, options)
		},
	})

	return nil
}

func defaultIngressRules(ingress *networkingv1.Ingress, options Options) (*component.Table, error) {
	return createIngressRulesView(ingress, options)
}
