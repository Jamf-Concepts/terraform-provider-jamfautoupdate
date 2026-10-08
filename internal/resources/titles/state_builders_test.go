// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/terraform-provider-jamfautoupdate/internal/client"
)

func TestBuildTitleModelsFromResponse_SingleTitle(t *testing.T) {
	titles := []client.Title{
		{
			TitleName:        new("GoogleChrome"),
			TitleDisplayName: new("Google Chrome"),
			TitleVersion:     new("120.0"),
			MinimumOS:        new("12.0"),
			MaximumOS:        new("15.0"),
			PatchDefinition: client.PatchDefinition{
				Requirements: []client.Requirement{
					{Name: new("Application Bundle ID"), Value: new("com.google.Chrome")},
				},
			},
		},
	}

	models, err := buildTitleModelsFromResponse(titles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}

	m := models[0]
	if m.TitleName.ValueString() != "GoogleChrome" {
		t.Errorf("expected GoogleChrome, got %s", m.TitleName.ValueString())
	}
	if m.TitleDisplayName.ValueString() != "Google Chrome" {
		t.Errorf("expected Google Chrome, got %s", m.TitleDisplayName.ValueString())
	}
	if m.AppBundleID.ValueString() != "com.google.Chrome" {
		t.Errorf("expected com.google.Chrome, got %s", m.AppBundleID.ValueString())
	}
}

func TestBuildTitleModelsFromResponse_NilFields(t *testing.T) {
	titles := []client.Title{
		{
			TitleName:       new("TestApp"),
			PatchDefinition: client.PatchDefinition{},
		},
	}

	models, err := buildTitleModelsFromResponse(titles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}

	m := models[0]
	if !m.TitleDisplayName.IsNull() {
		t.Error("expected null TitleDisplayName")
	}
	if !m.IconBase64.IsNull() {
		t.Error("expected null IconBase64")
	}
	if !m.AppBundleID.IsNull() {
		t.Error("expected null AppBundleID")
	}
}

func TestBuildTitleModelsFromResponse_EmptySlice(t *testing.T) {
	models, err := buildTitleModelsFromResponse([]client.Title{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != 0 {
		t.Errorf("expected 0 models, got %d", len(models))
	}
}

func TestBuildTitleModelsFromResponse_BundleIDExtraction(t *testing.T) {
	titles := []client.Title{
		{
			TitleName: new("TestApp"),
			PatchDefinition: client.PatchDefinition{
				Requirements: []client.Requirement{
					{Name: new("OS Version"), Value: new("12.0")},
					{Name: new("Application Bundle ID"), Value: new("com.test.app")},
					{Name: new("Processor"), Value: new("arm64")},
				},
			},
		},
	}

	models, err := buildTitleModelsFromResponse(titles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if models[0].AppBundleID.ValueString() != "com.test.app" {
		t.Errorf("expected com.test.app, got %s", models[0].AppBundleID.ValueString())
	}
}

func TestBuildTitleModelsFromResponse_NoBundleID(t *testing.T) {
	titles := []client.Title{
		{
			TitleName: new("TestApp"),
			PatchDefinition: client.PatchDefinition{
				Requirements: []client.Requirement{
					{Name: new("OS Version"), Value: new("12.0")},
				},
			},
		},
	}

	models, err := buildTitleModelsFromResponse(titles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !models[0].AppBundleID.IsNull() {
		t.Error("expected null AppBundleID when no bundle ID requirement exists")
	}
}

func TestBuildTitleModelsFromResponse_NoRequirements(t *testing.T) {
	titles := []client.Title{
		{
			TitleName:       new("TestApp"),
			PatchDefinition: client.PatchDefinition{},
		},
	}

	models, err := buildTitleModelsFromResponse(titles)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !models[0].AppBundleID.IsNull() {
		t.Error("expected null AppBundleID when no requirements exist")
	}
}

func TestExtractBundleID_Found(t *testing.T) {
	reqs := []client.Requirement{
		{Name: new("Application Bundle ID"), Value: new("com.example.app")},
	}
	result := extractBundleID(reqs)
	if result == nil || *result != "com.example.app" {
		t.Errorf("expected com.example.app, got %v", result)
	}
}

func TestExtractBundleID_NotFound(t *testing.T) {
	reqs := []client.Requirement{
		{Name: new("OS Version"), Value: new("12.0")},
	}
	result := extractBundleID(reqs)
	if result != nil {
		t.Errorf("expected nil, got %s", *result)
	}
}

func TestExtractBundleID_EmptySlice(t *testing.T) {
	result := extractBundleID([]client.Requirement{})
	if result != nil {
		t.Errorf("expected nil, got %s", *result)
	}
}

func TestExtractBundleID_NilSlice(t *testing.T) {
	result := extractBundleID(nil)
	if result != nil {
		t.Errorf("expected nil, got %s", *result)
	}
}

func TestBuildTitleModelsFromResponse_AllProfilesMapped(t *testing.T) {
	modelType := reflect.TypeFor[TitleModel]()
	payload := map[string]string{"title_name": "TestApp"}
	var profileTags []string
	for field := range modelType.Fields() {
		tag := field.Tag.Get("tfsdk")
		if strings.HasSuffix(tag, "_profile") {
			payload[tag] = "value-" + tag
			profileTags = append(profileTags, tag)
		}
	}
	if len(profileTags) != 31 {
		t.Fatalf("expected 31 profile fields on the model, got %d", len(profileTags))
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	var title client.Title
	if err := json.Unmarshal(raw, &title); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	models, err := buildTitleModelsFromResponse([]client.Title{title})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	modelValue := reflect.ValueOf(models[0])
	for i := range modelType.NumField() {
		tag := modelType.Field(i).Tag.Get("tfsdk")
		if !strings.HasSuffix(tag, "_profile") {
			continue
		}
		got := modelValue.Field(i).Interface().(types.String)
		if got.ValueString() != "value-"+tag {
			t.Errorf("%s: expected %q, got %q", tag, "value-"+tag, got.ValueString())
		}
	}
}

func TestBuildTitleModelsFromResponse_NullProfiles(t *testing.T) {
	var title client.Title
	if err := json.Unmarshal([]byte(`{"title_name":"TestApp","accessibility_profile":null,"camera_profile":null}`), &title); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	models, err := buildTitleModelsFromResponse([]client.Title{title})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !models[0].AccessibilityProfile.IsNull() || !models[0].CameraProfile.IsNull() {
		t.Error("expected null profiles for JSON null values")
	}
}
