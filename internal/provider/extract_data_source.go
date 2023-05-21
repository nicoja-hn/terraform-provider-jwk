// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	// "github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ExtractDataSource{}

func NewExtractDataSource() datasource.DataSource {
	return &ExtractDataSource{}
}

// ExtractDataSource defines the data source implementation.
type ExtractDataSource struct {
	jwk types.String
}

// ExtractDataSourceModel describes the data source data model.
type ExtractDataSourceModel struct {
	PublicCertificate types.String `tfsdk:"public_certificate"`
	SigningAlgorithm  types.String `tfsdk:"signing_algorithm"`
	Jwk               types.String `tfsdk:"jwk"`
	Id                types.String `tfsdk:"id"`
}

func (d *ExtractDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extract"
}

func (d *ExtractDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Extract data source",

		Attributes: map[string]schema.Attribute{
			"public_certificate": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
			},
			"signing_algorithm": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Required:            true,
				Validators: []validator.String{
					// These are example validators from terraform-plugin-framework-validators
					// stringvalidator.LengthBetween(0, 6),
					// stringvalidator.OneOfCaseInsensitive("rs256"),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(r|R|e|E)(s|S)(256|384|512)$`),
						"must contain only lowercase alphanumeric characters",
					),
				},
			},
			"jwk": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Computed:            true,
				Sensitive:           true,
			},
			"jwk": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Computed:            true,
				Sensitive:           true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Computed:            true,
			},
		},
	}
}

func (d *ExtractDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ExtractDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var content string = data.PublicCertificate.ValueString()
	output, err := readKey(content)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	d.jwk = types.StringValue(string(output))
	data.Jwk = d.jwk
	data.Id = types.StringValue(fmt.Sprint(hash(string(output))))

	// // Write logs using the tflog package
	// // Documentation: https://terraform.io/plugin/log
	// tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
