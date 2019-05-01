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

func TestFeed_Unsubscribe(t *testing.T) {
	c := make(chan []byte)
	q := make(chan bool)

	s := Subscription{
		Key: "1",
		Messages: c,
	}

	f := Feed{
		quit: q,
		subscriptions: map[string]chan<- []byte{
			s.Key: c,
			"2": make(chan<- []byte),
		},
	}

	f.Unsubscribe(s)

	subCount := len(f.subscriptions)

	if subCount != 1 {
		t.Errorf("feed should have 0 active subscriptions, still has %d", subCount)
	}
}

func TestFeed_Close(t *testing.T) {
	q := make(chan bool)

	f := Feed{
		quit: q,
		subscriptions: map[string]chan<- []byte{
			"1": make(chan<- []byte),
			"2": make(chan<- []byte),
		},
	}

	go f.Close()

	sig := <- q

	if !sig {
		t.Error("feed should send true bool value to quit channel during close")
	}

	if len(f.subscriptions) > 0 {
		t.Error("feed should not have any subscriptions left over after close")
	}

	if f.quit != nil {
		t.Error("feed's quit channel should be unset during close")
	}
}