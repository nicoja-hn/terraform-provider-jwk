// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	// "github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure JWKProvider satisfies various provider interfaces.
var _ provider.Provider = &JWKProvider{}

// JWKProvider defines the provider implementation.
type JWKProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// JWKProviderModel describes the provider data model.
type JWKProviderModel struct {}

func (p *JWKProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jwk"
	resp.Version = p.version
}

func (p *JWKProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
		},
	}
}

func (p *JWKProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JWKProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// // Example client configuration for data sources and resources
	// client := http.DefaultClient
	// resp.DataSourceData = client
	// resp.ResourceData = client
}

func (p *JWKProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExtractResource,
	}
}

func (p *JWKProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExtractDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JWKProvider{
			version: version,
		}
	}
}
