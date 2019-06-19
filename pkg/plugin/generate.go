/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package plugin

//go:generate mockgen -source=runner.go -destination=./fake/mock_runners.go -package=fake github.com/heptio/developer-dash/pkg/plugin Runners
//go:generate mockgen -destination=./fake/mock_manager_store.go -package=fake github.com/heptio/developer-dash/pkg/plugin ManagerStore
//go:generate mockgen -destination=./fake/mock_client_factory.go -package=fake github.com/heptio/developer-dash/pkg/plugin ClientFactory
//go:generate mockgen -source=client.go -destination=./fake/mock_service.go -package=fake github.com/heptio/developer-dash/pkg/plugin Service
//go:generate mockgen -source=broker.go -destination=./fake/mock_broker.go -package=fake github.com/heptio/developer-dash/pkg/plugin Broker
//go:generate mockgen -source=dashboard/dashboard.pb.go -destination=./fake/mock_plugin_client.go -package=fake github.com/heptio/developer-dash/pkg/plugin/dashboard PluginClient
//go:generate mockgen -source=../../vendor/github.com/hashicorp/go-plugin/protocol.go -destination=./fake/mock_client_protocol.go -package=fake github.com/hashicorp/go-plugin ClientProtocol
