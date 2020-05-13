/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kLabels "k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

// NetworkPolicyListHandler is a printFunc that prints network policies
func NetworkPolicyListHandler(ctx context.Context, list *networkingv1.NetworkPolicyList, options Options) (component.Component, error) {
	if list == nil {
		return nil, errors.New("network policy list is nil")
	}

	cols := component.NewTableCols("Name", "Labels", "Age")
	ot := NewObjectTable("Network Policies", "We couldn't find any network policies!", cols, options.DashConfig.ObjectStore())

	for _, networkPolicy := range list.Items {
		row := component.TableRow{}
		nameLink, err := options.Link.ForObject(&networkPolicy, networkPolicy.Name)
		if err != nil {
			return nil, err
		}

		row["Name"] = nameLink
		row["Labels"] = component.NewLabels(networkPolicy.Labels)
		row["Age"] = component.NewTimestamp(networkPolicy.CreationTimestamp.Time)

		if err := ot.AddRowForObject(ctx, &networkPolicy, row); err != nil {
			return nil, fmt.Errorf("add row for object: %w", err)
		}
	}

	return ot.ToComponent()
}

// NetworkPolicyHandler is a printFunc that prints NetworkPolicies
func NetworkPolicyHandler(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, options Options) (component.Component, error) {
	o := NewObject(networkPolicy)
	o.EnableEvents()

	np, err := newNetworkPolicyHander(networkPolicy, o)
	if err != nil {
		return nil, err
	}

	if err := np.Config(); err != nil {
		return nil, errors.Wrap(err, "print networkPolicy configuration")
	}
	if err := np.Status(); err != nil {
		return nil, errors.Wrap(err, "print networkPolicy status")
	}

	if err := np.Pods(ctx, networkPolicy, options); err != nil {
		return nil, errors.Wrap(err, "print networkPolicy pods")
	}

	return o.ToComponent(ctx, options)
}

type networkPolicyObject interface {
	Config() error
	Status() error
	Pods(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, options Options) error
}

type networkPolicyHandler struct {
	networkPolicy *networkingv1.NetworkPolicy
	configFunc    func(*networkingv1.NetworkPolicy) (*component.Summary, error)
	summaryFunc   func(*networkingv1.NetworkPolicy) (*component.Summary, error)
	podFunc       func(context.Context, *networkingv1.NetworkPolicy, Options) (component.Component, error)
	object        *Object
}

var _ networkPolicyObject = (*networkPolicyHandler)(nil)

func newNetworkPolicyHander(networkPolicy *networkingv1.NetworkPolicy, object *Object) (*networkPolicyHandler, error) {
	if networkPolicy == nil {
		return nil, errors.New("can't print a nil network policy")
	}

	if object == nil {
		return nil, errors.New("can't print network policy using a nil object printer")
	}

	nph := &networkPolicyHandler{
		networkPolicy: networkPolicy,
		configFunc:    defaultNetworkPolicyConfig,
		summaryFunc:   defaultNetWorkPolicySummary,
		podFunc:       defaultNetworkPolicyPods,
		object:        object,
	}

	return nph, nil
}

func (n *networkPolicyHandler) Config() error {
	out, err := n.configFunc(n.networkPolicy)
	if err != nil {
		return err
	}

	n.object.RegisterConfig(out)
	return nil
}

func defaultNetworkPolicyConfig(networkPolicy *networkingv1.NetworkPolicy) (*component.Summary, error) {
	return NewNetworkPolicyConfiguration(networkPolicy).Create()
}

func (n *networkPolicyHandler) Status() error {
	out, err := n.summaryFunc(n.networkPolicy)
	if err != nil {
		return err
	}

	n.object.RegisterSummary(out)
	return nil
}

func defaultNetWorkPolicySummary(networkPolicy *networkingv1.NetworkPolicy) (*component.Summary, error) {
	return createNetworkPolicySummaryStatus(networkPolicy)
}

