CONTAINER_RUNTIME ?= $(shell command -v podman 2> /dev/null || echo docker)

.PHONY: test
default: test

# Start host containers used for playground and testing
hosts: clean
	$(CONTAINER_RUNTIME) network create remote
	$(CONTAINER_RUNTIME) build -t remotehost tests
	$(CONTAINER_RUNTIME) run --rm -d --net remote --name remotehost -p 8022:22 remotehost
	$(CONTAINER_RUNTIME) run --rm -d --net remote --name remotehost2 -p 8023:22 remotehost

# Stop containers
clean:
	$(CONTAINER_RUNTIME) rm -f remotehost remotehost2 || true
	$(CONTAINER_RUNTIME) network rm remote || true

# Run acceptance tests
test: hosts
	$(CONTAINER_RUNTIME) run --rm --net remote -v ~/go:/go:z -v $(PWD):/provider:z --workdir /provider \
	-e "TF_LOG=INFO" -e "TF_ACC=1" -e "TF_ACC_TERRAFORM_VERSION=1.0.11" -e "TESTARGS=$(TESTARGS)" \
	golang:1.16 bash tests/test.sh

# Install provider in playground
INSTALL_DIR=playground
TARGET_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PROVIDER_PATH=.terraform/providers/registry.terraform.io/tenstad/remote/99.0.0/$(TARGET_ARCH)
BIN_PATH=$(INSTALL_DIR)/$(PROVIDER_PATH)/terraform-provider-remote_v99.0.0
install:
	mkdir -p $(INSTALL_DIR)/$(PROVIDER_PATH)
	go build -ldflags="-s -w -X main.version=99.0.0" -o $(BIN_PATH)

doc:
	go generate
