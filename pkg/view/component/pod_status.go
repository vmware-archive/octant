/*
Copyright (c) 2019 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package component

import "github.com/vmware-tanzu/octant/internal/util/json"

// PodSummary is a status summary for a pod.
type PodSummary struct {
	Details    []Component `json:"details,omitempty"`
	Properties []Property  `json:"properties,omitempty"`
	Status     NodeStatus  `json:"status,omitempty"`
}

// PodStatusConfig is config for PodStatus.
type PodStatusConfig struct {
	Pods map[string]PodSummary `json:"pods,omitempty"`
}

// PodStatus represents the status for a group of pods.
//
// +octant:component
type PodStatus struct {
	Base
	Config PodStatusConfig `json:"config"`
}

var _ Component = (*PodStatus)(nil)

// NewPodStatus creates a PodStatus.
func NewPodStatus() *PodStatus {
	return &PodStatus{
		Base: newBase(TypePodStatus, nil),
		Config: PodStatusConfig{
			Pods: make(map[string]PodSummary),
		},
	}
}

type podStatusMarshal PodStatus

// MarshalJSON implements json.Marshaler.
func (ps *PodStatus) MarshalJSON() ([]byte, error) {
	m := podStatusMarshal(*ps)
	m.Metadata.Type = TypePodStatus
	return json.Marshal(&m)
}

// AddSummary adds summary for a pod.
func (ps *PodStatus) AddSummary(name string, details []Component, properties []Property, status NodeStatus) {
	ps.Config.Pods[name] = PodSummary{
		Details:    details,
		Properties: properties,
		Status:     status,
	}
}

func (ps *PodStatus) Status() NodeStatus {

	tally := make(map[NodeStatus]int)

	for _, summary := range ps.Config.Pods {
		tally[summary.Status]++
	}

	if errors, ok := tally[NodeStatusError]; ok && errors > 0 {
		return NodeStatusError
	} else if warnings, ok := tally[NodeStatusWarning]; ok && warnings > 0 {
		return NodeStatusWarning
	} else {
		return NodeStatusOK
	}
}

func (podSummary *PodSummary) UnmarshalJSON(data []byte) error {
	stage := struct {
		Details    []TypedObject    `json:"details,omitempty"`
		Properties []PropertyObject `json:"properties,omitempty"`
		Status     NodeStatus       `json:"status,omitempty"`
	}{}

	if err := json.Unmarshal(data, &stage); err != nil {
		return err
	}

	podSummary.Status = stage.Status

	for _, to := range stage.Details {
		status, err := to.ToComponent()
		if err != nil {
			return err
		}

		podSummary.Details = append(podSummary.Details, status)
	}

	for _, to := range stage.Properties {
		val, err := to.Value.ToComponent()
		if err != nil {
			return err
		}

		podSummary.Properties = append(podSummary.Properties, Property{Label: to.Label, Value: val})
	}

	return nil
}
