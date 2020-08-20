package main

import (
	"fmt"
	"sort"
	"strings"
)

// LDFlags represents ldflags
type LDFlags map[string][]string

// String returns the ldflags as a string
func (l LDFlags) String() string {
	var o []string
	keys := make([]string, 0, len(l))
	for k := range l {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		ss := l[k]
		if len(ss) == 0 {
			o = append(o, "-"+k)
			continue
		}
		for _, s := range ss {
			o = append(o, fmt.Sprintf(`-%s %s`, k, s))
		}
	}
	return strings.Join(o, " ")
}

// Set allows setting the values for use by flag
func (l LDFlags) Set(s string) error {
	segments := strings.SplitN(s, ":", 2)
	flag := segments[0]

	val := l[flag]
	if len(segments) == 2 {
		val = strings.Split(segments[1], ",")
	}
	l[flag] = append(l[flag], val...)

	return nil
}
