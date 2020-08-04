/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func unmarshal(to TypedObject) (Component, error) {
	var o Component
	var err error

	switch to.Metadata.Type {
	case typeAnnotations:
		t := &Annotations{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal annotations config")
		o = t
	case typeButtonGroup:
		t := &ButtonGroup{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal buttonGroup config")
		o = t
	case typeCard:
		t := &Card{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal card config")
		o = t
	case typeCardList:
		t := &CardList{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal cardList config")
		o = t
	case typeCodeBlock:
		t := &Code{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal code config")
		o = t
	case typeContainers:
		t := &Containers{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal containers config")
		o = t
	case typeDonutChart:
		t := &DonutChart{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal donutChart config")
		o = t
	case typeEditor:
		t := &Editor{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal editor config")
		o = t
	case typeError:
		t := &Error{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal error config")
		o = t
	case typeExpressionSelector:
		t := &ExpressionSelector{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expressionSelector config")
		o = t
	case typeFlexLayout:
		t := &FlexLayout{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal flexlayout config")
		o = t
	case typeGraphviz:
		t := &Graphviz{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal graphviz config")
		o = t
	case typeGridActions:
		t := &GridActions{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal gridActions config")
		o = t
	case typeIFrame:
		t := &IFrame{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal iframe config")
		o = t
	case typeLabels:
		t := &Labels{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labels config")
		o = t
	case typeLabelSelector:
		t := &LabelSelector{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labelSelector config")
		o = t
	case typeLoading:
		t := &Loading{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal loading config")
		o = t
	case typeLink:
		t := &Link{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal link config")
		o = t
	case typeList:
		t := &List{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal list config")
		o = t
	case typeLogs:
		t := &Logs{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal logs config")
		o = t
	case typeQuadrant:
		t := &Quadrant{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal quadrant config")
		o = t
	case typeResourceViewer:
		t := &ResourceViewer{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal resourceViewer config")
		o = t
	case typeSelectors:
		t := &Selectors{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal selectors config")
		o = t
	case typeSingleStat:
		t := &SingleStat{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal singleStat config")
		o = t
	case typeStepper:
		t := &Stepper{base: base{Metadata: to.Metadata}}
		if uErr := json.Unmarshal(to.Config, &t.Config); uErr != nil {
			err = fmt.Errorf("unmarshal stepper config: %w", uErr)
		}
		o = t
	case typeSummary:
		t := &Summary{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal summary config")
		o = t
	case typeTable:
		t := &Table{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal table config")
		o = t
	case typeText:
		t := &Text{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal text config")
		o = t
	case typeTimestamp:
		t := &Timestamp{base: base{Metadata: to.Metadata}}
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
