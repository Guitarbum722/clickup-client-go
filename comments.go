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
	"net/url"
	"strconv"
)

type List struct {
	List string `json:"list"`
}

type CodeBlock struct {
	CodeBlock string `json:"code-block"`
}

type Attributes struct {
	Bold      bool       `json:"bold"`
	Italic    bool       `json:"italic"`
	Code      bool       `json:"code"`
	CodeBlock *CodeBlock `json:"code-block,omitempty"`
	List      *List      `json:"list"`
}

type Emoticon struct {
	Code string `json:"code"`
}

type ComplexComment struct {
	Text       string      `json:"text"`
	Type       string      `json:"type,omitempty"`
	Attributes *Attributes `json:"attributes,omitempty"`
	Emoticon   *Emoticon   `json:"emoticon,omitempty"`
}

type CreateCommentRequest struct {
	CommentText string           `json:"comment_text,omitempty"` // plain text
	Comment     []ComplexComment `json:"comment,omitempty"`
	Assignee    int              `json:"assignee,omitempty"`
	NotifyAll   bool             `json:"notify_all,omitempty"`
}

func (c *CreateCommentRequest) BulletedListItem(text string, attributes *Attributes) {
	if c.Comment == nil {
		c.Comment = make([]ComplexComment, 0, 2)
	}
	comment := []ComplexComment{
		{
			Text:       text,
			Attributes: attributes,
		},
		{
			Text: "\n",
			Attributes: &Attributes{
				List: &List{
					List: "bullet",
				},
			},
		},
	}
	c.Comment = append(c.Comment, comment...)
}

func (c *CreateCommentRequest) NumberedListItem(text string, attributes *Attributes) {
	if c.Comment == nil {
		c.Comment = make([]ComplexComment, 0, 2)
	}
	comment := []ComplexComment{
		{
			Text:       text,
			Attributes: attributes,
		},
		{
			Text: "\n",
			Attributes: &Attributes{
				List: &List{
					List: "ordered",
				},
			},
		},
	}
	c.Comment = append(c.Comment, comment...)
}

func (c *CreateCommentRequest) ChecklistItem(text string, checked bool, attributes *Attributes) {
	if c.Comment == nil {
		c.Comment = make([]ComplexComment, 0, 2)
	}
	isCheckedVal := "unchecked"
	if checked {
		isCheckedVal = "checked"
	}
	comment := []ComplexComment{
		{
			Text:       text,
			Attributes: attributes,
		},
		{
			Text: "\n",
			Attributes: &Attributes{
				List: &List{
					List: isCheckedVal,
				},
			},
		},
	}
	c.Comment = append(c.Comment, comment...)
}

type CreateTaskCommentRequest struct {
	CreateCommentRequest
	TaskID           string `json:"-"`
	UseCustomTaskIDs bool   `json:"-"`
	WorkspaceID      string `json:"-"`
}

func NewCreateTaskCommentRequest(taskID string, useCustomTaskIDs bool, workspaceID string) *CreateTaskCommentRequest {
	return &CreateTaskCommentRequest{
		TaskID:           taskID,
		UseCustomTaskIDs: useCustomTaskIDs,
		WorkspaceID:      workspaceID,
	}
}

type CreateCommentResponse struct {
	ID        int    `json:"id"`
	HistoryID string `json:"hist_id"`
	Date      int    `json:"date"`
}

type CreateTaskCommentResponse struct {
	CreateCommentResponse
}

func (c *Client) CreateTaskComment(ctx context.Context, comment CreateTaskCommentRequest) (*CreateTaskCommentResponse, error) {
	if comment.TaskID == "" {
		return nil, fmt.Errorf("must provide a task id to create a task comment: %w", ErrValidation)
	}
	if comment.UseCustomTaskIDs && comment.WorkspaceID == "" {
		return nil, fmt.Errorf("must provide a workspace id for a new task comment if using custom task ID: %w", ErrValidation)
	}

	b, err := json.Marshal(comment)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new task comment: %w", err)
	}

	buf := bytes.NewBuffer(b)

	urlValues := url.Values{}
	urlValues.Set("custom_task_ids", strconv.FormatBool(comment.UseCustomTaskIDs))
	urlValues.Add("team_id", comment.WorkspaceID)

	endpoint := fmt.Sprintf("/task/%s/comment/?%s", comment.TaskID, urlValues.Encode())

	var commentResponse CreateTaskCommentResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &commentResponse); err != nil {
		return nil, ErrCall
	}

	return &commentResponse, nil
}

