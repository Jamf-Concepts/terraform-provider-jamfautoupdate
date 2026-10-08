// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
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

	titlesAttr, ok := attrs["titles"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("expected titles to be a ListNestedAttribute, got %T", attrs["titles"])
	}

	nestedAttrs := titlesAttr.NestedObject.Attributes
	modelType := reflect.TypeFor[TitleModel]()
	for i := range modelType.NumField() {
		tag := modelType.Field(i).Tag.Get("tfsdk")
		if _, ok := nestedAttrs[tag]; !ok {
			t.Errorf("model field %q has no matching schema attribute", tag)
		}
	}
	if len(nestedAttrs) != modelType.NumField() {
		t.Errorf("schema has %d nested attributes, model has %d fields", len(nestedAttrs), modelType.NumField())
	}
}
