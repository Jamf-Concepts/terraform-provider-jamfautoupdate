// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"fmt"
	"io"
)

// closeWithLog closes the given closer and logs any error using the client's logger.
func (c *Client) closeWithLog(ctx context.Context, closer io.Closer, name string) {
	if err := closer.Close(); err != nil {
		if c.logger != nil {
			c.logger.LogAuth(ctx, fmt.Sprintf("Failed to close %s", name), map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("warning: failed to close %s: %v\n", name, err)
		}
	}
}

// titlesMissing returns the list of requested title names that are not present in the given titles slice.
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
