// Copyright 2025 Jamf Software LLC.

package titles

import (
	"context"
	"fmt"
	"io"
	"strings"

	"encoding/json"
	"os"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TitlesDataSource struct {
	client          *client.Client
	imageProcessor  *Processor
	definitionsURL  string
	definitionsFile string
}

type TitlesDataSourceModel struct {
	TitleNames types.List   `tfsdk:"title_names"`
	Titles     []TitleModel `tfsdk:"titles"`
}

type TitleModel struct {
	TitleName                types.String `tfsdk:"title_name"`
	TitleDisplayName         types.String `tfsdk:"title_display_name"`
	TitleDescription         types.String `tfsdk:"title_description"`
	TitleVersion             types.String `tfsdk:"title_version"`
	MinimumOS                types.String `tfsdk:"minimum_os"`
	MaximumOS                types.String `tfsdk:"maximum_os"`
	IconBase64               types.String `tfsdk:"icon_base64"`
	UninstallIconBase64      types.String `tfsdk:"uninstall_icon_base64"`
	ExtensionAttribute       types.String `tfsdk:"extension_attribute"`
	ContentFilterProfile     types.String `tfsdk:"content_filter_profile"`
	KernelExtensionProfile   types.String `tfsdk:"kernel_extension_profile"`
	ManagedLoginItemsProfile types.String `tfsdk:"managed_login_items_profile"`
	NotificationsProfile     types.String `tfsdk:"notifications_profile"`
	PPPCPProfile             types.String `tfsdk:"pppcp_profile"`
	ScreenRecordingProfile   types.String `tfsdk:"screen_recording_profile"`
	SystemExtensionProfile   types.String `tfsdk:"system_extension_profile"`
	AppBundleID              types.String `tfsdk:"app_bundle_id"`
}

// NewTitlesDataSource creates a new instance of the titles data source with initialized image processor.
// It implements the datasource.DataSource interface for Terraform provider use.
func NewTitlesDataSource() datasource.DataSource {
	return &TitlesDataSource{
		imageProcessor: NewProcessor(),
	}
}

// Metadata returns the full name of the data source by combining the provider name with "_titles".
// This name is used to identify the data source in Terraform configurations.
func (d *TitlesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_titles"
}

// Schema defines the schema for the titles data source.
// It specifies all available attributes and their types, including the optional title_names filter
// and the computed titles list containing detailed information about each title.
func (d *TitlesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about Jamf Auto Update titles.",
		Attributes: map[string]schema.Attribute{
			"title_names": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Optional list of specific title names to retrieve. If not provided, returns all titles.",
			},
			"titles": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of titles and their details",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"title_name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the title",
						},
						"title_display_name": schema.StringAttribute{
							Computed:    true,
							Description: "The display name of the title",
						},
						"title_description": schema.StringAttribute{
							Computed:    true,
							Description: "The description of the title",
						},
						"title_version": schema.StringAttribute{
							Computed:    true,
							Description: "The version of the title",
						},
						"minimum_os": schema.StringAttribute{
							Computed:    true,
							Description: "Minimum OS version required",
						},
						"maximum_os": schema.StringAttribute{
							Computed:    true,
							Description: "Maximum OS version supported",
						},
						"icon_base64": schema.StringAttribute{
							Computed:    true,
							Description: "The icon in base64 format",
						},
						"uninstall_icon_base64": schema.StringAttribute{
							Computed:    true,
							Description: "The uninstall icon in base64 format",
						},
						"extension_attribute": schema.StringAttribute{
							Computed:    true,
							Description: "Extension attribute data",
						},
						"content_filter_profile": schema.StringAttribute{
							Computed:    true,
							Description: "Content filter profile data",
						},
						"kernel_extension_profile": schema.StringAttribute{
							Computed:    true,
							Description: "Kernel extension profile data",
						},
						"managed_login_items_profile": schema.StringAttribute{
							Computed:    true,
							Description: "Managed login items profile data",
						},
						"notifications_profile": schema.StringAttribute{
							Computed:    true,
							Description: "Notifications profile data",
						},
						"pppcp_profile": schema.StringAttribute{
							Computed:    true,
							Description: "PPPCP profile data",
						},
						"screen_recording_profile": schema.StringAttribute{
							Computed:    true,
							Description: "Screen recording profile data",
						},
						"system_extension_profile": schema.StringAttribute{
							Computed:    true,
							Description: "System extension profile data",
						},
						"app_bundle_id": schema.StringAttribute{
							Computed:    true,
							Description: "The application bundle identifier",
						},
					},
				},
			},
		},
	}
}

