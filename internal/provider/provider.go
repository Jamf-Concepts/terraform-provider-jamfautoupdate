// Copyright 2025 Jamf Software LLC.

package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/client"
	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/resources/titles"
)

// Constants for environment variable names.
const (
	envDefinitionsURL  = "JAMF_AUTO_UPDATE_DEFINITIONS_URL"
	envDefinitionsFile = "JAMF_AUTO_UPDATE_DEFINITIONS_FILE"
)

// Ensure JamfAutoUpdateProvider satisfies various provider interfaces.
var _ provider.Provider = &JamfAutoUpdateProvider{}

// JamfAutoUpdateProvider defines the provider implementation.
type JamfAutoUpdateProvider struct {
	client  *client.Client
	version string
}

// JamfAutoUpdateProvider describes the provider data model.
type JamfAutoUpdateProviderModel struct {
	DefinitionsURL  types.String `tfsdk:"definitions_url"`
	DefinitionsFile types.String `tfsdk:"definitions_file"`
}

func (p *JamfAutoUpdateProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfautoupdate"
	resp.Version = p.version
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

func (p *JamfAutoUpdateProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JamfAutoUpdateProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	definitionsURL := data.DefinitionsURL.ValueString()
	if definitionsURL == "" {
		definitionsURL = getenv(envDefinitionsURL)
	}
	definitionsFile := data.DefinitionsFile.ValueString()
	if definitionsFile == "" {
		definitionsFile = getenv(envDefinitionsFile)
	}

	urlSet := definitionsURL != ""
	fileSet := definitionsFile != ""

	if urlSet == fileSet {
		resp.Diagnostics.AddError(
			"Invalid provider configuration",
			"Exactly one of definitions_url or definitions_file must be set.",
		)
		return
	}

	var clientObj *client.Client
	if urlSet {
		clientObj = client.NewClient(definitionsURL, "")
	} else {
		clientObj = client.NewClient("", definitionsFile)
	}

	clientObj.SetLogger(NewTerraformLogger())

	p.client = clientObj
	resp.DataSourceData = clientObj
	resp.ResourceData = clientObj
}

func (p *JamfAutoUpdateProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *JamfAutoUpdateProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		titles.NewTitlesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JamfAutoUpdateProvider{
			version: version,
		}
	}
}

// getenv is a helper to get an environment variable, returns empty string if not set.
func getenv(key string) string {
	v, _ := os.LookupEnv(key)
	return v
}
