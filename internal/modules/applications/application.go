/*
 * Copyright (c) 2019 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package applications

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/vmware-tanzu/octant/pkg/store"
)

type descriptor struct {
	name     string
	instance string
	version  string
}

func descriptorFromFields(fields map[string]string) (descriptor, error) {
	if fields == nil {
		return descriptor{}, errors.New("fields is nil")
	}

	name := fields["name"]
	if name == "" {
		return descriptor{}, errors.New("name is blank")
	}

	instance := fields["instance"]
	if instance == "" {
		return descriptor{}, errors.New("instance is blank")
	}

	version := fields["version"]
	if version == "" {
		return descriptor{}, errors.New("version is blank")
	}

	return descriptor{
		name:     name,
		instance: instance,
		version:  version,
	}, nil
}

func (d *descriptor) applicationTitle() string {
	return fmt.Sprintf("%s@%s %s", d.name, d.instance, d.version)
}

type metadata struct {
	podCount int
}

type application struct {
	Name     string
	Instance string
	Version  string
	State    string
	PodCount int
}

func (a *application) Title() string {
	return fmt.Sprintf("%s@%s %s", a.Name, a.Instance, a.Version)
}

func (a *application) Path(prefix, namespace string) string {
	return filepath.Join("/", prefix, "namespace", namespace, a.Name, a.Instance, a.Version)
}

func listApplications(ctx context.Context, objectStore store.Store, namespace string) ([]application, error) {
	key := store.Key{
		Namespace:  namespace,
		APIVersion: "v1",
		Kind:       "Pod",
	}

	pods, _, err := objectStore.List(ctx, key)
	if err != nil {
		return nil, err
	}

	apps := make(map[descriptor]metadata)
	for i := range pods.Items {
		d, belongs, err := podBelongsToApplication(&pods.Items[i])
		if err != nil {
			return nil, err
		}

		if !belongs {
			continue
		}

		m := apps[d]
		m.podCount++
		apps[d] = m
	}

	var list []application

	for _, d := range sortedDescriptorList(apps) {
		a := application{
			Name:     d.name,
			Instance: d.instance,
			Version:  d.version,
			State:    "state",
			PodCount: apps[d].podCount,
		}
		list = append(list, a)
	}

	return list, nil
}

func sortedDescriptorList(apps map[descriptor]metadata) []descriptor {
	var descriptors []descriptor
	for d := range apps {
		descriptors = append(descriptors, d)
	}
	sort.Slice(descriptors, func(i, j int) bool {
		if descriptors[i].name < descriptors[j].name {
			return true
		}
		if descriptors[i].name > descriptors[j].name {
			return false
		}
		if descriptors[i].instance < descriptors[j].instance {
			return true
		}
		if descriptors[i].instance > descriptors[j].instance {
			return false
		}
		return descriptors[i].version < descriptors[j].version
	})
	return descriptors
}

func podBelongsToApplication(object *unstructured.Unstructured) (descriptor, bool, error) {
	if object == nil {
		return descriptor{}, false, errors.New("object is nil")
	}

	apiVersion, kind := object.GetObjectKind().GroupVersionKind().ToAPIVersionAndKind()
	if apiVersion != "v1" && kind != "Pod" {
		return descriptor{}, false, errors.New("object is not a v1 pod")
	}

	appName, err := getLabel(object, appLabelName)
	if err != nil {
		return descriptor{}, false, nil
	}
	appInstance, err := getLabel(object, appLabelInstance)
	if err != nil {
		return descriptor{}, false, nil
	}
	appVersion, err := getLabel(object, appLabelVersion)
	if err != nil {
		return descriptor{}, false, nil
	}

	d := descriptor{
		name:     appName,
		instance: appInstance,
		version:  appVersion,
	}

	return d, true, nil
}

func getLabel(object *unstructured.Unstructured, key string) (string, error) {
	if object == nil {
		return "", errors.New("object is nil")
	}

	if key == "" {
		return "", errors.New("key is blank")
	}

	value := object.GetLabels()[key]
	if value == "" {
		return "", errors.Errorf("value for key %s is blank", key)
	}

	return value, nil
}