// Configure sets up the data source with the provider-level client.
// It is called by the Terraform provider framework to initialize the data source
// with provider-level configuration data.
func (d *TitlesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	// ProviderData is a map[string]string
	providerData, ok := req.ProviderData.(map[string]string)
	if !ok {
		resp.Diagnostics.AddError("Provider data error", "Could not parse provider data.")
		return
	}
	d.definitionsURL = providerData["definitions_url"]
	d.definitionsFile = providerData["definitions_file"]
	if d.definitionsFile != "" {
		d.client = nil // not used
	} else {
		d.client = client.NewClient(d.definitionsURL)
	}
}

// Read fetches title data from the Jamf Auto Update API and processes it into the Terraform state.
// It handles both filtered and unfiltered requests, processes icons for uninstall operations,
// and extracts bundle IDs from patch definitions. Any errors during processing are reported
// through the diagnostics system.
func (d *TitlesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TitlesDataSourceModel

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var titleNames []string
	if !state.TitleNames.IsNull() {
		diags = state.TitleNames.ElementsAs(ctx, &titleNames, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !state.TitleNames.IsUnknown() && len(titleNames) == 0 {
		state.Titles = []TitleModel{}
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	var titles []client.Title
	var err error
	if d.definitionsFile != "" {
		// Read from local file
		file, err := os.Open(d.definitionsFile)
		if err != nil {
			resp.Diagnostics.AddError("Unable to open definitions file", err.Error())
			return
		}
		defer func() {
			if cerr := file.Close(); cerr != nil {
				resp.Diagnostics.AddError("Error closing definitions file", cerr.Error())
			}
		}()
		bytes, err := io.ReadAll(file)
		if err != nil {
			resp.Diagnostics.AddError("Unable to read definitions file", err.Error())
			return
		}
		var allTitles []client.Title
		if err := json.Unmarshal(bytes, &allTitles); err != nil {
			resp.Diagnostics.AddError("Unable to parse definitions file", err.Error())
			return
		}
		// Filter if titleNames is set
		if len(titleNames) > 0 {
			for _, t := range allTitles {
				for _, name := range titleNames {
					if t.TitleName != nil && *t.TitleName == name {
						titles = append(titles, t)
						break
					}
				}
			}
		} else {
			titles = allTitles
		}
	} else {
		// Use API client as before
		titles, err = d.client.GetTitles(ctx, titleNames...)
		if err != nil {
			if titlesErr, ok := err.(*client.TitlesNotFoundError); ok {
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
	}

	for _, title := range titles {
		var bundleID *string
		if len(title.PatchDefinition.Requirements) > 0 {
			for _, req := range title.PatchDefinition.Requirements {
				if *req.Name == "Application Bundle ID" {
					bundleID = req.Value
					break
				}
			}
		}

		var uninstallIcon *string
		if title.IconHiRes != nil {
			var err error
			uninstallIcon, err = d.imageProcessor.ProcessUninstallIcon(*title.IconHiRes)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error processing icon",
					fmt.Sprintf("Could not process icon: %s", err),
				)
				return
			}
		}

		titleModel := TitleModel{
			TitleName:                types.StringPointerValue(title.TitleName),
			TitleDisplayName:         types.StringPointerValue(title.TitleDisplayName),
			TitleDescription:         types.StringPointerValue(title.TitleDescription),
			TitleVersion:             types.StringPointerValue(title.TitleVersion),
			MinimumOS:                types.StringPointerValue(title.MinimumOS),
			MaximumOS:                types.StringPointerValue(title.MaximumOS),
			IconBase64:               types.StringPointerValue(title.IconHiRes),
			UninstallIconBase64:      types.StringPointerValue(uninstallIcon),
			ExtensionAttribute:       types.StringPointerValue(title.GetExtensionAttribute()),
			ContentFilterProfile:     types.StringPointerValue(title.GetContentFilterProfile()),
			KernelExtensionProfile:   types.StringPointerValue(title.GetKernelExtensionProfile()),
			ManagedLoginItemsProfile: types.StringPointerValue(title.GetManagedLoginItemsProfile()),
			NotificationsProfile:     types.StringPointerValue(title.GetNotificationsProfile()),
			PPPCPProfile:             types.StringPointerValue(title.GetPPPCPProfile()),
			ScreenRecordingProfile:   types.StringPointerValue(title.GetScreenRecordingProfile()),
			SystemExtensionProfile:   types.StringPointerValue(title.GetSystemExtensionProfile()),
			AppBundleID:              types.StringPointerValue(bundleID),
		}
		state.Titles = append(state.Titles, titleModel)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
