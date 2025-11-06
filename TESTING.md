# Local Testing Guide for JWK Terraform Provider

This guide shows you how to develop and test the provider locally.

## Prerequisites

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0 or [OpenTofu](https://opentofu.org/docs/intro/install/) >= 1.6

**Note:** This provider works with both Terraform and OpenTofu. Throughout this guide, you can replace `terraform` commands with `tofu` commands.

## 1. Build the Provider Locally

### Using Make (Recommended)

```bash
# See all available commands
make help

# Build for local development with dev_overrides
make install-local

# Build and install to Terraform plugin directory
make install
```

### Using Go directly

```bash
# In the project directory
go build -o terraform-provider-jwk
```

## 2. Install Provider for Local Development

### Option A: Using Make (Recommended)

```bash
# Install to Terraform plugin directory
make install

# This installs to: ~/.terraform.d/plugins/registry.terraform.io/nicoja-hn/jwk/1.0.0/<os>_<arch>/
```

After installation, **remove dev_overrides** from your config file and run `terraform init` or `tofu init`.

### Option B: Using `go install`

```bash
go install
```

This installs the provider in `$GOPATH/bin`.

### Option C: Manual Installation with Development Override

Create or edit the CLI configuration file:
- **Terraform**: `~/.terraformrc` (Linux/Mac) or `%APPDATA%\terraform.rc` (Windows)
- **OpenTofu**: `~/.tofurc` (Linux/Mac) or `%APPDATA%\tofu.rc` (Windows)

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

**Important Notes about dev_overrides:**
- Terraform/OpenTofu ignores the `required_providers` version - you must manually rebuild the provider when you make changes
- **Skip `terraform init` / `tofu init`** when using dev_overrides - go directly to `plan` or `apply`
- The init command will try to download from the registry and may fail, but this is expected and can be ignored

## 3. Create Test Configuration

You can use the examples from the `examples/` directory, or create your own `test.tf`:

```hcl
terraform {
  required_providers {
    jwk = {
      source = "nicoja-hn/jwk"
    }
  }
}

provider "jwk" {
  # No configuration needed
}

# Extract JWK from a public certificate
data "jwk_extract" "example" {
  public_certificate = <<-EOT
    -----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
    4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
    +qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
    kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
    0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
    cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
    mwIDAQAB
    -----END PUBLIC KEY-----
  EOT

  signing_algorithm = "RS256"
}

output "jwk" {
  value     = data.jwk_extract.example.jwk
  sensitive = true
}
```

**Available Data Sources:**
- `jwk_extract` - Convert PEM public certificates to JWK format (supports RS256/384/512 and ES256/384/512)

See `examples/README.md` for more detailed usage examples.

## 4. Run Terraform/OpenTofu

### With dev_overrides (Recommended for Development)

```bash
# SKIP init when using dev_overrides!
# The provider is already available locally

# Show plan
terraform plan
# or with OpenTofu:
tofu plan

# Apply changes
terraform apply
# or with OpenTofu:
tofu apply

# Clean up
terraform destroy
# or with OpenTofu:
tofu destroy
```

### Without dev_overrides (Standard Installation)

```bash
# Initialize Terraform/OpenTofu
terraform init

# Show plan
terraform plan

# Apply changes
terraform apply

# Clean up
terraform destroy
```

**Note:** If you see an error during `init` about the provider not being found in the registry while using dev_overrides, that's expected - just skip `init` and go directly to `plan`/`apply`.

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