type CreateChatViewCommentRequest struct {
	CreateCommentRequest
	ViewID string `json:"-"`
}

func NewCreateChatViewCommentRequest(viewID string) *CreateChatViewCommentRequest {
	return &CreateChatViewCommentRequest{
		ViewID: viewID,
	}
}

type CreateChatViewCommentResponse struct {
	CreateCommentResponse
}

func (c *Client) CreateChatViewComment(ctx context.Context, comment CreateChatViewCommentRequest) (*CreateChatViewCommentResponse, error) {
	if comment.ViewID == "" {
		return nil, fmt.Errorf("must provide a view id to create a view comment: %w", ErrValidation)
	}

	b, err := json.Marshal(comment)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new comment: %w", err)
	}

	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/view/%s/comment", comment.ViewID)

	var commentResponse CreateChatViewCommentResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &commentResponse); err != nil {
		return nil, ErrCall
	}

	return &commentResponse, nil

}

type CreateListCommentRequest struct {
	CreateCommentRequest
	ListID string `json:"-"`
}

func NewCreateListCommentRequest(listID string) *CreateListCommentRequest {
	return &CreateListCommentRequest{
		ListID: listID,
	}
}

type CreateListCommentResponse struct {
	CreateCommentResponse
}

func (c *Client) CreateListComment(ctx context.Context, comment CreateListCommentRequest) (*CreateListCommentResponse, error) {
	if comment.ListID == "" {
		return nil, fmt.Errorf("must provide a list id to create a view comment: %w", ErrValidation)
	}

	b, err := json.Marshal(comment)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize new comment: %w", err)
	}

	buf := bytes.NewBuffer(b)

	endpoint := fmt.Sprintf("/list/%s/comment", comment.ListID)

	var commentResponse CreateListCommentResponse

	if err := c.call(ctx, http.MethodPost, endpoint, buf, &commentResponse); err != nil {
		return nil, ErrCall
	}

	return &commentResponse, nil
}

type CommentsResponse struct {
	Comments []struct {
		ID          string           `json:"id"`
		Comment     []ComplexComment `json:"comment"`
		CommentText string           `json:"comment_text"`
		User        *TeamUser        `json:"user"`
		Assignee    *TeamUser        `json:"assignee"`
		AssignedBy  *TeamUser        `json:"assigned_by"`
		Reactions   []struct {
			Reaction string   `json:"reaction"`
			Date     string   `json:"date"`
			User     TeamUser `json:"user"`
		} `json:"reactions"`
		Date string `json:"date"`
	} `json:"comments"`
}

type CommentsQuery struct {
	TaskID           string
	UseCustomTaskIDs bool
	WorkspaceID      string
	ListID           string
	ViewID           string
}

type CommentsForTaskQuery struct {
	CommentsQuery
}

func (c *Client) TaskComments(ctx context.Context, query CommentsForTaskQuery) (CommentsResponse, error) {
	panic("TODO")
}

type CommentsForTaskViewQuery struct {
	CommentsQuery
}

func (c *Client) ChatViewComments(ctx context.Context, query CommentsForTaskViewQuery) (CommentsResponse, error) {
	panic("TODO")
}

type CommentsForListQuery struct {
	CommentsQuery
}

func (c *Client) ListComments(ctx context.Context, query CommentsForListQuery) (CommentsResponse, error) {
	panic("TODO")
}

type UpdateCommentRequest struct {
	CommentText string `json:"comment_text"` // plain text
	Comment     []struct {
		Text       string `json:"text"`
		Type       string `json:"type"`
		Attributes struct {
			Bold      bool `json:"bold"`
			Italic    bool `json:"italic"`
			Code      bool `json:"code"`
			CodeBlock struct {
				CodeBlock string `json:"code-block"`
			} `json:"code-block,omitempty"`
		} `json:"attributes"`
		Emoticon struct {
			Code string `json:"code"`
		} `json:"emoticon"`
	} `json:"comment"`
	Assignee  int  `json:"assignee"`
	NotifyAll bool `json:"notify_all"`
	Resolved  bool `json:"resolved"`
}

func (c *Client) UpdateComment(ctx context.Context, comment UpdateCommentRequest) error {
	panic("TODO")
}

func (c *Client) DeleteComment(ctx context.Context, commentID string) error {
	panic("TODO")
}
