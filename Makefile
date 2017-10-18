SHELL = bash
GOTOOLS = \
	github.com/mitchellh/gox \
	github.com/tcnksm/ghr

GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

GIT_LATEST_TAG=$(shell git describe --abbrev=0 --tags)
VERSION_IMPORT=github.com/andrexus/terraform-provider-proxmox/proxmox
GOLDFLAGS=-X $(VERSION_IMPORT).providerVersion=$(GIT_LATEST_TAG)
OSARCH=darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64
DIST_USER=andrexus

export GOLDFLAGS

all: bin dist

bin: tools
	@echo "==> Building..."
	gox -ldflags "${GOLDFLAGS}" -osarch "${OSARCH}" -output "build/{{.OS}}_{{.Arch}}_{{.Dir}}"

dist:
	ghr -u ${DIST_USER} --token ${GITHUB_TOKEN} --replace --prerelease ${GIT_LATEST_TAG} build/

tools:
	go get -u -v $(GOTOOLS)

.PHONY: all bin dist tools