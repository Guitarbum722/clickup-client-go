// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import "context"

type CreateCommentRequest struct {
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
}

type CreateTaskCommentRequest struct {
	CreateCommentRequest
	TaskID           string
	UseCustomTaskIDs bool
	WorkspaceID      string
}

type CreateCommentResponse struct {
	ID        string `json:"id"`
	HistoryID string `json:"hist_id"`
	Date      int    `json:"date"`
}

type CreateTaskCommentResponse struct {
	CreateCommentResponse
}

func (c *Client) CreateTaskComment(ctx context.Context, comment CreateTaskCommentRequest) (*CreateTaskCommentResponse, error) {
	panic("TODO")
}

type CreateChatViewCommentRequest struct {
	CreateCommentRequest
	ViewID string
}

type CreateChatViewCommentResponse struct {
	CreateCommentResponse
}

func (c *Client) CreateChatViewComment(ctx context.Context, comment CreateChatViewCommentRequest) (*CreateChatViewCommentResponse, error) {
	panic("TODO")
}

type CreateListCommentRequest struct {
	CreateCommentRequest
	ListID string
}

type CreateListCommentResponse struct {
	CreateCommentResponse
}

func (c *Client) CreateListComment(ctx context.Context, commen CreateListCommentRequest) (*CreateListCommentResponse, error) {
	panic("TODO")
}

type CommentsResponse struct {
	Comments []struct {
		ID      string `json:"id"`
		Comment []struct {
			Text       string `json:"text"`
			Attributes struct {
				CodeBlock struct {
					CodeBlock string `json:"code-block"`
				} `json:"code-block"`
			} `json:"attributes,omitempty"`
		} `json:"comment"`
		CommentText string   `json:"comment_text"`
		User        TeamUser `json:"user"`
		Assignee    TeamUser `json:"assignee"`
		AssignedBy  TeamUser `json:"assigned_by"`
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
