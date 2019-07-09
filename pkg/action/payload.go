/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package action

import (
	"strconv"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Payload is an action payload.
type Payload map[string]interface{}

// GroupVersionKind extracts a GroupVersionKind from a payload.
func (p Payload) GroupVersionKind() (schema.GroupVersionKind, error) {
	group, err := p.String("group")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	version, err := p.String("version")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	kind, err := p.String("kind")
	if err != nil {
		return schema.GroupVersionKind{}, err
	}

	return schema.GroupVersionKind{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, nil
}

// String returns a string from the payload.
func (p Payload) String(key string) (string, error) {
	s, ok := p[key].(string)
	if !ok {
		return "", errors.Errorf("payload does not contain %q", key)
	}

	return s, nil
}

// StringSlice returns a string slice from the payload.
func (p Payload) StringSlice(key string) ([]string, error) {
	sli, ok := p[key].([]interface{})
	if !ok {
		return nil, errors.Errorf("payload does not contain %q", key)
	}

	var list []string
	for i := range sli {
		s, ok := sli[i].(string)
		if !ok {
			return nil, errors.New("could not convert slice entry to string")
		}

		list = append(list, s)
	}

	return list, nil
}

// Float64 returns a float64 from the payload.
func (p Payload) Float64(key string) (float64, error) {
	switch v := p[key].(type) {
	case string:
		return strconv.ParseFloat(v, 64)
	case float64:
		return v, nil
	default:
		return 0, errors.Errorf("unable to handle type %T for %q; got %#v", p[key], key, v)
	}
}
