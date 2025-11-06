terraform {
  required_providers {
    jwk = {
      source = "nicoja-hn/jwk"
    }
  }
}

provider "jwk" {
  # Provider configuration (currently no configuration needed)
}

# Example: Extract JWK from a TLS public certificate
#
# This data source converts a PEM-encoded public certificate
# to JWK (JSON Web Key) format for use in JWT/JWS operations

# Example with RSA key
data "jwk_extract" "rsa_example" {
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

# Output the generated JWK
output "rsa_jwk" {
  value       = data.jwk_extract.rsa_example.jwk
  description = "The RSA public key in JWK format"
  sensitive   = true
}

# Example with ECDSA P-256 key
data "jwk_extract" "ec_example" {
  public_certificate = <<-EOT
    -----BEGIN PUBLIC KEY-----
    MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYD54V/vp+54P9DXarYqx4MPcm+HK
    RIQzNasYSoRQHQ/6S6Ps8tpMcT+KvIIC8W/e9k0W7Cm72M1P9jU7SLf/vg==
    -----END PUBLIC KEY-----
  EOT

  signing_algorithm = "ES256"
}

output "ec_jwk" {
  value       = data.jwk_extract.ec_example.jwk
  description = "The ECDSA public key in JWK format"
  sensitive   = true
}
