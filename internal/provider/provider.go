// Copyright 2025 Jamf Software LLC.

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/resources/titles"
)

type JamfAutoUpdateProvider struct {
	definitionsURL  string
	definitionsFile string
}

type JamfAutoUpdateProviderModel struct {
	DefinitionsURL  types.String `tfsdk:"definitions_url"`
	DefinitionsFile types.String `tfsdk:"definitions_file"`
}

// New creates a new instance of the Jamf Auto Update provider.
func New() provider.Provider {
	return &JamfAutoUpdateProvider{}
}

// Metadata sets the metadata for the Jamf Auto Update provider.
func (p *JamfAutoUpdateProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfautoupdate"
}

// Schema defines the schema for the Jamf Auto Update provider.
func (p *JamfAutoUpdateProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"definitions_url": schema.StringAttribute{
				Optional:    true,
				Description: "The baseURL of the Definitions API. Mutually exclusive with definitions_file.",
			},
			"definitions_file": schema.StringAttribute{
				Optional:    true,
				Description: "Path to a local JSON file containing definitions. Mutually exclusive with definitions_url.",
			},
		},
	}
}

// Resources returns an empty slice as there are no resources defined for this provider.
func (p *JamfAutoUpdateProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns the data sources available in the Jamf Auto Update provider.
func (p *JamfAutoUpdateProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		titles.NewTitlesDataSource,
	}
}

// Configure initializes the provider with the given configuration.
func (p *JamfAutoUpdateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JamfAutoUpdateProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	urlSet := !config.DefinitionsURL.IsNull() && config.DefinitionsURL.ValueString() != ""
	fileSet := !config.DefinitionsFile.IsNull() && config.DefinitionsFile.ValueString() != ""

	if urlSet == fileSet {
		resp.Diagnostics.AddError(
			"Invalid provider configuration",
			"Exactly one of definitions_url or definitions_file must be set.",
		)
		return
	}

	p.definitionsURL = config.DefinitionsURL.ValueString()
	p.definitionsFile = config.DefinitionsFile.ValueString()

	// Pass both to the data source
	resp.DataSourceData = map[string]string{
		"definitions_url":  p.definitionsURL,
		"definitions_file": p.definitionsFile,
	}
}
