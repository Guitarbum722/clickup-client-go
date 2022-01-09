# clickup-client-go

### Get Started

Create a Clickup Client by providing a `ClientOpts`.  The default `Doer` is an `http.Client` with a `20` second timeout.

Use the `APITokenAuthenticator` for a simple authentication mechanism and provide your Clickup user's API Key.
If you want to implement `Authenticator` in different ways (eg. OAuth flow) then provide your implementation to the client.

```go
	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Args[1],
		},
	})
```

### Tasks

```go
	queryOpts := &clickup.TaskQueryOptions{
		IncludeArchived: false,
	}

	for {
		tasks, _ := client.TasksForList("list-id", queryOpts)

		for _, task := range tasks.Tasks {
			fmt.Println("Task: ", task.CustomID, task.Name)
		}
		if len(tasks.Tasks) < clickup.MaxPageSize {
			return
		} else {
			queryOpts.Page++
		}
	}
}
```

### Create and get webhooks

Create a webhook and listen for Task Updated Events for a particular list.

```go
	newWebhook, _ := client.CreateWebhook(workspaceID, &clickup.CreateWebhookRequest{
		Endpoint: "https://your-webhook.site/myhook",
		Events: []clickup.WebhookEvent{
			clickup.EventTaskUpdated,
		},
		ListID: "list-id",
	})

	fmt.Println("New Webhook ID: ", newWebhook.ID)
	fmt.Println("Health: ", newWebhook.Webhook.Health.FailCount, newWebhook.Webhook.Health.Status)
```

Get webhooks for a workspace.

```go
	webhooks, _ := client.WebhooksFor(workspaceID)

	fmt.Println("Existing webhooks:")
	for _, v := range webhooks.Webhooks {
		fmt.Println("ID:", v.ID)
		fmt.Println("Status:", v.Health.Status, v.Health.FailCount)

		fmt.Println("deleting webhook for:", v.ID)
		client.DeleteWebhook(v.ID)
	}
```

Other webhook events...

```
	EventAll                     
	EventTaskCreated             
	EventTaskUpdated             
	EventTaskDeleted             
	EventTaskPriorityUpdated     
	EventTaskStatusUpdated       
	EventTaskAssigneeUpdated     
	EventTaskDueDateUpdated      
	EventTaskTagUpdated          
	EventTaskMoved               
	EventTaskCommentPosted       
	EventTaskCommentUpdated      
	EventTaskTimeEstimateUpdated 
	EventTaskTimeTrackedUpdated  
	EventListCreated             
	EventListUpdated             
	EventListDeleted             
	EventFolderCreated           
	EventFolderUpdated           
	EventFolderDeleted           
	EventSpaceCreated            
	EventSpaceUpdated            
	EventSpaceDeleted            
	EventGoalCreated             
	EventGoalUpdated             
	EventGoalDeleted             
	EventKeyResultCreated        
	EventKeyResultUpdated        
	EventKeyResultDeleted        
```

### Comments

Comments are not very intuitive via Clickup's API (IMO). This library provides some helpers to construct a comment request (builder).
The comment request mode is completely exported, so feel free to construct it yourself if desired.

```go
	comment := clickup.NewCreateTaskCommentRequest(
		os.Getenv("CLICKUP_TASK_ID"),
		true,
		os.Getenv("CLICKUP_WORKSPACE_ID"),
	)
	comment.BulletedListItem("Bullet Item 4asdf", nil)
	comment.BulletedListItem("Bullet Item 5", nil)
	comment.BulletedListItem("Bullet Item 6", &clickup.Attributes{Italic: true})
	comment.NumberedListItem("Numbered Item 1", nil)
	comment.ChecklistItem("Checklist item 1", false, nil)
	comment.ChecklistItem("Checklist item 2", true, nil)

	res, err := client.CreateTaskComment(context.Background(), *comment)
	if err != nil {
		panic(err)
	}
```


### Pagination

The clickup API is a little inconsistent with pagination.  This client library will aim to document behavior as well as it can.  For example, use the `Page` attribute in `TaskQueryOptions` and call `TasksForList()` again.  

Unfortunately, the GET Tasks operation returns up to 100 tasks and the caller must know that the last page was reached only if there are less than 100.

### Client Library Progress

✅️ Implemented or partially implemented

🙅️ Not implemented


***

✅️ Attachments
✅️ Authorization (API Key supported "out of box." See `Authenticator` interface to implement OAuth, etc.)
✅️ Checklists
✅️ Comments
✅️ Folders
✅️ Goals
🙅️ Guests
✅️ Lists
🙅️ Members
🙅️ Shared Hierarchy
✅️ Spaces
✅️ Tags
✅️ Tasks
✅️ Task Templates
✅️ Team
✅️ Team
🙅️ Time Tracking
🙅️ Users
✅️ Views
✅️ Webhooks 


### Contributions

Open a PR with your fork or open an issue.  Open to help!

