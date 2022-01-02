// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"reflect"
	"testing"
)

func TestCreateCommentRequest_BulletedListItem(t *testing.T) {
	type fields struct {
		CommentText string
		Comment     []ComplexComment
		Assignee    int
		NotifyAll   bool
	}
	type args struct {
		text       string
		attributes *Attributes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateCommentRequest
	}{
		{
			name: "Successful build bulleted list starting with nil comment items",
			fields: fields{
				Comment:     nil,
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Bullet item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build bulleted list starting with zero len comment items",
			fields: fields{
				Comment:     []ComplexComment{},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Bullet item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build bulleted list with existing comment items",
			fields: fields{
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Bullet item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Bullet item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
					{
						Text:       "Bullet item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "bullet",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateCommentRequest{
				CommentText: tt.fields.CommentText,
				Comment:     tt.fields.Comment,
				Assignee:    tt.fields.Assignee,
				NotifyAll:   tt.fields.NotifyAll,
			}
			if c.BulletedListItem(tt.args.text, tt.args.attributes); !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CreateCommentRequest.BulletedListItem() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestCreateCommentRequest_NumberedListItem(t *testing.T) {
	type fields struct {
		CommentText string
		Comment     []ComplexComment
		Assignee    int
		NotifyAll   bool
	}
	type args struct {
		text       string
		attributes *Attributes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateCommentRequest
	}{
		{
			name: "Successful build numbered list starting with nil comment items",
			fields: fields{
				Comment:     nil,
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Ordered item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build ordered list starting with zero len comment items",
			fields: fields{
				Comment:     []ComplexComment{},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Ordered item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build ordered list with existing comment items",
			fields: fields{
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "Ordered item 1",
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "Ordered item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
					{
						Text:       "Ordered item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "ordered",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateCommentRequest{
				CommentText: tt.fields.CommentText,
				Comment:     tt.fields.Comment,
				Assignee:    tt.fields.Assignee,
				NotifyAll:   tt.fields.NotifyAll,
			}
			if c.NumberedListItem(tt.args.text, tt.args.attributes); !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CreateCommentRequest.NumberedListItem() = %v, want %v", c, tt.want)
			}
		})
	}
}

func TestCreateCommentRequest_ChecklistItem(t *testing.T) {
	type fields struct {
		CommentText string
		Comment     []ComplexComment
		Assignee    int
		NotifyAll   bool
	}
	type args struct {
		text       string
		checked    bool
		attributes *Attributes
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *CreateCommentRequest
	}{
		{
			name: "Successful build checklist item with nil comment items",
			fields: fields{
				Comment:     nil,
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "checklist item 1",
				checked:    true,
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "checklist item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "checked",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build checklist item starting with zero len comment items",
			fields: fields{
				Comment:     []ComplexComment{},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "checklist item 1",
				checked:    false,
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "checklist item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
				},
			},
		},
		{
			name: "Successful build checklist item with existing comment items",
			fields: fields{
				Comment: []ComplexComment{
					{
						Text:       "checklist item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
				},
				CommentText: "",
				Assignee:    0,
				NotifyAll:   false,
			},
			args: args{
				text:       "checklist item 1",
				checked:    false,
				attributes: nil,
			},
			want: &CreateCommentRequest{
				CommentText: "",
				Comment: []ComplexComment{
					{
						Text:       "checklist item 0",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
					{
						Text:       "checklist item 1",
						Attributes: nil,
					},
					{
						Text: "\n",
						Attributes: &Attributes{
							List: &List{
								List: "unchecked",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CreateCommentRequest{
				CommentText: tt.fields.CommentText,
				Comment:     tt.fields.Comment,
				Assignee:    tt.fields.Assignee,
				NotifyAll:   tt.fields.NotifyAll,
			}
			if c.ChecklistItem(tt.args.text, tt.args.checked, tt.args.attributes); !reflect.DeepEqual(c, tt.want) {
				t.Errorf("CreateCommentRequest.ChecklistItem() = %v, want %v", c, tt.want)
			}
		})
	}
}
