/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 */

package dash_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/dash"
	"github.com/vmware-tanzu/octant/pkg/dash/fake"
)

func TestConfigWatcher_Add(t *testing.T) {
	tests := []struct {
		name         string
		filenames    []string
		filenamesErr error
	}{
		{
			name:      "in general",
			filenames: []string{"kubeconfig"},
		},
		{
			name:         "file add error",
			filenames:    []string{"kubeconfig"},
			filenamesErr: fmt.Errorf("boom"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			watcherConfig := fake.NewMockWatcherConfig(controller)
			fileWatcher := fake.NewMockFileWatcher(controller)

			for _, f := range tt.filenames {
				fileWatcher.EXPECT().Add(f).Return(tt.filenamesErr)
			}

			fileWatcherOption := dash.ConfigWatcherFileWatcher(fileWatcher)

			ctx := context.Background()
			cw, err := dash.NewConfigWatcher(watcherConfig, fileWatcherOption)
			require.NoError(t, err)

			err = cw.Add(ctx, tt.filenames...)
			if tt.filenamesErr != nil {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestConfigWatcher_Watch(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	watcherConfig := fake.NewMockWatcherConfig(controller)

	watcherConfig.EXPECT().CurrentContext().Return("name")

	ch := make(chan bool, 1)
	watcherConfig.EXPECT().UseContext(gomock.Any(), "name").
		DoAndReturn(func(_ context.Context, _ string) error {
			ch <- true
			return nil
		})

	fileWatcher := fake.NewMockFileWatcher(controller)

	eventCh := make(chan fsnotify.Event)
	fileWatcher.EXPECT().Events().Return(eventCh).AnyTimes()

	errCh := make(chan error)
	fileWatcher.EXPECT().Errors().Return(errCh).AnyTimes()

	fileWatcherOption := dash.ConfigWatcherFileWatcher(fileWatcher)

	ctx, cancel := context.WithCancel(context.Background())
	cw, err := dash.NewConfigWatcher(watcherConfig, fileWatcherOption)
	require.NoError(t, err)

	go cw.Watch(ctx)

	event := fsnotify.Event{
		Name: "event",
		Op:   0,
	}

	eventCh <- event
	<-ch
	cancel()
}
