GOCMD=go
GOBUILD=$(GOCMD) build

hcli-dev:
	@mkdir -p ./build
	@$(GOBUILD) -o build/hcli ./cmd/hcli

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none npm start

web-deps:
	@cd web; npm i

web-build: web-deps
	@cd web; yarn build
	@go generate ./web