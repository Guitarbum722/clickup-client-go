package main

import (
	"fmt"
	"os"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {
	if os.Args[1] == "" || os.Args[2] == "" {
		panic("missing api key or team id")
	}
	apiKey := os.Args[1]
	folderID := os.Args[2]

	client := clickup.NewClient(&clickup.ClientOpts{
		APIToken: apiKey,
		Doer:     nil,
	})

	lists, err := client.ListsForFolder(folderID, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, folder := range lists.Lists {
		fmt.Printf("List ID: %s\nName: %s\n\n", folder.ID, folder.Name)
	}

	list, err := client.ListByID(os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Single Folder\n\nSpace ID: %s\nName: %s\n", list.ID, list.Name)
}
