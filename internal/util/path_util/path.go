/*
Copyright (c) 2021 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package path_util

import (
	"path"
	"strings"
)

// NamespacedPath generates the URL for namespaced path
// by joining base url, namespace and additional path segments.
//
func NamespacedPath(base, namespace string, paths ...string) string {
	return path.Join(append([]string{base, "namespace", namespace}, paths...)...)
}

// PrefixedPath ensures that provided url starts with slash ("/")
//
func PrefixedPath(url string) string {
	if !strings.HasPrefix(url, "/") {
		url = path.Join("/", url)
	}
	return url
}
