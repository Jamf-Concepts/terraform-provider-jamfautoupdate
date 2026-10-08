// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// Logger is an interface for logging HTTP requests and responses.
type Logger interface {
	LogRequest(ctx context.Context, method, url string, body []byte)
	LogResponse(ctx context.Context, statusCode int, headers http.Header, body []byte)
	LogAuth(ctx context.Context, message string, fields map[string]any)
}

// Title represents a title retrieved from the Jamf Auto Update API.
type Title struct {
	TitleName                           *string         `json:"title_name"`
	TitleDisplayName                    *string         `json:"title_display_name"`
	TitleDescription                    *string         `json:"title_description"`
	TitleVersion                        *string         `json:"title_version"`
	MinimumOS                           *string         `json:"minimum_os"`
	MaximumOS                           *string         `json:"maximum_os"`
	IconHiRes                           *string         `json:"icon_hires"`
	ExtensionAttribute                  *string         `json:"extension_attribute"`
	ContentFilterProfile                *string         `json:"content_filter_profile"`
	KernelExtensionProfile              *string         `json:"kernel_extension_profile"`
	ManagedLoginItemsProfile            *string         `json:"managed_login_items_profile"`
	NotificationsProfile                *string         `json:"notifications_profile"`
	PPPCPProfile                        *string         `json:"pppcp_profile"`
	ScreenRecordingProfile              *string         `json:"screen_recording_profile"`
	SystemExtensionProfile              *string         `json:"system_extension_profile"`
	AccessibilityProfile                *string         `json:"accessibility_profile"`
	AddressBookProfile                  *string         `json:"address_book_profile"`
	AppleEventsProfile                  *string         `json:"apple_events_profile"`
	BluetoothAlwaysProfile              *string         `json:"bluetooth_always_profile"`
	CalendarProfile                     *string         `json:"calendar_profile"`
	CameraProfile                       *string         `json:"camera_profile"`
	FileProviderExtensionProfile        *string         `json:"file_provider_extension_profile"`
	FileProviderPresenceProfile         *string         `json:"file_provider_presence_profile"`
	ListenEventProfile                  *string         `json:"listen_event_profile"`
	MediaLibraryProfile                 *string         `json:"media_library_profile"`
	MicrophoneProfile                   *string         `json:"microphone_profile"`
	PhotosProfile                       *string         `json:"photos_profile"`
	PostEventProfile                    *string         `json:"post_event_profile"`
	RemindersProfile                    *string         `json:"reminders_profile"`
	SpeechRecognitionProfile            *string         `json:"speech_recognition_profile"`
	SystemPolicyAllFilesProfile         *string         `json:"system_policy_all_files_profile"`
	SystemPolicyAppBundlesProfile       *string         `json:"system_policy_app_bundles_profile"`
	SystemPolicyAppDataProfile          *string         `json:"system_policy_app_data_profile"`
	SystemPolicyDesktopFolderProfile    *string         `json:"system_policy_desktop_folder_profile"`
	SystemPolicyDocumentsFolderProfile  *string         `json:"system_policy_documents_folder_profile"`
	SystemPolicyDownloadsFolderProfile  *string         `json:"system_policy_downloads_folder_profile"`
	SystemPolicyNetworkVolumesProfile   *string         `json:"system_policy_network_volumes_profile"`
	SystemPolicyRemovableVolumesProfile *string         `json:"system_policy_removable_volumes_profile"`
	SystemPolicySysAdminFilesProfile    *string         `json:"system_policy_sys_admin_files_profile"`
	PatchDefinition                     PatchDefinition `json:"patch_definition"`
}

// PatchDefinition represents the patch definition of a title.
type PatchDefinition struct {
	Requirements []Requirement `json:"requirements"`
}

// Requirement represents a requirement in the patch definition.
type Requirement struct {
	Name  *string `json:"name"`
	Value *string `json:"value"`
}

// TitlesNotFoundError is returned when one or more requested titles are not found.
type TitlesNotFoundError struct {
	MissingTitles []string
}

// Error returns a formatted string listing the titles that were not found.
func (e *TitlesNotFoundError) Error() string {
	return fmt.Sprintf("The following titles were not found: %s", strings.Join(e.MissingTitles, ", "))
}
