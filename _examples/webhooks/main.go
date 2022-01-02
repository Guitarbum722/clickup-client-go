// Copyright (c) 2022, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {
	if os.Args[1] == "" || os.Args[2] == "" {
		panic("missing api key or team id")
	}
	apiKey := os.Args[1]
	workspaceID := os.Args[2]
	listID := os.Args[3]

	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: apiKey,
		},
	})

	newWebhook, err := client.CreateWebhook(context.Background(), workspaceID, &clickup.CreateWebhookRequest{
		Endpoint: "https://your-webhook.site/myhook",
		Events: []clickup.WebhookEvent{
			clickup.EventTaskUpdated,
		},
		ListID: listID,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("New Webhook ID: ", newWebhook.ID)
	fmt.Println("Health: ", newWebhook.Webhook.Health.FailCount, newWebhook.Webhook.Health.Status)

	webhooks, err := client.WebhooksFor(context.Background(), workspaceID)
	if err != nil {
		panic(err)
	}

	fmt.Println("Existing webhooks:")
	for _, v := range webhooks.Webhooks {
		fmt.Println("ID:", v.ID)
		fmt.Println("Status:", v.Health.Status, v.Health.FailCount)

		fmt.Println("deleting webhook for:", v.ID)
		if err := client.DeleteWebhook(context.Background(), v.ID); err != nil {
			panic(err)
		}
	}
}
