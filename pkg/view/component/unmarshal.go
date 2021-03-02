/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"github.com/pkg/errors"
)

func unmarshal(to TypedObject) (Component, error) {
	var o Component
	var err error

	switch to.Metadata.Type {
	case TypeAnnotations:
		t := &Annotations{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal annotations config")
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
	case TypeIFrame:
		t := &IFrame{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal iframe config")
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
	case TypeLoading:
		t := &Loading{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal loading config")
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
	case TypeModal:
		t := &Modal{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal modal config")
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
	case TypeSelectors:
		t := &Selectors{Base: Base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal selectors config")
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

	default:
		return nil, errors.Errorf("unknown view component %q", to.Metadata.Type)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal component")
	}

	return o, nil
}
