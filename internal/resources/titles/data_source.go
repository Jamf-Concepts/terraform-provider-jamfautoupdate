// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// defaultReadTimeout is the default timeout duration for reading titles from the API.
const defaultReadTimeout = 90 * time.Second

var _ datasource.DataSource = &TitlesDataSource{}

// NewTitlesDataSource returns a new instance of the titles data source.
func NewTitlesDataSource() datasource.DataSource {
	return &TitlesDataSource{}
}

// TitlesDataSource defines the data source implementation.
type TitlesDataSource struct {
	client *client.Client
}

func (d *TitlesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_titles"
}

func (d *TitlesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about Jamf Auto Update titles. Available titles are shown in the [Jamf Auto Update Catalog Browser](https://support.datajar.co.uk/hc/en-us/articles/4409234438161-Jamf-Auto-Update-Catalog-Browser-User-Guide)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx),
			"title_names": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "List of specific title names to retrieve.",
			},
			"titles": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of titles and their details",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"title_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The name of the title",
						},
						"title_display_name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The display name of the title",
						},
						"title_description": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The description of the title",
						},
						"title_version": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The version of the title",
						},
						"minimum_os": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Minimum OS version required",
						},
						"maximum_os": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Maximum OS version supported",
						},
						"icon_base64": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The icon in base64 format",
						},
						"uninstall_icon_base64": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The uninstall icon in base64 format",
						},
						"extension_attribute": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Extension attribute data",
						},
						"content_filter_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Content filter profile data",
						},
						"kernel_extension_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Kernel extension profile data",
						},
						"managed_login_items_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Managed login items profile data",
						},
						"notifications_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Notifications profile data",
						},
						"pppcp_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "PPPCP profile data",
						},
						"screen_recording_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Screen recording profile data",
						},
						"system_extension_profile": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "System extension profile data",
						},
						"app_bundle_id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The application bundle identifier",
						},
					},
				},
			},
		},
	}
}

func (d *TitlesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = c
}

func (d *TitlesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data TitlesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var titleNames []string
	if !data.TitleNames.IsNull() {
		resp.Diagnostics.Append(data.TitleNames.ElementsAs(ctx, &titleNames, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !data.TitleNames.IsUnknown() && len(titleNames) == 0 {
		data.Titles = []TitleModel{}
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	readTimeout := defaultReadTimeout
	if !data.Timeouts.IsNull() && !data.Timeouts.IsUnknown() {
		configuredTimeout, timeoutDiags := data.Timeouts.Read(ctx, defaultReadTimeout)
		resp.Diagnostics.Append(timeoutDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
		readTimeout = configuredTimeout
	}

	readCtx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	titles, err := d.client.GetTitles(readCtx, titleNames...)
	if err != nil {
		if titlesErr, ok := errors.AsType[*client.TitlesNotFoundError](err); ok {
			resp.Diagnostics.AddError(
				"Requested titles not found",
				fmt.Sprintf("The following titles do not exist: %s",
					strings.Join(titlesErr.MissingTitles, ", ")),
			)
			return
		}
		resp.Diagnostics.AddError(
			"Unable to read Jamf Auto Update titles",
			err.Error(),
		)
		return
	}

	models, err := buildTitleModelsFromResponse(titles)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error processing title data",
			err.Error(),
		)
		return
	}
	data.Titles = models

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d titles from Jamf Auto Update API", len(data.Titles)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
