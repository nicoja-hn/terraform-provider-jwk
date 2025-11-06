# Example Quick Start Guide

This directory contains example Terraform configurations for the JWK provider.

## Current Status

⚠️ **Note**: The provider currently contains scaffolding/example resources (`jwk_example`) as placeholders. Real JWK functionality needs to be implemented.

## Using the Examples

### 1. Basic Provider Setup

See `provider/provider.tf` for the basic provider configuration.

### 2. Testing Locally

If you're developing the provider locally:

```bash
# Build and install the provider
cd /path/to/terraform-provider-jwk
make install

# Navigate to an example
cd examples/provider

# Initialize and plan (or just plan if using dev_overrides)
terraform plan
terraform apply
```

## Planned JWK Resources (To Be Implemented)

The following resources would be typical for a JWK provider:

### Resources

- `jwk_rsa_key` - Generate RSA key pairs as JWK
- `jwk_ec_key` - Generate Elliptic Curve key pairs as JWK
- `jwk_symmetric_key` - Generate symmetric keys as JWK
- `jwk_key_set` - Manage JWK Sets (JWKS)

### Data Sources

- `jwk_key_set` - Read JWK Set from a URL or file
- `jwk_public_key` - Extract public key from a JWK

### Example Use Cases

1. **Generate signing keys for JWT tokens**
2. **Rotate keys automatically**
3. **Manage key sets for multiple environments**
4. **Import existing keys**
5. **Export public keys for verification**

## Contributing

To implement real JWK resources:

1. Review the scaffolding in `internal/provider/example_resource.go`
2. Create new resource files like `internal/provider/rsa_key_resource.go`
3. Implement the JWK generation/management logic
4. Update these examples with real usage

See `TESTING.md` in the repository root for development guidelines.
