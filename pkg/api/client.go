package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

// DefaultURL is the url used by a API client if none is given
const DefaultURL = "https://graphigo.prd.dlive.tv/"
const DefaultURLWebsocket = "wss://graphigostream.prd.dlive.tv/"

// Client is used to send requests to DLive's API
type Client struct {
	URL  string // The URL for DLive's API
	Auth string // An authorization token to send along with requests
}

type responseError struct {
	Message string
}

func (re responseError) Error() string {
	return "GraphQL API Error: " + re.Message
}

type response struct {
	Data   interface{}
	Errors []responseError
}

type request struct {
	Query         string                 `json:"query"`
	Vars          map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName"`
}

type webSocketRequest struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Payload request `json:"payload"`
}

// GlobalInformation fetches language information about DLive
func (c *Client) GlobalInformation() (interface{}, error) {
	req := request{
		Query: GlobalInformationQuery(),
		Vars:  map[string]interface{}{},
	}
	return c.run(req)
}

// Query Methods
func (c *Client) LivestreamPage(displayName string, add bool, isLoggedIn bool) (interface{}, error) {
	req := request{
		Query: LivestreamPageQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"add":         add,
			"isLoggedIn":  isLoggedIn,
		},
	}
	return c.run(req)
}

func (c *Client) LivestreamChatRoomInfo(displayName string, isLoggedIn bool, limit string) (interface{}, error) {
	req := request{
		Query: LivestreamChatRoomInfoQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"isLoggedIn":  isLoggedIn,
			"limit":       limit,
		},
	}
	return c.run(req)
}

func (c *Client) LivestreamProfileVideos(displayName string, first string) (interface{}, error) {
	req := request{
		Query: LivestreamProfileVideoQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
		},
	}
	return c.run(req)
}

func (c *Client) LivestreamProfileReplays(displayName string, first string) (interface{}, error) {
	req := request{
		Query: LivestreamProfileReplayQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
		},
	}
	return c.run(req)
}

func (c *Client) LivestreamProfileFollowers(displayName string, sortBy string, first string, isLoggedIn bool) (interface{}, error) {
	req := request{
		Query: LivestreamProfileFollowersQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"sortBy":      sortBy,
			"first":       first,
			"isLoggedIn":  false,
		},
	}
	return c.run(req)
}

func (c *Client) LivestreamProfileFollowing(displayName string, sortBy string, first string, isLoggedIn bool) (interface{}, error) {
	req := request{
		Query: LivestreamProfileFollowingQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"sortBy":      sortBy,
			"first":       first,
			"isLoggedIn":  false,
		},
	}
	return c.run(req)
}

func (c *Client) LivestreamProfileWallet(displayName string, first string, isLoggedIn bool) (interface{}, error) {
	req := request{
		Query: LivestreamProfileWalletQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
			"isLoggedIn":  false,
		},
	}
	return c.run(req)
}

func (c *Client) TopContributors(displayName string, rule string, first string, queryStream bool) (interface{}, error) {
	req := request{
		Query: TopContributorsQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"rule":        rule,
			"first":       first,
			"queryStream": false,
		},
	}
	return c.run(req)
}

func (c *Client) StreamChatBannedUsers(displayName string, first string, search string) (interface{}, error) {
	req := request{
		Query: StreamChatBannedUsersQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
			"search":      search,
		},
	}
	return c.run(req)
}

func (c *Client) StreamChatModerators(displayName string, first string, search string) (interface{}, error) {
	req := request{
		Query: StreamChatModeratorsQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
			"search":      search,
		},
	}
	return c.run(req)
}

func (c *Client) AllowedActions(username string, streamer string) (interface{}, error) {
	req := request{
		Query: AllowedActionsQuery(),
		Vars: map[string]interface{}{
			"username": username,
			"streamer": streamer,
		},
	}
	return c.run(req)
}

// Mutation Methods

// Subscription Methods
func (c *Client) StreamMessageFeed(streamer string, messages chan []byte) error {
	req := webSocketRequest{
		ID:   "1",
		Type: "start",
		Payload: request{
			Query: StreamMessageSubscription(),
			Vars: map[string]interface{}{
				"streamer": "TheHighlord",
			},
			OperationName: "StreamMessageSubscription",
		},
	}
	return c.connectWebsocket(req, messages)
}

func (c *Client) run(req request) (interface{}, error) {
	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return "", err
	}

	resp, err := http.Post(c.URL, "application/json", &body)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return "", err
	}

	var data response

	if err := json.NewDecoder(&buf).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Errors) > 0 {
		return "", data.Errors[0]
	}

	return data.Data, nil
}

func (c *Client) connectWebsocket(req webSocketRequest, messages chan []byte) error {
	go func(req webSocketRequest, messages chan []byte) {
		conn, _, err := websocket.DefaultDialer.Dial(DefaultURLWebsocket, http.Header{
			"Sec-WebSocket-Protocol": []string{"graphql-ws"},
			"Sec-WebSocket-Version": []string{"13"},
		})

		if err != nil {
			log.Fatal("Dial:", err)
			return
		}

		err = conn.WriteJSON(struct {
			Type    string      `json:"type"`
			Payload interface{} `json:"payload"`
		}{
			Type:    "connection_init",
			Payload: map[string]interface{}{},
		})

		if err != nil {
			log.Fatal("Connection Init:", err)
			return
		}

		err = conn.WriteJSON(req)

		if err != nil {
			log.Fatal("GraphQL Subscription Start:", err)
			return
		}

		defer conn.Close()

		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				log.Fatal("Read:", err)
				return
			}
			messages <- m
		}
	}(req, messages)
	return nil
}
