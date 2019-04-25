package api

import (
	"errors"
	"github.com/nu7hatch/gouuid"
)

type FeedSetupFunc func(chan<- []byte, <-chan bool, webSocketRequest)

type Feed struct {
	FeedType  string // The type of subscription this Feed utilizes
	SetupFunc FeedSetupFunc
	quit      chan<- bool              // The channel used to terminate the goroutine writing to the input channel
	input     <-chan []byte            // The channel used to read data from goroutine that is streaming data from a remote source
	outputs   map[string]chan<- []byte // A group of channels interested in this websocket connection's Feed
}

// publish will send the data given to all output channels it currently knows of
// Returns the length of data, times the number of output channels data was written to
// Returns an error if Feed has no output channels
func (f *Feed) publish(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	if len(f.outputs) == 0 {
		return 0, errors.New("no output channels to write to")
	}

	for _, c := range f.outputs {
		go func(data []byte, c chan<- []byte) {
			c <- data
			return
		}(p, c)
	}

	return len(f.outputs) * len(p), nil
}

func (f *Feed) Subscribe(output chan<- []byte) (string, error) {
	id, err := uuid.NewV4()

	if err != nil {
		return "", nil
	}

	f.outputs[id.String()] = output

	return id.String(), nil
}

func (f *Feed) Unsubscribe(id string) error {
	if c, ok := f.outputs[id]; ok {
		close(c)
		delete(f.outputs, id)

		if len(f.outputs) == 0 {
			 if err := f.Close(); err != nil {
			 	return err
			 }
		}

		return nil
	}

	return errors.New("output has already been closed")
}

func (f *Feed) Close() error {
	// Send termination signal to goroutine
	f.quit <- true

	// Close termination channel
	close(f.quit)

	// Close any existing output channels
	for k, v := range f.outputs {
		close(v)
		delete(f.outputs, k)
	}

	return nil
}
