package api

import (
	"testing"
)

func TestFeed_String(t *testing.T) {
	f := Feed{
		key: "test",
		subscriptions: make(map[string]chan<- []byte),
	}

	expected := "feed(test) -- subscription count(0) -- active (false)"

	result := f.String()

	if result != expected {
		t.Errorf("returned --%s--, should have been --%s--", result, expected)
	}
}

func TestFeed_Active(t *testing.T) {
	f := Feed{
		quit: make(chan<- bool),
	}

	active := f.Active()

	if !active {
		t.Errorf("returned %t, should be %t", active, true)
	}
}

func TestFeed_Publish(t *testing.T) {
	chan1 := make(chan []byte)
	chan2 := make(chan []byte)

	f := Feed{
		subscriptions: map[string]chan<- []byte{
			"one": chan1,
			"two": chan2,
		},
	}

	p := []byte{'h','e','l','l','o'}

	res, err := f.Publish(p)

	if err != nil {
		t.Error("got error: ", err)
	}

	expected := len(p) * 2

	if res != expected {
		t.Errorf("expected return length of %d, got %d", expected, res)
	}
}

func TestFeed_Subscribe(t *testing.T) {
	f := Feed{
		subscriptions: make(map[string]chan<- []byte),
	}

	s, err := f.Subscribe()

	if err != nil {
		t.Error("error when subscribing to feed: ", err)
	}

	if s.feed != &f {
		t.Error("subscription has incorrect feed reference")
	}
}