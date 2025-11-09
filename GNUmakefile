# Detect OS and Architecture
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
PROVIDER_NAME := jwk
NAMESPACE := nicoja-hn
VERSION := 1.0.0

# Detect if OpenTofu or Terraform is being used
HAS_TOFU := $(shell command -v tofu 2> /dev/null)
ifdef HAS_TOFU
    REGISTRY := registry.opentofu.org
else
    REGISTRY := registry.terraform.io
endif

# Installation paths
PLUGIN_DIR := $(HOME)/.terraform.d/plugins/$(REGISTRY)/$(NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS)_$(ARCH)
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
	@echo "Creating configuration file..."
	@$(MAKE) -s create-config
	@echo ""
	@echo "Installation complete! To use the provider:"
	@echo "  1. Run: terraform init (or tofu init)"
	@echo "  2. Then: terraform plan (or tofu plan)"

.PHONY: create-config
create-config: ## Create or update Terraform/OpenTofu config file
	@echo "Checking for Terraform/OpenTofu config..."
	@if command -v tofu >/dev/null 2>&1; then \
		CONFIG_FILE="$(HOME)/.tofurc"; \
		REGISTRY="registry.opentofu.org"; \
		echo "OpenTofu detected, creating $$CONFIG_FILE for registry.opentofu.org"; \
	elif command -v terraform >/dev/null 2>&1; then \
		CONFIG_FILE="$(HOME)/.terraformrc"; \
		REGISTRY="registry.terraform.io"; \
		echo "Terraform detected, creating $$CONFIG_FILE for registry.terraform.io"; \
	else \
		CONFIG_FILE="$(HOME)/.terraformrc"; \
		REGISTRY="registry.terraform.io"; \
		echo "Creating default config at $$CONFIG_FILE"; \
	fi; \
	if [ -f "$$CONFIG_FILE" ]; then \
		echo "Backing up existing config to $$CONFIG_FILE.backup"; \
		cp "$$CONFIG_FILE" "$$CONFIG_FILE.backup"; \
	fi; \
	echo "provider_installation {" > "$$CONFIG_FILE"; \
	echo "  # For all other providers, use the normal registry" >> "$$CONFIG_FILE"; \
	echo "  direct {" >> "$$CONFIG_FILE"; \
	echo "    exclude = [\"$(NAMESPACE)/*\"]" >> "$$CONFIG_FILE"; \
	echo "  }" >> "$$CONFIG_FILE"; \
	echo "" >> "$$CONFIG_FILE"; \
	echo "  filesystem_mirror {" >> "$$CONFIG_FILE"; \
	echo "    path    = \"$(HOME)/.terraform.d/plugins\"" >> "$$CONFIG_FILE"; \
	echo "    include = [\"$(NAMESPACE)/*\"]" >> "$$CONFIG_FILE"; \
	echo "  }" >> "$$CONFIG_FILE"; \
	echo "}" >> "$$CONFIG_FILE"; \
	echo "Config file created/updated: $$CONFIG_FILE"; \
	echo ""; \
	echo "Using registry: $$REGISTRY"; \
	echo "Provider installed at: $(PLUGIN_DIR)"

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

.PHONY: show-config
show-config: ## Show what the config file would contain
	@echo "Configuration that will be created:"
	@echo ""
	@echo "provider_installation {"
	@echo "  # For all other providers, use the normal registry"
	@echo "  direct {"
	@echo "    exclude = [\"$(NAMESPACE)/*\"]"
	@echo "  }"
	@echo ""
	@echo "  filesystem_mirror {"
	@echo "    path    = \"$(HOME)/.terraform.d/plugins\""
	@echo "    include = [\"$(NAMESPACE)/*\"]"
	@echo "  }"
	@echo "}"
	@echo ""
	@echo "Provider will be installed to:"
	@echo "  $(PLUGIN_DIR)"
	@echo ""
	@if command -v tofu >/dev/null 2>&1; then \
		echo "Registry: registry.opentofu.org (OpenTofu detected)"; \
	else \
		echo "Registry: registry.terraform.io (Terraform detected)"; \
	fi
