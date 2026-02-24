// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// buildTitleModelsFromResponse converts a slice of client.Title API responses into TitleModel state values.
func buildTitleModelsFromResponse(titles []client.Title) ([]TitleModel, error) {
	models := make([]TitleModel, 0, len(titles))

	for _, title := range titles {
		bundleID := extractBundleID(title.PatchDefinition.Requirements)

		var uninstallIcon *string
		if title.IconHiRes != nil {
			var err error
			uninstallIcon, err = processUninstallIcon(*title.IconHiRes)
			if err != nil {
				return nil, err
			}
		}

		model := TitleModel{
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
		models = append(models, model)
	}

	return models, nil
}

// extractBundleID finds the Application Bundle ID from a title's patch definition requirements.
func extractBundleID(requirements []client.Requirement) *string {
	for _, req := range requirements {
		if req.Name != nil && *req.Name == "Application Bundle ID" {
			return req.Value
		}
	}
	return nil
}
