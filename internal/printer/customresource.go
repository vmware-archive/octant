/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/util/jsonpath"

	"github.com/vmware-tanzu/octant/internal/link"
	"github.com/vmware-tanzu/octant/internal/octant"
	octantStrings "github.com/vmware-tanzu/octant/internal/util/strings"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// CustomResourceListHandler prints a list of custom resources with
// optional custom columns.
func CustomResourceListHandler(crdObject *unstructured.Unstructured, resources *unstructured.UnstructuredList, version string, linkGenerator link.Interface) (component.Component, error) {
	if crdObject == nil {
		return nil, fmt.Errorf("custom resource definition is nil")
	}

	tableName := fmt.Sprintf("%s/%s", crdObject.GetName(), version)
	table := component.NewTable(tableName, "We couldn't find any custom resources!", component.NewTableCols("Name", "Labels"))

	crd, err := octant.NewCustomResourceDefinition(crdObject)
	if err != nil {
		return nil, fmt.Errorf("create custom resource definition parse tool: %w", err)
	}

	if len(resources.Items) > 0 {
		versionName := resources.Items[0].GroupVersionKind().Version
		version, err := crd.Version(versionName)
		if err != nil {
			return nil, fmt.Errorf("get version '%s' from crd %s: %w", versionName, crdObject.GetName(), err)
		}

		for _, column := range version.PrinterColumns {
			name := column.Name
			if octantStrings.Contains(column.Name, []string{"Name", "Labels", "Age"}) {
				name = fmt.Sprintf("Resource %s", column.Name)
			}
			table.AddColumn(name)
		}
		table.AddColumn("Age")
	}

	for i := range resources.Items {
		versionName := resources.Items[i].GroupVersionKind().Version
		version, err := crd.Version(versionName)
		if err != nil {
			return nil, fmt.Errorf("get version '%s' from crd '%s': %w", versionName, crdObject.GetName(), err)
		}

		cr := resources.Items[i]
		row := component.TableRow{}

		name, err := linkGenerator.ForObject(&cr, cr.GetName())
		if err != nil {
			return nil, err
		}

		row["Name"] = name
		row["Labels"] = component.NewLabels(cr.GetLabels())
		row["Age"] = component.NewTimestamp(cr.GetCreationTimestamp().Time)

		for _, column := range version.PrinterColumns {
			s, err := printCustomColumn(cr.Object, column)
			if err != nil {
				return nil, fmt.Errorf("print custom column %q in CRD %q: %w",
					column.Name, crdObject.GetName(), err)
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

func printCustomColumn(m interface{}, column octant.CustomResourceDefinitionPrinterColumn) (string, error) {
	j := jsonpath.New(column.Name)
	buf := bytes.Buffer{}

	s := strings.Replace(column.JSONPath, "\\", "", -1)

	if err := j.Parse(fmt.Sprintf("{%s}", s)); err != nil {
		return "", fmt.Errorf("json path parse error for '%s': %w", s, err)
	}
	if err := j.Execute(&buf, m); err != nil {
		// inspecting the error string because jsonpath doesn't do typed errors
		if strings.Contains(err.Error(), "is not found") {
			return "<not found>", nil
		}

		return "", fmt.Errorf("json path execute error: %w", err)
	}

	return buf.String(), nil
}

// CustomResourceHandler prints custom resource objects. If the
// object has columns specified, it will print those columns as well.
func CustomResourceHandler(ctx context.Context, crd, cr *unstructured.Unstructured, options Options) (component.Component, error) {
	object := NewObject(cr)
	object.EnableEvents()

	h, err := newCustomResourceHandler(crd, cr, object)
	if err != nil {
		return nil, err
	}

	if err := h.Config(); err != nil {
		return nil, fmt.Errorf("print custom resource configuration: %w", err)
	}

	if err := h.Status(); err != nil {
		return nil, fmt.Errorf("print custom resource status: %w", err)
	}

	view, err := object.ToComponent(ctx, options)
	if err != nil {
		return nil, fmt.Errorf("print custom resource: %w", err)
	}

	return view, nil
}

type customResourceObject interface {
	Config() error
	Status() error
}

type customResourceHandler struct {
	statusFunc func(crd, cr *unstructured.Unstructured) (*component.Summary, error)
	configFunc func(crd, cr *unstructured.Unstructured) (*component.Summary, error)
	crd        *unstructured.Unstructured
	cr         *unstructured.Unstructured
	object     *Object
}

var _ customResourceObject = (*customResourceHandler)(nil)

func newCustomResourceHandler(crd, u *unstructured.Unstructured, object *Object) (*customResourceHandler, error) {
	if crd == nil {
		return nil, fmt.Errorf("custom resource definition is nil")
	}

	if u == nil {
		return nil, fmt.Errorf("custom resource is nil")
	}

	if object == nil {
		return nil, fmt.Errorf("can't print custom resource using a nil object printer")
	}

	h := &customResourceHandler{
		crd:        crd,
		cr:         u,
		statusFunc: printCustomResourceStatus,
		configFunc: printCustomResourceConfig,
		object:     object,
	}

	return h, nil
}

func (c *customResourceHandler) Config() error {
	out, err := c.configFunc(c.crd, c.cr)
	if err != nil {
		return err
	}
	c.object.RegisterConfig(out)
	return nil
}

func printCustomResourceConfig(crd, cr *unstructured.Unstructured) (*component.Summary, error) {
	return printCustomResourceSummaryWithPrefix(crd, cr, "Configuration", ".spec")
}

func (c *customResourceHandler) Status() error {
	out, err := c.statusFunc(c.crd, c.cr)
	if err != nil {
		return err
	}
	c.object.RegisterSummary(out)
	return nil
}

func printCustomResourceStatus(crd, cr *unstructured.Unstructured) (*component.Summary, error) {
	return printCustomResourceSummaryWithPrefix(crd, cr, "Status", ".status")
}

func printCustomResourceSummaryWithPrefix(crd, cr *unstructured.Unstructured, title, prefix string) (*component.Summary, error) {
	crdVersion, err := crdVersion(crd, cr)
	if err != nil {
		return nil, fmt.Errorf("fetch crd version: %w", err)
	}

	summary := component.NewSummary(title)

	sections := component.SummarySections{}

	for _, column := range crdVersion.PrinterColumns {
		if strings.HasPrefix(column.JSONPath, prefix) {
			s, err := printCustomColumn(cr.Object, column)
			if err != nil {
				return nil, fmt.Errorf("print custom column '%s': %w", column.Name, err)
			}

			if s != "" {
				sections.AddText(column.Name, s)
			}
		}
	}

	summary.Add(sections...)

	return summary, nil
}

func crdVersion(crd, cr *unstructured.Unstructured) (octant.CustomResourceDefinitionVersion, error) {
	if crd == nil {
		return octant.CustomResourceDefinitionVersion{}, fmt.Errorf("custom resource definition is nil")
	}

	if cr == nil {
		return octant.CustomResourceDefinitionVersion{}, fmt.Errorf("custom resource is nil")
	}

	octantCRD, err := octant.NewCustomResourceDefinition(crd)
	if err != nil {
		return octant.CustomResourceDefinitionVersion{}, fmt.Errorf("create octant CRD: %w", err)
	}

	crdVersion, err := octantCRD.Version(cr.GroupVersionKind().Version)
	if err != nil {
		return octant.CustomResourceDefinitionVersion{}, fmt.Errorf("get version for custom resource: %w", err)
	}

	return crdVersion, nil
}
