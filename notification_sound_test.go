package main

import (
	"testing"
	"time"
)

// The reason this test exist is to help with development (e.g.: test if sound
// is working). It is useless outside of development purposes and needs a proper
// sound environment to work, and this is the reason why it is not run in CI.
func TestPlayNotificationSound(t *testing.T) {
	err := initBeep()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	<-playNotificationSound()
	// it takes a while until the sound finishes playing
	time.Sleep(5 * time.Second)
}
