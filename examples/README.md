# JWK Terraform Provider Examples

This directory contains example configurations for the JWK provider.

## Available Resources

### Data Sources

#### `jwk_extract`

Converts PEM-encoded public certificates to JWK (JSON Web Key) format.

**Supported algorithms:**
- RSA: RS256, RS384, RS512
- ECDSA: ES256, ES384, ES512 (for P-256, P-384, P-521 curves)

**Example:** See `provider/provider.tf`

## Using the Examples

### Local Development

```bash
# Build and install the provider
cd /path/to/terraform-provider-jwk
make install

# Navigate to examples
cd examples/provider

# Run Terraform/OpenTofu
terraform plan
terraform apply
```

### With dev_overrides

If using dev_overrides in `~/.terraformrc` or `~/.tofurc`:

```bash
# Skip init, go directly to plan
terraform plan
terraform apply
```

## Common Use Cases

### 1. Convert TLS Certificates for JWT Verification

```hcl
data "jwk_extract" "api_key" {
  public_certificate = file("path/to/public-key.pem")
  signing_algorithm  = "RS256"
}

# Use in your application configuration
output "jwks" {
  value = jsonencode({
    keys = [jsondecode(data.jwk_extract.api_key.jwk)]
  })
}
```

### 2. Create JWKS for Multiple Keys

```hcl
data "jwk_extract" "key1" {
  public_certificate = file("key1.pem")
  signing_algorithm  = "RS256"
}

data "jwk_extract" "key2" {
  public_certificate = file("key2.pem")
  signing_algorithm  = "ES256"
}

locals {
  jwks = {
    keys = concat(
      jsondecode(data.jwk_extract.key1.jwk).keys,
      jsondecode(data.jwk_extract.key2.jwk).keys
    )
  }
}

output "jwks_json" {
  value = jsonencode(local.jwks)
}
```

## Documentation

For the full provider documentation, see:
- Provider configuration: `docs/index.md`
- Data source `jwk_extract`: `docs/data-sources/extract.md`

## Contributing

To add more examples or improve existing ones, please submit a pull request!
