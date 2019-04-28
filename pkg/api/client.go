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
	Endpoint          string          // The endpoint for DLive's API
	WebsocketEndpoint string          // The endpoint used for making websocket connections
	Auth              string          // An authorization token to send along with requests
	Feeds             map[string]Feed // Any active websocket streams the client is consuming
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
	req := Request{
		Query: GlobalInformationQuery(),
	}
	return c.Send(req)
}

// Query Methods
func (c *Client) LivestreamPage(args LivestreamPageArgs) (Response, error) {
	req := Request{
		Query: LivestreamPageQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) LivestreamChatRoomInfo(args LivestreamChatRoomInfoArgs) (Response, error) {
	req := Request{
		Query: LivestreamChatRoomInfoQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) LivestreamProfileVideos(args LivestreamProfileVideoArgs) (Response, error) {
	req := Request{
		Query: LivestreamProfileVideoQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) LivestreamProfileReplays(args LivestreamProfileReplayArgs) (Response, error) {
	req := Request{
		Query: LivestreamProfileReplayQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) LivestreamProfileFollowers(args LivestreamProfileFollowersArgs) (Response, error) {
	req := Request{
		Query: LivestreamProfileFollowersQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) LivestreamProfileFollowing(args LivestreamProfileFollowingArgs) (Response, error) {
	req := Request{
		Query: LivestreamProfileFollowingQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) LivestreamProfileWallet(args LivestreamProfileWalletArgs) (Response, error) {
	req := Request{
		Query: LivestreamProfileWalletQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) TopContributors(args TopContributorsArgs) (Response, error) {
	req := Request{
		Query: TopContributorsQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) StreamChatBannedUsers(args StreamChatBannedUsersArgs) (Response, error) {
	req := Request{
		Query: StreamChatBannedUsersQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) StreamChatModerators(args StreamChatModeratorsArgs) (Response, error) {
	req := Request{
		Query: StreamChatModeratorsQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

func (c *Client) AllowedActions(args AllowedActionsArgs) (Response, error) {
	req := Request{
		Query: AllowedActionsQuery(),
		Vars:  args,
	}
	return c.Send(req)
}

// Mutation Methods
func (c *Client) SendStreamChat(args SendStreamChatMessageArgs) (Response, error) {
	req := Request{
		Query: SendStreamChatMessageMutation(),
		Vars:  args,
	}
	return c.Send(req)
}

// Subscription Methods
func (c *Client) StreamMessageFeed(args StreamMessageFeedArgs) (*Subscription, error) {
	k := "StreamMessageFeed:" + args.Streamer

	if f, ok := c.Feeds[k]; ok {
		if f.Active() {
			return f.Subscribe()
		}
	} else {
		c.Feeds[k] = Feed{
			key: k,
			subscriptions: make(map[string]chan<- []byte),
		}
	}

	f := c.Feeds[k]

	r := WebSocketRequest{
		ID:   "1",
		Type: "start",
		Payload: Request{
			Query: StreamMessageSubscription(),
			Vars:  args,
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

// Send takes the provided request, sends it to the DLive API endpoint, then returns the decoded JSON response
func (c *Client) Send(req Request) (Response, error) {
	client := http.Client{}
	var body bytes.Buffer
	var data Response

	if err := json.NewEncoder(&body).Encode(req); err != nil {
		return data, err
	}

	r, err := http.NewRequest(http.MethodPost, c.Endpoint, &body)

	if err != nil {
		return data, err
	}

	// API Client was given an auth token, set it to request header
	if c.Auth != "" {
		r.Header.Set("authorization", c.Auth)
	}

	r.Header.Set("content-type", "application/json")

	resp, err := client.Do(r)

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

// setupWebsocket is the default func used to setup a websocket connection for a feed
func (c *Client) setupWebsocket(req WebSocketRequest) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.WebsocketEndpoint, http.Header{
		"Sec-WebSocket-Protocol": []string{"graphql-ws"},
		"Sec-WebSocket-Version":  []string{"13"},
	})

	if err != nil {
		log.Println("Dial:", err)
		return nil, err
	}

	init := WebSocketRequest{
		Type:    connectionInit,
		Payload: Request{},
	}

	err = conn.WriteJSON(init)

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
