/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/gob"
	"log"
	"os"
	"reflect"

	"github.com/vmware-tanzu/octant/pkg/view/component"
	"github.com/vmware-tanzu/octant/pkg/view/tsgen"
)

func main() {
	m := &tsgen.Model{
		ComponentNames: []string{"component.Link", "component.Text"},
	}

	jobs := []struct {
		Name   string
		TSName string
		Type   reflect.Type
	}{
		{
			Name:   "Link",
			TSName: component.TypeLink,
			Type:   reflect.TypeOf(component.LinkConfig{}),
		},
		{
			Name:   "Text",
			TSName: component.TypeText,
			Type:   reflect.TypeOf(component.TextConfig{}),
		},
	}

	for _, job := range jobs {
		c := tsgen.Component{
			Name:   job.Name,
			TSName: job.TSName,
		}

		xType := job.Type

		for i := 0; i < xType.NumField(); i++ {
			f, ok := tsgen.ConvertField(xType, i, m.ComponentNames)
			if !ok {
				continue
			}
			c.Fields = append(c.Fields, f)
		}

		m.Components = append(m.Components, c)
	}

	enc := gob.NewEncoder(os.Stdout)
	if err := enc.Encode(m); err != nil {
		log.Fatal("encode error:", err)
	}
}
