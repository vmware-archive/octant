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
			"unmarshal containers config in %q", t.Metadata.Title)
		o = t
	case "expressionSelector":
		t := &ExpressionSelector{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expressionSelector config in %q", t.Metadata.Title)
		o = t
	case "grid":
		t := &Grid{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal grid config in %q", t.Metadata.Title)
		o = t
	case "labels":
		t := &Labels{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labels config in %q", t.Metadata.Title)
		o = t
	case "labelSelector":
		t := &LabelSelector{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labelSelector config in %q", t.Metadata.Title)
		o = t
	case "link":
		t := &Link{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal link config in %q", t.Metadata.Title)
		o = t
	case "list":
		t := &List{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal list config in %q", t.Metadata.Title)
		o = t
	case "panel":
		t := &Panel{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal panel config in %q", t.Metadata.Title)
		o = t
	case "quadrant":
		t := &Quadrant{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal quadrant config in %q", t.Metadata.Title)
		o = t
	case "selectors":
		t := &Selectors{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal selectors config in %q", t.Metadata.Title)
		o = t
	case "summary":
		t := &Summary{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal summary config in %q", t.Metadata.Title)
		o = t
	case "table":
		t := &Table{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal table config in %q", t.Metadata.Title)
		o = t
	case "text":
		t := &Text{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal text config in %q", t.Metadata.Title)
		o = t
	case "timestamp":
		t := &Timestamp{Metadata: to.Metadata}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal timestamp config in %q", t.Metadata.Title)
		o = t

	default:
		return nil, errors.Errorf("unknown view component %q", to.Metadata.Type)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal component")
	}

	return o, nil
}
