package clickup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SingleList struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Orderindex int    `json:"orderindex"`
	Status     struct {
	} `json:"status"`
	Priority struct {
	} `json:"priority"`
	Assignee struct {
	} `json:"assignee"`
	TaskCount int `json:"task_count"`
	DueDate   struct {
	} `json:"due_date"`
	StartDate struct {
	} `json:"start_date"`
	Folder struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Hidden bool   `json:"hidden"`
		Access bool   `json:"access"`
	} `json:"folder"`
	Space struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Access bool   `json:"access"`
	} `json:"space"`
	Archived         bool   `json:"archived"`
	OverrideStatuses bool   `json:"override_statuses"`
	PermissionLevel  string `json:"permission_level"`
}

type ListsResponse struct {
	Lists []SingleList `json:"lists"`
}

func (c *Client) ListsForFolder(folderID string, includeArchived bool) (*ListsResponse, error) {

	urlValues := url.Values{}
	urlValues.Set("archived", strconv.FormatBool(includeArchived))

	endpoint := fmt.Sprintf("%s/folder/%s/list/?%s", c.baseURL, folderID, urlValues.Encode())

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("lists for folder request failed: %w", err)
	}
	req.Header.Add("Authorization", c.opts.APIToken)

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make lists request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		var errResponse ErrClickupResponse
		if err := decoder.Decode(&errResponse); err != nil {
			return nil, &HTTPError{
				Status:     res.Status,
				StatusCode: res.StatusCode,
				URL:        res.Request.URL.String(),
			}
		}
		errResponse.StatusCode = res.StatusCode
		errResponse.Status = res.Status
		return nil, &errResponse
	}

	var lists ListsResponse

	if err := decoder.Decode(&lists); err != nil {
		return nil, fmt.Errorf("failed to parse lists: %w", err)
	}

	return &lists, nil
}

func (c *Client) ListByID(listID string) (*SingleList, error) {

	endpoint := fmt.Sprintf("%s/list/%s", c.baseURL, listID)

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("list request failed: %w", err)
	}
	req.Header.Add("Authorization", c.opts.APIToken)

	res, err := c.doer.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make list request: %w", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		var errResponse ErrClickupResponse
		if err := decoder.Decode(&errResponse); err != nil {
			return nil, &HTTPError{
				Status:     res.Status,
				StatusCode: res.StatusCode,
				URL:        res.Request.URL.String(),
			}
		}
		errResponse.StatusCode = res.StatusCode
		errResponse.Status = res.Status
		return nil, &errResponse
	}

	var list SingleList

	if err := decoder.Decode(&list); err != nil {
		return nil, fmt.Errorf("failed parse to list: %w", err)
	}

	return &list, nil
}
