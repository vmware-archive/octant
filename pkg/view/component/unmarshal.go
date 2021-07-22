/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"github.com/vmware-tanzu/octant/internal/util/json"

	"github.com/pkg/errors"
)

func unmarshal(to TypedObject) (Component, error) {
	var o Component
	var err error

	switch to.Metadata.Type {
	case TypeAccordion:
		t := &Accordion{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal accordion config")
		o = t
	case TypeAnnotations:
		t := &Annotations{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal annotations config")
		o = t
	case TypeButton:
		t := &Button{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal button config")
		o = t
	case TypeButtonGroup:
		t := &ButtonGroup{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal buttonGroup config")
		o = t
	case TypeCard:
		t := &Card{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal card config")
		o = t
	case TypeCardList:
		t := &CardList{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal cardList config")
		o = t
	case TypeCode:
		t := &Code{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal code config")
		o = t
	case TypeContainers:
		t := &Containers{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal containers config")
		o = t
	case TypeDonutChart:
		t := &DonutChart{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal donutChart config")
		o = t
	case TypeDropdown:
		t := &Dropdown{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal dropdown config")
		o = t
	case TypeEditor:
		t := &Editor{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal editor config")
		o = t
	case TypeError:
		t := &Error{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal error config")
		o = t
	case TypeExpandableRowDetail:
		t := &ExpandableRowDetail{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expandable row detail config")
		o = t
	case TypeExtension:
		t := &Extension{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal extension config")
		o = t
	case TypeExpressionSelector:
		t := &ExpressionSelector{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expressionSelector config")
		o = t
	case TypeFlexLayout:
		t := &FlexLayout{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal flexlayout config")
		o = t
	case TypeFormField:
		t := &FormField{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal form field config")
		o = t
	case TypeGraphviz:
		t := &Graphviz{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal graphviz config")
		o = t
	case TypeGridActions:
		t := &GridActions{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal gridActions config")
		o = t
	case TypeIcon:
		t := &Icon{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal icon config")
		o = t
	case TypeIFrame:
		t := &IFrame{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal iframe config")
		o = t
	case TypeJSONEditor:
		t := &JSONEditor{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal jsonEditor config")
		o = t
	case TypeLabels:
		t := &Labels{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labels config")
		o = t
	case TypeLabelSelector:
		t := &LabelSelector{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labelSelector config")
		o = t
	case TypeLink:
		t := &Link{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal link config")
		o = t
	case TypeList:
		t := &List{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal list config")
		o = t
	case TypeLogs:
		t := &Logs{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal logs config")
		o = t
	case TypeLoading:
		t := &Loading{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal loading config")
		o = t
	case TypeModal:
		t := &Modal{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal modal config")
		o = t
	case TypePodStatus:
		t := &PodStatus{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal pod status config")
		o = t
	case TypePort:
		t := &Port{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal port config")
		o = t
	case TypePorts:
		t := &Ports{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal ports config")
		o = t
	case TypeQuadrant:
		t := &Quadrant{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal quadrant config")
		o = t
	case TypeResourceViewer:
		t := &ResourceViewer{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal resourceViewer config")
		o = t
	case TypeSelectFile:
		t := &SelectFile{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal select file config")
		o = t
	case TypeSelectors:
		t := &Selectors{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal selectors config")
		o = t
	case TypeSignpost:
		t := &Signpost{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal signpost config")
		o = t
	case TypeSingleStat:
		t := &SingleStat{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal singleStat config")
		o = t
	case TypeStepper:
		t := &Stepper{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal stepper config")
		o = t
	case TypeSummary:
		t := &Summary{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal summary config")
		o = t
	case TypeTable:
		t := &Table{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal table config")
		o = t
	case TypeTabsView:
		t := &TabsView{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal tabs config")
		o = t
	case TypeTerminal:
		t := &Terminal{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal terminal config")
		o = t
	case TypeText:
		t := &Text{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal text config")
		o = t
	case TypeTimeline:
		t := &Timeline{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal timeline config")
		o = t
	case TypeTimestamp:
		t := &Timestamp{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal timestamp config")
		o = t
	case TypeYAML:
		t := &YAML{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal YAML config")
		o = t
	default:
		return nil, errors.Errorf("unknown view component %q", to.Metadata.Type)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal component")
	}

	return o, nil
}
