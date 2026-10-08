// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strings"
)

// closeWithLog closes the given closer and logs any error using the client's logger.
func (c *Client) closeWithLog(ctx context.Context, closer io.Closer, name string) {
	if err := closer.Close(); err != nil {
		if c.logger != nil {
			c.logger.LogAuth(ctx, fmt.Sprintf("Failed to close %s", name), map[string]any{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("warning: failed to close %s: %v\n", name, err)
		}
	}
}

// titlesMissing returns the list of requested title names that are not present in the given titles slice.
// Names are compared case-insensitively, matching the behavior of the definitions API.
func titlesMissing(titles []Title, requested []string) []string {
	found := make(map[string]struct{}, len(titles))
	for _, title := range titles {
		if title.TitleName != nil {
			found[strings.ToLower(*title.TitleName)] = struct{}{}
		}
	}

	var missing []string
	for _, name := range requested {
		if _, ok := found[strings.ToLower(name)]; !ok {
			missing = append(missing, name)
		}
	}

	return missing
}

// chunkTitleNames splits names into consecutive groups whose escaped, comma-joined length does not
// exceed maxLength. A name longer than maxLength is placed in a group of its own.
func chunkTitleNames(names []string, maxLength int) [][]string {
	var chunks [][]string
	var current []string
	currentLength := 0

	for _, name := range names {
		nameLength := len(url.PathEscape(name))
		added := nameLength
		if len(current) > 0 {
			added++
		}

		if len(current) > 0 && currentLength+added > maxLength {
			chunks = append(chunks, current)
			current = nil
			currentLength = 0
			added = nameLength
		}

		current = append(current, name)
		currentLength += added
	}

	if len(current) > 0 {
		chunks = append(chunks, current)
	}

	return chunks
}
