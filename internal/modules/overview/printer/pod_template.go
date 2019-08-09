/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package printer

import (
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/pkg/view/component"
	"github.com/vmware/octant/pkg/view/flexlayout"
)

type podTemplateLayoutOptions struct {
	parent          runtime.Object
	containers      []corev1.Container
	podTemplateSpec corev1.PodTemplateSpec
	isInit          bool
	printOptions    Options
}

type podTemplateFunc func(fl *flexlayout.FlexLayout, options podTemplateLayoutOptions) error

type PodTemplate struct {
	parent          runtime.Object
	podTemplateSpec corev1.PodTemplateSpec

	podTemplateHeaderFunc           podTemplateFunc
	podTemplateInitContainersFunc   podTemplateFunc
	podTemplateContainersFunc       podTemplateFunc
	podTemplatePodConfigurationFunc podTemplateFunc
}

func NewPodTemplate(parent runtime.Object, podTemplateSpec corev1.PodTemplateSpec) *PodTemplate {
	return &PodTemplate{
		parent:                          parent,
		podTemplateSpec:                 podTemplateSpec,
		podTemplateHeaderFunc:           podTemplateHeader,
		podTemplateInitContainersFunc:   podTemplateContainers,
		podTemplateContainersFunc:       podTemplateContainers,
		podTemplatePodConfigurationFunc: podTemplatePodConfiguration,
	}
}

func (pt *PodTemplate) AddToFlexLayout(fl *flexlayout.FlexLayout, options Options) error {
	if fl == nil {
		return errors.New("flex layout is nil")
	}

	baseOptions := podTemplateLayoutOptions{
		parent:          pt.parent,
		podTemplateSpec: pt.podTemplateSpec,
		printOptions:    options,
	}

	if err := pt.podTemplateHeaderFunc(fl, baseOptions); err != nil {
		return errors.Wrap(err, "pod template header")
	}

	initContainerOptions := baseOptions
	initContainerOptions.containers = pt.podTemplateSpec.Spec.InitContainers
	initContainerOptions.isInit = true

	if err := pt.podTemplateInitContainersFunc(fl, initContainerOptions); err != nil {
		return errors.Wrap(err, "pod template init containers")
	}

	containerOptions := baseOptions
	containerOptions.containers = pt.podTemplateSpec.Spec.Containers
	containerOptions.isInit = false

	if err := pt.podTemplateContainersFunc(fl, containerOptions); err != nil {
		return errors.Wrap(err, "pod template containers")
	}

	if err := pt.podTemplatePodConfigurationFunc(fl, baseOptions); err != nil {
		return errors.Wrap(err, "pod template pod configuration")
	}

	return nil
}

func podTemplateHeader(fl *flexlayout.FlexLayout, options podTemplateLayoutOptions) error {
	headerSection := fl.AddSection()
	podTemplateHeader := NewPodTemplateHeader(options.podTemplateSpec.ObjectMeta.Labels)
	headerLabels := podTemplateHeader.Create()

	if err := headerSection.Add(headerLabels, component.WidthFull); err != nil {
		return errors.Wrap(err, "add pod template header")
	}

	return nil
}

func podTemplateContainers(fl *flexlayout.FlexLayout, options podTemplateLayoutOptions) error {
	if len(options.containers) < 1 {
		return nil
	}

	portForwarder := options.printOptions.DashConfig.PortForwarder()

	containerSection := fl.AddSection()

	for _, container := range options.containers {
		containerConfig := NewContainerConfiguration(options.parent, &container, portForwarder, options.isInit, options.printOptions)
		summary, err := containerConfig.Create()
		if err != nil {
			return err
		}

		if err := containerSection.Add(summary, component.WidthHalf); err != nil {
			return errors.Wrap(err, "add container")
		}
	}

	return nil
}

func podTemplatePodConfiguration(fl *flexlayout.FlexLayout, options podTemplateLayoutOptions) error {
	podSection := fl.AddSection()

	volumeTable, err := printVolumes(options.podTemplateSpec.Spec.Volumes)
	if err != nil {
		return errors.Wrap(err, "print volumes")
	}
	if !volumeTable.IsEmpty() {
		if err := podSection.Add(volumeTable, component.WidthHalf); err != nil {
			return err
		}
	}

	tolerationList, err := printTolerations(options.podTemplateSpec.Spec)
	if err != nil {
		return errors.Wrap(err, "print tolerations")
	}
	if !tolerationList.IsEmpty() {
		if err := podSection.Add(tolerationList, component.WidthHalf); err != nil {
			return err
		}
	}

	affinityList, err := printAffinity(options.podTemplateSpec.Spec)
	if err != nil {
		return errors.Wrap(err, "print affinities")
	}
	if !affinityList.IsEmpty() {
		if err := podSection.Add(affinityList, component.WidthHalf); err != nil {
			return err
		}
	}

	return nil
}

// PodTemplateHeader creates a pod template header. It consists of a
// selectors component with title `Pod Template` and the associated
// match selectors.
type PodTemplateHeader struct {
	labels map[string]string
}

// NewPodTemplateHeader creates an instance of PodTemplateHeader.
func NewPodTemplateHeader(labels map[string]string) *PodTemplateHeader {
	return &PodTemplateHeader{
		labels: labels,
	}
}

// Create creates a labels component.
func (pth *PodTemplateHeader) Create() *component.Labels {
	view := component.NewLabels(pth.labels)
	view.Metadata.SetTitleText("Pod Template")

	return view
}
