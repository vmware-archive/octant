package printer

import (
	"reflect"
	"strings"

	rbacv1 "k8s.io/api/rbac/v1"
)

type simpleResource struct {
	Group             string
	Resource          string
	ResourceNameExist bool
	ResourceName      string
}

// compactRules combines rules that contain a single APIGroup/Resource, differ only by verb, and contain no other attributes.
// this is a fast check, and works well with the decomposed "missing rules" list from a Covers check.
func compactRules(rules []rbacv1.PolicyRule) ([]rbacv1.PolicyRule, error) {
	compacted := make([]rbacv1.PolicyRule, 0, len(rules))

	simpleRules := map[simpleResource]*rbacv1.PolicyRule{}
	for _, rule := range rules {
		if resource, isSimple := isSimpleResourceRule(&rule); isSimple {
			if existingRule, ok := simpleRules[resource]; ok {
				// Add the new verbs to the existing simple resource rule
				if existingRule.Verbs == nil {
					existingRule.Verbs = []string{}
				}
				existingRule.Verbs = append(existingRule.Verbs, rule.Verbs...)
			} else {
				// Copy the rule to accumulate matching simple resource rules into
				simpleRules[resource] = rule.DeepCopy()
			}
		} else {
			compacted = append(compacted, rule)
		}
	}

	// Once we've consolidated the simple resource rules, add them to the compacted list
	for _, simpleRule := range simpleRules {
		compacted = append(compacted, *simpleRule)
	}

	return compacted, nil
}

// isSimpleResourceRule returns true if the given rule contains verbs, a single resource, a single API group, at most one Resource Name, and no other values
func isSimpleResourceRule(rule *rbacv1.PolicyRule) (simpleResource, bool) {
	resource := simpleResource{}

	// If we have "complex" rule attributes, return early without allocations or expensive comparisons
	if len(rule.ResourceNames) > 1 || len(rule.NonResourceURLs) > 0 {
		return resource, false
	}
	// If we have multiple api groups or resources, return early
	if len(rule.APIGroups) != 1 || len(rule.Resources) != 1 {
		return resource, false
	}

	// Test if this rule only contains APIGroups/Resources/Verbs/ResourceNames
	simpleRule := &rbacv1.PolicyRule{APIGroups: rule.APIGroups, Resources: rule.Resources, Verbs: rule.Verbs, ResourceNames: rule.ResourceNames}
	if !reflect.DeepEqual(simpleRule, rule) {
		return resource, false
	}

	if len(rule.ResourceNames) == 0 {
		resource = simpleResource{Group: rule.APIGroups[0], Resource: rule.Resources[0], ResourceNameExist: false}
	} else {
		resource = simpleResource{Group: rule.APIGroups[0], Resource: rule.Resources[0], ResourceNameExist: true, ResourceName: rule.ResourceNames[0]}
	}

	return resource, true
}

// BreakdownRule takes a rule and builds an equivalent list of rules that each have at most one verb, one
// resource, and one resource name
func BreakdownRule(rule rbacv1.PolicyRule) []rbacv1.PolicyRule {
	subrules := []rbacv1.PolicyRule{}
	for _, group := range rule.APIGroups {
		for _, resource := range rule.Resources {
			for _, verb := range rule.Verbs {
				if len(rule.ResourceNames) > 0 {
					for _, resourceName := range rule.ResourceNames {
						subrules = append(subrules, rbacv1.PolicyRule{APIGroups: []string{group}, Resources: []string{resource}, Verbs: []string{verb}, ResourceNames: []string{resourceName}})
					}

				} else {
					subrules = append(subrules, rbacv1.PolicyRule{APIGroups: []string{group}, Resources: []string{resource}, Verbs: []string{verb}})
				}

			}
		}
	}

	// Non-resource URLs are unique because they only combine with verbs.
	for _, nonResourceURL := range rule.NonResourceURLs {
		for _, verb := range rule.Verbs {
			subrules = append(subrules, rbacv1.PolicyRule{NonResourceURLs: []string{nonResourceURL}, Verbs: []string{verb}})
		}
	}

	return subrules
}

func CombineResourceGroup(resource, group []string) string {
	if len(resource) == 0 {
		return ""
	}
	parts := strings.SplitN(resource[0], "/", 2)
	combine := parts[0]

	if len(group) > 0 && group[0] != "" {
		combine = combine + "." + group[0]
	}

	if len(parts) == 2 {
		combine = combine + "/" + parts[1]
	}
	return combine
}
