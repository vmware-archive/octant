/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GroupVersionParserFunc is a function that parses a string group/version into a schema GroupVersion.
// If it is unable, it returns an error.
type GroupVersionParserFunc func(groupVersion string) (schema.GroupVersion, error)

// DiscoveryResourceInfo returns information about resources.
type DiscoveryResourceInfo struct {
	resourceLists      []*metav1.APIResourceList
	groupVersionParser GroupVersionParserFunc
}

// NewDiscoveryResourceInfo creates an instance of DiscoveryResourceInfo.
func NewDiscoveryResourceInfo(resourceLists []*metav1.APIResourceList, optionList ...Option) *DiscoveryResourceInfo {
	opts := buildOptions(optionList...)

	dri := &DiscoveryResourceInfo{
		resourceLists:      resourceLists,
		groupVersionParser: opts.groupVersionParser,
	}

	return dri
}

// PreferredVersion returns the preferred version for a group/kind pair.
func (dri *DiscoveryResourceInfo) PreferredVersion(groupKind schema.GroupKind) (string, error) {
	for _, resourceList := range dri.resourceLists {
		groupVersion, err := dri.groupVersionParser(resourceList.GroupVersion)
		if err != nil {
			return "", fmt.Errorf("parse group version %s: %w", resourceList.GroupVersion, err)
		}

		if groupVersion.Group != groupKind.Group {
			continue
		}

		for _, apiResource := range resourceList.APIResources {
			if apiResource.Kind != groupKind.Kind {
				continue
			}

			return groupVersion.Version, nil
		}
	}

	return "", fmt.Errorf("unknown version for %s", groupKind)
}
