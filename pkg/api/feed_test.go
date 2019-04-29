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