func (n *networkPolicyHandler) Pods(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, options Options) error {
	n.object.RegisterItems(ItemDescriptor{
		Width: component.WidthFull,
		Func: func() (component.Component, error) {
			return n.podFunc(ctx, networkPolicy, options)
		},
	})
	return nil
}

func defaultNetworkPolicyPods(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, options Options) (component.Component, error) {
	return createNetworkPodListView(ctx, networkPolicy, options)
}

// NetworkPolicyConfiguration generates networkPolicy configuration
type NetworkPolicyConfiguration struct {
	networkPolicy *networkingv1.NetworkPolicy
}

// NewNetworkPolicyConfiguration creates an instance of NetworkPolicyConfiguration
func NewNetworkPolicyConfiguration(n *networkingv1.NetworkPolicy) *NetworkPolicyConfiguration {
	return &NetworkPolicyConfiguration{
		networkPolicy: n,
	}
}

// Create creates a networkPolicy configuration summary
func (n *NetworkPolicyConfiguration) Create() (*component.Summary, error) {
	if n.networkPolicy == nil {
		return nil, errors.New("network policy is nil")
	}

	sections := make([]component.SummarySection, 0)

	networkPolicy := n.networkPolicy

	if networkPolicy.Spec.PolicyTypes != nil {
		var policyText []string
		for _, policy := range networkPolicy.Spec.PolicyTypes {
			policyText = append(policyText, string(policy))
		}

		sections = append(sections, component.SummarySection{
			Header:  "Policy Types",
			Content: component.NewText(strings.Join(policyText, ", ")),
		})
	}

	if &networkPolicy.Spec.PodSelector != nil {
		selectors, err := selectorToComponent(&networkPolicy.Spec.PodSelector)
		if err != nil {
			return nil, err
		}

		sections = append(sections, component.SummarySection{
			Header:  "Selectors",
			Content: selectors,
		})
	}

	summary := component.NewSummary("Configuration", sections...)

	return summary, nil
}

func selectorToComponent(selector *metav1.LabelSelector) (*component.Selectors, error) {
	if selector == nil {
		return nil, errors.New("nil label selector")
	}

	var selectors []component.Selector

	for _, lsr := range selector.MatchExpressions {
		o, err := component.MatchOperator(string(lsr.Operator))
		if err != nil {
			return nil, err
		}

		es := component.NewExpressionSelector(lsr.Key, o, lsr.Values)
		selectors = append(selectors, es)
	}

	for k, v := range selector.MatchLabels {
		ls := component.NewLabelSelector(k, v)
		selectors = append(selectors, ls)
	}

	return component.NewSelectors(selectors), nil
}

func createNetworkPolicySummaryStatus(networkPolicy *networkingv1.NetworkPolicy) (*component.Summary, error) {
	if networkPolicy == nil {
		return nil, errors.New("unable to generate status from a nil network policy")
	}

	ingressTable, err := createIngressRules(networkPolicy.Spec.Ingress)
	if err != nil {
		return nil, errors.Wrapf(err, "create ingress rules")
	}

	egressTable, err := createEgressRules(networkPolicy.Spec.Egress)
	if err != nil {
		return nil, errors.Wrapf(err, "create egress rules")
	}

	summary := component.NewSummary("Status", []component.SummarySection{
		{
			Header:  "Summary",
			Content: policyDescriber(networkPolicy),
		},
		{
			Header:  "Allowing ingress traffic",
			Content: ingressTable,
		},
		{
			Header:  "Allowing egress traffic",
			Content: egressTable,
		},
	}...)

	return summary, nil
}

