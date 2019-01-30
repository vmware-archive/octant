package printer

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/heptio/developer-dash/internal/view/component"
)

type AffinityDescriber struct {
	podSpec corev1.PodSpec
}

func NewAffinityDescriber(podSpec corev1.PodSpec) *AffinityDescriber {
	return &AffinityDescriber{
		podSpec: podSpec,
	}
}

func (ad *AffinityDescriber) Create() (*component.List, error) {
	var items []component.ViewComponent

	if affinity := ad.podSpec.Affinity; affinity == nil {
		items = append(items, component.NewText("", "Pod affinity is not configured."))
	} else {
		items = append(items, ad.nodeAffinity(*affinity)...)
		items = append(items, ad.podAffinity(*affinity)...)
	}

	list := component.NewList("", items)

	return list, nil
}

type podAffinityOptions struct {
	weight     int32
	isRequired bool
	anti       bool
}

func (ad *AffinityDescriber) podAffinity(affinity corev1.Affinity) []component.ViewComponent {
	var items []component.ViewComponent

	if podAffinity := affinity.PodAffinity; podAffinity != nil {
		requiredOptions := podAffinityOptions{isRequired: true}
		items = append(items,
			ad.podAffinityTerms(
				podAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
				requiredOptions)...)
		for _, weighted := range podAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
			preferredOptions := podAffinityOptions{weight: weighted.Weight}
			items = append(items,
				ad.podAffinityTerms(
					[]corev1.PodAffinityTerm{weighted.PodAffinityTerm},
					preferredOptions)...)
		}
	}

	if podAntiAffinity := affinity.PodAntiAffinity; podAntiAffinity != nil {
		requiredOptions := podAffinityOptions{isRequired: true, anti: true}
		items = append(items,
			ad.podAffinityTerms(
				podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
				requiredOptions)...)
		for _, weighted := range podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
			preferredOptions := podAffinityOptions{weight: weighted.Weight, anti: true}
			items = append(items,
				ad.podAffinityTerms(
					[]corev1.PodAffinityTerm{weighted.PodAffinityTerm},
					preferredOptions)...)
		}
	}

	return items
}

func (ad *AffinityDescriber) podAffinityTerms(terms []corev1.PodAffinityTerm, options podAffinityOptions) []component.ViewComponent {
	var items []component.ViewComponent

	for _, term := range terms {
		var b strings.Builder
		switch {
		case options.isRequired && !options.anti:
			b.WriteString("Schedule with pod")
		case options.isRequired && options.anti:
			b.WriteString("Do not schedule with pod")
		case !options.isRequired && !options.anti:
			b.WriteString("Prefer to schedule with pod")
		case !options.isRequired && options.anti:
			b.WriteString("Prefer to not schedule with pod")
		}

		if term.LabelSelector != nil {
			matchLabels := printMatchLabels(term.LabelSelector.MatchLabels)

			if matchLabels != "" {
				b.WriteString(fmt.Sprintf(" labeled %s", matchLabels))
			}

			matchExpressions := printLabelSelectorRequirement(term.LabelSelector.MatchExpressions)
			if matchExpressions != "" {
				b.WriteString(fmt.Sprintf(" where %s", matchExpressions))
			}
		}

		b.WriteString(fmt.Sprintf(" in topology %s.", term.TopologyKey))

		if options.weight > 0 {
			b.WriteString(fmt.Sprintf(" Weight %d.", options.weight))
		}

		items = append(items, component.NewText("", b.String()))
	}

	return items
}

func (ad *AffinityDescriber) nodeAffinity(affinity corev1.Affinity) []component.ViewComponent {
	var items []component.ViewComponent

	if nodeAffinity := affinity.NodeAffinity; nodeAffinity != nil {
		for _, preferred := range nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {
			options := nodeSelectorRequirementOptions{
				weight: preferred.Weight,
			}
			items = append(items, ad.nodeSelectorTerms([]corev1.NodeSelectorTerm{preferred.Preference}, options)...)
		}

		if required := nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution; required != nil {
			options := nodeSelectorRequirementOptions{
				isRequired: true,
			}
			items = append(items, ad.nodeSelectorTerms(required.NodeSelectorTerms, options)...)
		}
	}

	return items
}

type nodeSelectorRequirementOptions struct {
	isRequired bool
	weight     int32
}

func (ad *AffinityDescriber) nodeSelectorTerms(terms []corev1.NodeSelectorTerm, options nodeSelectorRequirementOptions) []component.ViewComponent {
	var items []component.ViewComponent

	for _, term := range terms {
		items = append(items, ad.nodeSelectorRequirement("label", options, term.MatchExpressions)...)
		items = append(items, ad.nodeSelectorRequirement("field", options, term.MatchFields)...)
	}

	return items
}

func (ad *AffinityDescriber) nodeSelectorRequirement(itemType string, options nodeSelectorRequirementOptions, nodeSelectorRequirements []corev1.NodeSelectorRequirement) []component.ViewComponent {
	preamble := "Prefer to schedule on nodes"
	if options.isRequired {
		preamble = "Schedule on nodes"
	}

	var items []component.ViewComponent
	for _, nsr := range nodeSelectorRequirements {
		var b strings.Builder

		switch nsr.Operator {
		case corev1.NodeSelectorOpIn:
			b.WriteString(fmt.Sprintf("%s with %s %s with values %s.",
				preamble, itemType, nsr.Key, strings.Join(nsr.Values, ", ")))
		case corev1.NodeSelectorOpNotIn:
			b.WriteString(fmt.Sprintf("%s with %s %s without values %s.",
				preamble, itemType, nsr.Key, strings.Join(nsr.Values, ", ")))
		case corev1.NodeSelectorOpExists:
			b.WriteString(fmt.Sprintf("%s where %s %s exists.",
				preamble, itemType, nsr.Key))
		case corev1.NodeSelectorOpDoesNotExist:
			b.WriteString(fmt.Sprintf("%s where %s %s does not exist.",
				preamble, itemType, nsr.Key))
		case corev1.NodeSelectorOpGt:
			b.WriteString(fmt.Sprintf("%s where %s %s is greater than %s.",
				preamble, itemType, nsr.Key, nsr.Values[0]))
		case corev1.NodeSelectorOpLt:
			b.WriteString(fmt.Sprintf("%s where %s %s is less than %s.",
				preamble, itemType, nsr.Key, nsr.Values[0]))
		default:
			continue
		}

		if options.weight > 0 {
			b.WriteString(fmt.Sprintf(" Weight %d.", options.weight))
		}

		items = append(items, component.NewText("", b.String()))
	}

	return items
}
