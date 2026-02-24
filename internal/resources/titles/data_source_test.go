// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestTitlesDataSource_Metadata(t *testing.T) {
	ds := &TitlesDataSource{}
	req := datasource.MetadataRequest{
		ProviderTypeName: "jamfautoupdate",
	}
	resp := &datasource.MetadataResponse{}

	ds.Metadata(context.Background(), req, resp)

	if resp.TypeName != "jamfautoupdate_titles" {
		t.Errorf("expected jamfautoupdate_titles, got %s", resp.TypeName)
	}
}

func TestTitlesDataSource_Schema(t *testing.T) {
	ds := &TitlesDataSource{}
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	ds.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema errors: %v", resp.Diagnostics)
	}

	attrs := resp.Schema.Attributes
	if attrs == nil {
		t.Fatal("expected non-nil schema attributes")
	}

	expectedAttrs := []string{"timeouts", "title_names", "titles"}
	for _, name := range expectedAttrs {
		if _, ok := attrs[name]; !ok {
			t.Errorf("expected attribute %q in schema", name)
		}
	}

	titlesAttr, ok := attrs["titles"]
	if !ok {
		t.Fatal("missing titles attribute")
	}

	nested, ok := titlesAttr.(interface {
		GetNestedObject() interface{ GetAttributes() map[string]interface{} }
	})
	_ = nested
	_ = ok

	expectedNestedAttrs := []string{
		"title_name", "title_display_name", "title_description", "title_version",
		"minimum_os", "maximum_os", "icon_base64", "uninstall_icon_base64",
		"extension_attribute", "content_filter_profile", "kernel_extension_profile",
		"managed_login_items_profile", "notifications_profile", "pppcp_profile",
		"screen_recording_profile", "system_extension_profile", "app_bundle_id",
	}
	if len(expectedNestedAttrs) != 17 {
		t.Errorf("expected 17 nested attributes, listed %d", len(expectedNestedAttrs))
	}
}
