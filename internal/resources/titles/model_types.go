// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
