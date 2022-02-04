.PHONY: test
default: test

# Start host containers used for playground and testing
hosts: clean
	docker network create remote
	docker build -t remotehost tests
	docker run --rm -d --net remote --name remotehost -p 8022:22 remotehost
	docker run --rm -d --net remote --name remotehost2 -p 8023:22 remotehost

# Stop containers
clean:
	docker rm -f remotehost remotehost2
	docker network rm remote || true

# Run acceptance tests
test: hosts
	docker run --rm --net remote -v $(PWD):/app -v ~/go:/go --workdir /app \
	-e "TF_ACC=1" -e "TF_ACC_TERRAFORM_VERSION=1.0.11" -e "TESTARGS=$(TESTARGS)" \
	golang:1.16 bash tests/test.sh

# Install provider in playground
INSTALL_DIR=playground
TARGET_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PROVIDER_PATH=.terraform/providers/registry.terraform.io/tenstad/remote/99.0.0/$(TARGET_ARCH)
BIN_PATH=$(INSTALL_DIR)/$(PROVIDER_PATH)/terraform-provider-remote_v99.0.0
install:
	mkdir -p $(INSTALL_DIR)/$(PROVIDER_PATH)
	go build -ldflags="-s -w -X main.version=99.0.0" -o $(BIN_PATH)
