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
	@cd web; npm i

web-build: web-deps
	@cd web; npm build
	@go generate ./web