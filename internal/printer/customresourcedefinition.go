/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package printer

import (
	"context"
	"fmt"
	"strings"

	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"

	"github.com/pkg/errors"

	"github.com/vmware-tanzu/octant/internal/cluster"
	oerrors "github.com/vmware-tanzu/octant/internal/errors"
	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// CustomResourceDefinitionListHandler is a printFunc that lists custom resource definitions.
func CustomResourceDefinitionListHandler(
	ctx context.Context,
	list *apiextv1.CustomResourceDefinitionList,
	opts Options) (component.Component, error) {
	if list == nil {
		return nil, fmt.Errorf("daemon set list is nil")
	}

	cols := component.NewTableCols("Name", "Age")
	ot := NewObjectTable(
		"Custom Resource Definitions",
		"We couldn't find any custom resource definitions!",
		cols,
		opts.DashConfig.ObjectStore())

	for _, crd := range list.Items {
		row := component.TableRow{}

		nameLink, err := opts.Link.ForObject(&crd, crd.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Age"] = component.NewTimestamp(crd.ObjectMeta.CreationTimestamp.Time)

		if err := ot.AddRowForObject(ctx, &crd, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// CustomResourceDefinitionHandler is a print func that prints a custom resource definition.
func CustomResourceDefinitionHandler(
	ctx context.Context,
	crd *apiextv1.CustomResourceDefinition,
	options Options) (component.Component, error) {
	o, err := options.ObjectFactory.Factory(crd, options)
	if err != nil {
		return nil, fmt.Errorf("create object factory: %w", err)
	}

	return o.ToComponent(ctx, options)
}

// CustomResourceDefinitionObjectFactory is an object factory that describes a crd summary.
func CustomResourceDefinitionObjectFactory(object runtime.Object, options Options) (*Object, error) {
	crd, ok := object.(*apiextv1.CustomResourceDefinition)
	if !ok {
		return nil, fmt.Errorf("expected object with gvk %s; got object with gvk %s",
			gvk.CustomResourceDefinition, object.GetObjectKind().GroupVersionKind())
	}

	printObject := NewObject(object)
	printObject.EnableEvents()

	h := NewCustomResourceDefinitionSummary(crd, printObject)
	if err := h.BuildConfig(); err != nil {
		return nil, fmt.Errorf("print crd configuration: %w", err)
	}

	if err := h.BuildItems(options); err != nil {
		return nil, fmt.Errorf("print crd additional items")
	}

	return printObject, nil
}

// CustomResourceDefinitionVersionList lists crd versions.
func CustomResourceDefinitionVersionList(
	ctx context.Context,
	crd *unstructured.Unstructured,
	namespace string,
	options Options) (component.Component, error) {
	crdTool, err := octant.NewCustomResourceDefinitionTool(crd)
	if err != nil {
		return nil, err
	}

	discoveryClient, err := options.DashConfig.ClusterClient().DiscoveryClient()
	if err != nil {
		return nil, fmt.Errorf("discovery client: %w", err)
	}

	resourceLists, err := discoveryClient.ServerPreferredResources()
	if err != nil {
		//TODO: determine the best way to handle these types of errors for all resources, not just metrics.
		if discovery.IsGroupDiscoveryFailedError(err) && strings.Contains(err.Error(), "metrics") {
			logger := log.From(ctx)
			logger.Debugf("preferred resources: %w", err)
		} else {
			return nil, fmt.Errorf("preferred resources: %w", err)
		}
	}

	groupKind, err := crdTool.GroupKind()
	if err != nil {
		return nil, fmt.Errorf("crd group kind: %w", err)
	}

	dri := cluster.NewDiscoveryResourceInfo(resourceLists)
	resourceVersion, err := dri.PreferredVersion(groupKind)
	if err != nil {
		return nil, fmt.Errorf("preferred version for %s: %w", groupKind, err)
	}

	view, err := genViewForCRDVersion(ctx, crd, resourceVersion, namespace, options)
	if err != nil {
		return nil, fmt.Errorf("generate view for CRD %q version %q: %w", crd.GetName(), resourceVersion, err)
	}

	list := component.NewList(nil, []component.Component{view})
	return list, nil
}

func printErrorCard(name, version string, err error) (component.Component, error) {
	errStr := fmt.Sprintf("%s", err)
	title := fmt.Sprintf("%s/%s", name, version)

	var ae *oerrors.AccessError
	if errors.As(err, &ae) {
		errStr = fmt.Sprintf("Access Error, failed to %s: %s", ae.Verb(), ae.Key())
	}

	c := component.NewCard(component.TitleFromString(title))
	c.SetBody(component.NewText(errStr))
	return c, nil
}

func genViewForCRDVersion(
	ctx context.Context,
	crd *unstructured.Unstructured,
	version, namespace string,
	options Options) (component.Component, error) {
	objectStore := options.DashConfig.ObjectStore()

	crGVK, err := gvk.CustomResource(crd, version)
	if err != nil {
		return nil, err
	}

	key := store.KeyFromGroupVersionKind(crGVK)
	key.Namespace = namespace

	customResources, _, err := objectStore.List(ctx, key)
	if err != nil {
		return printErrorCard(crd.GetName(), crd.GetAPIVersion(), err)
	}

	lister := NewCustomResourceLister()
	return lister.List(crd, customResources, version, options.Link)
}

type CRDSummaryFunc func(*apiextv1.CustomResourceDefinition, Options) (component.Component, error)

// CustomResourceDefinitionSummary creates a crd summary.
type CustomResourceDefinitionSummary struct {
	object          ObjectInterface
	crd             *apiextv1.CustomResourceDefinition
	additionalFuncs []CRDSummaryFunc
}

// CustomResourceDefinitionSummaryOption is an option for configuring CustomResourceDefinitionSummary.
type CustomResourceDefinitionSummaryOption func(s *CustomResourceDefinitionSummary)

func CustomResourceDefinitionSummaryItems(funcs ...CRDSummaryFunc) CustomResourceDefinitionSummaryOption {
	return func(s *CustomResourceDefinitionSummary) {
		s.additionalFuncs = funcs
	}
}

// NewCustomResourceDefinitionSummary creates an instance of CustomResourceDefinitionSummary.
func NewCustomResourceDefinitionSummary(
	crd *apiextv1.CustomResourceDefinition,
	object ObjectInterface,
	options ...CustomResourceDefinitionSummaryOption) *CustomResourceDefinitionSummary {
	h := CustomResourceDefinitionSummary{
		crd:             crd,
		object:          object,
		additionalFuncs: defaultCustomResourceDefinitionAdditionalItems,
	}

	for _, option := range options {
		option(&h)
	}

	return &h
}

var defaultCustomResourceDefinitionAdditionalItems = []CRDSummaryFunc{
	func(crd *apiextv1.CustomResourceDefinition, options Options) (component.Component, error) {
		return CreateCRDConditionsTable(crd)
	},
}

var (
	// CRDConditionsColumns are columns for a crd conditions table.
	CRDConditionsColumns = component.NewTableCols("Type", "Status", "Last Transition Time", "Message", "Reason")
)

// CreateCRDConditionsTable creates a crd conditions table.
func CreateCRDConditionsTable(crd *apiextv1.CustomResourceDefinition) (*component.Table, error) {
	table := component.NewTable("Conditions", "", CRDConditionsColumns)

	for _, condition := range crd.Status.Conditions {
		row := component.TableRow{}

		row["Type"] = component.NewText(string(condition.Type))
		row["Status"] = component.NewText(string(condition.Status))
		row["Last Transition Time"] = component.NewTimestamp(condition.LastTransitionTime.Time)
		row["Message"] = component.NewText(condition.Message)
		row["Reason"] = component.NewText(condition.Reason)

		table.Add(row)
	}

	return table, nil
}

// BuildConfig adds configuration data for a crd to the summary.
func (h *CustomResourceDefinitionSummary) BuildConfig() error {
	crd := h.crd
	sections := component.SummarySections{}

	if crd.Spec.Conversion != nil {
		sections.AddText("Conversion Strategy", string(crd.Spec.Conversion.Strategy))
	}

	sections.AddText("Group", crd.Spec.Group)
	sections.AddText("Kind", crd.Spec.Names.Kind)
	sections.AddText("List Kind", crd.Spec.Names.ListKind)
	sections.AddText("Plural", crd.Spec.Names.Plural)
	sections.AddText("Singular", crd.Spec.Names.Singular)
	sections.AddText("Short Names", strings.Join(crd.Spec.Names.ShortNames, ", "))

	if crd.Spec.Names.Categories != nil {
		sections.AddText("Categories", strings.Join(crd.Spec.Names.Categories, ", "))
	}

	summary := component.NewSummary("Configuration", sections...)

	h.object.RegisterConfig(summary)

	return nil
}

// BuildItems adds additional items to the crd summary.
func (h *CustomResourceDefinitionSummary) BuildItems(options Options) error {
	var itemDescriptors []ItemDescriptor

	for i := range h.additionalFuncs {
		c, err := h.additionalFuncs[i](h.crd, options)
		if err != nil {
			return fmt.Errorf("crd item failed: %w", err)
		}

		itemDescriptors = append(itemDescriptors, ItemDescriptor{
			Component: c,
			Width:     component.WidthFull,
		})
	}

	h.object.RegisterItems(itemDescriptors...)

	return nil
}
