package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// DefaultURL is the url used by a API client if none is given
const DefaultURL = "https://graphigo.prd.dlive.tv/"
const DefaultURLWebsocket = "wss://graphigostream.prd.dlive.tv/"

// Client is used to send requests to DLive's API
type Client struct {
	Endpoint string          // The endpoint for DLive's API
	Auth     string          // An authorization token to send along with requests
	Feeds    map[string]Feed // Any active websocket streams the client is consuming
}

func (c *Client) Feed(key string) (*Feed, error) {
	if f, ok := c.Feeds[key]; ok {
		return &f, nil
	} else {
		return nil, errors.New(fmt.Sprintf("no active feed found with key (%s)", key))
	}
}

func (c *Client) FeedCount() int {
	return len(c.Feeds)
}

// GlobalInformation fetches language information about DLive
func (c *Client) GlobalInformation() (Response, error) {
	req := request{
		Query: GlobalInformationQuery(),
	}
	return c.sendQuery(req)
}

// Query Methods
func (c *Client) LivestreamPage(args LivestreamPageArgs) (interface{}, error) {
	req := request{
		Query: LivestreamPageQuery(),
		Vars:  args,
	}
	return c.sendQuery(req)
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
	return c.sendQuery(req)
}

func (c *Client) LivestreamProfileVideos(displayName string, first string) (interface{}, error) {
	req := request{
		Query: LivestreamProfileVideoQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
		},
	}
	return c.sendQuery(req)
}

func (c *Client) LivestreamProfileReplays(displayName string, first string) (interface{}, error) {
	req := request{
		Query: LivestreamProfileReplayQuery(),
		Vars: map[string]interface{}{
			"displayname": displayName,
			"first":       first,
		},
	}
	return c.sendQuery(req)
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
	return c.sendQuery(req)
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
	return c.sendQuery(req)
}

func (c *Client) LivestreamProfileWallet(args LivestreamProfileWalletArgs) (Response, error) {
	req := request{
		Query: LivestreamProfileWalletQuery(),
		Vars:  args,
	}
	return c.sendQuery(req)
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
	return c.sendQuery(req)
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
	return c.sendQuery(req)
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
	return c.sendQuery(req)
}

func (c *Client) AllowedActions(username string, streamer string) (interface{}, error) {
	req := request{
		Query: AllowedActionsQuery(),
		Vars: map[string]interface{}{
			"username": username,
			"streamer": streamer,
		},
	}
	return c.sendQuery(req)
}

// Mutation Methods
func (c *Client) SendStreamChat(input StreamChatInput) error {
	req := request{
		Query: SendStreamChatMessageMutation(),
		Vars: map[string]interface{}{
			"stream": input,
		},
	}
	return c.sendMutation(req)
}

// Subscription Methods
func (c *Client) StreamMessageFeed(streamer string) (*Subscription, error) {
	k := "StreamMessageFeed:" + streamer

	if f, ok := c.Feeds[k]; ok {
		if f.Active() {
			return f.Subscribe()
		}
	} else {
		c.Feeds[k] = Feed{
			subscriptions: make(map[string]chan<- []byte),
		}
	}

	f := c.Feeds[k]

	r := webSocketRequest{
		ID:   "1",
		Type: "start",
		Payload: request{
			Query: StreamMessageSubscription(),
			Vars: map[string]interface{}{
				"streamer": streamer,
			},
		},
	}

	err := f.Start(r, c.setupWebsocket)

	if err != nil {
		return nil, err
	}

	s, err := f.Subscribe()

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (c *Client) sendQuery(req request) (Response, error) {
	var body bytes.Buffer
	var data Response

	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return data, err
	}

	resp, err := http.Post(c.Endpoint, "application/json", &body)

	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return data, err
	}

	if err := json.NewDecoder(&buf).Decode(&data); err != nil {
		return data, err
	}

	if len(data.Errors) > 0 {
		return data, data.Errors[0]
	}

	return data, nil
}

func (c *Client) sendMutation(req request) error {
	var body bytes.Buffer

	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return err
	}

	resp, err := http.Post(c.Endpoint, "application/json", &body)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var buf bytes.Buffer

	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return err
	}

	var data Response

	if err := json.NewDecoder(&buf).Decode(&data); err != nil {
		return err
	}

	if len(data.Errors) > 0 {
		return data.Errors[0]
	}

	return nil
}

func (c *Client) setupWebsocket(req webSocketRequest) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(DefaultURLWebsocket, http.Header{
		"Sec-WebSocket-Protocol": []string{"graphql-ws"},
		"Sec-WebSocket-Version":  []string{"13"},
	})

	if err != nil {
		log.Println("Dial:", err)
		return nil, err
	}

	err = conn.WriteJSON(struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}{
		Type:    "connection_init",
		Payload: map[string]interface{}{},
	})

	if err != nil {
		log.Println("Connection Init:", err)
		return nil, err
	}

	err = conn.WriteJSON(req)

	if err != nil {
		log.Println("GraphQL Subscription Start:", err)
		return nil, err
	}

	return conn, nil
}
