package clickup

import (
	"context"
	"fmt"
	"net/http"
)

type SharedHierarchyResponse struct {
	Shared struct {
		Tasks []struct {
			ID       string `json:"id"`
			CustomID struct {
			} `json:"custom_id"`
			Name   string `json:"name"`
			Status struct {
				Status     string `json:"status"`
				Color      string `json:"color"`
				Type       string `json:"type"`
				Orderindex int    `json:"orderindex"`
			} `json:"status"`
			Orderindex  string `json:"orderindex"`
			DateCreated string `json:"date_created"`
			DateUpdated string `json:"date_updated"`
			DateClosed  string `json:"date_closed"`
			Archived    bool   `json:"archived"`
			Creator     struct {
				ID             int    `json:"id"`
				Username       string `json:"username"`
				Color          string `json:"color"`
				Email          string `json:"email"`
				ProfilePicture string `json:"profilePicture"`
			} `json:"creator"`
			Assignees []struct {
				ID             int    `json:"id"`
				Username       string `json:"username"`
				Color          string `json:"color"`
				Initials       string `json:"initials"`
				Email          string `json:"email"`
				ProfilePicture string `json:"profilePicture"`
			} `json:"assignees"`
			Parent struct {
			} `json:"parent"`
			Priority struct {
			} `json:"priority"`
			DueDate struct {
			} `json:"due_date"`
			StartDate struct {
			} `json:"start_date"`
			Points struct {
			} `json:"points"`
			TimeEstimate struct {
			} `json:"time_estimate"`
			CustomFields []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Type       string `json:"type"`
				TypeConfig struct {
					NewDropDown bool `json:"new_drop_down"`
					Options     []struct {
						ID         string `json:"id"`
						Name       string `json:"name"`
						Color      string `json:"color"`
						Orderindex int    `json:"orderindex"`
					} `json:"options"`
				} `json:"type_config"`
				DateCreated    string `json:"date_created"`
				HideFromGuests bool   `json:"hide_from_guests"`
				Required       bool   `json:"required"`
			} `json:"custom_fields"`
			TeamID          string `json:"team_id"`
			URL             string `json:"url"`
			PermissionLevel string `json:"permission_level"`
			List            struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Access bool   `json:"access"`
			} `json:"list"`
			Project struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Hidden bool   `json:"hidden"`
				Access bool   `json:"access"`
			} `json:"project"`
			Folder struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Hidden bool   `json:"hidden"`
				Access bool   `json:"access"`
			} `json:"folder"`
			Space struct {
				ID string `json:"id"`
			} `json:"space"`
		} `json:"tasks"`
		Lists []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Orderindex int    `json:"orderindex"`
			TaskCount  string `json:"task_count"`
			Archived   bool   `json:"archived"`
		} `json:"lists"`
		Folders []struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			Orderindex int    `json:"orderindex"`
			TaskCount  string `json:"task_count"`
			Archived   bool   `json:"archived"`
		} `json:"folders"`
	} `json:"shared"`
}

// SharedHierarchy returns resources that the authenticated user has access to, but not to its parent.
// From the ClickUp documentation:
// "Returns all resources you have access to where you don't have access to its parent. For example,
// if you have a access to a shared task, but don't have access to its parent list, it will come back in this request."
func (c *Client) SharedHierarchy(ctx context.Context, workspaceID string) (*SharedHierarchyResponse, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id to query shared hierarchy: %w", ErrValidation)
	}
	endpoint := fmt.Sprintf("/team/%s/shared", workspaceID)

	var sharedHierarchyResponse SharedHierarchyResponse

	if err := c.call(ctx, http.MethodGet, endpoint, nil, &sharedHierarchyResponse); err != nil {
		return nil, fmt.Errorf("failed to make clickup request: %w", err)
	}

	return &sharedHierarchyResponse, nil
}
