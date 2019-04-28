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

	rc := make(chan api.Response)

	go func(responseChan chan<- api.Response, client api.Client) {
		resp, err := c.GlobalInformation()

		if err != nil {
			log.Fatalf("error when sending query, %s\n", err)
		}

		responseChan <- resp

		close(responseChan)

		return
	}(rc, c)

	resp := <- rc

	prettyResponse, err := json.MarshalIndent(resp, "", " ")

	if err != nil {
		log.Fatal("could not pretty print response")
	}

	fmt.Println(string(prettyResponse))
}
