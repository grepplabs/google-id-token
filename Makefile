.DEFAULT_GOAL := build

.PHONY: clean build build.local build.linux build.osx

BINARY        ?= google-id-token
SOURCES        = $(shell find . -name '*.go' | grep -v /vendor/)
GOPKGS        = $(shell go list ./... | grep -v /vendor/)
BUILD_FLAGS   ?=
LDFLAGS       ?= -w -s
TAG           ?= "v0.0.1"

default: build.local

test:
	go test -v -race `go list ./...`

fmt:
	go fmt $(GOPKGS)

check:
	golint $(GOPKGS)
	go vet $(GOPKGS)

build.local: build/$(BINARY)
build.linux: build/linux/$(BINARY)
build.osx: build/osx/$(BINARY)

build: build/$(BINARY)

build/$(BINARY): $(SOURCES)
	CGO_ENABLED=0 go build -o build/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

build/linux/$(BINARY): $(SOURCES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o build/linux/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

build/osx/$(BINARY): $(SOURCES)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o build/osx/$(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

release: clean build.linux build/osx/$(BINARY)
	git tag $(TAG) && git push --tags
	github-release release -u grepplabs -r $(BINARY) --tag $(TAG)
	github-release upload -u grepplabs -r $(BINARY) -t $(TAG) -f build/linux/$(BINARY) -n linux/amd64/$(BINARY)
	github-release upload -u grepplabs -r $(BINARY) -t $(TAG) -f build/osx/$(BINARY) -n darwin/amd64/$(BINARY)
	github-release info -u grepplabs -r $(BINARY)

clean:
	@rm -rf build
