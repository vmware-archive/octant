/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package api_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/vmware-tanzu/octant/internal/api"
	"github.com/vmware-tanzu/octant/internal/api/fake"
	"github.com/vmware-tanzu/octant/internal/octant"

	"k8s.io/client-go/tools/clientcmd"
)

func Test_watchConfig(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	octantClient := fake.NewMockOctantClient(controller)

	octantClient.EXPECT().Send(octant.Event{
		Type: octant.EventTypeRefresh,
	})

	fs := afero.NewMemMapFs()
	config := clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	afero.WriteFile(fs, config, []byte(""), 0755)

	kubeConfigPath := make(chan string, 1)

	manager := api.NewLoadingManager()
	manager.WatchConfig(kubeConfigPath, octantClient, fs)

	select {
	case path := <-kubeConfigPath:
		assert.Equal(t, path, config)
	case <-time.After(3 * time.Second):
		t.FailNow()
	}
}
