package printer

import (
	"context"
	"fmt"
	"strings"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/sets"
)

func IngressListHandler(ctx context.Context, list *extv1beta1.IngressList, options Options) (component.ViewComponent, error) {
	if list == nil {
		return nil, errors.New("ingress list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Hosts", "Address", "Ports", "Age")
	table := component.NewTable("Ingresses", cols)

	for _, ingress := range list.Items {
		ports := "80"
		if len(ingress.Spec.TLS) > 0 {
			ports = "80, 443"
		}

		row := component.TableRow{}
		row["Name"] = link.ForObject(&ingress, ingress.Name)
		row["Labels"] = component.NewLabels(ingress.Labels)
		row["Hosts"] = component.NewText(formatIngressHosts(ingress.Spec.Rules))
		row["Address"] = component.NewText(loadBalancerStatusStringer(ingress.Status.LoadBalancer))
		row["Ports"] = component.NewText(ports)
		row["Age"] = component.NewTimestamp(ingress.CreationTimestamp.Time)

		table.Add(row)
	}

	return table, nil
}

func IngressHandler(ctx context.Context, ingress *extv1beta1.Ingress, options Options) (component.ViewComponent, error) {
	o := NewObject(ingress)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return printIngressConfig(ingress)
	}, 16)

	o.RegisterItems(ItemDescriptor{
		Func: func() (component.ViewComponent, error) {
			return printRulesForIngress(ingress)
		},
		Width: component.WidthFull,
	})

	o.EnableEvents()

	return o.ToComponent(ctx, options)
}

func printIngressConfig(ingress *extv1beta1.Ingress) (component.ViewComponent, error) {
	if ingress == nil {
		return nil, errors.New("ingress is nil")
	}

	var sections component.SummarySections

	if backend := ingress.Spec.Backend; backend != nil {
		backendPath := link.ForGVK(ingress.Namespace, "v1", "Service",
			backend.ServiceName, backendStringer(backend))

		sections.Add("Default Backend", backendPath)
	} else {
		sections.AddText("Default Backend", "Default is not configured")
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func printRulesForIngress(ingress *extv1beta1.Ingress) (component.ViewComponent, error) {
	if ingress == nil {
		return nil, errors.New("ingress is nil")
	}

	cols := component.NewTableCols("Host", "Path", "Backends")
	table := component.NewTable("Rules", cols)

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
			servicePath := link.ForGVK(ingress.Namespace, "v1", "Service",
				path.Backend.ServiceName, backendStringer(&path.Backend))

			table.Add(component.TableRow{
				"Host":     component.NewText(host),
				"Path":     component.NewText(path.Path),
				"Backends": servicePath,
			})
		}
	}

	if backend := ingress.Spec.Backend; ruleCount == 0 && backend != nil {
		servicePath := link.ForGVK(ingress.Namespace, "v1", "Service",
			backend.ServiceName, backendStringer(backend))

		table.Add(component.TableRow{
			"Host":     component.NewText("*"),
			"Path":     component.NewText("*"),
			"Backends": servicePath,
		})

	}

	return table, nil
}

// backendStringer behaves just like a string interface and converts the given backend to a string.
func backendStringer(backend *extv1beta1.IngressBackend) string {
	if backend == nil {
		return ""
	}
	return fmt.Sprintf("%v:%v", backend.ServiceName, backend.ServicePort.String())
}

func formatIngressHosts(rules []extv1beta1.IngressRule) string {
	list := []string{}
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