func createIngressRules(ingressRules []networkingv1.NetworkPolicyIngressRule) (*component.Table, error) {
	cols := component.NewTableCols("Policy", "Value")
	ingressRuleTable := component.NewTable("", "No rules found", cols)

	if ingressRules == nil {
		return ingressRuleTable, nil
	}

	var portText []string

	for _, rule := range ingressRules {
		if rule.Ports != nil {
			for _, port := range rule.Ports {
				protocol := string(*port.Protocol)
				portText = append(portText, port.Port.String()+"/"+protocol)
			}

			row := component.TableRow{}
			row["Policy"] = component.NewText("To Port")
			row["Value"] = component.NewText(strings.Join(portText, ", "))
			ingressRuleTable.Add(row)
		}
		if rule.From != nil {
			for _, peer := range rule.From {
				if peer.IPBlock != nil {
					row := component.TableRow{}
					row["Policy"] = component.NewText("From IP Block")
					row["Value"] = component.NewText(describeIPBlock(peer.IPBlock))

					ingressRuleTable.Add(row)
				}

				if peer.NamespaceSelector != nil {
					selectors, err := selectorToComponent(peer.NamespaceSelector)
					if err != nil {
						return nil, err
					}

					row := component.TableRow{}
					row["Policy"] = component.NewText("From Namespace Selector")
					row["Value"] = selectors
					ingressRuleTable.Add(row)
				}

				if peer.PodSelector != nil {
					selectors, err := selectorToComponent(peer.PodSelector)
					if err != nil {
						return nil, err
					}

					row := component.TableRow{}
					row["Policy"] = component.NewText("From Pod Selector")
					row["Value"] = selectors
					ingressRuleTable.Add(row)
				}
			}
		}
	}
	return ingressRuleTable, nil
}

func createEgressRules(egressRules []networkingv1.NetworkPolicyEgressRule) (*component.Table, error) {
	cols := component.NewTableCols("Policy", "Value")
	egressRuleTable := component.NewTable("", "No rules found", cols)

	if egressRules == nil {
		return egressRuleTable, nil
	}

	var portText []string

	for _, rule := range egressRules {
		if rule.Ports != nil {
			for _, port := range rule.Ports {
				protocol := string(*port.Protocol)
				portText = append(portText, port.Port.String()+"/"+protocol)
			}

			row := component.TableRow{}
			row["Policy"] = component.NewText("To Port")
			row["Value"] = component.NewText(strings.Join(portText, ", "))
			egressRuleTable.Add(row)
		}
		if rule.To != nil {
			for _, peer := range rule.To {
				if peer.IPBlock != nil {
					row := component.TableRow{}
					row["Policy"] = component.NewText("From IP Block")
					row["Value"] = component.NewText(describeIPBlock(peer.IPBlock))
					egressRuleTable.Add(row)
				}

				if peer.NamespaceSelector != nil {
					selectors, err := selectorToComponent(peer.NamespaceSelector)
					if err != nil {
						return nil, err
					}

					row := component.TableRow{}
					row["Policy"] = component.NewText("From Namespace Selector")
					row["Value"] = selectors
					egressRuleTable.Add(row)
				}

				if peer.PodSelector != nil {
					selectors, err := selectorToComponent(peer.PodSelector)
					if err != nil {
						return nil, err
					}

					row := component.TableRow{}
					row["Policy"] = component.NewText("From Pod Selector")
					row["Value"] = selectors
					egressRuleTable.Add(row)
				}
			}
		}
	}
	return egressRuleTable, nil
}

