/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"time"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/config"
)

type buildInfoResponse struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Time    string `json:"time"`
}

type kubeConfigPathResponse struct {
	Path string `json:"path"`
}

type HelperGeneratorOption func(generator *HelperGenerator)

type HelperGenerator struct {
	DashConfig config.Dash
}

var _ octant.Generator = (*HelperGenerator)(nil)

func NewHelperGenerator(dashConfig config.Dash, options ...HelperGeneratorOption) *HelperGenerator {
	hg := &HelperGenerator{
		DashConfig: dashConfig,
	}

	for _, option := range options {
		option(hg)
	}

	return hg
}

func (h *HelperGenerator) Events(ctx context.Context) ([]event.Event, error) {
	version, commit, time := h.DashConfig.BuildInfo()

	resp := buildInfoResponse{
		Version: version,
		Commit:  commit,
		Time:    time,
	}

	buildInfoEvent := event.Event{
		Type: event.EventTypeBuildInfo,
		Data: resp,
	}

	kubeConfigPathEvent := event.Event{
		Type: event.EventTypeKubeConfigPath,
		Data: kubeConfigPathResponse{h.DashConfig.KubeConfigPath()},
	}

	return []event.Event{buildInfoEvent, kubeConfigPathEvent}, nil
}

func (HelperGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

func (HelperGenerator) Name() string {
	return "buildInfo"
}
