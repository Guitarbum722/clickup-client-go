package main

import (
	"fmt"
	"os"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {
	client := clickup.NewClient(&clickup.ClientOpts{
		Doer: nil,
		Authenticator: &clickup.APITokenAuthenticator{
			APIToken: os.Args[1],
		},
	})

	f, err := os.Open("examples/create-attachment/attachment_sample.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	params := clickup.AttachmentParams{
		FileName: f.Name(),
		Reader:   f,
	}

	attachment, err := client.CreateTaskAttachment(os.Args[2], os.Args[3], true, &params)
	if err != nil {
		panic(err)
	}

	fmt.Println(attachment)
}
