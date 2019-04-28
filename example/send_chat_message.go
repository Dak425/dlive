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
		Auth: "ADD AUTH TOKEN HERE",
	}

	args := api.SendStreamChatMessageArgs{
		Input: api.SendStreamChatMessageInput{
			Message: "MESSAGE HERE",
			RoomRole: "ROOM ROLE CONSTANT",
			Streamer: "STREAMER ID HERE (Lino Account ID)",
			Subscribing: true,
		},
	}

	resp, err := c.SendStreamChat(args)

	if err != nil {
		log.Fatal(err)
	}

	prettyResponse, err := json.MarshalIndent(resp, "", " ")

	fmt.Println(string(prettyResponse))
}
