package overview

import (
	"context"
	"strings"

	"github.com/heptio/developer-dash/internal/cache"
	"github.com/heptio/developer-dash/internal/view"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/clock"
)

type IngressSummary struct{}

var _ view.View = (*IngressSummary)(nil)

func NewIngressSummary(prefix, namespace string, c clock.Clock) view.View {
	return &IngressSummary{}
}

func (js *IngressSummary) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	ingress, err := retrieveIngress(object)
	if err != nil {
		return nil, err
	}

	detail, err := printIngressSummary(ingress)
	if err != nil {
		return nil, err
	}

	summary := content.NewSummary("Details", []content.Section{detail})
	return []content.Content{
		&summary,
	}, nil
}

type IngressDetails struct{}

var _ view.View = (*IngressDetails)(nil)

func NewIngressDetails(prefix, namespace string, c clock.Clock) view.View {
	return &IngressDetails{}
}

func (ing *IngressDetails) Content(ctx context.Context, object runtime.Object, c cache.Cache) ([]content.Content, error) {
	ingress, err := retrieveIngress(object)
	if err != nil {
		return nil, err
	}

	return []content.Content{
		ingressTLSTable(ingress),
		ingressRulesTable(ingress),
	}, nil
}

func ingressTLSTable(ingress *v1beta1.Ingress) *content.Table {
	table := content.NewTable("TLS", "TLS is not configured for this Ingress")

	table.Columns = []content.TableColumn{
		view.TableCol("Secret"),
		view.TableCol("Hosts"),
	}

	for _, tls := range ingress.Spec.TLS {
		table.AddRow(content.TableRow{
			"Secret": content.NewStringText(tls.SecretName),
			"Hosts":  content.NewStringText(strings.Join(tls.Hosts, ", ")),
		})
	}

	return &table
}

func ingressRulesTable(ingress *v1beta1.Ingress) *content.Table {
	table := content.NewTable("Rules", "Rules are not configured for this Ingress")

	table.Columns = view.TableCols("Host", "Path", "Backend")

	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				backendText := backendStringer(&path.Backend)
				table.AddRow(content.TableRow{
					"Host":    content.NewStringText(rule.Host),
					"Path":    content.NewStringText(path.Path),
					"Backend": content.NewLinkText(backendText, gvkPath("v1", "Service", path.Backend.ServiceName)),
				})
			}
		}
	}

	return &table
}

func retrieveIngress(object runtime.Object) (*v1beta1.Ingress, error) {
	ingress, ok := object.(*v1beta1.Ingress)
	if !ok {
		return nil, errors.Errorf("expected object to be an Ingress, it was %T", object)
	}

	return ingress, nil
}
