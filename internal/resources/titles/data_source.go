// Copyright 2025 Jamf Software LLC.

package titles

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"strings"
	"time"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/image/draw"
)

const (
	OverlayImageBase64 = "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"
	BaseImageSize      = 512
	OverlaySize        = 128
	TargetImageSize    = 512
)

const defaultReadTimeout = 90 * time.Second

var _ datasource.DataSource = &TitlesDataSource{}

func NewTitlesDataSource() datasource.DataSource {
	return &TitlesDataSource{}
}

// TitlesDataSource defines the data source implementation.
type TitlesDataSource struct {
	client *client.Client
}

// TitlesDataSourceModel describes the data source data model.
type TitlesDataSourceModel struct {
	TitleNames types.List     `tfsdk:"title_names"`
	Timeouts   timeouts.Value `tfsdk:"timeouts"`
	Titles     []TitleModel   `tfsdk:"titles"`
}

// TitleModel describes the structure of a title in the data source.
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

func (d *TitlesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_titles"
}

func (d *TitlesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about Jamf Auto Update titles. Available titles are shown in the [Jamf Auto Update Catalog Browser](https://definitions-admin.datajar.mobi/v2/titles)",
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

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
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
			uninstallIcon, err = processUninstallIcon(*title.IconHiRes)
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
			ExtensionAttribute:       types.StringPointerValue(title.ExtensionAttribute),
			ContentFilterProfile:     types.StringPointerValue(title.ContentFilterProfile),
			KernelExtensionProfile:   types.StringPointerValue(title.KernelExtensionProfile),
			ManagedLoginItemsProfile: types.StringPointerValue(title.ManagedLoginItemsProfile),
			NotificationsProfile:     types.StringPointerValue(title.NotificationsProfile),
			PPPCPProfile:             types.StringPointerValue(title.PPPCPProfile),
			ScreenRecordingProfile:   types.StringPointerValue(title.ScreenRecordingProfile),
			SystemExtensionProfile:   types.StringPointerValue(title.SystemExtensionProfile),
			AppBundleID:              types.StringPointerValue(bundleID),
		}
		data.Titles = append(data.Titles, titleModel)
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d titles from Jamf Auto Update API", len(data.Titles)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// processUninstallIcon processes a base64 encoded image by resizing it and adding an overlay.
// It takes a base64 encoded string of the original image and returns a base64 encoded string
// of the processed image. The processing includes resizing to the standard size and adding
// an uninstall overlay to the bottom right corner.
// Returns an error if any step of the image processing fails.
func processUninstallIcon(baseImageB64 string) (*string, error) {
	baseImageBytes, err := base64.StdEncoding.DecodeString(baseImageB64)
	if err != nil {
		return nil, fmt.Errorf("error decoding base image: %w", err)
	}

	baseImg, _, err := image.Decode(bytes.NewReader(baseImageBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding base image bytes: %w", err)
	}

	resizedBase := image.NewRGBA(image.Rect(0, 0, BaseImageSize, BaseImageSize))
	draw.CatmullRom.Scale(resizedBase, resizedBase.Bounds(), baseImg, baseImg.Bounds(), draw.Over, nil)
	baseImg = resizedBase

	bounds := baseImg.Bounds()
	rgba := image.NewRGBA(bounds)

	draw.Draw(rgba, bounds, baseImg, image.Point{}, draw.Src)

	overlayBytes, err := base64.StdEncoding.DecodeString(OverlayImageBase64)
	if err != nil {
		return nil, fmt.Errorf("error decoding overlay image: %w", err)
	}

	overlayImg, _, err := image.Decode(bytes.NewReader(overlayBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding overlay image bytes: %w", err)
	}

	resizedOverlay := image.NewRGBA(image.Rect(0, 0, OverlaySize, OverlaySize))
	draw.CatmullRom.Scale(resizedOverlay, resizedOverlay.Bounds(), overlayImg, overlayImg.Bounds(), draw.Over, nil)
	overlayImg = resizedOverlay
	overlayBounds := overlayImg.Bounds()

	x := bounds.Max.X - overlayBounds.Dx()
	y := bounds.Max.Y - overlayBounds.Dy()

	offset := image.Point{X: x, Y: y}

	draw.Draw(rgba, image.Rectangle{
		Min: offset,
		Max: offset.Add(overlayBounds.Size()),
	}, overlayImg, image.Point{}, draw.Over)

	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		return nil, fmt.Errorf("error encoding processed image: %w", err)
	}

	result := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &result, nil
}
