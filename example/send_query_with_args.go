package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dak425/dlive/pkg/api"
)

func main() {
	c := api.Client{
		Endpoint: api.DefaultURL,
		Feeds:    make(map[string]api.Feed),
	}

	a := api.LivestreamProfileWalletArgs{
		DisplayName: "TheHighlord",
		First:       20,
		IsLoggedIn:  false,
	}

	resp, err := c.LivestreamProfileWallet(a)

	if err != nil {
		log.Fatal(err)
	}

	prettyResponse, err := json.MarshalIndent(resp, "", " ")

	fmt.Println(string(prettyResponse))
}
