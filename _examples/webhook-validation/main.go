package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Guitarbum722/clickup-client-go"
)

func main() {

	secret := "imiO3dJZfIlyykAG"
	data := `{"event":"taskUpdated"}`

	req := &http.Request{
		Header: http.Header{
			"X-Signature": []string{"2831500d379c7e90a2c8b3ff55dec81a42889b8a91f6b97f8513d98ebb6b23bf"},
		},
		Body: io.NopCloser(strings.NewReader(data)),
	}

	result, err := clickup.VerifyWebhookSignature(req, secret)
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Valid(), result.SignatureFromClickup(), result.SignatureGenerated())
}
