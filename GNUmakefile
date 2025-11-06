# Detect OS and Architecture
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
PROVIDER_NAME := jwk
NAMESPACE := nicoja-hn
VERSION := 1.0.0

# Installation paths
PLUGIN_DIR := $(HOME)/.terraform.d/plugins/registry.terraform.io/$(NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS)_$(ARCH)
BINARY_NAME := terraform-provider-$(PROVIDER_NAME)_v$(VERSION)

default: testacc

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the provider binary
	go build -o $(BINARY_NAME)

.PHONY: install
install: build ## Build and install provider to Terraform plugin directory
	@echo "Installing provider to $(PLUGIN_DIR)"
	@mkdir -p $(PLUGIN_DIR)
	@cp $(BINARY_NAME) $(PLUGIN_DIR)/
	@echo "Provider installed successfully!"
	@echo ""
	@echo "To use this installation:"
	@echo "  1. Remove or comment out dev_overrides from ~/.terraformrc or ~/.tofurc"
	@echo "  2. Run: terraform init (or tofu init)"

.PHONY: install-local
install-local: ## Build for dev_overrides (in current directory)
	@echo "Building provider for dev_overrides"
	@go build -o terraform-provider-$(PROVIDER_NAME)
	@echo "Provider built as terraform-provider-$(PROVIDER_NAME)"
	@echo ""
	@echo "Make sure your ~/.terraformrc or ~/.tofurc has:"
	@echo "  dev_overrides {"
	@echo "    \"$(NAMESPACE)/$(PROVIDER_NAME)\" = \"$(shell pwd)\""
	@echo "  }"

.PHONY: test
test: ## Run unit tests
	go test -v ./... -short

# Run acceptance tests
.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: clean
clean: ## Clean build artifacts
	@rm -f terraform-provider-$(PROVIDER_NAME)
	@rm -f $(BINARY_NAME)
	@echo "Cleaned build artifacts"

.PHONY: uninstall
uninstall: ## Remove provider from Terraform plugin directory
	@rm -rf ~/.terraform.d/plugins/registry.terraform.io/$(NAMESPACE)/$(PROVIDER_NAME)
	@echo "Provider uninstalled from Terraform plugin directory"

.PHONY: docs
docs: ## Generate documentation
	go generate

.PHONY: show-paths
show-paths: ## Show installation paths
	@echo "OS/Arch: $(OS)/$(ARCH)"
	@echo "Plugin directory: $(PLUGIN_DIR)"
	@echo "Binary name: $(BINARY_NAME)"
	@echo "Current directory: $(shell pwd)"
