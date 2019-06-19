# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL=/bin/bash
BUILD_TIME=$(shell date -u +%Y-%m-%dT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

VERSION ?= v0.3.0

.PHONY: version
version:
	@echo $(VERSION)

# Run all tests
.PHONY: test
test: generate
	@echo "-> $@"
	@env go test -v ./internal/... ./pkg/...

# Run govet
.PHONY: vet
vet:
	@echo "-> $@"
	@env go vet  ./internal/... ./pkg/...

octant-dev:
	@mkdir -p ./build
	@env $(GOBUILD) -o build/octant $(GO_FLAGS) -v ./cmd/octant

octant-docker:
	@env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o /octant $(GO_FLAGS) -v ./cmd/octant

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none npm start

generate:
	@echo "-> $@"
	@go generate -v ./pkg/store ./pkg/plugin/api/proto ./pkg/plugin/dashboard ./pkg/plugin/api ./pkg/plugin ./internal/...

go-install:
	$(GOINSTALL) ./vendor/github.com/GeertJohan/go.rice
	$(GOINSTALL) ./vendor/github.com/GeertJohan/go.rice/rice
	$(GOINSTALL) ./vendor/github.com/asticode/go-astilectron-bundler/...
	$(GOINSTALL) ./vendor/github.com/golang/mock/gomock
	$(GOINSTALL) ./vendor/github.com/golang/mock/mockgen
	$(GOINSTALL) ./vendor/github.com/golang/protobuf/protoc-gen-go

# Remove all generated fakes
.PHONY: clean
clean:
	@rm -rf ./internal/portforward/fake
	@rm -rf ./internal/objectstore/fake
	@rm -rf ./internal/queryer/fake
	@rm -rf ./internal/cluster/fake
	@rm -rf ./internal/overview/printer/fake
	@rm -rf ./pkg/plugin/fake
	@rm -rf ./pkg/plugin/api/fake

web-deps:
	@cd web; npm ci

web-build: web-deps
	@cd web; npm run build
	@go generate ./web

web-test: web-deps
	@cd web; npm run test:headless

ui-server:
	OCTANT_DISABLE_OPEN_BROWSER=false OCTANT_LISTENER_ADDR=localhost:3001 $(GOCMD) run ./cmd/octant/main.go $(OCTANT_FLAGS)

ui-client:
	cd web; API_BASE=http://localhost:3001 npm run start

gen-electron:
	@GOCACHE=${HOME}/cache/go-build astilectron-bundler -v -c configs/electron/bundler.json

.PHONY: changelogs
changelogs:
	hacks/changelogs.sh

.PHONY: release
release:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --follow-tags

.PHONY: ci
ci: test vet web-test web-build octant-dev

.PHONY: ci-quick
ci-quick:
	@cd web; npm run build
	@go generate ./web
	make octant-dev

install-test-plugin:
	mkdir -p ~/.config/vmdash/plugins
	go build -o ~/.config/vmdash/plugins/pluginstub github.com/vmware/octant/cmd/pluginstub

.PHONY:
build-deps:
