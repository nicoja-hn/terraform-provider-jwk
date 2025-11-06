terraform {
  required_providers {
    jwk = {
      source = "nicoja-hn/jwk"
    }
  }
}

provider "jwk" {
  # Provider configuration (if needed)
}

# Example: Currently the provider has scaffolding resources
# These resources are named "jwk_example" and are placeholders
# for real JWK functionality to be implemented

# Data source example (placeholder)
data "jwk_example" "test" {
  configurable_attribute = "example-value"
}

# Resource example (placeholder)
resource "jwk_example" "test" {
  configurable_attribute = "example-value"
}

# Output the data source ID
output "data_source_id" {
  value       = data.jwk_example.test.id
  description = "The ID from the example data source"
}

# Output the resource ID
output "resource_id" {
  value       = jwk_example.test.id
  description = "The ID from the example resource"
}

# Note: In a real JWK provider implementation, you would have resources like:
#
# Generate an RSA key pair as JWK
# resource "jwk_rsa_key" "example" {
#   key_size = 2048
#   algorithm = "RS256"
#   key_id = "my-key-1"
# }
#
# Generate an ECDSA key pair as JWK
# resource "jwk_ec_key" "example" {
#   curve = "P-256"
#   algorithm = "ES256"
#   key_id = "my-ec-key-1"
# }
#
# Use a JWK from an external source
# data "jwk_key_set" "external" {
#   url = "https://example.com/.well-known/jwks.json"
# }
