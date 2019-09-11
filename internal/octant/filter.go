/*
 * Copyright (c) 2019 VMware, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package octant

import "fmt"

// Filter is used to filter queries for objects. Typically,
// the filter is an object's label.
type Filter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ToQueryParam converts the filter to a query parameter.
func (f *Filter) ToQueryParam() string {
	return fmt.Sprintf("%s:%s", f.Key, f.Value)
}

// IsEqual returns true if the filter equals the other filter.
func (f *Filter) IsEqual(other Filter) bool {
	return f.Key == other.Key && f.Value == other.Value
}
