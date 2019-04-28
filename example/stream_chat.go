package main

import (
	"fmt"
	"github.com/Dak425/dlive/pkg/api"
	"log"
)

func main() {
	c := api.Client{
		Endpoint: api.DefaultURL,
		Feeds:    make(map[string]api.Feed),
	}

	streamer := "dlive-21641280"

	s, err := c.StreamMessageFeed(streamer)

	if err != nil {
		log.Fatalf("unable to subscribe to %s's chat: %s\n", streamer, err)
	}

	count := 0

	for m := range s.Messages {
		count++
		fmt.Println(string(m))
		if count >= 5 {
			log.Println("closing subscription...")
			s.Close()
		}
	}
}
