// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestProviderMetadata(t *testing.T) {
	p := &JamfAutoUpdateProvider{version: "1.0.0"}
	req := provider.MetadataRequest{}
	resp := &provider.MetadataResponse{}

	p.Metadata(context.Background(), req, resp)

	if resp.TypeName != "jamfautoupdate" {
		t.Errorf("expected type name jamfautoupdate, got %s", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", resp.Version)
	}
}

func TestProviderSchema(t *testing.T) {
	p := &JamfAutoUpdateProvider{}
	req := provider.SchemaRequest{}
	resp := &provider.SchemaResponse{}

	p.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema errors: %v", resp.Diagnostics)
	}

	attrs := resp.Schema.Attributes
	if attrs == nil {
		t.Fatal("expected non-nil schema attributes")
	}

	if _, ok := attrs["definitions_url"]; !ok {
		t.Error("expected definitions_url attribute in schema")
	}
	if _, ok := attrs["definitions_file"]; !ok {
		t.Error("expected definitions_file attribute in schema")
	}
}

func TestProviderDataSources(t *testing.T) {
	p := &JamfAutoUpdateProvider{}
	dataSources := p.DataSources(context.Background())
	if len(dataSources) != 1 {
		t.Errorf("expected 1 data source, got %d", len(dataSources))
	}
}

func TestProviderResources(t *testing.T) {
	p := &JamfAutoUpdateProvider{}
	resources := p.Resources(context.Background())
	if len(resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(resources))
	}
}
