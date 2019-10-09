/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package strings

// Contains returns true if a slice contains a string and false if not.
func Contains(s string, sl []string) bool {
	for i := range sl {
		if s == sl[i] {
			return true
		}
	}

	return false
}

// Deduplicate removes any duplicated values and returns a new slice, keeping the order unchanged
func Deduplicate(s []string) []string {
	encountered := map[string]bool{}
	ret := make([]string, 0)
	for i := range s {
		if encountered[s[i]] {
			continue
		}
		encountered[s[i]] = true
		ret = append(ret, s[i])
	}
	return ret
}
