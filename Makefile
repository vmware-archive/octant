

all: web-build update

build:
	go run build.go build
	cp build/octant /Users/jstrachan/bin/octant

update:
	go generate ./web
	go run build.go build
	cp build/octant /Users/jstrachan/bin/octant

web-build:
	cd web && npm run build
	go generate ./web