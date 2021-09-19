package clickup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
)

type CreateAttachmentResponse struct {
	ID             string `json:"id"`
	Version        string `json:"version"`
	Date           int    `json:"date"`
	Title          string `json:"title"`
	Extension      string `json:"extension"`
	ThumbnailSmall string `json:"thumbnail_small"`
	ThumbnailLarge string `json:"thumbnail_large"`
	URL            string `json:"url"`
}

type AttachmentParams struct {
	FileName string
	Reader   io.Reader
}

func (c *Client) CreateTaskAttachment(taskID, workspaceID string, useCustomTaskIDs bool, params *AttachmentParams) (*CreateAttachmentResponse, error) {
	if useCustomTaskIDs && workspaceID == "" {
		return nil, fmt.Errorf("workspaceID must be provided if querying by custom task id: %w", ErrValidation)
	}

	contents, err := ioutil.ReadAll(params.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of Reader: %w", err)
	}

	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)
	part, err := multipartWriter.CreateFormFile("attachment", filepath.Base(params.FileName)) // must be "attachment"
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart field: %w", err)
	}
	part.Write(contents)

	if err := multipartWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(useCustomTaskIDs))
	urlValues.Add("team_id", workspaceID)

	endpoint := fmt.Sprintf("%s/task/%s/attachment/?%s", c.baseURL, taskID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodPost, endpoint, &body)
	if err != nil {
		return nil, fmt.Errorf("create attachment request failed: %w", err)
	}
	if err := c.AuthenticateFor(req); err != nil {
		return nil, fmt.Errorf("failed to authenticate client: %w", err)
	}
	req.Header.Add("Content-Type", multipartWriter.FormDataContentType())

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make create attachment request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errorFromResponse(res, decoder)
	}

	var attachmentResponse CreateAttachmentResponse

	if err := decoder.Decode(&attachmentResponse); err != nil {
		return nil, fmt.Errorf("failed to parse attachment response: %w", err)
	}

	return &attachmentResponse, nil
}
