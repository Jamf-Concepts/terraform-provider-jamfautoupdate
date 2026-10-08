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
	TitleName                           types.String `tfsdk:"title_name"`
	TitleDisplayName                    types.String `tfsdk:"title_display_name"`
	TitleDescription                    types.String `tfsdk:"title_description"`
	TitleVersion                        types.String `tfsdk:"title_version"`
	MinimumOS                           types.String `tfsdk:"minimum_os"`
	MaximumOS                           types.String `tfsdk:"maximum_os"`
	IconBase64                          types.String `tfsdk:"icon_base64"`
	UninstallIconBase64                 types.String `tfsdk:"uninstall_icon_base64"`
	ExtensionAttribute                  types.String `tfsdk:"extension_attribute"`
	ContentFilterProfile                types.String `tfsdk:"content_filter_profile"`
	KernelExtensionProfile              types.String `tfsdk:"kernel_extension_profile"`
	ManagedLoginItemsProfile            types.String `tfsdk:"managed_login_items_profile"`
	NotificationsProfile                types.String `tfsdk:"notifications_profile"`
	PPPCPProfile                        types.String `tfsdk:"pppcp_profile"`
	ScreenRecordingProfile              types.String `tfsdk:"screen_recording_profile"`
	SystemExtensionProfile              types.String `tfsdk:"system_extension_profile"`
	AccessibilityProfile                types.String `tfsdk:"accessibility_profile"`
	AddressBookProfile                  types.String `tfsdk:"address_book_profile"`
	AppleEventsProfile                  types.String `tfsdk:"apple_events_profile"`
	BluetoothAlwaysProfile              types.String `tfsdk:"bluetooth_always_profile"`
	CalendarProfile                     types.String `tfsdk:"calendar_profile"`
	CameraProfile                       types.String `tfsdk:"camera_profile"`
	FileProviderExtensionProfile        types.String `tfsdk:"file_provider_extension_profile"`
	FileProviderPresenceProfile         types.String `tfsdk:"file_provider_presence_profile"`
	ListenEventProfile                  types.String `tfsdk:"listen_event_profile"`
	MediaLibraryProfile                 types.String `tfsdk:"media_library_profile"`
	MicrophoneProfile                   types.String `tfsdk:"microphone_profile"`
	PhotosProfile                       types.String `tfsdk:"photos_profile"`
	PostEventProfile                    types.String `tfsdk:"post_event_profile"`
	RemindersProfile                    types.String `tfsdk:"reminders_profile"`
	SpeechRecognitionProfile            types.String `tfsdk:"speech_recognition_profile"`
	SystemPolicyAllFilesProfile         types.String `tfsdk:"system_policy_all_files_profile"`
	SystemPolicyAppBundlesProfile       types.String `tfsdk:"system_policy_app_bundles_profile"`
	SystemPolicyAppDataProfile          types.String `tfsdk:"system_policy_app_data_profile"`
	SystemPolicyDesktopFolderProfile    types.String `tfsdk:"system_policy_desktop_folder_profile"`
	SystemPolicyDocumentsFolderProfile  types.String `tfsdk:"system_policy_documents_folder_profile"`
	SystemPolicyDownloadsFolderProfile  types.String `tfsdk:"system_policy_downloads_folder_profile"`
	SystemPolicyNetworkVolumesProfile   types.String `tfsdk:"system_policy_network_volumes_profile"`
	SystemPolicyRemovableVolumesProfile types.String `tfsdk:"system_policy_removable_volumes_profile"`
	SystemPolicySysAdminFilesProfile    types.String `tfsdk:"system_policy_sys_admin_files_profile"`
	AppBundleID                         types.String `tfsdk:"app_bundle_id"`
}
