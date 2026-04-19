package utils

import (
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
	ch := GlobalBroadcaster.Subscribe()
	defer GlobalBroadcaster.Unsubscribe(ch)

	msg := "test message"
	GlobalBroadcaster.Broadcast(msg)

	select {
	case got := <-ch:
		if got != msg {
			t.Errorf("got %v, want %v", got, msg)
		}
	case <-time.After(1 * time.Second):
		t.Error("timeout waiting for broadcast")
	}
}
