/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package javascript

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime/schema"

	jsoniter "github.com/json-iterator/go"

	"github.com/vmware-tanzu/octant/pkg/view/component"
)

var json = jsoniter.ConfigFastest

// ConvertToComponent attempts to convert interface i to a Component.
func ConvertToComponent(name string, i interface{}) (component.Component, error) {
	rawComponent, ok := i.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to get %s map", name)
	}

	rawMetadata, ok := rawComponent["metadata"]
	if !ok {
		return nil, fmt.Errorf("unable to get metadata from %s", name)
	}

	metadataJSON, err := json.Marshal(rawMetadata)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal metadata from: %s: %w", name, err)
	}

	metadata := component.Metadata{}
	if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
		return nil, fmt.Errorf("unable to unmarhal metadata from %s: %w", name, err)
	}

	config, ok := rawComponent["config"]
	if !ok {
		return nil, fmt.Errorf("unable to get config from %s", name)
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal buttonGroup config: %w", err)
	}

	typedObject := component.TypedObject{
		Config:   configJSON,
		Metadata: metadata,
	}

	c, err := typedObject.ToComponent()
	if err != nil {
		return nil, fmt.Errorf("unable to convert buttonGroup to component: %w", err)
	}
	return c, nil
}

// ConvertToItems attempts to convert interface i to a list of FlexLayoutItem.
func ConvertToItems(name string, i interface{}) ([]component.FlexLayoutItem, error) {
	var items []component.FlexLayoutItem

	v, ok := i.([]interface{})
	if !ok {
		return items, fmt.Errorf("unable to parse printHandler %s summary items", name)
	}

	for idx, ii := range v {
		mapItem, ok := ii.(map[string]interface{})
		if !ok {
			return items, fmt.Errorf("unable to parse %s summary items", name)
		}
		flexLayoutItem := component.FlexLayoutItem{}
		jsonSS, err := json.Marshal(mapItem)
		if err != nil {
			return items, fmt.Errorf("unable to marshal json in position %d in %s", idx, name)
		}
		if err := json.Unmarshal(jsonSS, &flexLayoutItem); err != nil {
			return items, fmt.Errorf("unable to unmarshal json in position %d in %s", idx, name)
		}
		items = append(items, flexLayoutItem)
	}

	return items, nil
}

// ConvertToSections attempts to convert interface i to a list of SummarySection.
func ConvertToSections(name string, i interface{}) ([]component.SummarySection, error) {
	var sections []component.SummarySection

	v, ok := i.([]interface{})
	if !ok {
		return sections, fmt.Errorf("unable to parse printHandler %s summary sections", name)
	}

	for idx, ii := range v {
		mapSummarySection, ok := ii.(map[string]interface{})
		if !ok {
			return sections, fmt.Errorf("unable to parse %s summary section", name)
		}
		ss := component.SummarySection{}
		jsonSS, err := json.Marshal(mapSummarySection)
		if err != nil {
			return sections, fmt.Errorf("unable to marshal json GVK in position %d in %s", idx, name)
		}
		if err := json.Unmarshal(jsonSS, &ss); err != nil {
			return sections, fmt.Errorf("unable to unmarshal json GVK in position %d in %s", idx, name)
		}
		sections = append(sections, ss)
	}

	return sections, nil
}

// ConvertToActions attempts to convert interface i to a list of action names.
func ConvertToActions(i interface{}) ([]string, error) {
	actions, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to parse ActionNames")
	}
	actionNames := make([]string, len(actions))
	for i := 0; i < len(actions); i++ {
		actionNames[i] = actions[i].(string)
	}
	return actionNames, nil
}

// ConvertToGVKs attempts to convert interface i to a list of GroupVersionKind.
func ConvertToGVKs(name string, i interface{}) ([]schema.GroupVersionKind, error) {
	GVKs, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s: unable to parse GVK list for supportPrinterConfig", name)
	}
	var gvkList []schema.GroupVersionKind
	for i, ii := range GVKs {
		mapGvk, ok := ii.(map[string]interface{})
		if !ok {
			return gvkList, fmt.Errorf("%s: unable to parse GVK in position %d in supportPrinterConfig", name, i)
		}
		gvk := schema.GroupVersionKind{}

		jsonGvk, err := json.Marshal(mapGvk)
		if err != nil {
			return gvkList, fmt.Errorf("%s: unable to marshal json GVK in position %d in supportPrinterConfig", name, i)
		}

		if err := json.Unmarshal(jsonGvk, &gvk); err != nil {
			return gvkList, fmt.Errorf("%s: unable to unmarshal json GVK in position %d in supportPrinterConfig", name, i)
		}

		gvkList = append(gvkList, gvk)
	}
	return gvkList, nil
}
