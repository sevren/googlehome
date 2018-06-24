# Set an output prefix, which is the local directory if not specified
PREFIX?=$(shell pwd)

# Setup name variables for the package/tool
NAME := googlehomekodi-go
PKG := github.com/sevren/googlehome/$(NAME)

# Set any default go build tags
BUILDTAGS :=
DOCKERREG := registry.dockerhub.com/sevren/googlehomekodi

# Set the build dir, where built cross-compiled binaries will be output
BUILDDIR := ${PREFIX}/build

GO_LDFLAGS_STATIC=--ldflags "-linkmode external -extldflags -static"

# Set our default go compiler
GO := go

# List the GOOS and GOARCH to build
GOOSARCHES = linux/amd64

.PHONY: build
build: $(NAME) ## Builds a dynamic executable or package

$(NAME): $(wildcard *.go) $(wildcard */*.go) VERSION.txt
	@echo "+ $@"
	@dep ensure
	$(GO) build  -gcflags "-N -l" -tags "$(BUILDTAGS)" ${GO_LDFLAGS} -o $(NAME) .

.PHONY: docker
docker: static ## Builds a docker image
	@echo "+ $@"
	@dep ensure
	@docker build -t $(DOCKERREG)$(NAME) .

.PHONY: static
static: ## Builds a static executable
	@echo "+ $@"
	$(GO) build -tags netgo -installsuffix netgo ${GO_LDFLAGS_STATIC} -o $(NAME) .


.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

