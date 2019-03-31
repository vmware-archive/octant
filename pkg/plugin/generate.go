package plugin

//go:generate mockgen -destination=./fake/mock_runners.go -package=fake github.com/heptio/developer-dash/pkg/plugin Runners
//go:generate mockgen -destination=./fake/mock_manager_store.go -package=fake github.com/heptio/developer-dash/pkg/plugin ManagerStore
//go:generate mockgen -destination=./fake/mock_client_factory.go -package=fake github.com/heptio/developer-dash/pkg/plugin ClientFactory
//go:generate mockgen -destination=./fake/mock_service.go -package=fake github.com/heptio/developer-dash/pkg/plugin Service
//go:generate mockgen -destination=./fake/mock_broker.go -package=fake github.com/heptio/developer-dash/pkg/plugin Broker
//go:generate mockgen -destination=./fake/mock_plugin_client.go -package=fake github.com/heptio/developer-dash/pkg/plugin/proto PluginClient
//go:generate mockgen -destination=./fake/mock_client_protocol.go -package=fake github.com/hashicorp/go-plugin ClientProtocol
