// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"os"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "titles-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}
	return f.Name()
}

func TestGetTitlesFromFile_AllTitles(t *testing.T) {
	path := writeTempFile(t, testMultipleTitlesJSON)
	c := NewClient("", path)
	titles, err := c.GetTitles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 2 {
		t.Fatalf("expected 2 titles, got %d", len(titles))
	}
	if *titles[0].TitleName != "GoogleChrome" {
		t.Errorf("expected GoogleChrome, got %s", *titles[0].TitleName)
	}
	if *titles[1].TitleName != "Firefox" {
		t.Errorf("expected Firefox, got %s", *titles[1].TitleName)
	}
}

func TestGetTitlesFromFile_SpecificTitles(t *testing.T) {
	path := writeTempFile(t, testMultipleTitlesJSON)
	c := NewClient("", path)
	titles, err := c.GetTitles(context.Background(), "Firefox")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 1 {
		t.Fatalf("expected 1 title, got %d", len(titles))
	}
	if *titles[0].TitleName != "Firefox" {
		t.Errorf("expected Firefox, got %s", *titles[0].TitleName)
	}
}

func TestGetTitlesFromFile_MissingTitle(t *testing.T) {
	path := writeTempFile(t, testMultipleTitlesJSON)
	c := NewClient("", path)
	_, err := c.GetTitles(context.Background(), "NonExistent")
	if err == nil {
		t.Fatal("expected error for missing title")
	}
	notFoundErr, ok := err.(*TitlesNotFoundError)
	if !ok {
		t.Fatalf("expected TitlesNotFoundError, got %T: %v", err, err)
	}
	if len(notFoundErr.MissingTitles) != 1 || notFoundErr.MissingTitles[0] != "NonExistent" {
		t.Errorf("expected [NonExistent], got %v", notFoundErr.MissingTitles)
	}
}

func TestGetTitlesFromFile_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "[]")
	c := NewClient("", path)
	titles, err := c.GetTitles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 0 {
		t.Errorf("expected 0 titles, got %d", len(titles))
	}
}

func TestGetTitlesFromFile_InvalidJSON(t *testing.T) {
	path := writeTempFile(t, "not json")
	c := NewClient("", path)
	_, err := c.GetTitles(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGetTitlesFromFile_FileNotFound(t *testing.T) {
	c := NewClient("", "/nonexistent/path/titles.json")
	_, err := c.GetTitles(context.Background())
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestGetTitlesFromFile_NotArray(t *testing.T) {
	path := writeTempFile(t, `{"title_name":"AppA"}`)
	c := NewClient("", path)
	_, err := c.GetTitles(context.Background(), "AppA")
	if err == nil {
		t.Fatal("expected error for non-array JSON")
	}
}

func TestGetTitlesFromFile_NullTitleName(t *testing.T) {
	json := `[{"title_name":null,"title_version":"1.0","patch_definition":{"requirements":[]}},{"title_name":"AppA","title_version":"2.0","patch_definition":{"requirements":[]}}]`
	path := writeTempFile(t, json)
	c := NewClient("", path)
	titles, err := c.GetTitles(context.Background(), "AppA")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 1 {
		t.Fatalf("expected 1 title, got %d", len(titles))
	}
	if *titles[0].TitleName != "AppA" {
		t.Errorf("expected AppA, got %s", *titles[0].TitleName)
	}
}
