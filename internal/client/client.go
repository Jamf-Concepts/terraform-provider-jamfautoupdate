// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// defaultHTTPTimeout is the maximum duration for HTTP requests made by the client.
const defaultHTTPTimeout = 30 * time.Second

// Client is a Jamf Auto Update API client.
type Client struct {
	baseURL         string
	definitionsFile string
	httpClient      *http.Client
	logger          Logger
}

// NewClient creates a new Jamf Auto Update API client.
// If definitionsFile is not empty, it will read from the file instead of making HTTP requests.
func NewClient(baseURL string, definitionsFile string) *Client {
	return &Client{
		baseURL:         baseURL,
		definitionsFile: definitionsFile,
		httpClient:      &http.Client{Timeout: defaultHTTPTimeout},
	}
}

// SetLogger sets the logger for the client.
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	defer c.closeWithLog(ctx, resp.Body, "response body")

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

// maxLogBodySize is the maximum number of bytes read from a response body for logging purposes.
const maxLogBodySize = 1 << 20 // 1 MiB

// logHTTPResponse logs the HTTP response details using the client's logger.
// It reads the full body, logs up to maxLogBodySize bytes, and replaces resp.Body
// so subsequent readers still see the complete response.
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

	logBody := responseBody
	if len(logBody) > maxLogBodySize {
		logBody = logBody[:maxLogBodySize]
	}

	c.logger.LogResponse(ctx, resp.StatusCode, resp.Header, logBody)
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
}
