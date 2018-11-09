package overview

import (
	"context"
	"fmt"
	"strings"

	"github.com/heptio/developer-dash/internal/content"
	"github.com/pkg/errors"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type IngressSummary struct{}

var _ View = (*IngressSummary)(nil)

func NewIngressSummary() *IngressSummary {
	return &IngressSummary{}
}

func (js *IngressSummary) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
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

var _ View = (*IngressDetails)(nil)

func NewIngressDetails() *IngressDetails {
	return &IngressDetails{}
}

func (ing *IngressDetails) Content(ctx context.Context, object runtime.Object, c Cache) ([]content.Content, error) {
	ingress, ok := object.(*v1beta1.Ingress)
	if !ok {
		return nil, errors.Errorf("expected object to be Ingress, it was %T", object)
	}

	return []content.Content{
		ingressTLSTable(ingress),
		ingressRulesTable(ingress),
	}, nil
}

func ingressTLSTable(ingress *v1beta1.Ingress) *content.Table {
	table := content.NewTable("TLS")

	table.Columns = []content.TableColumn{
		tableCol("Secret"),
		tableCol("Hosts"),
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
	table := content.NewTable("Rules")

	table.Columns = []content.TableColumn{
		tableCol("Host"),
		tableCol("Path"),
		tableCol("Backend"),
	}

	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP != nil {
			for _, path := range rule.HTTP.Paths {
				backend := fmt.Sprintf("%s:%s", path.Backend.ServiceName, path.Backend.ServicePort.String())
				table.AddRow(content.TableRow{
					"Host":    content.NewStringText(rule.Host),
					"Path":    content.NewStringText(path.Path),
					"Backend": content.NewStringText(backend),
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
