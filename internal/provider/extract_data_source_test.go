// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExtractDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccExtractDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.jwk_extract.test", "jwk", testAccExtractDataSourceResult),
				),
			},
		},
	})
}

const testAccExtractDataSourceConfig = `
data "jwk_extract" "test" {
	public_certificate = <<EOT
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE6GDI8fW9IHNg4kcUfmQ/5xZB3kKW
wlIEoiwQPj4JkS72uOECKq/RiLP0T6niU83gkCaKJMvAyvtb4o4rb4AUYQ==
-----END PUBLIC KEY-----
EOT
}
`
const testAccExtractDataSourceResult = `{"keys":[{"use":"sig","kty":"EC","kid":"0n7nytRUxT_sx_gMeznMG7BYjtclTps52bmrK-1RgL4","crv":"P-256","alg":"ES256","x":"6GDI8fW9IHNg4kcUfmQ_5xZB3kKWwlIEoiwQPj4JkS4","y":"9rjhAiqv0Yiz9E-p4lPN4JAmiiTLwMr7W-KOK2-AFGE"}]}`
