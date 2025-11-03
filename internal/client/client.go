// Copyright 2025 Jamf Software LLC.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// Logger is an interface for logging HTTP requests and responses
type Logger interface {
	LogRequest(ctx context.Context, method, url string, body []byte)
	LogResponse(ctx context.Context, statusCode int, headers http.Header, body []byte)
	LogAuth(ctx context.Context, message string, fields map[string]interface{})
}

// Client is a Jamf Auto Update API client.
type Client struct {
	baseURL         string
	definitionsFile string
	httpClient      *http.Client
	logger          Logger
}

// Title represents a title retrieved from the Jamf Auto Update API.
type Title struct {
	TitleName                *string         `json:"title_name"`
	TitleDisplayName         *string         `json:"title_display_name"`
	TitleDescription         *string         `json:"title_description"`
	TitleVersion             *string         `json:"title_version"`
	MinimumOS                *string         `json:"minimum_os"`
	MaximumOS                *string         `json:"maximum_os"`
	IconHiRes                *string         `json:"icon_hires"`
	ExtensionAttribute       *string         `json:"extension_attribute"`
	ContentFilterProfile     *string         `json:"content_filter_profile"`
	KernelExtensionProfile   *string         `json:"kernel_extension_profile"`
	ManagedLoginItemsProfile *string         `json:"managed_login_items_profile"`
	NotificationsProfile     *string         `json:"notifications_profile"`
	PPPCPProfile             *string         `json:"pppcp_profile"`
	ScreenRecordingProfile   *string         `json:"screen_recording_profile"`
	SystemExtensionProfile   *string         `json:"system_extension_profile"`
	PatchDefinition          PatchDefinition `json:"patch_definition"`
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

// TitlesNotFoundError returns a list of titles that were not found.
func (e *TitlesNotFoundError) Error() string {
	return fmt.Sprintf("The following titles were not found: %s", strings.Join(e.MissingTitles, ", "))
}

// NewClient creates a new Jamf Auto Update API client.
// If definitionsFile is not empty, it will read from the file instead of making HTTP requests.
func NewClient(baseURL string, definitionsFile string) *Client {
	return &Client{
		baseURL:         baseURL,
		definitionsFile: definitionsFile,
		httpClient:      &http.Client{},
	}
}

// SetLogger sets the logger for the client
func (c *Client) SetLogger(logger Logger) {
	c.logger = logger
}

// GetTitles retrieves titles from the API or file. If titleNames is empty, it returns all titles.
// If titleNames contains one or more names, it returns data for those specific titles.
func (c *Client) GetTitles(ctx context.Context, titleNames ...string) ([]Title, error) {
	if c.definitionsFile != "" {
		return c.getTitlesFromFile(ctx, titleNames...)
	}

	url := c.baseURL
	if len(titleNames) > 0 {
		url = fmt.Sprintf("%s/%s", c.baseURL, strings.Join(titleNames, ","))
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if c.logger != nil {
		c.logger.LogRequest(ctx, req.Method, req.URL.String(), nil)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if c.logger != nil {
		c.logHTTPResponse(ctx, resp)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			if c.logger != nil {
				c.logger.LogAuth(ctx, "Failed to close response body", map[string]interface{}{
					"error": cerr.Error(),
				})
			} else {
				fmt.Printf("warning: failed to close response body: %v\n", cerr)
			}
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var titles []Title
	if err := json.NewDecoder(resp.Body).Decode(&titles); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(titleNames) > 0 {
		if missing := titlesMissing(titles, titleNames); len(missing) > 0 {
			return nil, &TitlesNotFoundError{MissingTitles: missing}
		}
	}

	return titles, nil
}

// logHTTPResponse logs the HTTP response details using the client's logger
func (c *Client) logHTTPResponse(ctx context.Context, resp *http.Response) {
	if resp.Body == nil {
		c.logger.LogResponse(ctx, resp.StatusCode, resp.Header, nil)
		return
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.LogAuth(ctx, "Failed to read response body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := resp.Body.Close(); err != nil {
		c.logger.LogAuth(ctx, "Failed to close response body", map[string]interface{}{
			"error": err.Error(),
		})
	}

	c.logger.LogResponse(ctx, resp.StatusCode, resp.Header, responseBody)
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
}

// getTitlesFromFile retrieves titles from a local JSON file.
func (c *Client) getTitlesFromFile(ctx context.Context, titleNames ...string) ([]Title, error) {
	if c.logger != nil {
		fields := map[string]interface{}{
			"definitions_file": c.definitionsFile,
		}
		if len(titleNames) > 0 {
			fields["requested_titles"] = titleNames
		}
		c.logger.LogAuth(ctx, "Reading titles from definitions file", fields)
	}

	file, err := os.Open(c.definitionsFile)
	if err != nil {
		return nil, fmt.Errorf("error opening definitions file: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			if c.logger != nil {
				c.logger.LogAuth(ctx, "Failed to close definitions file", map[string]interface{}{
					"error": cerr.Error(),
				})
			} else {
				fmt.Printf("warning: failed to close definitions file: %v\n", cerr)
			}
		}
	}()

	decoder := json.NewDecoder(file)

	if len(titleNames) == 0 {
		var titles []Title
		if err := decoder.Decode(&titles); err != nil {
			return nil, fmt.Errorf("error decoding definitions file: %w", err)
		}
		return titles, nil
	}

	wanted := make(map[string]struct{}, len(titleNames))
	for _, name := range titleNames {
		wanted[name] = struct{}{}
	}

	token, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("error reading definitions file: %w", err)
	}
	delim, ok := token.(json.Delim)
	if !ok || delim != '[' {
		return nil, fmt.Errorf("definitions file must contain a JSON array of titles")
	}

	var titles []Title
	for decoder.More() {
		var title Title
		if err := decoder.Decode(&title); err != nil {
			return nil, fmt.Errorf("error decoding definitions file: %w", err)
		}

		if title.TitleName == nil {
			continue
		}

		if _, ok := wanted[*title.TitleName]; ok {
			titles = append(titles, title)
			delete(wanted, *title.TitleName)

			if len(wanted) == 0 {
				break
			}
		}
	}

	if len(wanted) > 0 {
		missing := make([]string, 0, len(wanted))
		for name := range wanted {
			missing = append(missing, name)
		}
		return nil, &TitlesNotFoundError{MissingTitles: missing}
	}

	return titles, nil
}

func titlesMissing(titles []Title, requested []string) []string {
	found := make(map[string]struct{}, len(titles))
	for _, title := range titles {
		if title.TitleName != nil {
			found[*title.TitleName] = struct{}{}
		}
	}

	var missing []string
	for _, name := range requested {
		if _, ok := found[name]; !ok {
			missing = append(missing, name)
		}
	}

	return missing
}
