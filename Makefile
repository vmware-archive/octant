SHELL=/bin/bash
BUILD_TIME=$(shell date -u +%Y-%m-%dT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build

hcli-dev:
	@mkdir -p ./build
	@$(GOBUILD) -o build/hcli $(GO_FLAGS) ./cmd/hcli

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none npm start

web-deps:
	@cd web; npm ci

web-build: web-deps
	@cd web; npm run-script build
	@go generate ./web

embed-go: web-build
	@cd web; rice embed-go

ui-server:
	DASH_DISABLE_OPEN_BROWSER=false DASH_LISTENER_ADDR=localhost:3001 $(GOCMD) run ./cmd/hcli/main.go dash

ui-client:
	cd web; API_BASE=http://localhost:3001 npm run start
