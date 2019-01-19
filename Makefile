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
test:
	@echo "-> $@"
	@env GO111MODULE=on go test -v -mod=vendor ./internal/...

# Run govet
.PHONY: vet
vet:
	@echo "-> $@"
	@env GO111MODULE=on go vet -mod=vendor ./...

hcli-dev:
	@mkdir -p ./build
	@env GO111MODULE=on $(GOBUILD) -o build/hcli -mod=vendor $(GO_FLAGS) ./cmd/hcli

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none npm start

web-deps:
	@cd web; npm ci

web-build: web-deps
	@cd web; npm run-script build
	@go generate ./web

ui-server:
	DASH_DISABLE_OPEN_BROWSER=false DASH_LISTENER_ADDR=localhost:3001 $(GOCMD) run ./cmd/hcli/main.go dash $(DASH_FLAGS)

ui-client:
	cd web; API_BASE=http://localhost:3001 npm run start

gen-electron:
	@GOCACHE=${HOME}/cache/go-build astilectron-bundler -v -c configs/electron/bundler.json

.PHONY: release
release:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --follow-tags

.PHONY: ci
ci: gen-electron test vet web-build hcli-dev
