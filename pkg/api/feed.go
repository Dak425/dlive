package api

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
)

const connectionAckMessage = "connection_ack"
const connectionKeepAliveMessage = "ka"

// WebsocketFunc is the function used to setup the websocket used by a Feed
type WebsocketFunc func(request webSocketRequest) (*websocket.Conn, error)

// FeedMessage represents a GraphQL subscription message from DLive's API
type FeedMessage struct {
	MessageType string                 `json:"type"`    // The type of message sent from the API
	Payload     map[string]interface{} `json:"payload"` // The contents of the Response body
}

type Subscription struct {
	feed     *Feed         // The feed this subscription belongs to
	key      string        // The unique ID for this subscription for its feed
	Messages <-chan []byte // Channel that all new Response are written to
}

// Close removes this subscription from the feed
func (s Subscription) Close() {
	s.feed.Unsubscribe(s)
}

// Feed is a real-time data stream using a websocket
// When a feed receives data from its websocket, its writes that data to all its subscribers
type Feed struct {
	quit          chan<- bool              // The channel used to terminate the goroutine writing to the stream channel
	subscriptions map[string]chan<- []byte // A group of channels interested in this websocket connection's Feed
}

func (f *Feed) Active() bool {
	return f.quit == nil
}

// Publish will send the data given to all output channels it currently knows of
// Returns the length of data, times the number of output channels data was written to
// Returns an error if Feed has no output channels
func (f *Feed) Publish(p []byte) (int, error) {
	if len(p) == 0 {
		log.Println("publish -- payload has 0 len")
		return 0, nil
	}

	if len(f.subscriptions) == 0 {
		return 0, errors.New("no output channels to write to")
	}

	for _, c := range f.subscriptions {
		go func(data []byte, c chan<- []byte) {
			select {
			case c <- data:
				log.Println("Writing data to subscriber channel...")
				return
			default:
				log.Println("Waiting to write to subscriber channel...")
				time.Sleep(time.Second)
			}
		}(p, c)
	}

	return len(f.subscriptions) * len(p), nil
}

// Subscribe creates a new Subscription for the feed
func (f *Feed) Subscribe() (*Subscription, error) {
	var s Subscription

	id, err := uuid.NewV4()

	if err != nil {
		return &s, err
	}

	c := make(chan []byte)

	f.subscriptions[id.String()] = c

	s = Subscription{
		feed:     f,
		key:      id.String(),
		Messages: c,
	}

	return &s, nil
}

// Unsubscribe closes the subscription's channel and removes it from its map of subscribers
func (f *Feed) Unsubscribe(subscription Subscription) {
	if c, ok := f.subscriptions[subscription.key]; ok {
		close(c)
		delete(f.subscriptions, subscription.key)

		if len(f.subscriptions) == 0 {
			f.Close()
		}
	}
}

// Close sends a termination signal to consumer go routine, closes the termination signal channel, and closes all subscriptions
func (f *Feed) Close() {
	// Send termination signal to goroutine
	f.quit <- true

	// Close termination channel
	close(f.quit)

	// Unset websocket channels
	f.quit = nil

	// Close any existing output channels
	for k, v := range f.subscriptions {
		close(v)
		delete(f.subscriptions, k)
	}
}

// Start uses the provided request and websocketFunc to start a GraphQL websocket connection
// Returns an error if the feed already been started
func (f *Feed) Start(socketRequest webSocketRequest, websocketFunc WebsocketFunc) error {
	if f.quit != nil {
		return errors.New("feed has already been started")
	}

	// Setup websocket using provided func
	conn, err := websocketFunc(socketRequest)

	if err != nil {
		return err
	}

	q, s := f.Consume(conn)
	sq := make(chan bool)

	f.quit = sq

	go func(feed *Feed, quitFeed <-chan bool, quitConsume chan<- bool, streamConsume <-chan []byte) {
		for {
			select {
			case <-quitFeed:
				log.Println("termination signal received, terminating consumer...")
				quitConsume <- true
				close(quitConsume)
				return
			case m := <-streamConsume:
				if _, err := feed.Publish(m); err != nil {
					log.Println("error when publishing stream to subscribers: ", err)
				}
			default:
				log.Println("waiting on message from socket...")
				time.Sleep(time.Second)
			}
		}
	}(f, sq, q, s)

	return nil
}

// Consume uses the provided websocket to continuously read data from the socket and write it to the downstream channel
// Go routine will return when a termination signal is sent via the quit channel
func (f *Feed) Consume(conn *websocket.Conn) (chan<- bool, <-chan []byte) {
	q := make(chan bool)
	s := make(chan []byte)

	go func(socket *websocket.Conn, quit <-chan bool, stream chan<- []byte) {
		defer conn.Close()
		defer close(stream)

		for {
			_, m, err := conn.ReadMessage()

			if err != nil {
				log.Println("websocket read: ", err)
				return
			}

			var message FeedMessage

			err = json.Unmarshal(m, &message)

			if message.MessageType == connectionAckMessage || message.MessageType == connectionKeepAliveMessage {
				continue
			}

			select {
			case <-quit:
				log.Println("Termination signal received, closing websocket and channel...")
				return
			case stream <- m:
				log.Println("Writing stream to feed...")
			}
		}
	}(conn, q, s)

	return q, s
}
