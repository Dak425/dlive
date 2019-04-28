package main

import (
	"fmt"
	"log"

	"github.com/Dak425/dlive/pkg/api"
)

func main() {
	c := api.Client{
		Endpoint:          api.DefaultURL,
		WebsocketEndpoint: api.DefaultURLWebsocket,
		Feeds:             make(map[string]api.Feed),
	}

	args := api.StreamMessageFeedArgs{
		Streamer: "dlive-21641280",
	}

	s, err := c.StreamMessageFeed(args)

	if err != nil {
		log.Fatalf("unable to subscribe to %s's chat: %s\n", args.Streamer, err)
	}

	count := 0

	for m := range s.Messages {
		count++
		fmt.Println(string(m))
		if count >= 5 {
			log.Printf("closing subscription (%s)...\n", s.Key)
			s.Close()
		}
	}
}
