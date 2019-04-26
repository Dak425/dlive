package api

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WebsocketFunc func(request webSocketRequest) (*websocket.Conn, error)

type Subscription struct {
	key      string
	Messages <-chan []byte
}

type Feed struct {
	quit          chan<- bool              // The channel used to terminate the goroutine writing to the stream channel
	stream        <-chan []byte            // The channel used to read data from goroutine that is streaming data from a remote source
	subscriptions map[string]chan<- []byte // A group of channels interested in this websocket connection's Feed
}

// Publish will send the data given to all output channels it currently knows of
// Returns the length of data, times the number of output channels data was written to
// Returns an error if Feed has no output channels
func (f *Feed) Publish(p []byte) (int, error) {
	if len(p) == 0 {
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
	f.stream = nil

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

	go f.Watch(conn)

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
			log.Println(string(m))

			if err != nil {
				log.Println("websocket read: ", err)
				return
			}

			select {
			case <-quit:
				log.Println("Termination signal received, ending goroutine...")
				return
			case stream <- m:
				log.Println("Writing stream to feed...")
			default:
				log.Println("Waiting on message...")
				time.Sleep(time.Second)
			}
		}
	}(conn, q, s)

	return q, s
}

func (f *Feed) Watch(conn *websocket.Conn) {
	f.quit = make(chan bool)
	f.stream = make(chan []byte)

	q := make(chan bool)
	s := make(chan []byte)

	f.quit = q
	f.stream = s

	go func(quit <-chan bool, stream chan<- []byte, conn2 *websocket.Conn) {
		cq, cs:= f.Consume(conn2)

		for {
			select {
			case <-quit:
				log.Println("feed watcher terminating, sending termination to websocket consumer...")
				cq <- true
				close(cq)
				close(stream)
			case m := <- cs:
				log.Println("sending data from consumer to feed")
				if _, err := f.Publish(m); err != nil {
					log.Println("error when publishing stream to subscribers", err)
				}
			default:
				log.Println("waiting on consumer data...")
				time.Sleep(time.Second)
			}
		}
	}(q, s, conn)
}