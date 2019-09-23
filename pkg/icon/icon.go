/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package icon

import (
	"fmt"
	"io/ioutil"

	rice "github.com/GeertJohan/go.rice"
)

//go:generate rice embed-go

const (
	Applications            = "applications"
	ApplicationsApplication = "application"

	ClusterOverview                   = "cluster"
	ClusterOverviewClusterRole        = "c-role"
	ClusterOverviewClusterRoleBinding = "crb"
	ClusterOverviewNode               = "node"

	Configuration       = "cog"
	ConfigurationPlugin = "plugin"

	CustomResourceDefinition = "crd"

	Overview                      = "dashboard"
	OverviewConfigMap             = "cm"
	OverviewCronJob               = "cronjob"
	OverviewDaemonSet             = "ds"
	OverviewDeployment            = "deploy"
	OverviewIngress               = "ing"
	OverviewJob                   = "job"
	OverviewPersistentVolumeClaim = "pvc"
	OverviewPod                   = "pod"
	OverviewReplicaSet            = "rs"
	OverviewReplicationController = "deploy"
	OverviewRole                  = "role"
	OverviewRoleBinding           = "rb"
	OverviewSecret                = "secret"
	OverviewService               = "svc"
	OverviewServiceAccount        = "sa"
	OverviewStatefulSet           = "sts"
)

// LoadIcon loads an icon by name.
func LoadIcon(name string) (string, error) {
	box, err := rice.FindBox("svg")
	if err != nil {
		return "", err
	}

	p := fmt.Sprintf("%s.svg", name)

	f, err := box.Open(p)
	if err != nil {
		return "", err
	}

	defer func() {
		cErr := f.Close()
		if cErr != nil && err == nil {
			err = cErr
		}
	}()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
