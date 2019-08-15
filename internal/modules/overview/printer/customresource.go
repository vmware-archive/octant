/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"

	"github.com/vmware/octant/internal/link"
	dashstrings "github.com/vmware/octant/internal/util/strings"
	"github.com/vmware/octant/pkg/view/component"
)

// CustomResourceListHandler prints a list of custom resources with
// optional custom columns.
func CustomResourceListHandler(
	crdName string,
	crd *apiextv1beta1.CustomResourceDefinition,
	list *unstructured.UnstructuredList,
	linkGenerator link.Interface) (component.Component, error) {

	hasCustomColumns := len(crd.Spec.AdditionalPrinterColumns) > 0
	if hasCustomColumns {
		return printCustomCRDListTable(crdName, crd, list, linkGenerator)
	}

	return printGenericCRDTable(crdName, list, linkGenerator)
}

func printGenericCRDTable(crdName string, list *unstructured.UnstructuredList, linkGenerator link.Interface) (component.Component, error) {
	cols := component.NewTableCols("Name", "Labels", "Age")
	table := component.NewTable(crdName, "We couldn't find any custom resources!", cols)

	for i := range list.Items {
		cr := list.Items[i]
		row := component.TableRow{}

		name, err := linkGenerator.ForObject(&cr, cr.GetName())
		if err != nil {
			return nil, err
		}

		row["Name"] = name
		row["Labels"] = component.NewLabels(cr.GetLabels())
		row["Age"] = component.NewTimestamp(cr.GetCreationTimestamp().Time)

		table.Add(row)
	}

	table.Sort("Name", false)

	return table, nil
}

func printCustomCRDListTable(
	crdName string,
	crd *apiextv1beta1.CustomResourceDefinition,
	list *unstructured.UnstructuredList,
	linkGenerator link.Interface) (component.Component, error) {

	table := component.NewTable(crdName, "We couldn't find any custom resources!", component.NewTableCols("Name", "Labels"))
	for _, column := range crd.Spec.AdditionalPrinterColumns {
		name := column.Name
		if dashstrings.Contains(column.Name, []string{"Name", "Labels", "Age"}) {
			name = fmt.Sprintf("Resource %s", column.Name)
		}
		table.AddColumn(name)
	}

	table.AddColumn("Age")

	for i := range list.Items {
		cr := list.Items[i]
		row := component.TableRow{}

		name, err := linkGenerator.ForObject(&cr, cr.GetName())
		if err != nil {
			return nil, err
		}

		row["Name"] = name
		row["Labels"] = component.NewLabels(cr.GetLabels())
		row["Age"] = component.NewTimestamp(cr.GetCreationTimestamp().Time)

		for _, column := range crd.Spec.AdditionalPrinterColumns {
			s, err := printCustomColumn(cr.Object, column)
			if err != nil {
				return nil, errors.Wrapf(err, "print custom column %q in CRD %q",
					column.Name, crd.Name)
			}

			name := column.Name

			if _, ok := row[column.Name]; ok {
				name = fmt.Sprintf("Resource %s", column.Name)
			}

			row[name] = component.NewText(s)

		}

		table.Add(row)
	}

	table.Sort("Name", false)

	return table, nil
}

func printCustomColumn(m interface{}, column apiextv1beta1.CustomResourceColumnDefinition) (string, error) {
	j := jsonpath.New(column.Name)
	buf := bytes.Buffer{}

	s := strings.Replace(column.JSONPath, "\\", "", -1)

	if err := j.Parse(fmt.Sprintf("{%s}", s)); err != nil {
		return "", errors.Wrapf(err, "jsonpath parse error: %s", s)
	}
	if err := j.Execute(&buf, m); err != nil {
		// inspecting the error string because jsonpath doesn't do typed errors
		if strings.Contains(err.Error(), "is not found") {
			return "<not found>", nil
		}

		return "", errors.Wrapf(err, "jsonpath execute error")
	}

	return buf.String(), nil
}

// CustomResourceHandler prints custom resource objects. If the
// object has columns specified, it will print those columns as well.
func CustomResourceHandler(
	ctx context.Context,
	crd *apiextv1beta1.CustomResourceDefinition,
	object *unstructured.Unstructured,
	options Options) (component.Component, error) {
	o := NewObject(object)

	configSummary, err := printCustomResourceConfig(object, crd)
	if err != nil {
		return nil, err
	}

	statusSummary, err := printCustomResourceStatus(object, crd)
	if err != nil {
		return nil, err
	}

	o.RegisterConfig(configSummary)
	o.RegisterSummary(statusSummary)
	o.EnableEvents()

	view, err := o.ToComponent(ctx, options)
	if err != nil {
		return nil, errors.Wrap(err, "print custom resource")
	}

	return view, nil
}

func printCustomResourceConfig(u *unstructured.Unstructured, crd *apiextv1beta1.CustomResourceDefinition) (*component.Summary, error) {
	if crd == nil {
		return nil, errors.New("CRD is nil")
	}

	summary := component.NewSummary("Configuration")

	if len(crd.Spec.AdditionalPrinterColumns) < 1 {
		return summary, nil
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

	summary.Add(sections...)

	return summary, nil
}

func printCustomResourceStatus(u *unstructured.Unstructured, crd *apiextv1beta1.CustomResourceDefinition) (*component.Summary, error) {
	if crd == nil {
		return nil, errors.New("CRD is nil")
	}

	summary := component.NewSummary("Status")

	if len(crd.Spec.AdditionalPrinterColumns) < 1 {
		return summary, nil
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

	summary.Add(sections...)

	return summary, nil
}
