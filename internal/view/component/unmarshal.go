package component

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func unmarshal(to typedObject) (interface{}, error) {
	var o interface{}
	var err error

	switch to.Metadata.Type {
	case "containers":
		t := &Containers{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal containers config")
		o = t
	case "expressionSelector":
		t := &ExpressionSelector{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expressionSelector config")
		o = t
	case "grid":
		t := &Grid{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal grid config")
		o = t
	case "labels":
		t := &Labels{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labels config")
		o = t
	case "labelSelector":
		t := &LabelSelector{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labelSelector config")
		o = t
	case "link":
		t := &Link{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal link config")
		o = t
	case "list":
		t := &List{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal list config")
		o = t
	case "panel":
		t := &Panel{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal panel config")
		o = t
	case "quadrant":
		t := &Quadrant{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal quadrant config")
		o = t
	case "resourceViewer":
		t := &ResourceViewer{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal resourceViewer config")
		o = t
	case "selectors":
		t := &Selectors{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal selectors config")
		o = t
	case "summary":
		t := &Summary{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal summary config")
		o = t
	case "table":
		t := &Table{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal table config")
		o = t
	case "text":
		t := &Text{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal text config")
		o = t
	case "timestamp":
		t := &Timestamp{Metadata: to.Metadata}
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
