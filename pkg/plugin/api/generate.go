package api

//go:generate mockgen -source=server.go -destination=./fake/mock_dash_service.go -package=fake github.com/heptio/developer-dash/pkg/plugin/api Service
