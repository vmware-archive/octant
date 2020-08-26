/*
 * Copyright (c) 2020 the Octant contributors. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package javascript

import (
	"context"
	"testing"

	"github.com/vmware-tanzu/octant/pkg/action"

	"github.com/vmware-tanzu/octant/pkg/event"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/vmware-tanzu/octant/pkg/event/fake"
)

func TestDashboardSendEvent_Name(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	wsClient := fake.NewMockWSClientGetter(ctrl)

	d := NewDashboardSendEvent(wsClient)

	want := "SendEvent"
	got := d.Name()

	require.Equal(t, want, got)
}

func TestDashboardSendEvent_Call(t *testing.T) {
	type ctorArgs struct {
		wsClient func(ctx context.Context, ctrl *gomock.Controller) event.WSClientGetter
	}
	tests := []struct {
		name     string
		ctorArgs ctorArgs
		call     string
		wantErr  bool
	}{
		{
			name: "empty clientID",
			ctorArgs: ctorArgs{
				wsClient: func(ctx context.Context, ctrl *gomock.Controller) event.WSClientGetter {
					wsClient := fake.NewMockWSClientGetter(ctrl)
					return wsClient
				},
			},
			call:    `dashClient.SendEvent("", "event.octant.dev/alert", {namespace:'test', apiVersion: 'v1', kind:'Pod'})`,
			wantErr: true,
		},
		{
			name: "action sends",
			ctorArgs: ctorArgs{
				wsClient: func(ctx context.Context, ctrl *gomock.Controller) event.WSClientGetter {
					event := event.CreateEvent(event.EventTypeAlert, action.Payload{"namespace": "test", "apiVersion": "v1", "kind": "ReplicaSet"})

					wsClient := fake.NewMockWSClientGetter(ctrl)
					wsEventSender := fake.NewMockWSEventSender(ctrl)

					wsClient.EXPECT().Get("my-client-id").Return(wsEventSender)
					wsEventSender.EXPECT().Send(event)
					return wsClient
				},
			},
			call:    `dashClient.SendEvent("my-client-id", "event.octant.dev/alert", {namespace:'test', apiVersion: 'v1', kind:'ReplicaSet'})`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			d := NewDashboardSendEvent(tt.ctorArgs.wsClient(ctx, ctrl))

			runner := functionRunner{wantErr: tt.wantErr}
			runner.run(ctx, t, d, tt.call)

		})
	}
}
