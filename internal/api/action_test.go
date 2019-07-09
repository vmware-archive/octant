/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vmware/octant/internal/api/fake"
	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/mime"
	"github.com/vmware/octant/pkg/action"
)

func Test_action(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	payload := action.Payload{"foo": "bar", "action": "action"}
	actionDispatcher := fake.NewMockActionDispatcher(controller)
	actionDispatcher.EXPECT().
		Dispatch(gomock.Any(), "action", payload)

	logger := log.NopLogger()

	handler := newAction(logger, actionDispatcher)

	ts := httptest.NewServer(handler)
	defer ts.Close()

	client := ts.Client()

	req := updateRequest{
		Update: payload,
	}
	data, err := json.Marshal(&req)
	require.NoError(t, err)

	r := bytes.NewReader(data)

	res, err := client.Post(ts.URL, mime.JSONContentType, r)
	require.NoError(t, err)

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

}
