// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

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
	defer c.closeWithLog(ctx, file, "definitions file")

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
		slices.Sort(missing)
		return nil, &TitlesNotFoundError{MissingTitles: missing}
	}

	return titles, nil
}
