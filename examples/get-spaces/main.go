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
	teamID := os.Args[2]

	client := clickup.NewClient(&clickup.ClientOpts{
		APIToken:   apiKey,
		HTTPClient: nil,
	})

	spaces, err := client.GetSpaces(teamID, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, space := range spaces.Spaces {
		fmt.Printf("Space ID: %s\nName: %s\n\n", space.ID, space.Name)
	}

	singleSpace, err := client.GetSpace("14865529")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Single Space\n\nSpace ID: %s\nName: %s\n", singleSpace.ID, singleSpace.Name)
}
