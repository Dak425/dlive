package api

import (
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WebsocketFunc func(request webSocketRequest) (*websocket.Conn, error)

type Response struct {
	MessageType string                 `json:"type"`
	Payload     map[string]interface{} `json:"payload"`
}

type Subscription struct {
	feed     *Feed
	key      string
	Messages <-chan []byte
}

func (s Subscription) Close() error {
	return s.feed.Unsubscribe(s)
}

type Feed struct {
	quit          chan<- bool              // The channel used to terminate the goroutine writing to the stream channel
	subscriptions map[string]chan<- []byte // A group of channels interested in this websocket connection's Feed
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

func (f *Feed) Unsubscribe(subscription Subscription) error {
	if c, ok := f.subscriptions[subscription.key]; ok {
		close(c)
		delete(f.subscriptions, subscription.key)

		if len(f.subscriptions) == 0 {
			if err := f.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *Feed) Close() error {
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

	return nil
}

func (f *Feed) Start(socketRequest webSocketRequest, websocketFunc WebsocketFunc) error {
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

			var message Response

			err = json.Unmarshal(m, &message)

			if message.MessageType == "connection_ack" || message.MessageType == "ka" {
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
