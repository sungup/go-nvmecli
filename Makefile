OS=linux
ARCH=amd64

GOENV=env GOOS=$(OS) GOARCH=$(ARCH)

TEST_FLAGS =

GOCMD=go
GOBUILD=$(GOENV) $(GOCMD) build
GOTEST=$(GOCMD) test $(TEST_FLAGS)
GOVET=$(GOCMD) vet
GOFMT=$(GOCMD) fmt
GOCLEAN=$(GOCMD) clean
GOLIST=$(GOCMD) list -mod=vendor

PROJECT_PATH=github.com/sungup/go-nvme-ioctl

BUILD_LDFLAG="-extldflags=-static"

SOURCES=$(shell find . -name "*.go")

# build:

format:
	@$(GOFMT) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)

test: format
	@$(GOVET) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)
	@$(GOTEST) $(shell $(GOLIST) ./... | grep -v /$(PROJECT_PATH)/)

gomod-refresh:
	@rm -rf go.sum go.mod vendor
	@$(GOCMD) mod init
	@$(GOCMD) mod vendor -v

clean:
	@$(GOCLEAN)
	@$(GOCLEAN) -testcache
	# TODO add remove output files
	# @rm -f $(TARGETS)
