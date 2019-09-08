/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"fmt"
	"regexp"

	"github.com/davecgh/go-spew/spew"
)

type PathFilter struct {
	filterPath string
	Describer  Describer

	re *regexp.Regexp
}

// NewPathFilter creates a path filter.
func NewPathFilter(filterPath string, describer Describer) *PathFilter {

	// handle special case where you get `/namespace/default` by making
	// it optional in the regex.
	if filterPath == "/" {
		filterPath += "?"
	}
	re := regexp.MustCompile(fmt.Sprintf(`^(/namespace/(?P<namespace>[^/]+))?%s$`, filterPath))
	return &PathFilter{
		re:         re,
		filterPath: filterPath,
		Describer:  describer,
	}
}

func (pf *PathFilter) String() string {
	return pf.filterPath
}

// Match matches a contentPath against the filter.
//
// content paths look like:
//   /foo/bar
//   /namespace/default
//   /namespace/default/foo/bar
//   /
func (pf *PathFilter) Match(contentPath string) bool {
	return pf.re.MatchString(contentPath)
}

// Fields extracts parameters from the request path.
// In practice, this finds the field "name" for an object request.
func (pf *PathFilter) Fields(contentPath string) map[string]string {
	out := make(map[string]string)

	match := pf.re.FindStringSubmatch(contentPath)
	names := pf.re.SubexpNames()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("path filter fields is crashing the app")
			spew.Dump(contentPath, pf.filterPath, pf.re.String(), match, names)
			panic("i'm done")
		}
	}()

	for i, name := range names {
		if i != 0 && name != "" {
			out[name] = match[i]
		}
	}

	return out
}