func createNetworkPodListView(ctx context.Context, networkPolicy *networkingv1.NetworkPolicy, options Options) (component.Component, error) {
	options.DisableLabels = true
	podList := &corev1.PodList{}

	objectStore := options.DashConfig.ObjectStore()

	podSelectorList := []*metav1.LabelSelector{}
	selectorsList := []*metav1.LabelSelector{}
	keyList := []store.Key{}

	if networkPolicy.Spec.Ingress != nil {
		for _, rule := range networkPolicy.Spec.Ingress {
			if rule.From != nil {
				for _, peer := range rule.From {
					if peer.NamespaceSelector != nil {
						selectorsList = append(selectorsList, peer.NamespaceSelector)
					}

					if peer.PodSelector != nil {
						podSelectorList = append(podSelectorList, peer.PodSelector)
					}
				}
			}
		}
	}

	// Case with only pod selectors
	if len(selectorsList) == 0 {
		keyList = append(keyList, store.Key{Namespace: networkPolicy.Namespace, APIVersion: "v1", Kind: "Pod"})
	}
	// Case with namespace and pod selectors
	for _, selector := range selectorsList {
		s, err := metav1.LabelSelectorAsMap(selector)
		if err != nil {
			return nil, err
		}

		labelMap := kLabels.Set(s)

		ul, _, err := objectStore.List(ctx, store.Key{APIVersion: "v1", Kind: "Namespace", Selector: &labelMap})
		if err != nil {
			return nil, err
		}

		for i := range ul.Items {
			namespace := &corev1.Namespace{}
			err := runtime.DefaultUnstructuredConverter.FromUnstructured(ul.Items[i].Object, namespace)
			if err != nil {
				return nil, err
			}

			keyList = append(keyList, store.Key{Namespace: namespace.Name, APIVersion: "v1", Kind: "Pod"})
		}
	}

	// Get pods from namespaces filtering by pod selector
	for _, key := range keyList {
		for _, podSelector := range podSelectorList {
			pods, err := loadPods(ctx, key, objectStore, podSelector)
			if err != nil {
				return nil, err
			}

			for _, pod := range pods {
				podList.Items = append(podList.Items, *pod)
			}
		}
	}

	return PodListHandler(ctx, podList, options)
}

func policyDescriber(networkPolicy *networkingv1.NetworkPolicy) *component.Text {
	policyMap := map[*networkingv1.NetworkPolicySpec]string{
		&networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Ingress: []networkingv1.NetworkPolicyIngressRule{
				{},
			},
		}: "Deny all traffic from all other namespaces",
	}

	if &networkPolicy.Spec.PodSelector != nil {
		podSelectorPolicyMap := map[*networkingv1.NetworkPolicySpec]string{
			&networkingv1.NetworkPolicySpec{
				PodSelector: networkPolicy.Spec.PodSelector,
				Ingress:     []networkingv1.NetworkPolicyIngressRule{},
			}: "Deny all traffic to application",
			&networkingv1.NetworkPolicySpec{
				PodSelector: networkPolicy.Spec.PodSelector,
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{},
				},
			}: "Allow all traffic to application",
			&networkingv1.NetworkPolicySpec{
				PodSelector: networkPolicy.Spec.PodSelector,
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{
						From: []networkingv1.NetworkPolicyPeer{
							{
								NamespaceSelector: &metav1.LabelSelector{},
							},
						},
					},
				},
			}: "Allow traffic to application from all namespaces",
			&networkingv1.NetworkPolicySpec{
				PodSelector: networkPolicy.Spec.PodSelector,
				Ingress: []networkingv1.NetworkPolicyIngressRule{
					{
						From: []networkingv1.NetworkPolicyPeer{},
					},
				},
				PolicyTypes: []networkingv1.PolicyType{
					networkingv1.PolicyTypeIngress,
				},
			}: "Allow traffic from external clients",
		}

		for k, v := range podSelectorPolicyMap {
			policyMap[k] = v
		}

		if networkPolicy.Spec.Egress != nil {
			for _, egress := range networkPolicy.Spec.Egress {
				for _, peer := range egress.To {
					if peer.NamespaceSelector != nil {
						denyExternalSpec := &networkingv1.NetworkPolicySpec{
							PodSelector: networkPolicy.Spec.PodSelector,
							Egress:      networkPolicy.Spec.Egress,
							PolicyTypes: []networkingv1.PolicyType{
								networkingv1.PolicyTypeEgress,
							},
						}

						policyMap[denyExternalSpec] = "Deny external egress traffic"
					}
				}
			}
		}
	}

	for policySpec, policyDescription := range policyMap {
		if apiequality.Semantic.DeepEqual(&networkPolicy.Spec, policySpec) {
			return component.NewText(policyDescription)
		}
	}

	return component.NewText("---")
}

func describeIPBlock(source interface{}) string {
	data, _ := json.Marshal(source)
	return string(data)
}
