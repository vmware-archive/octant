/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"fmt"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func printMatchLabels(matchLabels map[string]string) string {
	var keys []string
	for k := range matchLabels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var matchLabelStrings []string
	for _, k := range keys {
		v := matchLabels[k]
		matchLabelStrings = append(matchLabelStrings, fmt.Sprintf("%s:%s", k, v))
	}

	return strings.Join(matchLabelStrings, ", ")
}

func printLabelSelectorRequirement(requirements []metav1.LabelSelectorRequirement) string {
	var sections []string

	for _, requirement := range requirements {
		switch requirement.Operator {
		case metav1.LabelSelectorOpIn:
			sections = append(sections, fmt.Sprintf("%s in [%s]",
				requirement.Key, strings.Join(requirement.Values, ", ")))
		case metav1.LabelSelectorOpNotIn:
			sections = append(sections, fmt.Sprintf("%s not in [%s]",
				requirement.Key, strings.Join(requirement.Values, ", ")))
		case metav1.LabelSelectorOpExists:
			sections = append(sections, fmt.Sprintf("%s exists",
				requirement.Key))
		case metav1.LabelSelectorOpDoesNotExist:
			sections = append(sections, fmt.Sprintf("%s does not exist",
				requirement.Key))
		}
	}

	return strings.Join(sections, ", ")
}
