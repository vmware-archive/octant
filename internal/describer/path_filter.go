/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"fmt"
	"regexp"
)

type PathFilter struct {
	path      string
	Describer Describer

	re *regexp.Regexp
}

// NewPathFilter creates a path filter.
func NewPathFilter(path string, describer Describer) *PathFilter {
	re := regexp.MustCompile(fmt.Sprintf("^%s/?$", path))

	return &PathFilter{
		re:        re,
		path:      path,
		Describer: describer,
	}
}

func (pf *PathFilter) String() string {
	return pf.path
}

func (pf *PathFilter) Match(path string) bool {
	return pf.re.MatchString(path)
}

// Fields extracts parameters from the request path.
// In practice, this finds the field "name" for an object request.
func (pf *PathFilter) Fields(path string) map[string]string {
	out := make(map[string]string)

	match := pf.re.FindStringSubmatch(path)
	for i, name := range pf.re.SubexpNames() {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}

	return out
}
