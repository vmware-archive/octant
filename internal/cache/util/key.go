package util

import (
	"fmt"
	"strings"

	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Key is a key for the cache.
type Key struct {
	Namespace  string
	APIVersion string
	Kind       string
	Name       string
	Selector   kLabels.Selector
}

func (k Key) String() string {
	var sb strings.Builder

	sb.WriteString("CacheKey[")
	if k.Namespace != "" {
		fmt.Fprintf(&sb, "Namespace=%q, ", k.Namespace)
	}
	fmt.Fprintf(&sb, "APIVersion=%q, ", k.APIVersion)
	fmt.Fprintf(&sb, "Kind=%q", k.Kind)

	if k.Name != "" {
		fmt.Fprintf(&sb, ", Name=%q", k.Name)
	}

	if k.Selector != nil {
		fmt.Fprintf(&sb, ", Selector=%q", k.Selector.String())
	}

	sb.WriteString("]")

	return sb.String()
}

// GroupVersionKind converts the Key to a GroupVersionKind.
func (k Key) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(k.APIVersion, k.Kind)
}
