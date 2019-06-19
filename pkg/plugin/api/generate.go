/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package api

//go:generate mockgen -source=server.go -destination=./fake/mock_dash_service.go -package=fake github.com/heptio/developer-dash/pkg/plugin/api Service
