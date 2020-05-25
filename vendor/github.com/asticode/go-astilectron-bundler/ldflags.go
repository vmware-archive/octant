package astibundler

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
			o = append(o, fmt.Sprintf(`-%s "%s"`, k, s))
		}
	}
	return strings.Join(o, " ")
}

// Merge merges ldflags
func (l LDFlags) Merge(r LDFlags) {
	for flag := range r {
		l[flag] = append(l[flag], r[flag]...)
	}
}
