// Copyright 2025 Jamf Software LLC.

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Title struct {
	TitleName                *string     `json:"title_name"`
	TitleDisplayName         *string     `json:"title_display_name"`
	TitleDescription         *string     `json:"title_description"`
	TitleVersion             *string     `json:"title_version"`
	MinimumOS                *string     `json:"minimum_os"`
	MaximumOS                *string     `json:"maximum_os"`
	IconHiRes                *string     `json:"icon_hires"`
	ExtensionAttribute       interface{} `json:"extension_attribute"`
	ContentFilterProfile     interface{} `json:"content_filter_profile"`
	KernelExtensionProfile   interface{} `json:"kernel_extension_profile"`
	ManagedLoginItemsProfile interface{} `json:"managed_login_items_profile"`
	NotificationsProfile     interface{} `json:"notifications_profile"`
	PPPCPProfile             interface{} `json:"pppcp_profile"`
	ScreenRecordingProfile   interface{} `json:"screen_recording_profile"`
	SystemExtensionProfile   interface{} `json:"system_extension_profile"`
	PatchDefinition          struct {
		Requirements []struct {
			Name  *string `json:"name"`
			Value *string `json:"value"`
		} `json:"requirements"`
	} `json:"patch_definition"`
}

type TitlesNotFoundError struct {
	MissingTitles []string
}

// TitlesNotFoundError returns a list of titles that were not found.
func (e *TitlesNotFoundError) Error() string {
	return fmt.Sprintf("The following titles were not found: %s", strings.Join(e.MissingTitles, ", "))
}

// NewClient creates a new Jamf Auto Update API client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// GetTitles retrieves titles from the API. If titleNames is empty, it returns all titles.
// If titleNames contains one or more names, it returns data for those specific titles.
func (c *Client) GetTitles(ctx context.Context, titleNames ...string) ([]Title, error) {
	url := c.baseURL
	if len(titleNames) > 0 {
		url = fmt.Sprintf("%s/%s", c.baseURL, strings.Join(titleNames, ","))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("warning: failed to close response body: %v\n", err)
		}
	}()

	var titles []Title
	if err := json.NewDecoder(resp.Body).Decode(&titles); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(titleNames) > 0 {
		foundTitles := make(map[string]bool)
		for _, title := range titles {
			if title.TitleName != nil {
				foundTitles[*title.TitleName] = true
			}
		}

		var missingTitles []string
		for _, requestedTitle := range titleNames {
			if !foundTitles[requestedTitle] {
				missingTitles = append(missingTitles, requestedTitle)
			}
		}

		if len(missingTitles) > 0 {
			return nil, &TitlesNotFoundError{MissingTitles: missingTitles}
		}
	}

	return titles, nil
}

func (t *Title) GetExtensionAttribute() *string {
	return interfaceToString(t.ExtensionAttribute)
}

func (t *Title) GetContentFilterProfile() *string {
	return interfaceToString(t.ContentFilterProfile)
}

func (t *Title) GetKernelExtensionProfile() *string {
	return interfaceToString(t.KernelExtensionProfile)
}

func (t *Title) GetManagedLoginItemsProfile() *string {
	return interfaceToString(t.ManagedLoginItemsProfile)
}

func (t *Title) GetNotificationsProfile() *string {
	return interfaceToString(t.NotificationsProfile)
}

func (t *Title) GetPPPCPProfile() *string {
	return interfaceToString(t.PPPCPProfile)
}

func (t *Title) GetScreenRecordingProfile() *string {
	return interfaceToString(t.ScreenRecordingProfile)
}

func (t *Title) GetSystemExtensionProfile() *string {
	return interfaceToString(t.SystemExtensionProfile)
}

// Helper function to convert interface{} to string.
func interfaceToString(v interface{}) *string {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case string:
		return &val
	default:
		return nil
	}
}
