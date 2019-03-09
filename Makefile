SHELL=/bin/bash
BUILD_TIME=$(shell date -u +%Y-%m-%dT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build

VERSION ?= v0.2.0

.PHONY: version
version:
	@echo $(VERSION)

# Run all tests
.PHONY: test
test: generate
	@echo "-> $@"
	@env go test -v ./internal/...

# Run govet
.PHONY: vet
vet:
	@echo "-> $@"
	@env go vet  ./internal/...

sugarloaf-dev:
	@mkdir -p ./build
	@env $(GOBUILD) -o build/sugarloaf $(GO_FLAGS) ./cmd/sugarloaf

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none npm start

generate:
	@go generate ./internal/...

web-deps:
	@cd web; npm ci

web-build: web-deps
	@cd web; npm run-script build
	@go generate ./web

web-test: web-deps
	@cd web; npm test

ui-server:
	DASH_DISABLE_OPEN_BROWSER=false DASH_LISTENER_ADDR=localhost:3001 $(GOCMD) run ./cmd/sugarloaf/main.go dash $(DASH_FLAGS)

ui-client:
	cd web; API_BASE=http://localhost:3001 npm run start

gen-electron:
	@GOCACHE=${HOME}/cache/go-build astilectron-bundler -v -c configs/electron/bundler.json

.PHONY: release
release:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --follow-tags

.PHONY: ci
ci: test vet web-test web-build sugarloaf-dev
