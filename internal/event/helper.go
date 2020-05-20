/*
Copyright (c) 2020 the Octant contributors. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"context"
	"time"

	"github.com/vmware-tanzu/octant/internal/config"
	"github.com/vmware-tanzu/octant/internal/octant"
)

type buildInfoResponse struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Time    string `json:"time"`
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

func (h *HelperGenerator) Event(ctx context.Context) (octant.Event, error) {
	version, commit, time := h.DashConfig.BuildInfo()

	resp := buildInfoResponse{
		Version: version,
		Commit:  commit,
		Time:    time,
	}

	e := octant.Event{
		Type: octant.EventTypeBuildInfo,
		Data: resp,
	}

	return e, nil
}

func (HelperGenerator) ScheduleDelay() time.Duration {
	return DefaultScheduleDelay
}

func (HelperGenerator) Name() string {
	return "buildInfo"
}
