/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/vmware-tanzu/octant/internal/event"
	"github.com/vmware-tanzu/octant/internal/kubeconfig"
	"github.com/vmware-tanzu/octant/internal/log"
	"github.com/vmware-tanzu/octant/internal/octant"
	"github.com/vmware-tanzu/octant/pkg/action"
)

const (
	UploadKubeConfig = "action.octant.dev/uploadKubeConfig"
)

type LoadingManager struct {
	client     OctantClient
	kubeConfig kubeconfig.TemporaryKubeConfig
	ctx        context.Context
	poller     Poller
}

var _ StateManager = (*LoadingManager)(nil)

func NewLoadingManager(temporaryKubeConfig kubeconfig.TemporaryKubeConfig) *LoadingManager {
	lm := &LoadingManager{
		kubeConfig: temporaryKubeConfig,
		poller:     NewInterruptiblePoller("loading"),
	}

	return lm
}

func (l *LoadingManager) Handlers() []octant.ClientRequestHandler {
	return []octant.ClientRequestHandler{
		{
			RequestType: UploadKubeConfig,
			Handler:     l.UploadKubeConfig,
		},
	}
}

func (l *LoadingManager) Start(ctx context.Context, state octant.State, client OctantClient) {
	l.client = client
	l.ctx = ctx

	l.poller.Run(ctx, nil, l.runUpdate(state, client), event.DefaultScheduleDelay)
}

func (l *LoadingManager) runUpdate(state octant.State, client OctantClient) PollerFunc {
	var previous []byte
	return func(ctx context.Context) bool {
		logger := log.From(ctx)

		if ctx.Err() == nil {
			event := octant.Event{
				Type: octant.EventTypeLoading,
			}

			cur, err := json.Marshal(event)
			if err != nil {
				logger.WithErr(err).Errorf("unable to marshal loading")
				return false
			}

			if bytes.Compare(previous, cur) != 0 {
				previous = cur
				client.Send(event)
			}
		}
		return false
	}
}

func (l *LoadingManager) UploadKubeConfig(state octant.State, payload action.Payload) error {
	kubeConfig64, err := payload.String("kubeConfig")
	if err != nil {
		return fmt.Errorf("getting kubeConfig from payload: %w", err)
	}

	kubeConfig, err := base64.StdEncoding.DecodeString(kubeConfig64)
	if err != nil {
		return fmt.Errorf("decode base64 kubeconfig: %w", err)
	}

	// Check if kube config can be converted to config object
	// TODO: Show error elsewhere instead of notifier
	if _, err := clientcmd.Load(kubeConfig); err != nil {
		message := fmt.Sprintf("Error parsing KubeConfig: %v", err)
		alert := action.CreateAlert(action.AlertTypeError, message, action.DefaultAlertExpiration)
		state.SendAlert(alert)
		return err
	}

	l.kubeConfig.KubeConfig <- string(kubeConfig)

	tempFile, err := ioutil.TempFile(os.TempDir(), "kubeconfig")
	if err != nil {
		return err
	}

	if _, err := tempFile.Write(kubeConfig); err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	l.client.Send(octant.Event{
		Type: octant.EventTypeRefresh,
	})

	l.kubeConfig.Path <- tempFile.Name()
	return nil
}
