// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testTitleJSON = `[{"title_name":"GoogleChrome","title_display_name":"Google Chrome","title_version":"1.0","patch_definition":{"requirements":[{"name":"Application Bundle ID","value":"com.google.Chrome"}]}}]`

const testMultipleTitlesJSON = `[{"title_name":"GoogleChrome","title_version":"1.0","patch_definition":{"requirements":[]}},{"title_name":"Firefox","title_version":"2.0","patch_definition":{"requirements":[]}}]`

func TestNewClient(t *testing.T) {
	c := NewClient("https://example.com", "")
	if c.baseURL != "https://example.com" {
		t.Errorf("expected baseURL https://example.com, got %s", c.baseURL)
	}
	if c.definitionsFile != "" {
		t.Errorf("expected empty definitionsFile, got %s", c.definitionsFile)
	}
	if c.httpClient == nil {
		t.Error("expected non-nil httpClient")
	}
}

func TestNewClient_FileMode(t *testing.T) {
	c := NewClient("", "/path/to/file.json")
	if c.baseURL != "" {
		t.Errorf("expected empty baseURL, got %s", c.baseURL)
	}
	if c.definitionsFile != "/path/to/file.json" {
		t.Errorf("expected /path/to/file.json, got %s", c.definitionsFile)
	}
}

func TestSetLogger(t *testing.T) {
	c := NewClient("https://example.com", "")
	if c.logger != nil {
		t.Error("expected nil logger before SetLogger")
	}
	logger := &mockLogger{}
	c.SetLogger(logger)
	if c.logger == nil {
		t.Error("expected non-nil logger after SetLogger")
	}
}

func TestGetTitles_AllTitles_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("expected path /, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testMultipleTitlesJSON))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
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

func TestGetTitles_SpecificTitles_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/GoogleChrome" {
			t.Errorf("expected path /GoogleChrome, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testTitleJSON))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	titles, err := c.GetTitles(context.Background(), "GoogleChrome")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 1 {
		t.Fatalf("expected 1 title, got %d", len(titles))
	}
	if *titles[0].TitleName != "GoogleChrome" {
		t.Errorf("expected GoogleChrome, got %s", *titles[0].TitleName)
	}
}

func TestGetTitles_SpecificTitles_URLPath(t *testing.T) {
	var capturedPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testMultipleTitlesJSON))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	_, err := c.GetTitles(context.Background(), "GoogleChrome", "Firefox")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedPath != "/GoogleChrome,Firefox" {
		t.Errorf("expected path /GoogleChrome,Firefox, got %s", capturedPath)
	}
}

func TestGetTitles_MissingTitle(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(testTitleJSON))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	_, err := c.GetTitles(context.Background(), "GoogleChrome", "NonExistent")
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

func TestGetTitles_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	titles, err := c.GetTitles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 0 {
		t.Errorf("expected 0 titles, got %d", len(titles))
	}
}

func TestGetTitles_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	_, err := c.GetTitles(context.Background())
	if err == nil {
		t.Fatal("expected error for HTTP 500")
	}
}

func TestGetTitles_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	_, err := c.GetTitles(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestGetTitles_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testTitleJSON))
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := NewClient(server.URL, "")
	_, err := c.GetTitles(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestGetTitles_NilPointerFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"title_name":"TestApp","patch_definition":{"requirements":[]}}]`))
	}))
	defer server.Close()

	c := NewClient(server.URL, "")
	titles, err := c.GetTitles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(titles) != 1 {
		t.Fatalf("expected 1 title, got %d", len(titles))
	}
	if titles[0].TitleDisplayName != nil {
		t.Error("expected nil TitleDisplayName")
	}
	if titles[0].IconHiRes != nil {
		t.Error("expected nil IconHiRes")
	}
}

// mockLogger implements Logger for testing.
type mockLogger struct {
	requestCount  int
	responseCount int
}

func (m *mockLogger) LogRequest(_ context.Context, _, _ string, _ []byte) {
	m.requestCount++
}

func (m *mockLogger) LogResponse(_ context.Context, _ int, _ http.Header, _ []byte) {
	m.responseCount++
}

func (m *mockLogger) LogAuth(_ context.Context, _ string, _ map[string]interface{}) {}
