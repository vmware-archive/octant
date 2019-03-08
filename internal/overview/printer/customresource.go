package printer

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"

	"github.com/heptio/developer-dash/internal/overview/link"
	"github.com/heptio/developer-dash/internal/view/component"
	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

// CustomResourceListHandler prints a list of custom resources with
// optional custom columns.
func CustomResourceListHandler(
	ctx context.Context,
	name, namespace string,
	crd *apiextv1beta1.CustomResourceDefinition,
	list []*unstructured.Unstructured) (component.ViewComponent, error) {

	hasCustomColumns := len(crd.Spec.AdditionalPrinterColumns) > 0
	if hasCustomColumns {
		return printCustomCRDListTable(name, namespace, crd, list)
	}

	return printGenericCRDTable(name, namespace, list)
}

func printGenericCRDTable(name, namespace string, list []*unstructured.Unstructured) (component.ViewComponent, error) {
	cols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable(name, cols)

	for _, cr := range list {
		row := component.TableRow{}

		row["Name"] = link.ForCustomResource(name, cr)
		row["Labels"] = component.NewLabels(cr.GetLabels())
		row["Age"] = component.NewTimestamp(cr.GetCreationTimestamp().Time)

		table.Add(row)
	}

	return table, nil
}

func printCustomCRDListTable(name, namespace string,
	crd *apiextv1beta1.CustomResourceDefinition,
	list []*unstructured.Unstructured) (component.ViewComponent, error) {

	table := component.NewTable(name, component.NewTableCols("Name", "Labels"))
	for _, column := range crd.Spec.AdditionalPrinterColumns {
		table.AddColumn(column.Name)
	}

	table.AddColumn("Age")

	for _, cr := range list {
		row := component.TableRow{}

		row["Name"] = link.ForCustomResource(name, cr)
		row["Labels"] = component.NewLabels(cr.GetLabels())
		row["Age"] = component.NewTimestamp(cr.GetCreationTimestamp().Time)

		for _, column := range crd.Spec.AdditionalPrinterColumns {
			s, err := printCustomColumn(cr.Object, column)
			if err != nil {
				return nil, errors.Wrapf(err, "print custom column %q in CRD %q",
					column.Name, crd.Name)
			}

			row[column.Name] = component.NewText(s)
		}

		table.Add(row)
	}

	return table, nil
}

func printCustomColumn(m map[string]interface{}, column apiextv1beta1.CustomResourceColumnDefinition) (string, error) {
	j := jsonpath.New(column.Name)
	if err := j.Parse(fmt.Sprintf("{%s}", column.JSONPath)); err != nil {
		return "", errors.Wrap(err, "parsing jsonpath")
	}

	var sb strings.Builder
	if err := j.Execute(&sb, m); err != nil {
		return "", nil
	}

	return sb.String(), nil
}

// CustomResourceHandler prints custom resource objects. If the
// object has columns specified, it will print those columns as well.
func CustomResourceHandler(
	ctx context.Context,
	crd *apiextv1beta1.CustomResourceDefinition,
	object *unstructured.Unstructured,
	options Options) (component.ViewComponent, error) {
	o := NewObject(object)

	o.RegisterConfig(func() (component.ViewComponent, error) {
		return printCustomResourceConfig(object, crd)
	}, 12)
	o.RegisterSummary(func() (component.ViewComponent, error) {
		return printCustomResourceStatus(object, crd)
	}, 12)

	o.EnableEvents()

	view, err := o.ToComponent(ctx, options)
	if err != nil {
		return nil, err
	}

	return view, nil
}

func printCustomResourceConfig(u *unstructured.Unstructured, crd *apiextv1beta1.CustomResourceDefinition) (component.ViewComponent, error) {
	if crd == nil {
		return nil, errors.New("CRD is nil")
	}

	if len(crd.Spec.AdditionalPrinterColumns) < 1 {
		// nothing to do
		return nil, nil
	}

	var sections component.SummarySections

	for _, column := range crd.Spec.AdditionalPrinterColumns {
		if strings.HasPrefix(column.JSONPath, ".spec") {
			s, err := printCustomColumn(u.Object, column)
			if err != nil {
				return nil, errors.Wrap(err, "print custom column")
			}

			if s != "" {
				sections.AddText(column.Name, s)
			}

		}
	}

	return component.NewSummary("Configuration", sections...), nil
}

func printCustomResourceStatus(u *unstructured.Unstructured, crd *apiextv1beta1.CustomResourceDefinition) (component.ViewComponent, error) {
	if crd == nil {
		return nil, errors.New("CRD is nil")
	}

	if len(crd.Spec.AdditionalPrinterColumns) < 1 {
		// nothing to do
		return nil, nil
	}

	var sections component.SummarySections

	for _, column := range crd.Spec.AdditionalPrinterColumns {
		if strings.HasPrefix(column.JSONPath, ".status") {
			s, err := printCustomColumn(u.Object, column)
			if err != nil {
				return nil, errors.Wrap(err, "print custom column")
			}

			sections.AddText(column.Name, s)
		}
	}

	return component.NewSummary("Status", sections...), nil
}
