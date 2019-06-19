/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

import (
	"encoding/json"

	"github.com/vmware/octant/pkg/plugin/dashboard"
	"github.com/vmware/octant/pkg/view/component"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func convertToCapabilities(in *dashboard.RegisterResponse_Capabilities) Capabilities {
	if in == nil {
		return Capabilities{}
	}

	c := Capabilities{
		SupportsPrinterStatus: convertToGroupVersionKindList(in.SupportsPrinterConfig),
		SupportsPrinterConfig: convertToGroupVersionKindList(in.SupportsPrinterConfig),
		SupportsPrinterItems:  convertToGroupVersionKindList(in.SupportsPrinterConfig),
		SupportsObjectStatus:  convertToGroupVersionKindList(in.SupportsPrinterConfig),
		SupportsTab:           convertToGroupVersionKindList(in.SupportsPrinterConfig),
	}

	return c
}

func convertFromCapabilities(in Capabilities) dashboard.RegisterResponse_Capabilities {
	c := dashboard.RegisterResponse_Capabilities{
		SupportsPrinterStatus: convertFromGroupVersionKindList(in.SupportsObjectStatus),
		SupportsPrinterConfig: convertFromGroupVersionKindList(in.SupportsPrinterConfig),
		SupportsPrinterItems:  convertFromGroupVersionKindList(in.SupportsPrinterItems),
		SupportsObjectStatus:  convertFromGroupVersionKindList(in.SupportsObjectStatus),
		SupportsTab:           convertFromGroupVersionKindList(in.SupportsTab),
	}

	return c
}

func convertToGroupVersionKindList(in []*dashboard.RegisterResponse_GroupVersionKind) []schema.GroupVersionKind {
	var list []schema.GroupVersionKind
	for i := range in {
		list = append(list, convertToGroupVersionKind(*in[i]))
	}

	return list
}

func convertFromGroupVersionKindList(in []schema.GroupVersionKind) []*dashboard.RegisterResponse_GroupVersionKind {
	var list []*dashboard.RegisterResponse_GroupVersionKind
	for i := range in {
		item := convertFromGroupVersionKind(in[i])
		list = append(list, &item)
	}

	return list
}

func convertToGroupVersionKind(in dashboard.RegisterResponse_GroupVersionKind) schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   in.Group,
		Version: in.Version,
		Kind:    in.Kind,
	}
}

func convertFromGroupVersionKind(in schema.GroupVersionKind) dashboard.RegisterResponse_GroupVersionKind {
	return dashboard.RegisterResponse_GroupVersionKind{
		Group:   in.Group,
		Version: in.Version,
		Kind:    in.Kind,
	}
}

func convertToSummarySections(in []*dashboard.PrintResponse_SummaryItem) ([]component.SummarySection, error) {
	var list []component.SummarySection

	for _, item := range in {
		converted, err := convertToSummarySection(*item)
		if err != nil {
			return nil, err
		}
		list = append(list, converted)
	}

	return list, nil
}

func convertFromSummarySections(in []component.SummarySection) ([]*dashboard.PrintResponse_SummaryItem, error) {
	var list []*dashboard.PrintResponse_SummaryItem

	for _, section := range in {
		converted, err := convertFromSummarySection(section)
		if err != nil {
			return nil, err
		}
		list = append(list, converted)
	}

	return list, nil
}

func convertToSummarySection(in dashboard.PrintResponse_SummaryItem) (component.SummarySection, error) {
	var typedObject component.TypedObject
	err := json.Unmarshal(in.Component, &typedObject)
	if err != nil {
		return component.SummarySection{}, err
	}

	view, err := typedObject.ToComponent()
	if err != nil {
		return component.SummarySection{}, err
	}

	return component.SummarySection{
		Header:  in.Header,
		Content: view,
	}, nil
}

func convertFromSummarySection(in component.SummarySection) (*dashboard.PrintResponse_SummaryItem, error) {
	data, err := json.Marshal(in.Content)
	if err != nil {
		return nil, err
	}

	return &dashboard.PrintResponse_SummaryItem{
		Header:    in.Header,
		Component: data,
	}, nil
}
