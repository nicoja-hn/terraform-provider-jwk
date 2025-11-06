# Local Testing Guide for JWK Terraform Provider

This guide shows you how to develop and test the provider locally.

## Prerequisites

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0

## 1. Build the Provider Locally

```bash
# In the project directory
go build -o terraform-provider-jwk
```

## 2. Install Provider for Local Development

### Option A: Using `go install`

```bash
go install
```

This installs the provider in `$GOPATH/bin`.

### Option B: Manual Installation with Development Override

Create or edit `~/.terraformrc` (Linux/Mac) or `%APPDATA%\terraform.rc` (Windows):

```hcl
provider_installation {
  dev_overrides {
    "nicoja-hn/jwk" = "/path/to/terraform-provider-jwk"
  }

  # For all other providers, use the normal registry
  direct {}
}
```

Replace `/path/to/terraform-provider-jwk` with the absolute path to your project directory.

**Important:** With dev_overrides, Terraform ignores the `required_providers` version - you must manually rebuild the provider when you make changes.

## 3. Create Test Configuration

Create a file `test.tf`:

```hcl
terraform {
  required_providers {
    jwk = {
      source = "nicoja-hn/jwk"
    }
  }
}

provider "jwk" {
  # Provider configuration here
}

# Example Data Source
data "jwk_example" "test" {
  configurable_attribute = "test"
}

output "example_output" {
  value = data.jwk_example.test
}

# Example Resource
resource "jwk_example" "test" {
  configurable_attribute = "test"
}
```

## 4. Initialize and Run Terraform

```bash
# Initialize Terraform
terraform init

# Show plan
terraform plan

# Apply changes
terraform apply

# Clean up
terraform destroy
```

## 5. Run Provider Tests

### Unit Tests

```bash
# All tests (without Acceptance Tests)
go test -v ./... -short

# Test specific package
go test -v ./internal/provider -short
```

### Acceptance Tests

Acceptance tests create real resources and may incur costs!

```bash
# All Acceptance Tests
TF_ACC=1 go test -v ./... -timeout 120m

# Specific test
TF_ACC=1 go test -v ./internal/provider -run TestAccExampleResource -timeout 10m
```

## 6. Check Code Quality

```bash
# Run linter
golangci-lint run

# Format code
go fmt ./...

# Run go vet
go vet ./...
```

## 7. Generate Documentation

```bash
# Generate Terraform docs
go generate

# This runs:
# - terraform fmt -recursive ./examples/
# - tfplugindocs (generates docs/)
```

**Note:** You need Terraform installed for formatting examples.

## Debugging

### With Debugger (e.g., Delve)

1. Start the provider in debug mode:

```bash
go build -o terraform-provider-jwk
./terraform-provider-jwk -debug
```

2. The provider outputs an address and environment variable set:

```
Provider started. To attach Terraform CLI, set the TF_REATTACH_PROVIDERS environment variable with the following:

TF_REATTACH_PROVIDERS='{"nicoja-hn/jwk":{"Protocol":"grpc","ProtocolVersion":6,"Pid":12345,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin12345.sock"}}}'
```

3. Export this variable and run Terraform:

```bash
export TF_REATTACH_PROVIDERS='...'
terraform plan
```

### With Logs

```bash
# Enable Terraform logs
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log

terraform apply

# Provider-specific logs
export TF_LOG_PROVIDER=TRACE
```

## Common Issues

### "Provider not found"

- Ensure `dev_overrides` has the correct path
- Verify the binary is named `terraform-provider-jwk`
- Rebuild after changes: `go build`

### "Invalid provider configuration"

- Check the provider schema definition in `internal/provider/provider.go`
- Ensure all required fields are set

### Tests Failing

- Ensure `TF_ACC=1` is only set for Acceptance Tests
- Verify Terraform is installed
- Check Terraform version: `terraform version`

## Additional Resources

- [Terraform Plugin Framework Docs](https://developer.hashicorp.com/terraform/plugin/framework)
- [HashiCorp Learn - Plugin Development](https://learn.hashicorp.com/collections/terraform/providers-plugin-framework)
- [Terraform Registry Publishing](https://www.terraform.io/docs/registry/providers/publishing.html)
