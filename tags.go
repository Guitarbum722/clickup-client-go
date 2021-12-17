// Copyright (c) 2021, John Moore
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

func (c *Client) TagsForSpace(ctx context.Context, spaceID string) (*TagsQueryResponse, error) {
	endpoint := fmt.Sprintf("%s/space/%s/tag", c.baseURL, spaceID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("tags for space request failed: %w", err)
	}
	if err := c.AuthenticateFor(ctx, req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make tags for spaces request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var tagsResponse TagsQueryResponse

	if err := decoder.Decode(&tagsResponse); err != nil {
		return nil, fmt.Errorf("failed parse to tags response: %w", err)
	}

	return &tagsResponse, nil
}

func (c *Client) CreateSpaceTag(ctx context.Context, spaceID string, tag Tag) error {
	if tag.Name == "" {
		return fmt.Errorf("must provide a name for new tag: %w", ErrValidation)
	}
	b, err := json.Marshal(tag)
	if err != nil {
		return fmt.Errorf("unable to serialize new tag: %w", err)
	}
	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("%s/space/%s/tag", c.baseURL, spaceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, buf)
	if err != nil {
		return fmt.Errorf("create tag request failed: %w", err)
	}
	if err := c.AuthenticateFor(ctx, req); err != nil {
		return fmt.Errorf("failed to authenticate client: %w", err)
	}

	res, err := c.doer.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make create tag request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return errorFromResponse(res, decoder)
	}

	return nil
}

func (c *Client) UpdateSpaceTag(ctx context.Context, spacID, tag Tag) error {
	return nil
}
