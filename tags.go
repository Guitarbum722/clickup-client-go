// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Tag struct {
	Name    string `json:"name"`
	TagFg   string `json:"tag_fg"`
	TagBg   string `json:"tag_bg"`
	Creator int    `json:"creator"`
}

type TagsQueryResponse struct {
	Tags []Tag `json:"tags"`
}

// TagsForSpace returns a listing of tags associated to the space using spaceID.
func (c *Client) TagsForSpace(ctx context.Context, spaceID string) (*TagsQueryResponse, error) {
	endpoint := fmt.Sprintf("/space/%s/tag", spaceID)

	var tagsResponse TagsQueryResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &tagsResponse); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &tagsResponse, nil
}

// CreateSpaceTag adds a new tag to the space with spaceID.
func (c *Client) CreateSpaceTag(ctx context.Context, spaceID string, tag Tag) error {
	if tag.Name == "" {
		return fmt.Errorf("must provide a name for new tag: %w", ErrValidation)
	}
	b, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("unable to serialize new tag: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/space/%s/tag", spaceID)

	return c.call(ctx, http.MethodPost, endpoint, buf, &struct{}{})
}

// UpdateSpaceTag updates an existing tag in the specified space with the parameters from tag.
func (c *Client) UpdateSpaceTag(ctx context.Context, spacID, tag Tag) error {
	panic("TODO not implemented")
}
