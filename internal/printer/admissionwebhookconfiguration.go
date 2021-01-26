/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"fmt"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/vmware-tanzu/octant/internal/gvk"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

// printYaml returns the object as yaml. Errors are returned as `<error>`.
func printYaml(obj interface{}) string {
	b, err := yaml.Marshal(obj)
	if err != nil {
		return "<error>"
	}
	return string(b)
}

func admissionWebhookRules(rules []admissionregistrationv1.RuleWithOperations, options Options) (component.Component, error) {
	columns := component.NewTableCols("API Groups", "API Versions", "Resources", "Operations", "Scope")
	table := component.NewTable("Rules", "There are no webhook rules!", columns)

	for _, rule := range rules {
		row := component.TableRow{}

		if len(rule.APIGroups) == 1 {
			row["API Groups"] = component.NewText(rule.APIGroups[0])
		} else {
			row["API Groups"] = component.NewMarkdownText(printYaml(rule.APIGroups))
		}
		if len(rule.APIVersions) == 1 {
			row["API Versions"] = component.NewText(rule.APIVersions[0])
		} else {
			row["API Versions"] = component.NewMarkdownText(printYaml(rule.APIVersions))
		}
		if len(rule.Resources) == 1 {
			row["Resources"] = component.NewText(rule.Resources[0])
		} else {
			row["Resources"] = component.NewMarkdownText(printYaml(rule.Resources))
		}
		if len(rule.Operations) == 1 {
			row["Operations"] = component.NewText(string(rule.Operations[0]))
		} else {
			row["Operations"] = component.NewMarkdownText(printYaml(rule.Operations))
		}
		if rule.Scope == nil {
			row["Scope"] = component.NewText(string(admissionregistrationv1.AllScopes))
		} else {
			row["Scope"] = component.NewText(string(*rule.Scope))
		}

		table.Add(row)
	}

	return table, nil
}

func admissionWebhookTimeout(timeoutSeconds *int32) component.Component {
	if timeoutSeconds == nil {
		return component.NewTextf("%ds", 10)
	}
	return component.NewTextf("%ds", *timeoutSeconds)
}

func admissionWebhookClientConfig(clientConfig admissionregistrationv1.WebhookClientConfig, options Options) (component.Component, error) {
	if clientConfig.Service != nil {
		return options.Link.ForGVK(
			clientConfig.Service.Namespace,
			gvk.Service.GroupVersion().String(),
			gvk.Service.Kind, clientConfig.Service.Name,
			fmt.Sprintf("%s/%s", clientConfig.Service.Namespace, clientConfig.Service.Name),
		)
	}
	if clientConfig.URL != nil {
		return component.NewText(*clientConfig.URL), nil
	}
	return component.NewText("unknown"), nil
}

func admissionWebhookFailurePolicy(failurePolicy *admissionregistrationv1.FailurePolicyType) component.Component {
	if failurePolicy == nil {
		return component.NewTextf("%s", admissionregistrationv1.Fail)
	}
	return component.NewTextf("%s", *failurePolicy)
}

func admissionWebhookMatchPolicy(matchPolicy *admissionregistrationv1.MatchPolicyType) component.Component {
	if matchPolicy == nil {
		return component.NewTextf("%s", admissionregistrationv1.Equivalent)
	}
	return component.NewTextf("%s", *matchPolicy)
}

func admissionWebhookSideEffects(sideEffects *admissionregistrationv1.SideEffectClass) component.Component {
	if sideEffects == nil {
		return component.NewTextf("%s", admissionregistrationv1.SideEffectClassUnknown)
	}
	return component.NewTextf("%s", *sideEffects)
}

func admissionWebhookLabelSelector(selector *metav1.LabelSelector) component.Component {
	if selector == nil {
		selector = &metav1.LabelSelector{}
	}
	matchLabels := printMatchLabels(selector.MatchLabels)
	matchExpressions := printLabelSelectorRequirement(selector.MatchExpressions)
	joiner := ""
	if matchLabels != "" && matchExpressions != "" {
		joiner = ", "
	} else if matchLabels == "" && matchExpressions == "" {
		joiner = "*"
	}
	return component.NewTextf("%s%s%s", matchLabels, joiner, matchExpressions)
}

func admissionWebhookAdmissionReviewVersions(admissionReviewVersions []string) component.Component {
	if len(admissionReviewVersions) == 1 {
		return component.NewText(admissionReviewVersions[0])
	}
	return component.NewMarkdownText(printYaml(admissionReviewVersions))
}

func admissionWebhookReinvocationPolicy(reinvocationPolicy *admissionregistrationv1.ReinvocationPolicyType) component.Component {
	if reinvocationPolicy == nil {
		return component.NewTextf("%s", admissionregistrationv1.NeverReinvocationPolicy)
	}
	return component.NewTextf("%s", *reinvocationPolicy)
}
